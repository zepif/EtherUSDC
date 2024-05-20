package pg

import (
    pgdb "gitlab.com/distributed_lab/kit/pgdb"
    "github.com/zepif/EtherUSDC/internal/data"
)

func NewMasterQ(db *pgdb.DB) data.MasterQ {
    return &masterQ{
        db: db.Clone(),
    }
}

type masterQ struct {
    db *pgdb.DB
}

func (m *masterQ) New() data.MasterQ {
    return NewMasterQ(m.db)
}

func (m *masterQ) Transaction(fn func(q data.MasterQ) error) error {
    return m.db.Transaction(func() error {
        return fn(m)
    })
}

func (m *masterQ) TransactionQ() data.TransactionQ {
    return newTransactionQ(m.db)
}
