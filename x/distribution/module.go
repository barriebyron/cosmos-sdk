package distribution

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// app module basics object
type AppModuleBasic struct{}

var _ sdk.AppModuleBasic = AppModuleBasic{}

// module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// module name
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// module name
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

var _ sdk.AppModule = AppModule{}

// register invariants
func (a AppModule) RegisterInvariants(ir sdk.InvariantRouter) {
	RegisterInvariants(ir, a.keeper)
}

// module message route name
func (AppModule) Route() string {
	return RouterKey
}

// module handler
func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(a.keeper)
}

// module querier route name
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// module querier
func (a AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(a.keeper)
}

// module init-genesis
func (a AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, a.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// module export genesis
func (a AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, a.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// module begin-block
func (a AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) sdk.Tags {
	BeginBlocker(ctx, req, a.keeper)
	return sdk.EmptyTags()
}

// module end-block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) ([]abci.ValidatorUpdate, sdk.Tags) {
	return []abci.ValidatorUpdate{}, sdk.EmptyTags()
}
