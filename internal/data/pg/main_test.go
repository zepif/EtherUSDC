package pg_test

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/kit/pgdb"
	"github.com/zepif/EtherUSDC/internal/data"
	"github.com/zepif/EtherUSDC/internal/data/pg"
)

func TestMasterQ_New(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	masterQ := pg.NewMasterQ(pgdb.New(db))
	newMasterQ := masterQ.New()

	assert.NotNil(t, newMasterQ)
	assert.IsType(t, &pg.MasterQ{}, newMasterQ)
}

func TestMasterQ_Transaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	masterQ := pg.NewMasterQ(pgdb.New(db))
	
    err = masterQ.Transaction(func(q data.MasterQ) error {
        tx := data.Transaction{
            TxHash:      "0xnewhash",
            FromAddress: "0xfromaddress",
            ToAddress:   "0xtoaddress",
            Value:       150.0,
            Timestamp:   time.Now().Unix(),
        }

        insertedTx, err := q.TransactionQ().Insert(tx)
        if err != nil {
            return err
        }

        fetchedTx, err := q.TransactionQ().Get(insertedTx.TxHash)
        if err != nil {
            return err
        }

        if fetchedTx.TxHash != insertedTx.TxHash {
            return errors.New("fetched transaction does not match inserted transaction")
        }

        return nil
    })

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMasterQ_TransactionQ(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	masterQ := pg.NewMasterQ(pgdb.New(db))
	transactionQ := masterQ.TransactionQ()

	assert.NotNil(t, transactionQ)
	assert.IsType(t, &pg.TransactionQ{}, transactionQ)
}

