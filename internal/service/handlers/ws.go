package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log := Log(r)
		log.WithError(err).Error("failed to upgrade connection to WebSocket")
		return
	}
	defer conn.Close()

	log := Log(r)
	db := DB(r).TransactionQ()

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
