package pg

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/zepif/EtherUSDC/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const usdcTransactionsTable = "usdcTransactions"

func newTransactionQ(db *pgdb.DB) data.TransactionQ {
	return &TransactionQ{
		db:  db,
		sql: sq.StatementBuilder.Select("*").From(usdcTransactionsTable),
	}
}

type TransactionQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *TransactionQ) GetByID(id int64) (*data.Transaction, error) {
	var tx data.Transaction
	query := q.sql.Where(sq.Eq{"id": id})
	err := q.db.Get(&tx, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (q *TransactionQ) Get(txHash string) ([]data.Transaction, error) {
	var txs []data.Transaction
	query := q.sql.Where(sq.Eq{"tx_hash": txHash})
	err := q.db.Select(&txs, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (q *TransactionQ) Select(filters ...data.TransactionFilter) ([]data.Transaction, error) {
	var txs []data.Transaction
	stmt := q.sql
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
	clauses := structs.Map(tx)
	delete(clauses, "id")
	stmt := sq.Insert(usdcTransactionsTable).SetMap(clauses).Suffix("RETURNING *")
	var result data.Transaction
	err := q.db.Get(&result, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert transaction to db")
	}
	return &result, nil
}

func (q *TransactionQ) FilterByFromAddress(addresses ...string) data.TransactionQ {
	q.sql = q.sql.Where(sq.Eq{"from_address": addresses})
	return q
}

func (q *TransactionQ) FilterByToAddress(addresses ...string) data.TransactionQ {
	q.sql = q.sql.Where(sq.Eq{"to_address": addresses})
	return q
}

func (q *TransactionQ) FilterByTimestampStart(start int64) data.TransactionQ {
	q.sql = q.sql.Where(sq.GtOrEq{"timestamp": start})
	return q
}

func (q *TransactionQ) FilterByTimestampEnd(end int64) data.TransactionQ {
	q.sql = q.sql.Where(sq.LtOrEq{"timestamp": end})
	return q
}

func (q *TransactionQ) FilterByTxHash(txHash string) data.TransactionQ {
	q.sql = q.sql.Where(sq.Eq{"tx_hash": txHash})
	return q
}

func (q *TransactionQ) FilterByBlockNumber(blockNumber int64) data.TransactionQ {
	q.sql = q.sql.Where(sq.GtOrEq{"block_number": blockNumber})
	return q
}

func (q *TransactionQ) Page(limit, offset uint64) data.TransactionQ {
	q.sql = q.sql.OrderBy("id").Limit(limit).Offset(offset)
	return q
}
