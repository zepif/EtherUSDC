package data_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/zepif/EtherUSDC/internal/data"
)

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

	mock.ExpectQuery(`INSERT INTO usdcTransactions`).
		WithArgs(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Value, mockTx.Timestamp).
		WillReturnRows(sqlmock.NewRows([]string{"txHash", "fromAddress", "toAddress", "value", "timestamp"}).
			AddRow(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Value, mockTx.Timestamp))

	q := pg.NewMasterQ(pgdb.New(db)).TransactionQ()
	result, err := q.Insert(mockTx)
	assert.NoError(t, err)
	assert.Equal(t, &mockTx, result)
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
			AddRow(expectedTx.TxHash, expectedTx.FromAddress, expectedTx.ToAddress, expectedTx.Value, expectedTx.Timestamp))

	q := pg.NewMasterQ(pgdb.New(db)).TransactionQ()
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
			AddRow(mockTx1.TxHash, mockTx1.FromAddress, mockTx1.ToAddress, mockTx1.Value, mockTx1.Timestamp).
			AddRow(mockTx2.TxHash, mockTx2.FromAddress, mockTx2.ToAddress, mockTx2.Value, mockTx2.Timestamp))

	q := pg.NewMasterQ(pgdb.New(db)).TransactionQ()
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
			AddRow(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Value, mockTx.Timestamp))

	q := pg.NewMasterQ(pgdb.New(db)).TransactionQ().FilterByFromAddress(mockTx.FromAddress)
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
			AddRow(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Value, mockTx.Timestamp))

	q := pg.NewMasterQ(pgdb.New(db)).TransactionQ().FilterByToAddress(mockTx.ToAddress)
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

	mock.ExpectQuery(`SELECT \* FROM usdcTransactions WHERE timestamp >= \AND timestamp <= \?`).
		WithArgs(startTimestamp, endTimestamp).
		WillReturnRows(sqlmock.NewRows([]string{"txHash", "fromAddress", "toAddress", "value", "timestamp"}).
			AddRow(mockTx.TxHash, mockTx.FromAddress, mockTx.ToAddress, mockTx.Value, mockTx.Timestamp))

	q := pg.NewMasterQ(pgdb.New(db)).TransactionQ().FilterByTimestamp(startTimestamp, endTimestamp)
	results, err := q.Select()
	assert.NoError(t, err)
	assert.Equal(t, []data.Transaction{mockTx}, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}


