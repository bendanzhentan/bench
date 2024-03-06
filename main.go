package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"keroro520/bench/core"
	"keroro520/bench/spec/simplecall"
	"keroro520/bench/spec/simpletransfer"
	"os"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "rpc-url",
				Usage:    "The RPC endpoint of the existing chain",
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
			configPath := context.String("config-path")
			_ = context.String("output-path")

			benchConfig, err := core.LoadConfig(configPath)
			if err != nil {
				log.Crit("Failed to load configuration file", "err", err)
			}

			client, err := ethclient.Dial(rpcURL)
			if err != nil {
				log.Crit("Failed to connect to the Ethereum client", "url", rpcURL, "err", err)
			}

			switch benchConfig.TxType {
			case "simpletransfer":
				err = simpletransfer.NewSpec(*benchConfig.SimpleTransfer).Run(context.Context, client)
				if err != nil {
					log.Crit("Failed to run simpletransfer", "err", err)
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
	}

	//go func() {
	//	stop := make(chan os.Signal, 1)
	//	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	//	<-stop
	//	log.Info("Received SIGINT or SIGTERM. Shutting down gracefully...")
	//	os.Exit(0)
	//}()

	err := app.Run(os.Args)
	if err != nil {
		log.Crit("Failed to run application", "err", err)
	}
}
