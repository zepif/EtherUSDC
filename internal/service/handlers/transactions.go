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
	Offset      int64   `url:"offset"`
	Limit       int64   `url:"limit"`
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	log := Log(r)
	d := DB(r).TransactionQ()

	request := TransactionsListRequest{Offset: 0, Limit: 50}

	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		log.WithError(err).Error("failed to decode request parameters")
		ape.RenderErr(w, problems.InternalError())
		return
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
	if request.StartTime != nil {
		d = d.FilterByTimestampStart(*request.StartTime)
	}
	if request.EndTime != nil {
		d = d.FilterByTimestampEnd(*request.EndTime)
	}
	d = d.Page(uint64(request.Limit), uint64(request.Offset))

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
		"data": txs,
	})
}
