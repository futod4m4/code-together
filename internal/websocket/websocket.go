package websocket

import (
	"github.com/futod4m4/m/config"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"net/http"
)

func NewUpgrader(cfg *config.Config) *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  cfg.WebSocketConfig.ReadBufferSize,
		WriteBufferSize: cfg.WebSocketConfig.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return cfg.WebSocketConfig.CheckOrigin
		},
	}
}

func HandleWebSocket(upgrader *websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("WebSocket Upgrade Error: ", err)
			return
		}
		defer conn.Close()

		log.Info("Client connected")

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("Read Error:", err)
				break
			}

			log.Printf("Received: %s\n", message)

			err = conn.WriteMessage(messageType, message)
			if err != nil {
				log.Error("Write Error:", err)
				break
			}
		}
	}
}
