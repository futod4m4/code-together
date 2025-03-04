package models

import "github.com/google/uuid"

type Session struct {
	SessionID string    `json:"session_id" redis:"session_id"`
	UserID    uuid.UUID `json:"user_id" redis:"user_id"`
	RoomID    uuid.UUID `json:"room_id" redis:"room_id"`
}
