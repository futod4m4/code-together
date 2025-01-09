package http

import (
	"fmt"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/roomCodes"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type roomCodeHandlers struct {
	cfg        *config.Config
	roomCodeUC roomCodes.RoomCodeUseCase
	logger     logger.Logger
}

func NewRoomCodeHandlers(cfg *config.Config, roomCodeUC roomCodes.RoomCodeUseCase, logger logger.Logger) roomCodes.HttpHandlers {
	return &roomCodeHandlers{
		cfg:        cfg,
		roomCodeUC: roomCodeUC,
		logger:     logger,
	}
}

func (h *roomCodeHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomCodeHandlers.Create")
		defer span.Finish()

		r := &models.RoomCode{}
		if err := c.Bind(r); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdRoomCode, err := h.roomCodeUC.CreateRoomCode(ctx, r)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdRoomCode)
	}
}

func (h *roomCodeHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomCodeHandlers.Update")
		defer span.Finish()

		fmt.Println("0")
		roomCodeUUID, err := uuid.Parse(c.Param("room_code_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		fmt.Println("1")
		r := &models.RoomCode{}
		if err = c.Bind(r); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		r.ID = roomCodeUUID

		fmt.Println("2")

		updatedRoomCode, err := h.roomCodeUC.UpdateRoomCode(ctx, r)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedRoomCode)
	}
}

func (h *roomCodeHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Create")
		defer span.Finish()

		roomCodeUUID, err := uuid.Parse(c.Param("room_code_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.roomCodeUC.DeleteRoomCode(ctx, roomCodeUUID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *roomCodeHandlers) GetRoomCodeByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomCodeHandlers.GetRoomCodeByID")
		defer span.Finish()

		roomCodeUUID, err := uuid.Parse(c.Param("room_code_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		roomCodeByID, err := h.roomCodeUC.GetRoomCodeByID(ctx, roomCodeUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, roomCodeByID)
	}
}

func (h *roomCodeHandlers) GetRoomCodeByRoomID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Create")
		defer span.Finish()

		roomID, err := uuid.Parse(c.Param("room_id"))

		roomCodeIDByRoomID, err := h.roomCodeUC.GetRoomCodeByRoomID(ctx, roomID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, roomCodeIDByRoomID)
	}
}
