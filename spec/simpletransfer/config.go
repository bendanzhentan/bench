package simpletransfer

type Config struct {
	TxnsCount          int      `toml:"txnscount"`
	InitializerPrivKey string   `toml:"initializerprivkey"`
	InitialBalance     uint64   `toml:"initialbalance"`
	PrivKeys           []string `toml:"privkeys"`
	TxValue            uint64   `toml:"txvalue"`
	TxGasPrice         uint64   `toml:"txgasprice"`
	TxGasLimit         uint64   `toml:"txgaslimit"`
}
