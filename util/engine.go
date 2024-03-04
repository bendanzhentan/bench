package util

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// NOTE: We don't care about the OP's derivation rules because it's consensus rules while we are testing the execution part.

func ProduceNextBlock(ctx context.Context, client *ethclient.Client, txs types.Transactions) (*common.Hash, error) {
	txsData := make([]hexutil.Bytes, len(txs), len(txs))
	for i, tx := range txs {
		txData, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		txsData[i] = txData
	}
	payloadAttributes := eth.PayloadAttributes{
		// Leave default values for now.
		Timestamp:             0,
		PrevRandao:            eth.Bytes32{},
		SuggestedFeeRecipient: common.Address{},
		Withdrawals:           nil,
		GasLimit:              nil,
		ParentBeaconBlockRoot: nil,

		// NoTxPool=true so that the block building process doesn't need to check the tx pool.
		NoTxPool: true,
		// Transactions to force into the block (always at the start of the transactions list).
		Transactions: txsData,
	}

	fcr, err := EngineForkchoiceUpdated(ctx, client, nil, &payloadAttributes)
	if err != nil {
		return nil, err
	}

	paylaod, err := EngineGetPayload(ctx, client, *fcr.PayloadID)
	if err != nil {
		return nil, err
	}

	status, err := EngineNewPayload(ctx, client, paylaod)
	if err != nil {
		return nil, err
	}

	fc := eth.ForkchoiceState{
		HeadBlockHash:      *status.LatestValidHash,
		SafeBlockHash:      *status.LatestValidHash,
		FinalizedBlockHash: *status.LatestValidHash,
	}
	_, err = EngineForkchoiceUpdated(ctx, client, &fc, nil)
	if err != nil {
		return nil, err
	}

	return status.LatestValidHash, nil
}

func EngineForkchoiceUpdated(ctx context.Context, client *ethclient.Client, fc *eth.ForkchoiceState, attributes *eth.PayloadAttributes) (*eth.ForkchoiceUpdatedResult, error) {
	fcCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	var result eth.ForkchoiceUpdatedResult
	err := client.Client().CallContext(fcCtx, &result, "engine_forkchoiceUpdatedV1", fc, attributes)
	if err != nil {
		return nil, err
	} else if result.PayloadStatus.Status != eth.ExecutionValid {
		return nil, fmt.Errorf("payload execution failed: %s, status: %s", *result.PayloadStatus.ValidationError, result.PayloadStatus.Status)
	}

	return &result, nil
}

func EngineGetPayload(ctx context.Context, client *ethclient.Client, payloadId eth.PayloadID) (*eth.ExecutionPayload, error) {
	var result eth.ExecutionPayload
	err := client.Client().CallContext(ctx, &result, "engine_getPayloadV1", payloadId)
	return &result, err
}

func EngineNewPayload(ctx context.Context, client *ethclient.Client, payload *eth.ExecutionPayload) (*eth.PayloadStatusV1, error) {
	execCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	var result eth.PayloadStatusV1
	err := client.Client().CallContext(execCtx, &result, "engine_newPayloadV1", payload)
	return &result, err
}
