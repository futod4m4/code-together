package utils

import (
	"context"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
)

func ValidateIsOwner(ctx context.Context, ownerID string, logger logger.Logger) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.UserID.String() != ownerID {
		logger.Errorf(
			"ValidateIsOwner, userID: %v, ownerID: %v",
			user.UserID.String(),
			ownerID,
		)
		return httpErrors.Forbidden
	}

	return nil
}
