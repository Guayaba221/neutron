package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	keeper2 "github.com/cosmos/interchain-security/v6/testutil/keeper"

	"github.com/neutron-org/neutron/v5/utils"

	"github.com/neutron-org/neutron/v5/app/config"

	"cosmossdk.io/log"
	cometbfttypes "github.com/cometbft/cometbft/abci/types"
	db2 "github.com/cosmos/cosmos-db"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	consumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	icssimapp "github.com/cosmos/interchain-security/v6/testutil/ibc_testing"
	"github.com/stretchr/testify/suite"

	appparams "github.com/neutron-org/neutron/v5/app/params"
	tokenfactorytypes "github.com/neutron-org/neutron/v5/x/tokenfactory/types"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	appProvider "github.com/cosmos/interchain-security/v6/app/provider"
	e2e "github.com/cosmos/interchain-security/v6/testutil/integration"

	"github.com/neutron-org/neutron/v5/app"
	ictxstypes "github.com/neutron-org/neutron/v5/x/interchaintxs/types"

	providertypes "github.com/cosmos/interchain-security/v6/x/ccv/provider/types"
	ccv "github.com/cosmos/interchain-security/v6/x/ccv/types"

	cmttypes "github.com/cometbft/cometbft/types"
)

var (
	// TestOwnerAddress defines a reusable bech32 address for testing purposes
	TestOwnerAddress = "neutron17dtl0mjt3t77kpuhg2edqzjpszulwhgzcdvagh"

	TestInterchainID = "owner_id"

	// provider-consumer connection takes connection-0
	ConnectionOne = "connection-1"

	// TestVersion defines a reusable interchainaccounts version string for testing purposes
	TestVersion = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: ConnectionOne,
		HostConnectionId:       ConnectionOne,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))
)

func init() {
	// ibctesting.DefaultTestingAppInit = SetupTestingApp()
	config.GetDefaultConfig()
	// Disable cache since enabled cache triggers test errors when `AccAddress.String()`
	// gets called before setting neutron bech32 prefix
	sdk.SetAddrCacheEnabled(false)
}

type IBCConnectionTestSuite struct {
	suite.Suite
	Coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	ChainProvider *ibctesting.TestChain
	ChainA        *ibctesting.TestChain
	ChainB        *ibctesting.TestChain

	ProviderApp e2e.ProviderApp
	ChainAApp   e2e.ConsumerApp
	ChainBApp   e2e.ConsumerApp

	CCVPathA     *ibctesting.Path
	CCVPathB     *ibctesting.Path
	Path         *ibctesting.Path
	TransferPath *ibctesting.Path
}

func GetTestConsumerAdditionProp(chain *ibctesting.TestChain) *providertypes.ConsumerAdditionProposal { //nolint:staticcheck
	prop := providertypes.NewConsumerAdditionProposal(
		chain.ChainID,
		"description",
		chain.ChainID,
		chain.LastHeader.GetHeight().(clienttypes.Height),
		[]byte("gen_hash"),
		[]byte("bin_hash"),
		time.Now(),
		ccv.DefaultConsumerRedistributeFrac,
		ccv.DefaultBlocksPerDistributionTransmission,
		"channel-0",
		ccv.DefaultHistoricalEntries,
		ccv.DefaultCCVTimeoutPeriod,
		ccv.DefaultTransferTimeoutPeriod,
		ccv.DefaultConsumerUnbondingPeriod,
		95,
		100,
		0,
		nil,
		nil,
		0,
		true,
	).(*providertypes.ConsumerAdditionProposal) //nolint:staticcheck

	return prop
}

