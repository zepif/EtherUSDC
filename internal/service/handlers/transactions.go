package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/zepif/EtherUSDC/internal/data"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func TransactionsByTime(w http.ResponseWriter, r *http.Request) {
	log := Log(r)

	startTimeStr := chi.URLParam(r, "startTime")
	endTimeStr := chi.URLParam(r, "endTime")

	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		log.WithError(err).Error("failed to parse start timestamp")
		return
	}

	endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
	if err != nil {
		log.WithError(err).Error("failed to parse end timestamp")
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
		"txHash":      (*tx).TxHash,
		"fromAddress": (*tx).FromAddress,
		"toAddress":   (*tx).ToAddress,
		"value":       (*tx).Values,
		"timestamp":   (*tx).Timestamp,
	}).Info("transaction retrieved")

	ape.Render(w, *tx)
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	log := Log(r)
	d := DB(r).TransactionQ()

	fromAddress := r.URL.Query().Get("fromAddress")
	toAddress := r.URL.Query().Get("toAddress")
	var filters []data.TransactionFilter
	if fromAddress != "" {
		d = d.FilterByFromAddress(fromAddress)
	}

	if toAddress != "" {
		d = d.FilterByToAddress(toAddress)
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
