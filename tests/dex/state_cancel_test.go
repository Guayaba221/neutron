package dex_state_test

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dextypes "github.com/neutron-org/neutron/v4/x/dex/types"
	"strconv"
	"testing"
	"time"
)

type cancelLimitOrderTestParams struct {
	// State Conditions
	SharedParams
	ExistingTokenAHolders string
	Filled                int
	WithdrawnCreator      bool
	WithdrawnOneOther     bool
	Expired               bool
	OrderType             int32 // JIT, GTT, GTC
}

func (p cancelLimitOrderTestParams) printTestInfo(t *testing.T) {
	t.Logf(`
		Existing Shareholders: %s
		Filled: %v
		WithdrawnCreator: %v
		WithdrawnOneOther: %t
		Expired: %t
		OrderType: %v`,
		p.ExistingTokenAHolders,
		p.Filled,
		p.WithdrawnCreator,
		p.WithdrawnOneOther,
		p.Expired,
		p.OrderType,
	)
}

func hydrateCancelLoTestCase(params map[string]string) cancelLimitOrderTestParams {
	return cancelLimitOrderTestParams{
		ExistingTokenAHolders: params["ExistingTokenAHolders"],
		Filled:                parseInt(params["Filled"]),
		WithdrawnCreator:      parseBool(params["WithdrawnCreator"]),
		WithdrawnOneOther:     parseBool(params["WithdrawnOneOther"]),
		Expired:               parseBool(params["Expired"]),
		OrderType:             dextypes.LimitOrderType_value[params["OrderType"]],
	}
}

func (s *DexStateTestSuite) setupCancelTest(params cancelLimitOrderTestParams) *dextypes.LimitOrderTranche {
	coinA := sdk.NewCoin(params.PairID.Token0, BaseTokenAmountInt)
	coinB := sdk.NewCoin(params.PairID.Token1, BaseTokenAmountInt)
	s.FundAcc(s.creator, sdk.NewCoins(coinA))
	var expTime *time.Time
	if params.OrderType == int32(dextypes.LimitOrderType_GOOD_TIL_TIME) {
		t := time.Now()
		expTime = &t
	}
	res := s.makePlaceLOSuccess(s.creator, coinA, coinB.Denom, DefaultSellPrice, dextypes.LimitOrderType(params.OrderType), expTime)

	totalDeposited := BaseTokenAmountInt
	if params.ExistingTokenAHolders == OneOtherAndCreatorLO {
		totalDeposited = totalDeposited.MulRaw(2)
		s.FundAcc(s.alice, sdk.NewCoins(coinA))
		s.makePlaceLOSuccess(s.alice, coinA, coinB.Denom, DefaultSellPrice, dextypes.LimitOrderType(params.OrderType), expTime)
	}

	if params.Filled > 0 {
		s.FundAcc(s.bob, sdk.NewCoins(coinB).MulInt(math.NewInt(10)))
		fillAmount := totalDeposited.MulRaw(int64(params.Filled)).QuoRaw(100)
		_, err := s.makePlaceTakerLO(s.bob, coinB, coinA.Denom, DefaultBuyPriceTaker, dextypes.LimitOrderType_IMMEDIATE_OR_CANCEL, &fillAmount)
		s.NoError(err)
	}

	if params.WithdrawnCreator {
		s.makeWithdrawFilledSuccess(s.creator, res.TrancheKey)
	}

	if params.WithdrawnOneOther {
		s.makeWithdrawFilledSuccess(s.alice, res.TrancheKey)
	}

	if params.Expired {
		s.App.DexKeeper.PurgeExpiredLimitOrders(s.Ctx, time.Now())
	}
	tick, err := dextypes.CalcTickIndexFromPrice(DefaultStartPrice)
	s.NoError(err)
	tranches, _ := s.App.DexKeeper.LimitOrderTrancheAll(s.Ctx, &dextypes.QueryAllLimitOrderTrancheRequest{
		PairId:     params.PairID.CanonicalString(),
		TokenIn:    params.PairID.Token0,
		Pagination: nil,
	})
	fmt.Println(tranches)
	req := dextypes.QueryGetLimitOrderTrancheRequest{
		PairId:     params.PairID.CanonicalString(),
		TickIndex:  -1 * tick,
		TokenIn:    params.PairID.Token0,
		TrancheKey: res.TrancheKey,
	}
	fmt.Println(req)
	tranchResp, err := s.App.DexKeeper.LimitOrderTranche(s.Ctx, &req)
	s.NoError(err)

	return tranchResp.LimitOrderTranche
}

