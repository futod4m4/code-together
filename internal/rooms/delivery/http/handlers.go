package http

import (
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/futod4m4/m/pkg/httpErrors"
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

type roomHandlers struct {
	cfg    *config.Config
	roomUC rooms.RoomUseCase
	logger logger.Logger
}

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
	docState  []byte
	broadcast chan []byte
	mu        sync.RWMutex
}

var (
	roomsMap = make(map[string]*RoomConnections)
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Разрешаем соединения с любого источника. Настройте это для безопасности.
			return true
		},
	}
)

func NewUpgrader(cfg *config.Config) *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  cfg.WebSocketConfig.ReadBufferSize,
		WriteBufferSize: cfg.WebSocketConfig.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func NewRoomHandlers(cfg *config.Config, roomUC rooms.RoomUseCase, logger logger.Logger) rooms.Handlers {
	return &roomHandlers{
		cfg:    cfg,
		roomUC: roomUC,
		logger: logger,
	}
}

func (h *roomHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Create")
		defer span.Finish()

		r := &models.Room{}
		if err := c.Bind(r); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdRoom, err := h.roomUC.CreateRoom(ctx, r)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdRoom)
	}
}

func (h *roomHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Create")
		defer span.Finish()

		roomUUID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		r := &models.Room{}
		if err = c.Bind(r); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		r.ID = roomUUID

		updatedRoom, err := h.roomUC.UpdateRoom(ctx, r)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedRoom)
	}
}

func (h *roomHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Create")
		defer span.Finish()

		roomUUID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.roomUC.DeleteRoom(ctx, roomUUID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *roomHandlers) GetRoomByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Create")
		defer span.Finish()

		roomUUID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		roomByID, err := h.roomUC.GetRoomByID(ctx, roomUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, roomByID)
	}
}

func (h *roomHandlers) Join() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Join")
		defer span.Finish()
		joinCode := c.Param("join_code")
		roomByJoinCode, err := h.roomUC.GetRoomByJoinCode(ctx, joinCode)

		roomID := roomByJoinCode.ID.String()
		if _, err := uuid.Parse(roomID); err != nil {
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

		sendCurrentState(safeConn, roomConn)

		// Слушаем сообщения от клиента
		go listenForUpdates(safeConn, roomConn)

		return nil
	}
}

func (h *roomHandlers) Leave() echo.HandlerFunc {
	// TODO: implement me
	return func(c echo.Context) error {
		return nil
	}
}

func getOrCreateRoom(roomID string) *RoomConnections {
	room, exists := roomsMap[roomID]
	if !exists {
		room = &RoomConnections{
			clients:   make(map[*SafeWebSocket]bool),
			docState:  []byte{}, // Изначально пустое состояние
			broadcast: make(chan []byte),
		}
		roomsMap[roomID] = room
		go broadcastUpdates(room) // Запуск обработчика рассылки
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
		for id, r := range roomsMap {
			if r == room {
				delete(roomsMap, id)
				log.Printf("Room %s closed", id)
				break
			}
		}
	}
}

func sendCurrentState(client *SafeWebSocket, room *RoomConnections) {
	room.mu.Lock()
	defer room.mu.Unlock()
	if len(room.docState) > 0 {
		err := client.WriteMessage(websocket.BinaryMessage, room.docState)
		if err != nil {
			log.Printf("Error sending current state: %v", err)
			client.conn.Close()
			delete(room.clients, client)
		}
	}
}

func listenForUpdates(client *SafeWebSocket, room *RoomConnections) {
	defer func() {
		removeClientFromRoom(client, room)
		client.conn.Close()
	}()

	for {
		// Чтение бинарного сообщения
		messageType, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		// Игнорируем не бинарные сообщения
		if messageType != websocket.BinaryMessage {
			log.Println("Ignoring non-binary message")
			continue
		}

		// Обновление состояния документа
		room.mu.Lock()
		room.docState = applyYUpdate(room.docState, message)
		room.mu.Unlock()

		// Передача обновления в канал
		room.broadcast <- message
	}
}

func applyYUpdate(currentState, update []byte) []byte {
	return update
}

func broadcastUpdates(room *RoomConnections) {
	for update := range room.broadcast {
		room.mu.Lock()
		for client := range room.clients {
			err := client.WriteMessage(websocket.BinaryMessage, update)
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				client.conn.Close()
				delete(room.clients, client)
			}
		}
		room.mu.Unlock()
	}
}
