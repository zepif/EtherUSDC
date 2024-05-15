package pg

import (
    "database/sql"
    sq "github.com/Masterminds/squirrel"
    "gitlab.com/distributed_lab/kit/pgdb"
    "gitlab.com/your-project/internal/data"
)

const uscTransactionsTable = "usdcTransactions"

func newTransactionQ(db *pgdb.DB) data.TransactionQ {
    return &usdcTransactionQ{
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
    err := q.db.Get(&tx, q.sql.Select("*").From(usdcTransactionsTable).Where(sq.Eq{"tx_hash": txHash}))
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
        stmt = filter(q).(*TransactionQ).sql
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
    stmt := q.sql.Insert(usdcTransactionsTable).SetMap(map[string]interface{}{
        "tx_hash":     tx.TxHash,
        "from_address": tx.FromAddress,
        "to_address":   tx.ToAddress,
        "value":        tx.Value,
        "timestamp":    tx.Timestamp,
    }).Suffix("RETURNING *")
    var result data.Transaction
    err := q.db.Get(&result, stmt)
    if err != nil {
        return nil, err
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
