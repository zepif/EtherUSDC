package data

type TransactionQ interface {
	Get(txHash string) ([]Transaction, error)
	Select(filters ...TransactionFilter) ([]Transaction, error)
	Insert(tx Transaction) (*Transaction, error)

	FilterByFromAddress(addresses ...string) TransactionQ
	FilterByToAddress(addresses ...string) TransactionQ
	FilterByTimestamp(start, end int64) TransactionQ
	FilterByTxHash(txHash string) TransactionQ
}

type Transaction struct {
	ID          int64   `db:"id" structs:"id"`
	TxHash      string  `db:"txhash" structs:"txhash"`
	FromAddress string  `db:"fromaddress" structs:"fromaddress"`
	ToAddress   string  `db:"toaddress" structs:"toaddress"`
	Values      float64 `db:"value" structs:"value"`
	Timestamp   int64   `db:"timestamp" structs:"timestamp"`
}

type TransactionFilter func(TransactionQ) TransactionQ
