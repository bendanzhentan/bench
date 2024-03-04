package core

import (
	"github.com/BurntSushi/toml"
	"keroro520/bench/spec/simplecall"
	"keroro520/bench/spec/simpletransfer"
)

type Config struct {
	TxType         string                `toml:"txtype"`
	SimpleTransfer simpletransfer.Config `toml:"simpletransfer"`
	SimpleCall     simplecall.Config     `toml:"simplecall"`
}

func LoadConfig(configFile string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(configFile, &c)
	return &c, err
}
