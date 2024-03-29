package simpletransfer

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"keroro520/bench/util"
	"math/big"
)

// TODO formalize the config when load
// TODO support 1559 transaction

type Spec struct {
	Config
	cursor int
}

func NewSpec(config Config) *Spec {
	return &Spec{cursor: 0, Config: config}
}

func (spec *Spec) Run(ctx context.Context, client *ethclient.Client, engine *ethclient.Client) ([]common.Hash, error) {
	shotTxs := make([]common.Hash, 0)
	for i := 0; i < spec.TxnsCount; i++ {
		txs, err := spec.NextTxs(ctx, client)
		if err != nil {
			return nil, err
		}

		txsData := make([][]byte, len(txs), len(txs))
		for j, tx := range txs {
			log.Info("Shot transaction", "hash", tx.Hash())

			txData, err := tx.MarshalBinary()
			if err != nil {
				return nil, err
			}
			txsData[j] = txData
		}

		_, err = util.ProduceNextBlock(ctx, engine, txsData, nil)
		if err != nil {
			return nil, err
		}

		for _, tx := range txs {
			shotTxs = append(shotTxs, tx.Hash())
		}
	}

	return shotTxs, nil
}

func (spec *Spec) NextTxs(ctx context.Context, client *ethclient.Client) (types.Transactions, error) {
	tx, err := spec.NextTx(ctx, client)
	if err != nil {
		return nil, err
	}
	return types.Transactions{tx}, nil
}

func (spec *Spec) NextTx(ctx context.Context, client *ethclient.Client) (*types.Transaction, error) {
	spec.cursor++
	return spec.nextTx(ctx, client)
}

func (spec *Spec) nextTx(ctx context.Context, client *ethclient.Client) (*types.Transaction, error) {
	from, err := util.NewAccountFromRaw(spec.PrivKeys[spec.cursor%len(spec.PrivKeys)])
	if err != nil {
		return nil, err
	}

	to, err := util.NewAccountFromRaw(spec.PrivKeys[(spec.cursor+1)%len(spec.PrivKeys)])
	if err != nil {
		return nil, err
	}
	toAddress := to.Address()

	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(ctx, from.Address())
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.LegacyTx{
		To:       &toAddress,
		Nonce:    nonce,
		GasPrice: big.NewInt(int64(spec.TxGasPrice)),
		Gas:      spec.TxGasLimit,
		Value:    big.NewInt(int64(spec.TxValue)),
	})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), &from.PrivateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
