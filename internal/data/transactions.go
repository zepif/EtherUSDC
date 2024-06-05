package data

type TransactionQ interface {
	Get(txHash string) ([]Transaction, error)
	Select(filters ...TransactionFilter) ([]Transaction, error)
	Insert(tx Transaction) (*Transaction, error)

	FilterByFromAddress(addresses ...string) TransactionQ
	FilterByToAddress(addresses ...string) TransactionQ
	FilterByTimestamp(start, end int64) TransactionQ
	FilterByTxHash(txHash string) TransactionQ
	FilterByBlockNumber(blockNumber int64) TransactionQ

	// Page(pageParams pgdb.OffsetPageParams) TransactionQ
	// GetTotalCount() (int64, error)
}

type Transaction struct {
	ID          int64   `db:"id" structs:"id"`
	TxHash      string  `db:"tx_hash" structs:"tx_hash"`
	FromAddress string  `db:"from_address" structs:"from_address"`
	ToAddress   string  `db:"to_address" structs:"to_address"`
	Values      float64 `db:"value" structs:"value"`
	Timestamp   int64   `db:"timestamp" structs:"timestamp"`
	BlockNumber int64   `db:"block_number" structs:"block_number"`
}

type TransactionFilter func(TransactionQ) TransactionQ
