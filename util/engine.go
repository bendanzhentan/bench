package util

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"
)

func ProduceNextBlock(ctx context.Context, engineClient *ethclient.Client, txsData [][]byte, parentBlockNumber *big.Int) (*common.Hash, error) {
	parent, err := engineClient.HeaderByNumber(ctx, parentBlockNumber)
	if err != nil {
		return nil, err
	}

	fc := engine.ForkchoiceStateV1{
		HeadBlockHash:      parent.Hash(),
		SafeBlockHash:      parent.Hash(),
		FinalizedBlockHash: parent.Hash(),
	}
	payloadAttributes := engine.PayloadAttributes{
		// Leave default values for now.
		Random:                common.Hash{},
		SuggestedFeeRecipient: common.Address{},
		Withdrawals:           []*types.Withdrawal{},

		// Post-shanghai beacon root is required.
		BeaconRoot: parent.ParentBeaconRoot,
		GasLimit:   &parent.GasLimit,

		// Timestamp interval is assumed to be 1 second.
		Timestamp: parent.Time + 1,
		// NoTxPool=true so that the block building process doesn't need to check the tx pool.
		NoTxPool: true,

		// Optimism: Transactions to force into the block (always at the start of the transactions list).
		//
		// NOTE: We don't care about the OP's derivation rules because it's consensus rules while we are testing
		// the execution part.
		// NOTE: The force-inclusion transactions feature is only available for Optimism. So please make sure that
		// the node is running using op-geth.
		Transactions: txsData,
	}

	fcr, err := EngineForkchoiceUpdated(ctx, engineClient, &fc, &payloadAttributes)
	if err != nil {
		return nil, err
	}

	payload, err := EngineGetPayload(ctx, engineClient, *fcr.PayloadID)
	if err != nil {
		return nil, err
	}

	status, err := EngineNewPayload(ctx, engineClient, *payload.ExecutionPayload)
	if err != nil {
		return nil, err
	}

	fc = engine.ForkchoiceStateV1{
		HeadBlockHash:      *status.LatestValidHash,
		SafeBlockHash:      *status.LatestValidHash,
		FinalizedBlockHash: *status.LatestValidHash,
	}
	_, err = EngineForkchoiceUpdated(ctx, engineClient, &fc, nil)
	if err != nil {
		return nil, err
	}

	log.Info("Shot block", "number", payload.ExecutionPayload.Number, "hash", status.LatestValidHash)

	return status.LatestValidHash, nil
}

func EngineForkchoiceUpdated(ctx context.Context, client *ethclient.Client, fc *engine.ForkchoiceStateV1, attributes *engine.PayloadAttributes) (*engine.ForkChoiceResponse, error) {
	fcCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	var result *engine.ForkChoiceResponse
	err := client.Client().CallContext(fcCtx, &result, "engine_forkchoiceUpdatedV2", fc, attributes)
	if err != nil {
		return nil, err
	} else if result.PayloadStatus.Status != "VALID" {
		return nil, fmt.Errorf("payload execution failed: %v", result)
	}

	return result, nil
}

func EngineGetPayload(ctx context.Context, client *ethclient.Client, payloadId engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error) {
	var result *engine.ExecutionPayloadEnvelope
	err := client.Client().CallContext(ctx, &result, "engine_getPayloadV2", payloadId)
	return result, err
}

func EngineNewPayload(ctx context.Context, client *ethclient.Client, payload engine.ExecutableData) (engine.PayloadStatusV1, error) {
	execCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	var result engine.PayloadStatusV1
	err := client.Client().CallContext(execCtx, &result, "engine_newPayloadV2", payload)
	return result, err
}