func (suite *IBCConnectionTestSuite) SetupTest() {
	// we need to redefine this variable to make tests work cause we use untrn as default bond denom in neutron
	sdk.DefaultBondDenom = appparams.DefaultDenom

	suite.Coordinator = NewProviderConsumerCoordinator(suite.T())
	suite.ChainProvider = suite.Coordinator.GetChain(ibctesting.GetChainID(1))
	suite.ChainA = suite.Coordinator.GetChain(ibctesting.GetChainID(2))
	suite.ChainB = suite.Coordinator.GetChain(ibctesting.GetChainID(3))
	suite.ProviderApp = suite.ChainProvider.App.(*appProvider.App)
	suite.ChainAApp = suite.ChainA.App.(*app.App)
	suite.ChainBApp = suite.ChainB.App.(*app.App)

	providerKeeper := suite.ProviderApp.GetProviderKeeper()
	consumerKeeperA := suite.ChainAApp.GetConsumerKeeper()
	consumerKeeperB := suite.ChainBApp.GetConsumerKeeper()

	// valsets must match
	providerValUpdates := cmttypes.TM2PB.ValidatorUpdates(suite.ChainProvider.Vals)
	consumerAValUpdates := cmttypes.TM2PB.ValidatorUpdates(suite.ChainA.Vals)
	consumerBValUpdates := cmttypes.TM2PB.ValidatorUpdates(suite.ChainB.Vals)
	suite.Require().True(len(providerValUpdates) == len(consumerAValUpdates), "initial valset not matching")
	suite.Require().True(len(providerValUpdates) == len(consumerBValUpdates), "initial valset not matching")

	for i := 0; i < len(providerValUpdates); i++ {
		addr1, _ := ccv.TMCryptoPublicKeyToConsAddr(providerValUpdates[i].PubKey)
		addr2, _ := ccv.TMCryptoPublicKeyToConsAddr(consumerAValUpdates[i].PubKey)
		addr3, _ := ccv.TMCryptoPublicKeyToConsAddr(consumerBValUpdates[i].PubKey)
		suite.Require().True(bytes.Equal(addr1, addr2), "validator mismatch")
		suite.Require().True(bytes.Equal(addr1, addr3), "validator mismatch")
	}

	ct := suite.ChainProvider.GetContext()
	// move chains to the next block
	suite.ChainProvider.NextBlock()
	suite.ChainA.NextBlock()
	suite.ChainB.NextBlock()

	initializationParameters := keeper2.GetTestInitializationParameters()
	// NOTE: we cannot use the time.Now() because the coordinator chooses a hardcoded start time
	// using time.Now() could set the spawn time to be too far in the past or too far in the future
	initializationParameters.SpawnTime = suite.Coordinator.CurrentTime
	// NOTE: the initial height passed to CreateConsumerClient
	// must be the height on the consumer when InitGenesis is called
	initializationParameters.InitialHeight = clienttypes.Height{RevisionNumber: 0, RevisionHeight: 2}

	// create consumer client on provider chain and set as consumer client for consumer chainID in provider keeper.
	prop1 := GetTestConsumerAdditionProp(suite.ChainA)

	providerKeeper.SetConsumerChainId(ct, prop1.ChainId, prop1.ChainId)
	err := providerKeeper.SetConsumerPowerShapingParameters(suite.ChainProvider.GetContext(), prop1.ChainId, keeper2.GetTestPowerShapingParameters())
	suite.Require().NoError(err)
	providerKeeper.SetConsumerPhase(ct, prop1.ChainId, providertypes.CONSUMER_PHASE_INITIALIZED)
	err = providerKeeper.SetConsumerInitializationParameters(ct, prop1.ChainId, initializationParameters)
	suite.Require().NoError(err)
	err = providerKeeper.SetConsumerMetadata(suite.ChainProvider.GetContext(), prop1.ChainId, keeper2.GetTestConsumerMetadata())
	suite.Require().NoError(err)
	err = providerKeeper.AppendConsumerToBeLaunched(suite.ChainProvider.GetContext(), prop1.ChainId, suite.Coordinator.CurrentTime)
	suite.Require().NoError(err)

	// opt-in all validators
	lastVals, err := providerKeeper.GetLastBondedValidators(suite.ChainProvider.GetContext())
	suite.Require().NoError(err)

	for _, v := range lastVals {
		consAddr, _ := v.GetConsAddr()
		providerKeeper.SetOptedIn(suite.ChainProvider.GetContext(), prop1.ChainId, providertypes.NewProviderConsAddress(consAddr))
	}

	prop2 := GetTestConsumerAdditionProp(suite.ChainB)

	providerKeeper.SetConsumerChainId(ct, prop2.ChainId, prop2.ChainId)
	err = providerKeeper.SetConsumerPowerShapingParameters(suite.ChainProvider.GetContext(), prop2.ChainId, keeper2.GetTestPowerShapingParameters())
	suite.Require().NoError(err)
	providerKeeper.SetConsumerPhase(ct, prop2.ChainId, providertypes.CONSUMER_PHASE_INITIALIZED)
	err = providerKeeper.SetConsumerInitializationParameters(ct, prop2.ChainId, initializationParameters)
	suite.Require().NoError(err)
	err = providerKeeper.SetConsumerMetadata(suite.ChainProvider.GetContext(), prop2.ChainId, keeper2.GetTestConsumerMetadata())
	suite.Require().NoError(err)
	err = providerKeeper.AppendConsumerToBeLaunched(suite.ChainProvider.GetContext(), prop2.ChainId, suite.Coordinator.CurrentTime)
	suite.Require().NoError(err)

	// opt-in all validators
	lastVals, err = providerKeeper.GetLastBondedValidators(suite.ChainProvider.GetContext())
	suite.Require().NoError(err)

	for _, v := range lastVals {
		consAddr, _ := v.GetConsAddr()
		providerKeeper.SetOptedIn(suite.ChainProvider.GetContext(), prop2.ChainId, providertypes.NewProviderConsAddress(consAddr))
	}

	// move provider to next block to commit the state
	suite.ChainProvider.NextBlock()

	// initialize the consumer chain with the genesis state stored on the provider
	consumerGenesisA, found := providerKeeper.GetConsumerGenesis(
		suite.ChainProvider.GetContext(),
		suite.ChainA.ChainID,
	)
	suite.Require().True(found, "consumer genesis not found")

	genesisStateA := consumertypes.GenesisState{
		Params:   consumerGenesisA.Params,
		Provider: consumerGenesisA.Provider,
		NewChain: consumerGenesisA.NewChain,
	}
	consumerKeeperA.InitGenesis(suite.ChainA.GetContext(), &genesisStateA)

	// initialize the consumer chain with the genesis state stored on the provider
	consumerGenesisB, found := providerKeeper.GetConsumerGenesis(
		suite.ChainProvider.GetContext(),
		suite.ChainB.ChainID,
	)
	suite.Require().True(found, "consumer genesis not found")

	genesisStateB := consumertypes.GenesisState{
		Params:   consumerGenesisB.Params,
		Provider: consumerGenesisB.Provider,
		NewChain: consumerGenesisB.NewChain,
	}
	consumerKeeperB.InitGenesis(suite.ChainB.GetContext(), &genesisStateB)

	suite.ChainA.NextBlock()
	suite.ChainB.NextBlock()

	// create paths for the CCV channel
	suite.CCVPathA = ibctesting.NewPath(suite.ChainA, suite.ChainProvider)
	suite.CCVPathB = ibctesting.NewPath(suite.ChainB, suite.ChainProvider)
	SetupCCVPath(suite.CCVPathA, suite)
	SetupCCVPath(suite.CCVPathB, suite)

	suite.SetupCCVChannels()

	suite.Path = NewICAPath(suite.ChainA, suite.ChainB, suite.ChainProvider)

	suite.Coordinator.SetupConnections(suite.Path)
}

