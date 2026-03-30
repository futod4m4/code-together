package http

import (
	"github.com/futod4m4/m/internal/chat"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strconv"
)

type chatHandlers struct {
	chatUC chat.UseCase
	logger logger.Logger
}

func NewChatHandlers(chatUC chat.UseCase, logger logger.Logger) chat.Handlers {
	return &chatHandlers{
		chatUC: chatUC,
		logger: logger,
	}
}

func (h *chatHandlers) CreateMessage() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "chatHandlers.CreateMessage")
		defer span.Finish()

		msg := &models.RoomMessage{}
		if err := c.Bind(msg); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
		}
		msg.UserID = user.UserID
		msg.Nickname = user.Nickname

		createdMsg, err := h.chatUC.CreateMessage(ctx, msg)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdMsg)
	}
}

func (h *chatHandlers) GetMessages() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "chatHandlers.GetMessages")
		defer span.Finish()

		roomID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 || limit > 100 {
			limit = 50
		}
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		if offset < 0 {
			offset = 0
		}

		messages, err := h.chatUC.GetMessagesByRoomID(ctx, roomID, limit, offset)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, messages)
	}
}
