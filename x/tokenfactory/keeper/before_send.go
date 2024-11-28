package keeper

import (
	"context"
	"encoding/json"
	"strings"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/neutron-org/neutron/v5/x/tokenfactory/types"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
)

func (k Keeper) setBeforeSendHook(ctx sdk.Context, denom, contractAddr string) error {
	// verify that denom is an x/tokenfactory denom
	_, _, err := types.DeconstructDenom(denom)
	if err != nil {
		return err
	}

	store := k.GetDenomPrefixStore(ctx, denom)

	// delete the store for denom prefix store when cosmwasm address is nil
	if contractAddr == "" {
		store.Delete([]byte(types.BeforeSendHookAddressPrefixKey))
		return nil
	}

	_, err = sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return err
	}

	store.Set([]byte(types.BeforeSendHookAddressPrefixKey), []byte(contractAddr))

	return nil
}

func (k Keeper) GetBeforeSendHook(ctx context.Context, denom string) string {
	store := k.GetDenomPrefixStore(ctx, denom)

	bz := store.Get([]byte(types.BeforeSendHookAddressPrefixKey))
	if bz == nil {
		return ""
	}

	return string(bz)
}

func CWCoinFromSDKCoin(in sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom:  in.GetDenom(),
		Amount: in.Amount.String(),
	}
}

// Hooks wrapper struct for bank keeper
type Hooks struct {
	k Keeper
}

var _ types.BankHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (k Keeper) AssertIsHookWhitelisted(ctx sdk.Context, denom string, contractAddress sdk.AccAddress) error {
	contractInfo := k.contractKeeper.GetContractInfo(ctx, contractAddress)
	if contractInfo == nil {
		return types.ErrBeforeSendHookNotWhitelisted.Wrapf("contract with address (%s) does not exist", contractAddress.String())
	}
	codeID := contractInfo.CodeID
	whitelistedHooks := k.GetParams(ctx).WhitelistedHooks
	denomCreator, _, err := types.DeconstructDenom(denom)
	if err != nil {
		return types.ErrBeforeSendHookNotWhitelisted.Wrapf("invalid denom: %s", denom)
	}

	for _, hook := range whitelistedHooks {
		if hook.CodeID == codeID && hook.DenomCreator == denomCreator {
			return nil
		}
	}

	return types.ErrBeforeSendHookNotWhitelisted.Wrapf("no whitelist for contract with codeID (%d) and denomCreator (%s) ", codeID, denomCreator)
}

// TrackBeforeSend calls the before send listener contract suppresses any errors
func (h Hooks) TrackBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	_ = h.k.callBeforeSendListener(ctx, from, to, amount, false)
}

// TrackBeforeSend calls the before send listener contract returns any errors
func (h Hooks) BlockBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	return h.k.callBeforeSendListener(ctx, from, to, amount, true)
}

// callBeforeSendListener iterates over each coin and sends corresponding sudo msg to the contract address stored in state.
// If blockBeforeSend is true, sudoMsg wraps BlockBeforeSendMsg, otherwise sudoMsg wraps TrackBeforeSendMsg.
// Note that we gas meter trackBeforeSend to prevent infinite contract calls.
// CONTRACT: this should not be called in beginBlock or endBlock since out of gas will cause this method to panic.
func (k Keeper) callBeforeSendListener(goCtx context.Context, from, to sdk.AccAddress, amount sdk.Coins, blockBeforeSend bool) (err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	defer func() {
		if r := recover(); r != nil {
			err = types.ErrTrackBeforeSendOutOfGas
		}
	}()

	for _, coin := range amount {
		contractAddr := k.GetBeforeSendHook(goCtx, coin.Denom)
		if contractAddr != "" {
			cwAddr, err := sdk.AccAddressFromBech32(contractAddr)
			if err != nil {
				return err
			}

			// Do not invoke hook if denom is not whitelisted
			// NOTE: hooks must already be whitelisted before they can be added, so under normal operation this check should never fail.
			// It is here as an emergency override if we want to shutoff a hook. We do not return the error because once it is removed from the whitelist
			// a hook should not be able to block a send.
			if err := k.AssertIsHookWhitelisted(ctx, coin.Denom, cwAddr); err != nil {
				ctx.Logger().Error(
					"Skipped hook execution due to missing whitelist",
					"err", err,
					"denom", coin.Denom,
					"contract", cwAddr.String(),
				)
				continue
			}

			var msgBz []byte

			// get msgBz, either BlockBeforeSend or TrackBeforeSend
			if blockBeforeSend {
				msg := types.BlockBeforeSendSudoMsg{
					BlockBeforeSend: types.BlockBeforeSendMsg{
						From:   from.String(),
						To:     to.String(),
						Amount: CWCoinFromSDKCoin(coin),
					},
				}
				msgBz, err = json.Marshal(msg)
			} else {
				msg := types.TrackBeforeSendSudoMsg{
					TrackBeforeSend: types.TrackBeforeSendMsg{
						From:   from.String(),
						To:     to.String(),
						Amount: CWCoinFromSDKCoin(coin),
					},
				}
				msgBz, err = json.Marshal(msg)
			}
			if err != nil {
				return err
			}

			em := sdk.NewEventManager()

			childCtx := ctx.WithGasMeter(storetypes.NewGasMeter(types.SudoHookGasLimit))
			_, err = k.contractKeeper.Sudo(childCtx.WithEventManager(em), cwAddr, msgBz)
			if err != nil {
				k.Logger(ctx).Debug("TokenFactory hooks: failed to sudo",
					"error", err, "contract_address", cwAddr)

				ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventFailedSudoCall,
					[]sdk.Attribute{sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
						sdk.NewAttribute(types.AttributeBeforeSendHookAddress, cwAddr.String()),
						sdk.NewAttribute(types.AttributeSudoErrorText, err.Error()),
					}...))

				// don't block or prevent transfer if there is no contract for some reason
				// It's not quite possible, but it's good to have such check just in case
				if strings.Contains(err.Error(), "no such contract") {
					return nil
				}

				// TF hooks should not block or prevent transfers from module accounts
				if k.isModuleAccount(ctx, from) {
					return nil
				}

				return errorsmod.Wrapf(err, "failed to call send hook for denom %s", coin.Denom)
			}

			// consume gas used for calling contract to the parent ctx
			ctx.GasMeter().ConsumeGas(childCtx.GasMeter().GasConsumed(), "track before send gas")
		}
	}
	return nil
}
