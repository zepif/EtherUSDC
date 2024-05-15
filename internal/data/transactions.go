package data

type TransactionQ interface {
    Get(txHash, string) (*Transaction, error)
    Select(filters ...USDCTransactionFilter) ([]USDCTransaction, error)
    Insert(tx USDCTransaction) (*USDCTransaction, error)

    FilterByFromAddress(addresses ...string) USDCTransactionQ
    FilterByToAddress(addresses ...string) USDCTransactionQ
    FilterByTimestamp(start, end int64) USDCTransactionQ
}

type Transaction struct {
    txHash      string  `db: "txHash"`
    fromAddress string  `db: "fromAddress"`
    toAddress   string  `db: "toAddress"`
    values      float64 `db: "value"`
    timestamp   int64   `db: "timestamp"`
}

type TransactionFilter func(TransactionQ) TransactionQ