func hydrateAllCancelLoTestCases(paramsList []map[string]string) []cancelLimitOrderTestParams {
	allTCs := make([]cancelLimitOrderTestParams, 0)
	for i, paramsRaw := range paramsList {
		pairID := generatePairID(i)
		tc := hydrateCancelLoTestCase(paramsRaw)
		tc.PairID = pairID
		allTCs = append(allTCs, tc)
	}

	return removeRedundantCancelLOTests(allTCs)
}

func removeRedundantCancelLOTests(params []cancelLimitOrderTestParams) []cancelLimitOrderTestParams {
	newParams := make([]cancelLimitOrderTestParams, 0)
	for _, p := range params {
		// it's impossible to withdraw 0 filled
		// error checks is not in a scope of the testcase (see withdraw filled test)
		if p.Filled == 0 && (p.WithdrawnOneOther || p.WithdrawnCreator) {
			continue
		}
		if p.Expired && p.OrderType == int32(dextypes.LimitOrderType_GOOD_TIL_CANCELLED) {
			continue
		}
		if p.WithdrawnOneOther && p.ExistingTokenAHolders == CreatorLO {
			continue
		}
		if p.ExistingTokenAHolders == OneOtherAndCreatorLO && p.OrderType != int32(dextypes.LimitOrderType_GOOD_TIL_CANCELLED) {
			// user tranches combined into tranches only for LimitOrderType_GOOD_TIL_CANCELLED
			// it does not make any sense to create two tranches
			continue
		}
		newParams = append(newParams, p)
	}
	return newParams
}

func (s *DexStateTestSuite) handleCancelErrors(params cancelLimitOrderTestParams, err error) {
	if params.Expired {
		if errors.Is(dextypes.ErrActiveLimitOrderNotFound, err) {
			s.T().Skip()
		}
	}
	s.NoError(err)
}

func TestCancel(t *testing.T) {
	testParams := []testParams{
		{field: "ExistingTokenAHolders", states: []string{CreatorLO, OneOtherAndCreatorLO}},
		{field: "Filled", states: []string{ZeroPCT, FiftyPCT, HundredPct}},
		{field: "WithdrawnCreator", states: []string{True, False}},
		{field: "WithdrawnOneOther", states: []string{True, False}},
		{field: "OrderType", states: []string{
			dextypes.LimitOrderType_name[int32(dextypes.LimitOrderType_GOOD_TIL_CANCELLED)],
			dextypes.LimitOrderType_name[int32(dextypes.LimitOrderType_GOOD_TIL_TIME)],
			dextypes.LimitOrderType_name[int32(dextypes.LimitOrderType_JUST_IN_TIME)],
		}},
		{field: "Expired", states: []string{True, False}},
	}
	testCasesRaw := generatePermutations(testParams)
	testCases := hydrateAllCancelLoTestCases(testCasesRaw)

	s := new(DexStateTestSuite)
	s.SetT(t)
	s.SetupTest()

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s.SetT(t)
			tc.printTestInfo(t)

			initialTrancheKey := s.setupCancelTest(tc)
			fmt.Println(initialTrancheKey)

			resp, err := s.makeCancel(s.creator, initialTrancheKey.Key.TrancheKey)
			s.handleCancelErrors(tc, err)
			fmt.Println("resp", resp)
			fmt.Println("err", err)

		})
	}
}
