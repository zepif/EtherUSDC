package data

type TransactionQ interface {
    Get(txHash string) (*Transaction, error)
    Select(filters ...TransactionFilter) ([]Transaction, error)
    Insert(tx Transaction) (*Transaction, error)

    FilterByFromAddress(addresses ...string) TransactionQ
    FilterByToAddress(addresses ...string) TransactionQ
    FilterByTimestamp(start, end int64) TransactionQ
}

type Transaction struct {
    txHash      string  `db:"txHash" structs:"txHash"`
    fromAddress string  `db:"fromAddress" structs:"fromAddress"`
    toAddress   string  `db:"toAddress" structs:"toAddress"`
    values      float64 `db:"value" structs:"value"`
    timestamp   int64   `db:"timestamp" structs:"timestamp"`
}

type TransactionFilter func(TransactionQ) TransactionQ