func (suite *IBCConnectionTestSuite) ConfigureTransferChannel() {
	suite.TransferPath = NewTransferPath(suite.ChainA, suite.ChainB, suite.ChainProvider)
	suite.Coordinator.SetupConnections(suite.TransferPath)
	err := SetupTransferPath(suite.TransferPath)
	suite.Require().NoError(err)
}

func (suite *IBCConnectionTestSuite) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	bankKeeper := suite.GetNeutronZoneApp(suite.ChainA).BankKeeper
	err := bankKeeper.MintCoins(suite.ChainA.GetContext(), tokenfactorytypes.ModuleName, amounts)
	suite.Require().NoError(err)

	err = bankKeeper.SendCoinsFromModuleToAccount(suite.ChainA.GetContext(), tokenfactorytypes.ModuleName, acc, amounts)
	suite.Require().NoError(err)
}

// update CCV path with correct info
func SetupCCVPath(path *ibctesting.Path, suite *IBCConnectionTestSuite) {
	// - set provider endpoint's clientID
	consumerClient, found := suite.ProviderApp.GetProviderKeeper().GetConsumerClientId(
		suite.ChainProvider.GetContext(),
		path.EndpointA.Chain.ChainID,
	)

	suite.Require().True(found, "consumer client not found")
	path.EndpointB.ClientID = consumerClient

	// - set consumer endpoint's clientID
	consumerKeeper := path.EndpointA.Chain.App.(*app.App).GetConsumerKeeper()
	providerClient, found := consumerKeeper.GetProviderClientID(path.EndpointA.Chain.GetContext())
	suite.Require().True(found, "provider client not found")
	path.EndpointA.ClientID = providerClient

	// - client config
	trustingPeriodFraction := suite.ProviderApp.GetProviderKeeper().GetTrustingPeriodFraction(suite.ChainProvider.GetContext())

	providerUnbondingPeriod, err := suite.ProviderApp.GetTestStakingKeeper().UnbondingTime(suite.ChainProvider.GetContext())
	suite.Require().NoError(err)
	path.EndpointB.ClientConfig.(*ibctesting.TendermintConfig).UnbondingPeriod = providerUnbondingPeriod
	path.EndpointB.ClientConfig.(*ibctesting.TendermintConfig).TrustingPeriod, _ = ccv.CalculateTrustPeriod(providerUnbondingPeriod, trustingPeriodFraction)
	consumerUnbondingPeriod := consumerKeeper.GetUnbondingPeriod(path.EndpointA.Chain.GetContext())
	path.EndpointA.ClientConfig.(*ibctesting.TendermintConfig).UnbondingPeriod = consumerUnbondingPeriod
	path.EndpointA.ClientConfig.(*ibctesting.TendermintConfig).TrustingPeriod, _ = ccv.CalculateTrustPeriod(consumerUnbondingPeriod, trustingPeriodFraction)
	// - channel config
	path.EndpointA.ChannelConfig.PortID = ccv.ConsumerPortID
	path.EndpointB.ChannelConfig.PortID = ccv.ProviderPortID
	path.EndpointA.ChannelConfig.Version = ccv.Version
	path.EndpointB.ChannelConfig.Version = ccv.Version
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
}

