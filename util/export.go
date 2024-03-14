package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Export(ctx context.Context, client *ethclient.Client, blockNumber uint64, txHashes []common.Hash, outputPath string) error {
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	marshaled, _ := json.Marshal(blockNumber)
	_, err = file.Write(marshaled)
	if err != nil {
		return err
	}
	file.WriteString("\n")

	for _, txHash := range txHashes {
		tx, isPending, err := client.TransactionByHash(ctx, txHash)
		if err != nil {
			return err
		}

		if tx == nil || isPending {
			return fmt.Errorf("transaction %s not found or pending", txHash.String())
		}

		marshaled, err := json.Marshal(tx)
		_, err = file.Write(marshaled)
		if err != nil {
			return err
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}
