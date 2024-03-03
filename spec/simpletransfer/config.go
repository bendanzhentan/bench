package simpletransfer

type Config struct {
	TxnsCount          int
	InitializerPrivKey string
	InitialBalance     uint64
	PrivKeys           []string
	TxValue            uint64
	TxGasPrice         uint64
	TxGasLimit         uint64
}
