package cron

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/neutron-org/neutron/x/cron/keeper"
	"github.com/neutron-org/neutron/x/cron/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the schedule
	for _, elem := range genState.ScheduleList {
		k.SetSchedule(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.ScheduleList = k.GetAllSchedules(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
