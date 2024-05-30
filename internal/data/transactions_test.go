package data_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zepif/EtherUSDC/internal/data"
	"github.com/zepif/EtherUSDC/internal/data/pg"
)

type MockDB struct {
	SQL *sql.DB
}

func NewMockDB(db *sql.DB) *MockDB {
	return &MockDB{SQL: db}
}

func (db *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.SQL.Query(query, args...)
}

func (db *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.SQL.Exec(query, args...)
}

func TestTransactionQ_Insert(t *testing.T) {
	mockTx := data.Transaction{
		TxHash:      "0x123",
		FromAddress: "0xabc",
		ToAddress:   "0xdef",
		Values:      100.0,
		Timestamp:   1234567890,
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`INSERT INTO usdcTransactions \(txHash, fromAddress, toAddress, value, timestamp\)`).
		WithArgs(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Values, mockTx.Timestamp).
		WillReturnResult(sqlmock.NewResult(1, 1))

	q := pg.NewMasterQ(NewMockDB(db)).TransactionQ()
	err = q.Insert(mockTx)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionQ_Get(t *testing.T) {
	txHash := "0x123"
	expectedTx := data.Transaction{
		TxHash:      txHash,
		FromAddress: "0xabc",
		ToAddress:   "0xdef",
		Values:      100.0,
		Timestamp:   1234567890,
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM usdcTransactions WHERE txHash = \?`).
		WithArgs(txHash).
		WillReturnRows(sqlmock.NewRows([]string{"txHash", "fromAddress", "toAddress", "value", "timestamp"}).
			AddRow(expectedTx.TxHash, expectedTx.FromAddress, expectedTx.ToAddress, expectedTx.Values, expectedTx.Timestamp))

	q := pg.NewMasterQ(NewMockDB(db)).TransactionQ()
	result, err := q.Get(txHash)
	assert.NoError(t, err)
	assert.Equal(t, &expectedTx, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionQ_Select(t *testing.T) {
	mockTx1 := data.Transaction{
		TxHash:      "0x123",
		FromAddress: "0xabc",
		ToAddress:   "0xdef",
		Values:      100.0,
		Timestamp:   1234567890,
	}
	mockTx2 := data.Transaction{
		TxHash:      "0x456",
		FromAddress: "0xghi",
		ToAddress:   "0xjkl",
		Values:      200.0,
		Timestamp:   1234567891,
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM usdcTransactions`).
		WillReturnRows(sqlmock.NewRows([]string{"txHash", "fromAddress", "toAddress", "value", "timestamp"}).
			AddRow(mockTx1.TxHash, mockTx1.FromAddress, mockTx1.ToAddress, mockTx1.Values, mockTx1.Timestamp).
			AddRow(mockTx2.TxHash, mockTx2.FromAddress, mockTx2.ToAddress, mockTx2.Values, mockTx2.Timestamp))

	q := pg.NewMasterQ(NewMockDB(db)).TransactionQ()
	results, err := q.Select()
	assert.NoError(t, err)
	assert.Equal(t, []data.Transaction{mockTx1, mockTx2}, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionQ_FilterByFromAddress(t *testing.T) {
	mockTx := data.Transaction{
		TxHash:      "0x123",
		FromAddress: "0xabc",
		ToAddress:   "0xdef",
		Values:      100.0,
		Timestamp:   1234567890,
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM usdcTransactions WHERE fromAddress = \?`).
		WithArgs(mockTx.FromAddress).
		WillReturnRows(sqlmock.NewRows([]string{"txHash", "fromAddress", "toAddress", "value", "timestamp"}).
			AddRow(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Values, mockTx.Timestamp))

	q := pg.NewMasterQ(NewMockDB(db)).TransactionQ().FilterByFromAddress(mockTx.FromAddress)
	results, err := q.Select()
	assert.NoError(t, err)
	assert.Equal(t, []data.Transaction{mockTx}, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionQ_FilterByToAddress(t *testing.T) {
	mockTx := data.Transaction{
		TxHash:      "0x123",
		FromAddress: "0xabc",
		ToAddress:   "0xdef",
		Values:      100.0,
		Timestamp:   1234567890,
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM usdcTransactions WHERE toAddress = \?`).
		WithArgs(mockTx.ToAddress).
		WillReturnRows(sqlmock.NewRows([]string{"txHash", "fromAddress", "toAddress", "value", "timestamp"}).
			AddRow(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Values, mockTx.Timestamp))

	q := pg.NewMasterQ(NewMockDB(db)).TransactionQ().FilterByToAddress(mockTx.ToAddress)
	results, err := q.Select()
	assert.NoError(t, err)
	assert.Equal(t, []data.Transaction{mockTx}, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionQ_FilterByTimestamp(t *testing.T) {
	mockTx := data.Transaction{
		TxHash:      "0x123",
		FromAddress: "0xabc",
		ToAddress:   "0xdef",
		Values:      100.0,
		Timestamp:   1234567890,
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	startTimestamp := int64(1234567890)
	endTimestamp := int64(1234567890)

	mock.ExpectQuery(`SELECT \* FROM usdcTransactions WHERE timestamp >= \? AND timestamp <= \?`).
		WithArgs(startTimestamp, endTimestamp).
		WillReturnRows(sqlmock.NewRows([]string{"txHash", "fromAddress", "toAddress", "value", "timestamp"}).
			AddRow(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Values, mockTx.Timestamp))

	q := pg.NewMasterQ(NewMockDB(db)).TransactionQ().FilterByTimestamp(startTimestamp, endTimestamp)
	results, err := q.Select()
	assert.NoError(t, err)
	assert.Equal(t, []data.Transaction{mockTx}, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

