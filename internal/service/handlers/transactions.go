package handlers

import (
    "net/http"

	"github.com/go-chi/chi"
    "gitlab.com/distributed_lab/ape"
    "gitlab.com/distributed_lab/ape/problems"
    "gitlab.com/distributed_lab/lorem"
    "gitlab.com/distributed_lab/logan/v3/errors"
)

type TransactionResponse struct {
    TxHash     string   `json:"tx_hash"`
    FromAddress string  `json:"from_address"`
    ToAddress   string  `json:"to_address"`
    Value       float64 `json:"value"`
    Timestamp   int64   `json:"timestamp"`
}

func parseTimestampRange(input string) (int64, int64, error) {
    if input == "" {
        return 0, 0, nil
    }

    times, err := lorem.ParseTimestampRange(input)
    if err != nil {
        return 0, 0, errors.Wrap(err, "failed to parse timestamp range")
    }

    return times[0].Unix(), times[1].Unix(), nil
}

func TransactionsByTime(w http.ResponseWriter, r *http.Request) {
    log := Log(r)
    startTime, endTime, err := parseTimestampRange(r)
    if err != nil {
        log.WithError(err).Error("failed to parse timestamp range")
        ape.RenderErr(w, problems.BadRequest(err))
        return
    }

    d := DB(r)
    TransactionQ := d.TransactionQ()
    
    log.WithFields(logan.F{
        "start_time": startTime,
        "end_time":   endTime,
    }).Info("retrieving transactions by time range")

    txs, err := TransactionQ().FilterByTimestamp(startTime, endTime).Select()
    if err != nil {
        log.WithError(err).Error("failed to get transactions by time range")
        ape.RenderErr(w, problems.InternalError(err))
        return
    }
    
    log.WithField("count", len(txs)).Info("transactions retrieved")
    ape.Render(w, txs)   
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
    log := Log(r)

    txHash := chi.URLParam(r, "txHash")

    d := DB(r)
    TransactionQ := d.TransactionQ()
    
    tx, err := TransactionQ().Get(txHash)
    if err != nil {
        if err == sql.ErrNoRows {
            log.WithError(err).Error("transaction not found")
            ape.RenderErr(w, problems.NotFound())
            return
        }
        log.WithError(err).Error("failed to get transaction")
        ape.RenderErr(w, problems.InternalError())
        return
    }
    
    log.WithFields(logan.F{
        "tx_hash":      tx.TxHash,
        "from_address": tx.FromAddress,
        "to_address":   tx.ToAddress,
        "value":        tx.Value,
        "timestamp":    tx.Timestamp,
    }).Info("transaction retrieved")

    resp := TransactionResponse(*tx)
    ape.Render(w, resp)
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
    log := Log(r)

    d := DB(r)
    TransactionQ := d.TransactionQ()

    from := r.Context().Value("from").(string)
    to := r.Context().Value("to").(string)

    var filters []data.TransactionFilter
    if from != "" {
        filters = append(filters, TransactionQ().FilterByFromAddress(from))
    }
    if to != "" {
        filters = append(filters, TransactionQ().FilterByToAddress(to))
    }

    log.WithFields(logan.F{
        "from":    from,
        "to":      to,
        "filters": len(filters),
    }).Info("retrieving transactions with filters")

    txs, err := TransactionQ().Select(filters...)
    if err != nil {
        log.WithError(err).Error("failed to get transactions with filters")
        ape.RenderErr(w, problems.InternalError(err))
        return
    }

    log.WithField("count", len(txs)).Info("transactions retrieved")
    ape.Render(w, txs)
}