func (suite *IBCConnectionTestSuite) SetupCCVChannels() {
	paths := []*ibctesting.Path{suite.CCVPathA, suite.CCVPathB}
	for _, path := range paths {
		suite.Coordinator.CreateConnections(path)

		err := path.EndpointA.ChanOpenInit()
		suite.Require().NoError(err)

		err = path.EndpointB.ChanOpenTry()
		suite.Require().NoError(err)

		err = path.EndpointA.ChanOpenAck()
		suite.Require().NoError(err)

		err = path.EndpointB.ChanOpenConfirm()
		suite.Require().NoError(err)

		err = path.EndpointA.UpdateClient()
		suite.Require().NoError(err)
	}
}

func testHomeDir(chainID string) string {
	projectRoot := utils.RootDir()
	return path.Join(projectRoot, ".testchains", chainID)
}

// NewCoordinator initializes Coordinator with interchain security dummy provider and 2 neutron consumer chains
func NewProviderConsumerCoordinator(t *testing.T) *ibctesting.Coordinator {
	coordinator := ibctesting.NewCoordinator(t, 0)
	chainID := ibctesting.GetChainID(1)

	ibctesting.DefaultTestingAppInit = icssimapp.ProviderAppIniter
	coordinator.Chains[chainID] = ibctesting.NewTestChain(t, coordinator, chainID)
	providerChain := coordinator.GetChain(chainID)

	_ = config.GetDefaultConfig()
	sdk.SetAddrCacheEnabled(false)
	chainID = ibctesting.GetChainID(2)
	ibctesting.DefaultTestingAppInit = SetupTestingApp(cmttypes.TM2PB.ValidatorUpdates(providerChain.Vals))
	coordinator.Chains[chainID] = ibctesting.NewTestChainWithValSet(t, coordinator,
		chainID, providerChain.Vals, providerChain.Signers)

	chainID = ibctesting.GetChainID(3)
	coordinator.Chains[chainID] = ibctesting.NewTestChainWithValSet(t, coordinator,
		chainID, providerChain.Vals, providerChain.Signers)

	return coordinator
}

func (suite *IBCConnectionTestSuite) GetNeutronZoneApp(chain *ibctesting.TestChain) *app.App {
	testApp, ok := chain.App.(*app.App)
	if !ok {
		panic("not NeutronZone app")
	}

	return testApp
}

func (suite *IBCConnectionTestSuite) StoreTestCode(ctx sdk.Context, addr sdk.AccAddress, path string) uint64 {
	// wasm file built with https://github.com/neutron-org/neutron-sdk/tree/main/contracts/reflect
	// wasm file built with https://github.com/neutron-org/neutron-dev-contracts/tree/feat/ica-register-fee-update/contracts/neutron_interchain_txs
	wasmCode, err := os.ReadFile(path)
	suite.Require().NoError(err)

	codeID, _, err := keeper.NewDefaultPermissionKeeper(suite.GetNeutronZoneApp(suite.ChainA).WasmKeeper).Create(ctx, addr, wasmCode, &wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeEverybody})
	suite.Require().NoError(err)

	return codeID
}

