package util

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
)

func Import(ctx context.Context, inputFile string, engineClient *ethclient.Client, txsPerBlock int) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	var parentBlockNumber uint64
	var firstLine = scanner.Bytes()
	err = json.Unmarshal(firstLine, &parentBlockNumber)
	if err != nil {
		return err
	}

	blockTxs := make([][]byte, 0)
	for scanner.Scan() {
		var tx types.Transaction
		encodedTransaction := scanner.Bytes()
		err = json.Unmarshal(encodedTransaction, &tx)
		if err != nil {
			return err
		}

		txData, err := tx.MarshalBinary()
		if err != nil {
			return err
		}

		blockTxs = append(blockTxs, txData)
		if len(blockTxs) < txsPerBlock {
			continue
		}

		_, err = ProduceNextBlock(ctx, engineClient, blockTxs, big.NewInt(int64(parentBlockNumber)))
		if err != nil {
			return err
		}

		parentBlockNumber += 1
		blockTxs = make([][]byte, 0)
	}

	if len(blockTxs) > 0 {
		_, err = ProduceNextBlock(ctx, engineClient, blockTxs, big.NewInt(int64(parentBlockNumber)))
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading file: %v", err)
	}
	return nil
}
