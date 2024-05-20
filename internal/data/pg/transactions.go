package pg

import (
    "database/sql"
    sq "github.com/Masterminds/squirrel"
    "github.com/fatih/structs"
    "gitlab.com/distributed_lab/kit/pgdb"
    "gitlab.com/distributed_lab/logan/v3/errors"
    "github.com/zepif/EtherUSDC/internal/data"
)

const usdcTransactionsTable = "usdcTransactions"

func newTransactionQ(db *pgdb.DB) data.TransactionQ {
    return &TransactionQ{
        db:  db,
        sql: sq.StatementBuilder,
    }
}

type TransactionQ struct {
    db  *pgdb.DB
    sql sq.StatementBuilderType
}

func (q *TransactionQ) Get(txHash string) (*data.Transaction, error) {
    var tx data.Transaction
    err := q.db.Get(&tx, q.sql.Select("*").From(usdcTransactionsTable).Where(sq.Eq{"txHash": txHash}))
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &tx, nil
}

func (q *TransactionQ) Select(filters ...data.TransactionFilter) ([]data.Transaction, error) {
    var txs []data.Transaction
    stmt := q.sql.Select("*").From(usdcTransactionsTable)
    for _, filter := range filters {
        stmt = filter(q).(*TransactionQ).sql.Select()
    }
    err := q.db.Select(&txs, stmt)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return txs, nil
}

func (q *TransactionQ) Insert(tx data.Transaction) (*data.Transaction, error) {
    clauses := structs.Map(tx)
    stmt := q.sql.Insert(usdcTransactionsTable).SetMap(clauses).Suffix("RETURNING *")
    var result data.Transaction
    err := q.db.Get(&result, stmt)
    if err != nil {
        return nil, errors.Wrap(err, "failed to insert nonce to db")
    }
    return &result, nil
}

func (q *TransactionQ) FilterByFromAddress(addresses ...string) data.TransactionQ {
    q.sql = q.sql.Where(sq.Eq{"fromAddress": addresses})
    return q
}

func (q *TransactionQ) FilterByToAddress(addresses ...string) data.TransactionQ {
    q.sql = q.sql.Where(sq.Eq{"toAddress": addresses})
    return q
}

func (q *TransactionQ) FilterByTimestamp(start, end int64) data.TransactionQ {
    q.sql = q.sql.Where(sq.And{
        sq.GtOrEq{"timestamp": start},
        sq.LtOrEq{"timestamp": end},
    })
    return q
}