func (suite *IBCConnectionTestSuite) InstantiateTestContract(ctx sdk.Context, funder sdk.AccAddress, codeID uint64) sdk.AccAddress {
	initMsgBz := []byte("{}")
	contractKeeper := keeper.NewDefaultPermissionKeeper(suite.GetNeutronZoneApp(suite.ChainA).WasmKeeper)
	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "demo contract", nil)
	suite.Require().NoError(err)

	return addr
}

func NewICAPath(chainA, chainB, chainProvider *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointB.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointA.ChannelConfig.Version = TestVersion
	path.EndpointB.ChannelConfig.Version = TestVersion

	trustingPeriodFraction := chainProvider.App.(*appProvider.App).GetProviderKeeper().GetTrustingPeriodFraction(chainProvider.GetContext())

	consumerUnbondingPeriodA := path.EndpointA.Chain.App.(*app.App).GetConsumerKeeper().GetUnbondingPeriod(path.EndpointA.Chain.GetContext())
	path.EndpointA.ClientConfig.(*ibctesting.TendermintConfig).UnbondingPeriod = consumerUnbondingPeriodA
	path.EndpointA.ClientConfig.(*ibctesting.TendermintConfig).TrustingPeriod, _ = ccv.CalculateTrustPeriod(consumerUnbondingPeriodA, trustingPeriodFraction)

	consumerUnbondingPeriodB := path.EndpointB.Chain.App.(*app.App).GetConsumerKeeper().GetUnbondingPeriod(path.EndpointB.Chain.GetContext())
	path.EndpointB.ClientConfig.(*ibctesting.TendermintConfig).UnbondingPeriod = consumerUnbondingPeriodB
	path.EndpointB.ClientConfig.(*ibctesting.TendermintConfig).TrustingPeriod, _ = ccv.CalculateTrustPeriod(consumerUnbondingPeriodB, trustingPeriodFraction)

	return path
}

// SetupICAPath invokes the InterchainAccounts entrypoint and subsequent channel handshake handlers
func SetupICAPath(path *ibctesting.Path, owner string) error {
	if err := RegisterInterchainAccount(path.EndpointA, owner); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenTry(); err != nil {
		return err
	}

	if err := path.EndpointA.ChanOpenAck(); err != nil {
		return err
	}

	return path.EndpointB.ChanOpenConfirm()
}

// RegisterInterchainAccount is a helper function for starting the channel handshake
func RegisterInterchainAccount(endpoint *ibctesting.Endpoint, owner string) error {
	icaOwner, _ := ictxstypes.NewICAOwner(owner, TestInterchainID)
	portID, err := icatypes.NewControllerPortID(icaOwner.String())
	if err != nil {
		return err
	}

	ctx := endpoint.Chain.GetContext()

	channelSequence := endpoint.Chain.App.GetIBCKeeper().ChannelKeeper.GetNextChannelSequence(ctx)

	a, ok := endpoint.Chain.App.(*app.App)
	if !ok {
		return fmt.Errorf("not NeutronZoneApp")
	}

	icaMsgServer := icacontrollerkeeper.NewMsgServerImpl(&a.ICAControllerKeeper)
	if _, err = icaMsgServer.RegisterInterchainAccount(ctx, &icacontrollertypes.MsgRegisterInterchainAccount{
		Owner:        icaOwner.String(),
		ConnectionId: endpoint.ConnectionID,
		Version:      TestVersion,
		Ordering:     channeltypes.ORDERED,
	}); err != nil {
		return err
	}

	// commit state changes for proof verification
	endpoint.Chain.NextBlock()

	// update port/channel ids
	endpoint.ChannelID = channeltypes.FormatChannelIdentifier(channelSequence)
	endpoint.ChannelConfig.PortID = portID

	return nil
}

