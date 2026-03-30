package sessions

import (
	"net/http"

	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type BanHandlers struct {
	db *sqlx.DB
}

func NewBanHandlers(db *sqlx.DB) *BanHandlers {
	return &BanHandlers{db: db}
}

// BanUser bans a registered user from a room
func (h *BanHandlers) BanUser() echo.HandlerFunc {
	type Req struct {
		RoomID string `json:"room_id"`
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}
	return func(c echo.Context) error {
		owner, err := utils.GetUserFromCtx(c.Request().Context())
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
		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
		}

		// Verify caller is owner
		var ownerID uuid.UUID
		err = h.db.GetContext(c.Request().Context(), &ownerID,
			`SELECT owner_id FROM rooms WHERE room_id = $1`, roomID)
		if err != nil || ownerID != owner.UserID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "only owner can ban"})
		}

		// Insert ban
		_, err = h.db.ExecContext(c.Request().Context(),
			`INSERT INTO banned_members (room_id, user_id, banned_by, reason)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (room_id, user_id) DO NOTHING`,
			roomID, userID, owner.UserID, req.Reason,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// Also remove from members
		_, _ = h.db.ExecContext(c.Request().Context(),
			`DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`, roomID, userID)

		return c.JSON(http.StatusOK, map[string]string{"status": "banned"})
	}
}

// BanIP bans a guest by IP address
func (h *BanHandlers) BanIP() echo.HandlerFunc {
	type Req struct {
		RoomID string `json:"room_id"`
		IP     string `json:"ip"`
		Reason string `json:"reason"`
	}
	return func(c echo.Context) error {
		owner, err := utils.GetUserFromCtx(c.Request().Context())
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

		// Verify owner
		var ownerID uuid.UUID
		err = h.db.GetContext(c.Request().Context(), &ownerID,
			`SELECT owner_id FROM rooms WHERE room_id = $1`, roomID)
		if err != nil || ownerID != owner.UserID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "only owner can ban"})
		}

		// Use a nil user_id with the IP
		_, err = h.db.ExecContext(c.Request().Context(),
			`INSERT INTO banned_members (room_id, banned_by, ip_address, reason, user_id)
			 VALUES ($1, $2, $3, $4, $2)
			 ON CONFLICT DO NOTHING`,
			roomID, owner.UserID, req.IP, req.Reason,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "ip_banned"})
	}
}

// CheckBan checks if current user is banned from a room
func (h *BanHandlers) CheckBan() echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room_id"})
		}

		// Check by user_id if authenticated
		user, userErr := utils.GetUserFromCtx(c.Request().Context())
		if userErr == nil {
			var count int
			_ = h.db.GetContext(c.Request().Context(), &count,
				`SELECT COUNT(*) FROM banned_members WHERE room_id = $1 AND user_id = $2`,
				roomID, user.UserID)
			if count > 0 {
				return c.JSON(http.StatusOK, map[string]bool{"banned": true})
			}
		}

		// Check by IP
		ip := c.RealIP()
		var count int
		_ = h.db.GetContext(c.Request().Context(), &count,
			`SELECT COUNT(*) FROM banned_members WHERE room_id = $1 AND ip_address = $2`,
			roomID, ip)

		return c.JSON(http.StatusOK, map[string]bool{"banned": count > 0})
	}
}

// GetBannedList returns banned users for a room
func (h *BanHandlers) GetBannedList() echo.HandlerFunc {
	type BannedEntry struct {
		BanID     uuid.UUID `json:"ban_id" db:"ban_id"`
		RoomID    uuid.UUID `json:"room_id" db:"room_id"`
		UserID    uuid.UUID `json:"user_id" db:"user_id"`
		IPAddress *string   `json:"ip_address" db:"ip_address"`
		Reason    string    `json:"reason" db:"reason"`
	}
	return func(c echo.Context) error {
		roomID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room_id"})
		}
		var bans []BannedEntry
		err = h.db.SelectContext(c.Request().Context(), &bans,
			`SELECT ban_id, room_id, user_id, ip_address, reason FROM banned_members WHERE room_id = $1`, roomID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if bans == nil {
			bans = []BannedEntry{}
		}
		return c.JSON(http.StatusOK, bans)
	}
}

// Unban removes a ban
func (h *BanHandlers) Unban() echo.HandlerFunc {
	return func(c echo.Context) error {
		banID, err := uuid.Parse(c.Param("ban_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ban_id"})
		}
		_, err = h.db.ExecContext(c.Request().Context(),
			`DELETE FROM banned_members WHERE ban_id = $1`, banID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusOK)
	}
}
