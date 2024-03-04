package simplecall

type Config struct {
	TxnsCount       int      `toml:"txnscount"`
	PrivKeys        []string `toml:"privkeys"`
	ContractAddress string   `toml:"contractaddress"`
	Calldata        string   `toml:"calldata"`
	TxValue         uint64   `toml:"txvalue"`
	TxGasPrice      uint64   `toml:"txgasprice"`
	TxGasLimit      uint64   `toml:"txgaslimit"`
}
