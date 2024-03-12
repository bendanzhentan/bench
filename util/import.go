package util

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"os"
)

// TODO 指定 block number
func Import(ctx context.Context, inputFile string, engineClient *ethclient.Client) error {
	// Open the file and read the file line by line
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// TODO parentBlockNumber
	parentBlockNumber := big.NewInt(3522)

	// TODO configured txs count
	blockTxsCount := 100
	blockTxs := make([][]byte, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read each line (RLP encoded transaction) from the file
		encodedTransaction := scanner.Bytes()

		var tx types.Transaction
		err = json.Unmarshal(encodedTransaction, &tx)
		if err != nil {
			return err
		}

		// TODO Txs count
		txData, err := tx.MarshalBinary()
		if err != nil {
			return err
		}

		blockTxs = append(blockTxs, txData)
		if len(blockTxs) < blockTxsCount {
			continue
		}

		blockHash, err := ProduceNextBlock(ctx, engineClient, blockTxs, parentBlockNumber)
		if err != nil {
			return err
		}

		blockTxs = make([][]byte, 0)
		parentBlockNumber = big.NewInt(0).Add(parentBlockNumber, big.NewInt(1))
		log.Info("Imported block", "number", parentBlockNumber, "hash", blockHash)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading file: %v", err)
	}
	return nil
}
