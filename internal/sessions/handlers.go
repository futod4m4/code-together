package sessions

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type CodingSession struct {
	ID         uuid.UUID  `json:"session_id" db:"session_id"`
	RoomID     uuid.UUID  `json:"room_id" db:"room_id"`
	Title      string     `json:"title" db:"title"`
	StartedBy  uuid.UUID  `json:"started_by" db:"started_by"`
	StartedAt  time.Time  `json:"started_at" db:"started_at"`
	EndedAt    *time.Time `json:"ended_at" db:"ended_at"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	MaxViewers int        `json:"max_viewers" db:"max_viewers"`
}

type Snapshot struct {
	ID          uuid.UUID `json:"snapshot_id" db:"snapshot_id"`
	SessionID   uuid.UUID `json:"session_id" db:"session_id"`
	Code        string    `json:"code" db:"code"`
	Language    string    `json:"language" db:"language"`
	Filename    string    `json:"filename" db:"filename"`
	TimestampMs int64     `json:"timestamp_ms" db:"timestamp_ms"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Handlers struct {
	db *sqlx.DB
}

func NewHandlers(db *sqlx.DB) *Handlers {
	return &Handlers{db: db}
}

func (h *Handlers) StartSession() echo.HandlerFunc {
	type Req struct {
		RoomID string `json:"room_id"`
		Title  string `json:"title"`
	}
	return func(c echo.Context) error {
		user, err := utils.GetUserFromCtx(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		req := &Req{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		roomID, err := uuid.Parse(req.RoomID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room_id"})
		}
		title := req.Title
		if title == "" {
			title = "Live Session"
		}
		sess := &CodingSession{}
		err = h.db.QueryRowxContext(c.Request().Context(),
			`INSERT INTO coding_sessions (room_id, title, started_by) VALUES ($1, $2, $3) RETURNING *`,
			roomID, title, user.UserID,
		).StructScan(sess)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusCreated, sess)
	}
}

func (h *Handlers) StopSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionID, err := uuid.Parse(c.Param("session_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session_id"})
		}
		_, err = h.db.ExecContext(c.Request().Context(),
			`UPDATE coding_sessions SET is_active = false, ended_at = CURRENT_TIMESTAMP WHERE session_id = $1`,
			sessionID,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusOK)
	}
}

func (h *Handlers) AddSnapshot() echo.HandlerFunc {
	type Req struct {
		SessionID   string `json:"session_id"`
		Code        string `json:"code"`
		Language    string `json:"language"`
		Filename    string `json:"filename"`
		TimestampMs int64  `json:"timestamp_ms"`
	}
	return func(c echo.Context) error {
		req := &Req{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		sessionID, err := uuid.Parse(req.SessionID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session_id"})
		}
		filename := req.Filename
		if filename == "" {
			filename = "main"
		}
		_, err = h.db.ExecContext(c.Request().Context(),
			`INSERT INTO session_snapshots (session_id, code, language, filename, timestamp_ms) VALUES ($1, $2, $3, $4, $5)`,
			sessionID, req.Code, req.Language, filename, req.TimestampMs,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusCreated)
	}
}

func (h *Handlers) GetRoomSessions() echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room_id"})
		}
		var sessions []CodingSession
		err = h.db.SelectContext(c.Request().Context(), &sessions,
			`SELECT * FROM coding_sessions WHERE room_id = $1 ORDER BY started_at DESC LIMIT 50`,
			roomID,
		)
		if err != nil && err != sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if sessions == nil {
			sessions = []CodingSession{}
		}
		return c.JSON(http.StatusOK, sessions)
	}
}

func (h *Handlers) GetSessionSnapshots() echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionID, err := uuid.Parse(c.Param("session_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session_id"})
		}
		var snapshots []Snapshot
		err = h.db.SelectContext(c.Request().Context(), &snapshots,
			`SELECT * FROM session_snapshots WHERE session_id = $1 ORDER BY timestamp_ms ASC`,
			sessionID,
		)
		if err != nil && err != sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if snapshots == nil {
			snapshots = []Snapshot{}
		}
		return c.JSON(http.StatusOK, snapshots)
	}
}

func (h *Handlers) UpdateViewerCount() echo.HandlerFunc {
	type Req struct {
		SessionID string `json:"session_id"`
		Count     int    `json:"count"`
	}
	return func(c echo.Context) error {
		req := &Req{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		sessionID, err := uuid.Parse(req.SessionID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session_id"})
		}
		_, err = h.db.ExecContext(c.Request().Context(),
			`UPDATE coding_sessions SET max_viewers = GREATEST(max_viewers, $1) WHERE session_id = $2`,
			req.Count, sessionID,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusOK)
	}
}
