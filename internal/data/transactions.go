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
    TxHash      string  `db:"txHash" structs:"txHash"`
    FromAddress string  `db:"fromAddress" structs:"fromAddress"`
    ToAddress   string  `db:"toAddress" structs:"toAddress"`
    Values      float64 `db:"value" structs:"value"`
    Timestamp   int64   `db:"timestamp" structs:"timestamp"`
}

type TransactionFilter func(TransactionQ) TransactionQ
