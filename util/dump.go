package util

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func DumpTransactions(ctx context.Context, client *ethclient.Client, txHashes []common.Hash, outputPath string) error {
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, txHash := range txHashes {
		tx, isPending, err := client.TransactionByHash(ctx, txHash)
		if err != nil {
			return err
		}

		if tx == nil || isPending {
			return fmt.Errorf("transaction %s not found or pending", txHash.String())
		}

		// Append the RLP transaction to the file
		encoded, err := rlp.EncodeToBytes(tx.Data())
		if err != nil {
			return err
		}
		_, err = file.Write(append(encoded, '\n'))
		if err != nil {
			return err
		}
	}

	return nil
}
