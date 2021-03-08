package keeper

import (
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	"github.com/bandprotocol/chain/pkg/bandrng"
	"github.com/bandprotocol/chain/x/oracle/types"
)

const FixedResolve = 4_100_000
const Factor = 8

func convertToOwasmGas(cosmos uint64) uint32 {
	return uint32(cosmos * Factor)
}

func convertToCosmosGas(owasm uint32) uint64 {
	return uint64(owasm / Factor)
}

// GetRandomValidators returns a pseudorandom subset of active validators. Each validator has
// chance of getting selected directly proportional to the amount of voting power it has.
func (k Keeper) GetRandomValidators(ctx sdk.Context, size int, id int64) ([]sdk.ValAddress, error) {
	valOperators := []sdk.ValAddress{}
	valPowers := []uint64{}
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
		func(idx int64, val exported.ValidatorI) (stop bool) {
			if k.GetValidatorStatus(ctx, val.GetOperator()).IsActive {
				valOperators = append(valOperators, val.GetOperator())
				valPowers = append(valPowers, val.GetTokens().Uint64())
			}
			return false
		})
	if len(valOperators) < size {
		return nil, sdkerrors.Wrapf(
			types.ErrInsufficientValidators, "%d < %d", len(valOperators), size)
	}
	rng, err := bandrng.NewRng(k.GetRollingSeed(ctx), sdk.Uint64ToBigEndian(uint64(id)), []byte(ctx.ChainID()))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrBadDrbgInitialization, err.Error())
	}
	tryCount := int(k.GetParam(ctx, types.KeySamplingTryCount))
	chosenValIndexes := bandrng.ChooseSomeMaxWeight(rng, valPowers, size, tryCount)
	validators := make([]sdk.ValAddress, size)
	for i, idx := range chosenValIndexes {
		validators[i] = valOperators[idx]
	}
	return validators, nil
}

// PrepareRequest takes an request specification object, performs the prepare call, and saves
// the request object to store. Also emits events related to the request.
func (k Keeper) PrepareRequest(ctx sdk.Context, r types.RequestSpec) error {
	start := time.Now()
	startGas := ctx.GasMeter().GasConsumed()
	askCount := r.GetAskCount()
	if askCount > k.GetParam(ctx, types.KeyMaxAskCount) {
		return sdkerrors.Wrapf(types.ErrInvalidAskCount, "got: %d, max: %d", askCount, k.GetParam(ctx, types.KeyMaxAskCount))
	}
	// Consume gas for data requests. We trust that we have reasonable params that don't cause overflow.
	ctx.GasMeter().ConsumeGas(askCount*k.GetParam(ctx, types.KeyPerValidatorRequestGas), "PER_VALIDATOR_REQUEST_FEE")
	// Get a random validator set to perform this request.
	validators, err := k.GetRandomValidators(ctx, int(askCount), k.GetRequestCount(ctx)+1)
	if err != nil {
		return err
	}
	// Create a request object. Note that RawRequestIDs will be populated after preparation is done.
	req := types.NewRequest(
		r.GetOracleScriptID(), r.GetCalldata(), validators, r.GetMinCount(),
		ctx.BlockHeight(), ctx.BlockTime(), r.GetClientID(), nil,
	)
	// Create an execution environment and call Owasm prepare function.
	env := types.NewPrepareEnv(req, int64(k.GetParam(ctx, types.KeyMaxRawRequestCount)))
	script, err := k.GetOracleScript(ctx, req.OracleScriptID)
	if err != nil {
		return err
	}
	code := k.GetFile(script.Filename)
	startPrepare := time.Now()
	output, err := k.owasmVM.Prepare(code, types.WasmPrepareGas, types.MaxDataSize, env)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrBadWasmExecution, err.Error())
	}
	k.Prepare(startPrepare, convertToCosmosGas(output.GasUsed))
	preparedTime := time.Since(startPrepare)
	ctx.GasMeter().ConsumeGas(convertToCosmosGas(output.GasUsed), "PREPARE_GAS")
	// Preparation complete! It's time to collect raw request ids.
	req.RawRequests = env.GetRawRequests()
	if len(req.RawRequests) == 0 {
		return types.ErrEmptyRawRequests
	}
	// We now have everything we need to the request, so let's add it to the store.
	id := k.AddRequest(ctx, req)
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequest)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyClientID, req.ClientID),
		sdk.NewAttribute(types.AttributeKeyOracleScriptID, fmt.Sprintf("%d", req.OracleScriptID)),
		sdk.NewAttribute(types.AttributeKeyCalldata, hex.EncodeToString(req.Calldata)),
		sdk.NewAttribute(types.AttributeKeyAskCount, fmt.Sprintf("%d", askCount)),
		sdk.NewAttribute(types.AttributeKeyMinCount, fmt.Sprintf("%d", req.MinCount)),
		sdk.NewAttribute(types.AttributeKeyGasUsed, fmt.Sprintf("%d", output.GasUsed)),
	)
	for _, val := range req.RequestedValidators {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyValidator, val.String()))
	}
	ctx.EventManager().EmitEvent(event)
	// Emit an event for each of the raw data requests.
	for _, rawReq := range env.GetRawRequests() {
		ds, err := k.GetDataSource(ctx, rawReq.DataSourceID)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, fmt.Sprintf("%d", rawReq.DataSourceID)),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, ds.Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, fmt.Sprintf("%d", rawReq.ExternalID)),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(rawReq.Calldata)),
		))
	}

	ctx.GasMeter().ConsumeGas(convertToCosmosGas(FixedResolve), "RESOLVE_RESERVATION")
	k.Request(start, ctx.GasMeter().GasConsumed()-startGas)
	k.NewRequest(
		ctx,
		int64(ctx.GasMeter().GasConsumed()-startGas),
		time.Since(start),
		int64(convertToCosmosGas(FixedResolve)),
		preparedTime,
	)
	return nil
}

// ResolveRequest resolves the given request and saves the result to the store. The function
// assumes that the given request is in a resolvable state with sufficient reporters.
func (k Keeper) ResolveRequest(ctx sdk.Context, reqID types.RequestID) {
	start := time.Now()
	gasStart := ctx.GasMeter().GasConsumed()
	req := k.MustGetRequest(ctx, reqID)
	env := types.NewExecuteEnv(req, k.GetReports(ctx, reqID))
	script := k.MustGetOracleScript(ctx, req.OracleScriptID)
	code := k.GetFile(script.Filename)
	startExec := time.Now()
	output, err := k.owasmVM.Execute(code, FixedResolve, types.MaxDataSize, env)
	k.Execute(startExec, uint64(output.GasUsed)/5)
	timeExec := time.Since(startExec)
	if err != nil {
		fmt.Printf("Fail to resolve %s\n", err.Error())
		k.ResolveFailure(ctx, reqID, err.Error())
	} else if env.Retdata == nil {
		k.ResolveFailure(ctx, reqID, "no return data")
	} else {
		k.ResolveSuccess(ctx, reqID, env.Retdata, output.GasUsed)
	}
	k.Resolve(start, ctx.GasMeter().GasConsumed()-gasStart)
	k.RecordResolve(
		ctx,
		int64(reqID),
		int64(ctx.GasMeter().GasConsumed()-gasStart),
		time.Since(start),
		int64(convertToCosmosGas(output.GasUsed)),
		timeExec,
	)
}