// SetupTestingApp initializes the IBC-go testing application
func SetupTestingApp(initValUpdates []cometbfttypes.ValidatorUpdate) func() (ibctesting.TestingApp, map[string]json.RawMessage) {
	return func() (ibctesting.TestingApp, map[string]json.RawMessage) {
		encoding := app.MakeEncodingConfig()
		db := db2.NewMemDB()
		homePath := testHomeDir("testchain-" + tmrand.NewRand().Str(6))
		testApp := app.New(
			log.NewNopLogger(),
			db,
			nil,
			false,
			map[int64]bool{},
			homePath,
			0,
			encoding,
			sims.EmptyAppOptions{},
			nil,
		)

		// we need to set up a TestInitChainer where we can redefine MaxBlockGas in ConsensusParamsKeeper
		testApp.SetInitChainer(testApp.TestInitChainer)
		// and then we manually init baseapp and load states
		testApp.LoadLatest()

		genesisState := app.NewDefaultGenesisState(testApp.AppCodec())

		// TODO: why isn't it in the `testApp.TestInitChainer`?
		// NOTE ibc-go/v7/testing.SetupWithGenesisValSet requires a staking module
		// genesisState or it panics. Feed a minimum one.
		genesisState[stakingtypes.ModuleName] = encoding.Marshaler.MustMarshalJSON(
			&stakingtypes.GenesisState{
				Params: stakingtypes.Params{BondDenom: sdk.DefaultBondDenom},
			},
		)

		var consumerGenesis ccv.ConsumerGenesisState
		encoding.Marshaler.MustUnmarshalJSON(genesisState[consumertypes.ModuleName], &consumerGenesis)
		consumerGenesis.Provider.InitialValSet = initValUpdates
		consumerGenesis.Params.Enabled = true
		genesisState[consumertypes.ModuleName] = encoding.Marshaler.MustMarshalJSON(&consumerGenesis)

		return testApp, genesisState
	}
}

// SetupValSetAppIniter is a simple wrapper for ICS e2e tests to satisfy interface
func SetupValSetAppIniter(initValUpdates []cometbfttypes.ValidatorUpdate) icssimapp.AppIniter {
	return SetupTestingApp(initValUpdates)
}

func NewTransferPath(chainA, chainB, chainProvider *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = types.PortID
	path.EndpointB.ChannelConfig.PortID = types.PortID
	path.EndpointA.ChannelConfig.Order = channeltypes.UNORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.UNORDERED
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	trustingPeriodFraction := chainProvider.App.(*appProvider.App).GetProviderKeeper().GetTrustingPeriodFraction(chainProvider.GetContext())
	consumerUnbondingPeriodA := path.EndpointA.Chain.App.(*app.App).GetConsumerKeeper().GetUnbondingPeriod(path.EndpointA.Chain.GetContext())
	path.EndpointA.ClientConfig.(*ibctesting.TendermintConfig).UnbondingPeriod = consumerUnbondingPeriodA
	path.EndpointA.ClientConfig.(*ibctesting.TendermintConfig).TrustingPeriod, _ = ccv.CalculateTrustPeriod(consumerUnbondingPeriodA, trustingPeriodFraction)

	consumerUnbondingPeriodB := path.EndpointB.Chain.App.(*app.App).GetConsumerKeeper().GetUnbondingPeriod(path.EndpointB.Chain.GetContext())
	path.EndpointB.ClientConfig.(*ibctesting.TendermintConfig).UnbondingPeriod = consumerUnbondingPeriodB
	path.EndpointB.ClientConfig.(*ibctesting.TendermintConfig).TrustingPeriod, _ = ccv.CalculateTrustPeriod(consumerUnbondingPeriodB, trustingPeriodFraction)

	return path
}

// SetupTransferPath
func SetupTransferPath(path *ibctesting.Path) error {
	channelSequence := path.EndpointA.Chain.App.GetIBCKeeper().ChannelKeeper.GetNextChannelSequence(path.EndpointA.Chain.GetContext())
	channelSequenceB := path.EndpointB.Chain.App.GetIBCKeeper().ChannelKeeper.GetNextChannelSequence(path.EndpointB.Chain.GetContext())

	// update port/channel ids
	path.EndpointA.ChannelID = channeltypes.FormatChannelIdentifier(channelSequence)
	path.EndpointB.ChannelID = channeltypes.FormatChannelIdentifier(channelSequenceB)

	if err := path.EndpointA.ChanOpenInit(); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenTry(); err != nil {
		return err
	}

	if err := path.EndpointA.ChanOpenAck(); err != nil {
		return err
	}

	return path.EndpointB.ChanOpenConfirm()
}
