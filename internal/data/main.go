package data

type MasterQ interface {
    New() MasterQ
    Transaction(fn func(db MasterQ) error) error
    TransactionQ() TransactionQ
}
