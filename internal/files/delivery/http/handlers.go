package http

import (
	"github.com/futod4m4/m/internal/files"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type fileHandlers struct {
	fileUC files.UseCase
	logger logger.Logger
}

func NewFileHandlers(fileUC files.UseCase, logger logger.Logger) files.Handlers {
	return &fileHandlers{
		fileUC: fileUC,
		logger: logger,
	}
}

func (h *fileHandlers) CreateFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "fileHandlers.CreateFile")
		defer span.Finish()

		f := &models.RoomFile{}
		if err := c.Bind(f); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		created, err := h.fileUC.CreateFile(ctx, f)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, created)
	}
}

func (h *fileHandlers) UpdateFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "fileHandlers.UpdateFile")
		defer span.Finish()

		fileID, err := uuid.Parse(c.Param("file_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		f := &models.RoomFile{}
		if err = c.Bind(f); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		f.ID = fileID

		updated, err := h.fileUC.UpdateFile(ctx, f)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updated)
	}
}

func (h *fileHandlers) DeleteFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "fileHandlers.DeleteFile")
		defer span.Finish()

		fileID, err := uuid.Parse(c.Param("file_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.fileUC.DeleteFile(ctx, fileID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *fileHandlers) GetFilesByRoomID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "fileHandlers.GetFilesByRoomID")
		defer span.Finish()

		roomID, err := uuid.Parse(c.Param("room_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		filesList, err := h.fileUC.GetFilesByRoomID(ctx, roomID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, filesList)
	}
}
