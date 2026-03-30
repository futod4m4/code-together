package sessions

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

const (
	maxFiles    = 50
	maxFileSize = 100 * 1024 // 100KB
)

var textExtensions = map[string]bool{
	".go": true, ".js": true, ".ts": true, ".jsx": true, ".tsx": true,
	".py": true, ".java": true, ".rs": true, ".php": true, ".rb": true,
	".c": true, ".cpp": true, ".h": true, ".hpp": true, ".cs": true,
	".html": true, ".css": true, ".scss": true, ".less": true,
	".json": true, ".xml": true, ".yaml": true, ".yml": true, ".toml": true,
	".md": true, ".txt": true, ".sh": true, ".bash": true,
	".sql": true, ".graphql": true, ".proto": true,
	".env": true, ".gitignore": true, ".dockerignore": true,
	".mk": true, ".mod": true, ".sum": true, ".cfg": true, ".ini": true,
	".vue": true, ".svelte": true,
}

type GitHubHandlers struct {
	db *sqlx.DB
}

func NewGitHubHandlers(db *sqlx.DB) *GitHubHandlers {
	return &GitHubHandlers{db: db}
}

func (h *GitHubHandlers) ImportRepo() echo.HandlerFunc {
	type Req struct {
		RoomID string `json:"room_id"`
		URL    string `json:"url"`
	}
	return func(c echo.Context) error {
		user, err := utils.GetUserFromCtx(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		_ = user

		req := &Req{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		roomID, err := uuid.Parse(req.RoomID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room_id"})
		}

		repoURL := req.URL
		if !strings.HasPrefix(repoURL, "https://") {
			repoURL = "https://github.com/" + repoURL
		}
		if !strings.HasSuffix(repoURL, ".git") {
			repoURL += ".git"
		}

		// Clone to temp dir (shallow, single branch)
		tmpDir, err := os.MkdirTemp("", "git-import-")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "temp dir failed"})
		}
		defer os.RemoveAll(tmpDir)

		cmd := exec.CommandContext(c.Request().Context(),
			"git", "clone", "--depth=1", "--single-branch", repoURL, tmpDir+"/repo",
		)
		var stderr strings.Builder
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("git clone failed: %s", stderr.String()),
			})
		}

		repoDir := tmpDir + "/repo"
		var importedFiles []map[string]string
		fileCount := 0

		err = filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == "vendor" || d.Name() == "__pycache__" || d.Name() == ".idea" || d.Name() == ".vscode" {
					return filepath.SkipDir
				}
				return nil
			}
			if fileCount >= maxFiles {
				return filepath.SkipAll
			}

			relPath, _ := filepath.Rel(repoDir, path)
			ext := strings.ToLower(filepath.Ext(relPath))

			// Skip non-text files
			base := filepath.Base(relPath)
			if !textExtensions[ext] && ext != "" && base != "Makefile" && base != "Dockerfile" {
				return nil
			}

			info, err := d.Info()
			if err != nil || info.Size() > maxFileSize {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			// Detect language from extension
			langMap := map[string]string{
				".go": "go", ".js": "javascript", ".ts": "typescript", ".py": "python",
				".java": "java", ".rs": "rust", ".php": "php", ".html": "html",
				".css": "css", ".json": "json", ".md": "markdown",
			}
			lang := langMap[ext]
			if lang == "" {
				lang = "plaintext"
			}

			_, execErr := h.db.ExecContext(c.Request().Context(),
				`INSERT INTO room_files (room_id, filename, language, content, is_entry_point)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (room_id, filename) DO UPDATE SET content = $4`,
				roomID, relPath, lang, string(content), fileCount == 0,
			)
			if execErr != nil {
				return nil
			}

			importedFiles = append(importedFiles, map[string]string{
				"filename": relPath,
				"language": lang,
			})
			fileCount++
			return nil
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"imported": len(importedFiles),
			"files":    importedFiles,
		})
	}
}
