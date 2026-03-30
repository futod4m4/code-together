package ws

import (
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"log"
	"net/http"
	"sync"
)

type SafeWebSocket struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (sw *SafeWebSocket) WriteMessage(messageType int, data []byte) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.conn.WriteMessage(messageType, data)
}

type RoomConnections struct {
	clients   map[*SafeWebSocket]bool
	broadcast chan broadcastMsg
	mu        sync.RWMutex
}

type broadcastMsg struct {
	data   []byte
	sender *SafeWebSocket
}

var (
	roomsMu  sync.Mutex
	roomsMap = make(map[string]*RoomConnections)
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type roomWSHandlers struct {
	cfg    *config.Config
	roomUC rooms.RoomUseCase
	logger logger.Logger
}

func NewRoomWSHandlers(cfg *config.Config, roomUC rooms.RoomUseCase, logger logger.Logger) rooms.WSHandlers {
	return &roomWSHandlers{
		cfg:    cfg,
		roomUC: roomUC,
		logger: logger,
	}
}

func NewUpgrader(cfg *config.Config) *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  cfg.WebSocketConfig.ReadBufferSize,
		WriteBufferSize: cfg.WebSocketConfig.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (h *roomWSHandlers) Join() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Join")
		defer span.Finish()
		joinCode := c.Param("join_code")
		roomByJoinCode, err := h.roomUC.GetRoomByJoinCode(ctx, joinCode)
		if err != nil {
			log.Println("Failed to get room:", err)
			return c.JSON(http.StatusBadRequest, "Invalid room")
		}

		roomID := roomByJoinCode.ID.String()
		if _, err = uuid.Parse(roomID); err != nil {
			log.Println("Invalid room ID:", err)
			return c.JSON(http.StatusBadRequest, "Invalid room ID")
		}

		roomConn := getOrCreateRoom(roomID)
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return c.JSON(http.StatusInternalServerError, "WebSocket upgrade failed")
		}

		safeConn := &SafeWebSocket{conn: conn}
		addClientToRoom(safeConn, roomConn)
		log.Printf("Client connected to room %s", roomID)

		go listenForUpdates(safeConn, roomConn)

		return nil
	}
}

func (h *roomWSHandlers) Leave() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func getOrCreateRoom(roomID string) *RoomConnections {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	room, exists := roomsMap[roomID]
	if !exists {
		room = &RoomConnections{
			clients:   make(map[*SafeWebSocket]bool),
			broadcast: make(chan broadcastMsg, 256),
		}
		roomsMap[roomID] = room
		go broadcastUpdates(room, roomID)
	}
	return room
}

func addClientToRoom(client *SafeWebSocket, roomConn *RoomConnections) {
	roomConn.mu.Lock()
	defer roomConn.mu.Unlock()
	roomConn.clients[client] = true
}

func removeClientFromRoom(client *SafeWebSocket, room *RoomConnections) {
	room.mu.Lock()
	defer room.mu.Unlock()
	delete(room.clients, client)
	if len(room.clients) == 0 {
		close(room.broadcast)
		roomsMu.Lock()
		for id, r := range roomsMap {
			if r == room {
				delete(roomsMap, id)
				log.Printf("Room %s closed", id)
				break
			}
		}
		roomsMu.Unlock()
	}
}

func listenForUpdates(client *SafeWebSocket, room *RoomConnections) {
	defer func() {
		removeClientFromRoom(client, room)
		client.conn.Close()
	}()

	for {
		messageType, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		if messageType != websocket.BinaryMessage {
			continue
		}

		room.broadcast <- broadcastMsg{data: message, sender: client}
	}
}

func broadcastUpdates(room *RoomConnections, roomID string) {
	for msg := range room.broadcast {
		room.mu.RLock()
		for client := range room.clients {
			if client == msg.sender {
				continue
			}
			err := client.WriteMessage(websocket.BinaryMessage, msg.data)
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				client.conn.Close()
			}
		}
		room.mu.RUnlock()
	}
}
