package usecase

import (
	"context"
	"github.com/futod4m4/m/internal/members"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type memberUC struct {
	memberRepo members.Repository
	logger     logger.Logger
}

func NewMemberUseCase(memberRepo members.Repository, logger logger.Logger) members.UseCase {
	return &memberUC{
		memberRepo: memberRepo,
		logger:     logger,
	}
}

func (u *memberUC) AddMember(ctx context.Context, member *models.RoomMember) (*models.RoomMember, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberUC.AddMember")
	defer span.Finish()

	validRoles := map[string]bool{"owner": true, "editor": true, "viewer": true}
	if !validRoles[member.Role] {
		return nil, httpErrors.NewRestError(http.StatusBadRequest, "Invalid role. Use: owner, editor, viewer", nil)
	}

	return u.memberRepo.AddMember(ctx, member)
}

func (u *memberUC) UpdateRole(ctx context.Context, roomID, userID uuid.UUID, role string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberUC.UpdateRole")
	defer span.Finish()

	validRoles := map[string]bool{"owner": true, "editor": true, "viewer": true}
	if !validRoles[role] {
		return httpErrors.NewRestError(http.StatusBadRequest, "Invalid role", nil)
	}

	return u.memberRepo.UpdateRole(ctx, roomID, userID, role)
}

func (u *memberUC) RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberUC.RemoveMember")
	defer span.Finish()

	return u.memberRepo.RemoveMember(ctx, roomID, userID)
}

func (u *memberUC) GetMembersByRoomID(ctx context.Context, roomID uuid.UUID) ([]*models.RoomMemberWithUser, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberUC.GetMembersByRoomID")
	defer span.Finish()

	return u.memberRepo.GetMembersByRoomID(ctx, roomID)
}

func (u *memberUC) GetMemberRole(ctx context.Context, roomID, userID uuid.UUID) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberUC.GetMemberRole")
	defer span.Finish()

	return u.memberRepo.GetMemberRole(ctx, roomID, userID)
}
