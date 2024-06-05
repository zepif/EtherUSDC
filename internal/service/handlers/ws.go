package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zepif/EtherUSDC/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type (
	ctxLog struct{}
	ctxDB  struct{}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log := r.Context().Value(ctxLog{}).(logan.Entry)
		log.WithError(err).Error("failed to upgrade connection to WebSocket")
		return
	}
	defer conn.Close()

	log := r.Context().Value(ctxLog{}).(logan.Entry)
	db := r.Context().Value(ctxDB{}).(data.TransactionQ)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.WithError(err).Error("failed to read WebSocket message")
			break
		}

		transactions, err := db.Select()
		if err != nil {
			log.WithError(err).Error("failed to get transactions from database")
			break
		}

		for _, tx := range transactions {
			txJSON, err := json.Marshal(tx)
			if err != nil {
				log.WithError(err).Error("failed to marshal transaction to JSON")
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, txJSON)
			if err != nil {
				log.WithError(err).Error("failed to send transaction data over WebSocket")
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}
