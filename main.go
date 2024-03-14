package main

import (
	"context"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/opio"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
	"keroro520/bench/core"
	"keroro520/bench/spec/simplecall"
	"keroro520/bench/spec/simpletransfer"
	"keroro520/bench/util"
	"os"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "export",
				Usage: "Export transactions from the chain",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "rpc-url",
						Usage:    "The RPC endpoint of the existing chain",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "engine-url",
						Usage:    "The Engine API endpoint of the existing chain",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "engine-jwt-secret",
						Usage:    "The JWT secret used to sign JWTs for the Engine API",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "config-path",
						Usage:    "The path to the configuration file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output-path",
						Usage:    "The file path for the output",
						Required: true,
					},
				},
				Action: func(context *cli.Context) error {
					rpcURL := context.String("rpc-url")
					engineURL := context.String("engine-url")
					engineJWTSecret := context.String("engine-jwt-secret")
					configPath := context.String("config-path")
					outputPath := context.String("output-path")

					benchConfig, err := core.LoadConfig(configPath)
					if err != nil {
						log.Crit("Failed to load configuration file", "err", err)
					}

					client, err := ethclient.Dial(rpcURL)
					if err != nil {
						log.Crit("Failed to connect to the Ethereum client", "url", rpcURL, "err", err)
					}

					var jwtSecret32 [32]byte
					copy(jwtSecret32[:], hexutils.HexToBytes(engineJWTSecret))
					c, err := rpc.DialOptions(context.Context, engineURL, rpc.WithHTTPAuth(node.NewJWTAuth(jwtSecret32)))
					if err != nil {
						log.Crit("Failed to create Engine client", "err", err)
					}
					engine := ethclient.NewClient(c)

					switch benchConfig.TxType {
					case "simpletransfer":
						spec := simpletransfer.NewSpec(*benchConfig.SimpleTransfer)
						blockNumber, err := client.BlockNumber(context.Context)
						if err != nil {
							log.Crit("Failed to get block number", "err", err)
						}
						txHashes, err := spec.Run(context.Context, client, engine)
						if err != nil {
							log.Crit("Failed to run simpletransfer", "err", err)
						}
						err = util.Export(context.Context, client, blockNumber, txHashes, outputPath)
						if err != nil {
							log.Crit("Failed to dump transactions", "err", err)
						}
					case "simplecall":
						err = simplecall.NewSpec(*benchConfig.SimpleCall).Run(context.Context, client)
						if err != nil {
							log.Crit("Failed to run simplecall", "err", err)
						}
					default:
						log.Crit("Unknown transaction type: %s", benchConfig.TxType)
					}

					return nil
				},
			},
			{
				Name:  "import",
				Usage: "import transactions to the chain",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "engine-url",
						Usage:    "The Engine API endpoint of the existing chain",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "engine-jwt-secret",
						Usage:    "The JWT secret used to sign JWTs for the Engine API",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "input-path",
						Usage:    "The file path for the input",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "block-txs-count",
						Usage:    "The number of transactions to include in each block",
						Required: true,
					},
				},
				Action: func(context *cli.Context) error {
					engineURL := context.String("engine-url")
					engineJWTSecret := context.String("engine-jwt-secret")
					inputPath := context.String("input-path")
					blockTxsCount := context.Int("block-txs-count")

					var jwtSecret32 [32]byte
					copy(jwtSecret32[:], hexutils.HexToBytes(engineJWTSecret))
					c, err := rpc.DialOptions(context.Context, engineURL, rpc.WithHTTPAuth(node.NewJWTAuth(jwtSecret32)))
					if err != nil {
						log.Crit("Failed to create Engine client", "err", err)
					}
					engine := ethclient.NewClient(c)

					err = util.Import(context.Context, inputPath, engine, blockTxsCount)
					if err != nil {
						log.Crit("Failed to import transactions", "err", err)
					}
					return nil
				},
			},
		},
	}

	// This is the most root context, used to propagate
	// cancellations to all spawned application-level goroutines
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		opio.BlockOnInterrupts()
		cancel()
	}()

	oplog.SetupDefaults()
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error("application failed", "err", err)
		os.Exit(1)
	}

	// TODO CTRL-C handling
	//go func() {
	//	stop := make(chan os.Signal, 1)
	//	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	//	<-stop
	//	log.Info("Received SIGINT or SIGTERM. Shutting down gracefully...")
	//	os.Exit(0)
	//}()
}
