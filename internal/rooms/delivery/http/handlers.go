package http

import (
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/members"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type roomHandlers struct {
	cfg      *config.Config
	roomUC   rooms.RoomUseCase
	memberUC members.UseCase
	logger   logger.Logger
}

func NewRoomHandlers(cfg *config.Config, roomUC rooms.RoomUseCase, memberUC members.UseCase, logger logger.Logger) rooms.HttpHandlers {
	return &roomHandlers{
		cfg:      cfg,
		roomUC:   roomUC,
		memberUC: memberUC,
		logger:   logger,
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

		// Auto-add owner as member
		if _, err := h.memberUC.AddMember(ctx, &models.RoomMember{
			RoomID: createdRoom.ID,
			UserID: createdRoom.OwnerID,
			Role:   "owner",
		}); err != nil {
			h.logger.Errorf("roomHandlers.Create.AddOwnerMember: %v", err)
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

func (h *roomHandlers) GetRoomByJoinCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Create")
		defer span.Finish()

		joinCode := c.Param("join_code")

		roomByJoinCode, err := h.roomUC.GetRoomByJoinCode(ctx, joinCode)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, roomByJoinCode)
	}
}

func (h *roomHandlers) GetMyRooms() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.GetMyRooms")
		defer span.Finish()

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
		}

		rooms, err := h.roomUC.GetRoomsByOwnerID(ctx, user.UserID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, rooms)
	}
}
