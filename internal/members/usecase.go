package members

import (
	"context"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	AddMember(ctx context.Context, member *models.RoomMember) (*models.RoomMember, error)
	UpdateRole(ctx context.Context, roomID, userID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error
	GetMembersByRoomID(ctx context.Context, roomID uuid.UUID) ([]*models.RoomMemberWithUser, error)
	GetMemberRole(ctx context.Context, roomID, userID uuid.UUID) (string, error)
}
