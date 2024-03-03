package core

import (
	"github.com/BurntSushi/toml"
	"keroro520/bench/spec/simplecall"
	"keroro520/bench/spec/simpletransfer"
)

type Config struct {
	TxType         string
	SimpleTransfer simpletransfer.Config
	SimpleCall     simplecall.Config
}

func LoadConfig(configFile string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(configFile, &c)
	return &c, err
}
