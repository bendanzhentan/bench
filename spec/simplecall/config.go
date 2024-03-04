package simplecall

type Config struct {
	TxnsCount       int
	PrivKeys        []string
	ContractAddress string
	Calldata        string
	TxValue         uint64
	TxGasPrice      uint64
	TxGasLimit      uint64
}
