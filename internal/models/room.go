package models

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Room struct {
	ID        uuid.UUID `json:"room_id" db:"room_id" redis:"room_id" validate:"omitempty"`
	Name      string    `json:"room_name" db:"room_name" redis:"room_name" validate:"omitempty"`
	JoinCode  string    `json:"join_code" db:"join_code" redis:"join_code" validate:"omitempty"`
	OwnerID   uuid.UUID `json:"owner_id" db:"owner_id" redis:"owner_id" validate:"omitempty"`
	Language  string    `json:"language" db:"language" redis:"language" validate:"omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at" redis:"created_at" validate:"omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" redis:"updated_at" validate:"omitempty"`
	Code      *RoomCode
}

type RoomCode struct {
	ID        int       `json:"room_code_id" db:"room_code_id" redis:"room_code_id" validate:"omitempty"`
	RoomID    uuid.UUID `json:"room_id" db:"room_id" redis:"room_id" validate:"omitempty"`
	Code      string    `json:"code" db:"code" redis:"code" validate:"omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at" redis:"created_at" validate:"omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" redis:"updated_at" validate:"omitempty"`
}

func (r *Room) GenJoinCode() error {
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		r.JoinCode = ""
		return err
	}

	code := base64.URLEncoding.EncodeToString(bytes)
	code = strings.ToUpper(strings.ReplaceAll(code, "-", ""))
	code = strings.ReplaceAll(code, "_", "")
	r.JoinCode = code[:12]
	return nil
}
