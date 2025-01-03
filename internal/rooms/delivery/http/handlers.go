package http

import (
	"github.com/futod4m4/m/config"
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
	cfg    *config.Config
	roomUC rooms.RoomUseCase
	logger logger.Logger
}

func NewRoomHandlers(cfg *config.Config, roomUC rooms.RoomUseCase, logger logger.Logger) rooms.Handlers {
	return &roomHandlers{
		cfg:    cfg,
		roomUC: roomUC,
		logger: logger,
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

func (h *roomHandlers) Join() echo.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (h *roomHandlers) Leave() echo.HandlerFunc {
	//TODO implement me
	panic("implement me")
}
