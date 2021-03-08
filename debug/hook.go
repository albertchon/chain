package debug

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

type HookBlock struct {
	Block      uint64 `json:"block"`
	GasWanted  uint64 `json:"gas_wanted"`
	GasUsed    uint64 `json:"gas_used"`
	BlockLimit uint64 `json:"block_limit"`
	Resolved   int    `json:"resolved"`
}

type Hook struct {
	db        *DB
	hookBlock HookBlock
	table     string
}

// AfterInitChain specify actions need to do after chain initialization (app.Hook interface).
func (h *Hook) AfterInitChain(ctx sdk.Context, req abci.RequestInitChain, res abci.ResponseInitChain) {
}

// AfterBeginBlock specify actions need to do after begin block period (app.Hook interface).
func (h *Hook) AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
	h.hookBlock.Block = uint64(ctx.BlockHeader().Height)
	h.hookBlock.GasUsed = 0
	h.hookBlock.GasWanted = 0
	h.hookBlock.BlockLimit = ctx.BlockGasMeter().Limit()
}

// AfterDeliverTx specify actions need to do after transaction has been processed (app.Hook interface).
func (h *Hook) AfterDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) {
	h.hookBlock.GasUsed += uint64(res.GasUsed)
	h.hookBlock.GasWanted += uint64(res.GasWanted)
}

// AfterEndBlock specify actions need to do after end block period (app.Hook interface).
func (h *Hook) AfterEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) {
	h.db.Insert("block-stats", h.hookBlock)
}

// ApplyQuery catch the custom query that matches specific paths (app.Hook interface).
func (h *Hook) ApplyQuery(req abci.RequestQuery) (res abci.ResponseQuery, stop bool) {
	return abci.ResponseQuery{}, false
}

// BeforeCommit specify actions need to do before commit block (app.Hook interface).
func (h *Hook) BeforeCommit() {
}

func NewHook(db *DB, table string) *Hook {
	return &Hook{
		db:    db,
		table: table,
	}
}
