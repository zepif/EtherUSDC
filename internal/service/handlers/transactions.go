package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/urlval"
)

type TransactionsListRequest struct {
	FromAddress *string `url:"fromAddress"`
	ToAddress   *string `url:"toAddress"`
	BlockNumber *int64  `url:"blockNumber"`
	TxHash      *string `url:"txHash"`
	StartTime   *int64  `url:"startTime"`
	EndTime     *int64  `url:"endTime"`
	Offset      *int    `url:"offset"`
	Limit       *int    `url:"limit"`
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	log := Log(r)
	d := DB(r).TransactionQ()

	var request TransactionsListRequest
	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		log.WithError(err).Error("failed to decode request parameters")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if request.Offset == nil {
		defaultOffset := 1
		request.Offset = &defaultOffset
	}
	if request.Limit == nil {
		defaultLimit := 10
		request.Limit = &defaultLimit
	}

	if request.FromAddress != nil {
		d = d.FilterByFromAddress(*request.FromAddress)
	}
	if request.ToAddress != nil {
		d = d.FilterByToAddress(*request.ToAddress)
	}
	if request.BlockNumber != nil {
		d = d.FilterByBlockNumber(*request.BlockNumber)
	}
	if request.TxHash != nil {
		d = d.FilterByTxHash(*request.TxHash)
	}
	if request.StartTime != nil && request.EndTime != nil {
		d = d.FilterByTimestamp(*request.StartTime, *request.EndTime)
	}
	d = d.Page(*request.Limit, *request.Offset)

	log.WithFields(logan.F{
		"from":        request.FromAddress,
		"to":          request.ToAddress,
		"blockNumber": request.BlockNumber,
		"txHash":      request.TxHash,
		"startTime":   request.StartTime,
		"endTime":     request.EndTime,
		"limit":       request.Limit,
		"offset":      request.Offset,
	}).Info("retrieving transactions with filters")

	txs, err := d.Select()
	if err != nil {
		log.WithError(err).Error("failed to get transactions with filters")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	log.WithField("count", len(txs)).Info("transactions retrieved")
	ape.Render(w, map[string]interface{}{
		"transactions": txs,
		"limit":        request.Limit,
		"offset":       request.Offset,
	})
}
