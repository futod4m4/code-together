package http

import (
	"github.com/futod4m4/m/internal/members"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type memberHandlers struct {
	memberUC members.UseCase
	logger   logger.Logger
}

func NewMemberHandlers(memberUC members.UseCase, logger logger.Logger) members.Handlers {
	return &memberHandlers{
		memberUC: memberUC,
		logger:   logger,
	}
}

func (h *memberHandlers) AddMember() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "memberHandlers.AddMember")
		defer span.Finish()

		m := &models.RoomMember{}
		if err := c.Bind(m); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		created, err := h.memberUC.AddMember(ctx, m)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, created)
	}
}

func (h *memberHandlers) UpdateRole() echo.HandlerFunc {
	type UpdateRoleRequest struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "memberHandlers.UpdateRole")
		defer span.Finish()

		roomID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		var req UpdateRoleRequest
		if err = c.Bind(&req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.memberUC.UpdateRole(ctx, roomID, userID, req.Role); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *memberHandlers) RemoveMember() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "memberHandlers.RemoveMember")
		defer span.Finish()

		roomID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.memberUC.RemoveMember(ctx, roomID, userID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *memberHandlers) GetMembers() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "memberHandlers.GetMembers")
		defer span.Finish()

		roomID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		membersList, err := h.memberUC.GetMembersByRoomID(ctx, roomID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, membersList)
	}
}
