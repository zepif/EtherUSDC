package workers

import (
    "context"
    "time"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/zepif/EtherUSDC/internal/data"
    "github.com/zepif/EtherUSDC/internal/service/eth"
    "gitlab.com/distributed_lab/logan/v3"
)

type TransactionWorker struct {
    log     *logan.Entry
    db      data.MasterQ
    client  *eth.EthClient
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewTransactionWorker(log *logan.Entry, db data.MasterQ, client *eth.EthClient) *TransactionWorker {
    ctx, cancel := context.WithCancel(context.Background())
    return &TransactionWorker{
        log:    log.WithField("worker", "usdctransactions"),
        db:     db,
        client: client,
        ctx:    ctx,
        cancel: cancel,
    }
}

func (w *TransactionWorker) Start() error {
    logs := make(chan types.Log)
    go w.consumeLogs(logs)

    err := w.client.ListenToEvents(w.ctx, logs)
    if err != nil {
        return err
    }

    return nil
}

func (w *TransactionWorker) Stop() {
    w.cancel()
}

func (w *TransactionWorker) consumeLogs(logs <-chan types.Log) {
    for {
        select {
        case <-w.ctx.Done():
            return
        case vLog := <-logs:
            event, err := w.client.ParseTransferEvent(vLog)
            if err != nil {
                w.log.WithError(err).Error("failed to parse transfer event")
                continue
            }
            w.saveTransaction(event)
        }
    }
}

func (w *TransactionWorker) saveTransaction(event *eth.TransferEvent) {
    tx := data.Transaction{
        TxHash:      vLog.TxHash.Hex(),
        FromAddress: event.From.Hex(),
        ToAddress:   event.To.Hex(),
        Value:       float64(event.Value.Int64()),
        Timestamp:   time.Now().Unix(),
    }

    _, err := w.db.TransactionQ().Insert(tx)
    if err != nil {
        w.log.WithError(err).Error("failed to save USDC transaction")
    }
}
