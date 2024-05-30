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
	"gitlab.com/distributed_lab/logan/v3"
)

type TransactionWorker struct {
	log    *logan.Entry
	db     data.MasterQ
	client *eth.EthClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTransactionWorker(log *logan.Entry, db data.MasterQ, client *eth.EthClient) *TransactionWorker {
	ctx, cancel := context.WithCancel(context.Background())
	log.WithField("context", fmt.Sprintf("%p", ctx)).Info("Created new context")
	return &TransactionWorker{
		log:    log.WithField("worker", "usdctransactions"),
		db:     db,
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *TransactionWorker) Start() error {
	logs := make(chan types.Log, 100)
	go w.consumeLogs(logs)
	err := w.client.ListenToEvents(w.ctx, logs)
	// w.consumeLogs(logs)
	if err != nil {
		w.log.WithError(err).Error("ListenToEvents failed")
		return err
	}

	w.log.Info("ListenToEvents started successfully")
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

			err = w.saveTransaction(vLog, event)
			if err != nil {
				w.log.WithError(err).Error("failed to save USDC transaction")
			}
		}
	}
}

func (w *TransactionWorker) saveTransaction(vLog types.Log, event *eth.TransferEvent) error {
	w.log.Info("Saving USDC transaction")

	tx := data.Transaction{
		TxHash:      vLog.TxHash.Hex(),
		FromAddress: event.From.Hex(),
		ToAddress:   event.To.Hex(),
		Values:      float64(event.Value.Int64()),
		Timestamp:   time.Now().Unix(),
	}

	_, err := w.db.TransactionQ().Insert(tx)
	if err != nil {
		w.log.WithError(err).WithFields(logan.F{
			"txHash":      tx.TxHash,
			"fromAddress": tx.FromAddress,
			"toAddress":   tx.ToAddress,
			"value":       tx.Values,
			"timestamp":   tx.Timestamp,
		}).Error("failed to save USDC transaction")

		return errors.Wrap(err, "failed to insert transaction into database")
	}

	return nil
}
