package http

import (
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

		roomCodeUUID, err := uuid.Parse(c.Param("room_code_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		r := &models.RoomCode{}
		if err = c.Bind(r); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		r.ID = roomCodeUUID

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
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		roomCodeIDByRoomID, err := h.roomCodeUC.GetRoomCodeByRoomID(ctx, roomID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, roomCodeIDByRoomID)
	}
}

func (h *roomCodeHandlers) Compile() echo.HandlerFunc {
	type ProjectFile struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}

	type CompileRequest struct {
		Code     string        `json:"code"`
		Language string        `json:"language"`
		Mode     string        `json:"mode"`
		TestCode string        `json:"test_code"`
		Files    []ProjectFile `json:"files"`
	}

	type CompileResponse struct {
		Output string `json:"output"`
		Error  string `json:"error,omitempty"`
	}

	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomHandlers.Compile")
		defer span.Finish()

		var req CompileRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		}

		// Convert files to map for the compiler
		filesMap := make(map[string]string)
		for _, f := range req.Files {
			filesMap[f.Filename] = f.Content
		}

		var output string
		var compileErr error

		if req.Mode == "test" && req.TestCode != "" {
			output, compileErr = utils.ExecuteCodeWithTests(ctx, req.Language, req.Code, req.TestCode)
		} else {
			output, compileErr = utils.ExecuteProject(ctx, req.Language, req.Code, filesMap)
		}

		if compileErr != nil {
			return c.JSON(http.StatusOK, CompileResponse{Output: output, Error: compileErr.Error()})
		}

		return c.JSON(http.StatusOK, CompileResponse{Output: output})
	}
}

func (h *roomCodeHandlers) DownloadCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "roomCodeHandlers.DownloadCode")
		defer span.Finish()

		roomCodeUUID, err := uuid.Parse(c.Param("room_code_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		roomCode, err := h.roomCodeUC.GetRoomCodeByID(ctx, roomCodeUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// Determine file extension based on room language (query param)
		lang := c.QueryParam("language")
		ext := ".txt"
		switch lang {
		case "javascript":
			ext = ".js"
		case "python":
			ext = ".py"
		case "java":
			ext = ".java"
		case "go":
			ext = ".go"
		case "rust":
			ext = ".rs"
		case "php":
			ext = ".php"
		}

		filename := "code" + ext
		c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
		c.Response().Header().Set("Content-Type", "text/plain")
		return c.String(http.StatusOK, roomCode.Code)
	}
}
