package handlers

import (
    "database/sql"
    "net/http"
    "time"
    
    "github.com/zepif/EtherUSDC/internal/data"
    "github.com/go-chi/chi"
    "gitlab.com/distributed_lab/ape"
    "gitlab.com/distributed_lab/ape/problems"
    "gitlab.com/distributed_lab/logan/v3"
    "gitlab.com/distributed_lab/logan/v3/errors"
)

func parseTimestampRange(input string) (int64, int64, error) {
    if input == "" {
        return 0, 0, nil
    }

    times, err := time.ParseInLocation("2006-01-02T15:04:05Z", input, time.UTC)
    if err != nil {
        return 0, 0, errors.Wrap(err, "failed to parse timestamp range")
    }

    startTime := times.UnixNano() / int64(time.Millisecond)
    endTime := startTime + int64(time.Hour/time.Millisecond)

    return startTime, endTime, nil
}

func TransactionsByTime(w http.ResponseWriter, r *http.Request) {
    log := Log(r)

    startTime, endTime, err := parseTimestampRange(r.URL.Query().Get("timestamp"))
    if err != nil {
        log.WithError(err).Error("failed to parse timestamp range")
        ape.RenderErr(w, problems.InternalError())
        return
    }

    d := DB(r)
    txs, err := d.TransactionQ().FilterByTimestamp(startTime, endTime).Select()
    if err != nil {
        log.WithError(err).Error("failed to get transactions by time range")
        ape.RenderErr(w, problems.InternalError())
        return
    }

    log.WithField("count", len(txs)).Info("transactions retrieved")
    ape.Render(w, txs)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
    log := Log(r)
    txHash := chi.URLParam(r, "txHash")

    db := DB(r)
    tx, err := db.TransactionQ().Get(txHash)
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
        "txHash":      (*tx).txHash,
        "fromAddress": (*tx).fromAddress,
        "toAddress":   (*tx).toAddress,
        "value":       (*tx).values,
        "timestamp":   (*tx).timestamp,
    }).Info("transaction retrieved")

    ape.Render(w, *tx)
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
    log := Log(r)
    d := DB(r)

    fromAddress := r.URL.Query().Get("fromAddress")
    toAddress := r.URL.Query().Get("toAddress")
    var filters []data.TransactionFilter
    if fromAddress != "" {
        d = d.TransactionQ().FilterByFromAddress(fromAddress)
    }

    if toAddress != "" {
        d = d.TransactionQ().FilterByToAddress(toAddress)
    }

    log.WithFields(logan.F{
        "from":    fromAddress,
        "to":      toAddress,
        "filters": len(filters),
    }).Info("retrieving transactions with filters")

    txs, err := d.Select(filters...)
    if err != nil {
        log.WithError(err).Error("failed to get transactions with filters")
        ape.RenderErr(w, problems.InternalError())
        return
    }

    log.WithField("count", len(txs)).Info("transactions retrieved")
    ape.Render(w, txs)
}
