package workers

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zepif/EtherUSDC/internal/data"
	"github.com/zepif/EtherUSDC/internal/service/eth"
	store "github.com/zepif/EtherUSDC/internal/store"
	"gitlab.com/distributed_lab/logan/v3"
)

type TransactionWorker struct {
	log        *logan.Entry
	db         data.MasterQ
	client     *eth.EthClient
	ctx        context.Context
	cancel     context.CancelFunc
	startBlock uint64
}

func NewTransactionWorker(log *logan.Entry, db data.MasterQ, client *eth.EthClient, startBlock uint64) *TransactionWorker {
	ctx, cancel := context.WithCancel(context.Background())
	log.WithField("context", fmt.Sprintf("%p", ctx)).Info("Created new context")
	return &TransactionWorker{
		log:        log.WithField("worker", "usdctransactions"),
		db:         db,
		client:     client,
		ctx:        ctx,
		cancel:     cancel,
		startBlock: startBlock,
	}
}

func (w *TransactionWorker) Start() error {
	logs := make(chan types.Log, 100)

	sub, err := w.client.SubscribeLogs(w.ctx, logs, w.startBlock)
	if err != nil {
		w.log.WithError(err).Error("Failed to subscribe to logs")
		return err
	}

	go w.consumeLogs(logs)

	w.log.Info("SubscribeLogs started successfully")

	go func() {
		for {
			select {
			case err := <-sub.Err():
				w.log.WithError(err).Error("Subscription error")
				return
			case <-w.ctx.Done():
				w.log.Info("Context canceled, stopping subscription")
				sub.Unsubscribe()
				return
			}
		}
	}()

	w.log.Info("LoadBlocks started successfully")
	return nil
}

func (w *TransactionWorker) Stop() {
	w.cancel()
	w.log.WithField("context", fmt.Sprintf("%p", w.ctx)).Info("Created new context")
}

func (w *TransactionWorker) consumeLogs(logs <-chan types.Log) {
	defer func() {
		if r := recover(); r != nil {
			w.log.WithField("panic", r).Error("Panic in consumeLogs goroutine")
			debug.PrintStack()
		}
	}()

	for {
		w.log.Info("Waiting for log event")
		select {
		case <-w.ctx.Done():
			w.log.Info("Saving USDC transaction")
			return
		case vLog := <-logs:
			w.log.WithFields(logan.F{
				"txHash":      vLog.TxHash.Hex(),
				"blockNumber": vLog.BlockNumber,
			}).Debug("Received log event")

			event, err := w.client.ParseTransferEvent(vLog)
			if err != nil {
				w.log.WithError(err).Error("failed to parse transfer event")
				continue
			}

			w.log.WithFields(logan.F{
				"from":  event.From.Hex(),
				"to":    event.To.Hex(),
				"value": event.Value.Int64(),
			}).Info("Saving USDC transaction")

			err = w.saveTransaction(vLog, *event)
			if err != nil {
				w.log.WithError(err).Error("failed to save USDC transaction")
			}
		}
	}
}

func (w *TransactionWorker) saveTransaction(vLog types.Log, event store.StoreTransfer) error {
	w.log.Info("Saving USDC transaction")

	tx := data.Transaction{
		TxHash:      vLog.TxHash.Hex(),
		FromAddress: event.From.Hex(),
		ToAddress:   event.To.Hex(),
		Values:      float64(event.Value.Int64()),
		Timestamp:   time.Now().Unix(),
		BlockNumber: int64(vLog.BlockNumber),
	}

	existingTxs, err := w.db.TransactionQ().Get(tx.TxHash)
	if err != nil {
		w.log.WithError(err).Error("failed to check for existing transactions")
		return errors.Wrap(err, "failed to check for existing transactions")
	}

	for _, existingTx := range existingTxs {
		if existingTx.FromAddress == tx.FromAddress && existingTx.ToAddress == tx.ToAddress &&
			existingTx.Values == tx.Values && existingTx.Timestamp == tx.Timestamp && existingTx.BlockNumber == tx.BlockNumber {
			w.log.WithField("txHash", tx.TxHash).Warn("transaction event already exists")
			return nil
		}
	}

	_, err = w.db.TransactionQ().Insert(tx)
	if err != nil {
		w.log.WithError(err).WithFields(logan.F{
			"txHash":      tx.TxHash,
			"fromAddress": tx.FromAddress,
			"toAddress":   tx.ToAddress,
			"value":       tx.Values,
			"timestamp":   tx.Timestamp,
			"blockNumber": tx.BlockNumber,
		}).Error("failed to save USDC transaction")

		return errors.Wrap(err, "failed to insert transaction into database")
	}

	return nil
}
