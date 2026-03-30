package repository

import (
	"context"
	"github.com/futod4m4/m/internal/members"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type memberRepo struct {
	db *sqlx.DB
}

func NewMemberRepository(db *sqlx.DB) members.Repository {
	return &memberRepo{db: db}
}

func (r *memberRepo) AddMember(ctx context.Context, member *models.RoomMember) (*models.RoomMember, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberRepo.AddMember")
	defer span.Finish()

	var m models.RoomMember
	if err := r.db.QueryRowxContext(ctx, addMember,
		&member.RoomID, &member.UserID, &member.Role,
	).StructScan(&m); err != nil {
		return nil, errors.Wrap(err, "memberRepo.AddMember.QueryRowxContext")
	}
	return &m, nil
}

func (r *memberRepo) UpdateRole(ctx context.Context, roomID, userID uuid.UUID, role string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberRepo.UpdateRole")
	defer span.Finish()

	_, err := r.db.ExecContext(ctx, updateRole, roomID, userID, role)
	if err != nil {
		return errors.Wrap(err, "memberRepo.UpdateRole.ExecContext")
	}
	return nil
}

func (r *memberRepo) RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberRepo.RemoveMember")
	defer span.Finish()

	_, err := r.db.ExecContext(ctx, removeMember, roomID, userID)
	if err != nil {
		return errors.Wrap(err, "memberRepo.RemoveMember.ExecContext")
	}
	return nil
}

func (r *memberRepo) GetMembersByRoomID(ctx context.Context, roomID uuid.UUID) ([]*models.RoomMemberWithUser, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberRepo.GetMembersByRoomID")
	defer span.Finish()

	var membersList []*models.RoomMemberWithUser
	if err := r.db.SelectContext(ctx, &membersList, getMembersByRoomID, roomID); err != nil {
		return nil, errors.Wrap(err, "memberRepo.GetMembersByRoomID.SelectContext")
	}
	return membersList, nil
}

func (r *memberRepo) GetMemberRole(ctx context.Context, roomID, userID uuid.UUID) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberRepo.GetMemberRole")
	defer span.Finish()

	var role string
	if err := r.db.GetContext(ctx, &role, getMemberRole, roomID, userID); err != nil {
		return "", errors.Wrap(err, "memberRepo.GetMemberRole.GetContext")
	}
	return role, nil
}

func (r *memberRepo) IsMember(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "memberRepo.IsMember")
	defer span.Finish()

	var exists bool
	if err := r.db.GetContext(ctx, &exists, isMember, roomID, userID); err != nil {
		return false, errors.Wrap(err, "memberRepo.IsMember.GetContext")
	}
	return exists, nil
}
