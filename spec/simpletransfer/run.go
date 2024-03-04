package simpletransfer

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Spec struct {
	Config
}

func NewSpec(config Config) *Spec {
	return &Spec{Config: config}
}

func (spec *Spec) Run(ctx context.Context, client *ethclient.Client) error {
	// TODO
	return nil
}
