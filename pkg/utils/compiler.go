package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const maxOutputSize = 64 * 1024 // 64KB
const execTimeout = 15 * time.Second

// ExecuteCode runs a single file (legacy, still used by tests)
func ExecuteCode(ctx context.Context, language, code string) (string, error) {
	return ExecuteProject(ctx, language, code, nil)
}

// ExecuteProject runs a project with multiple files
func ExecuteProject(ctx context.Context, language, mainCode string, files map[string]string) (string, error) {
	tempDir, err := os.MkdirTemp("", "code-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write all project files
	if len(files) > 0 {
		for filename, content := range files {
			fullPath := filepath.Join(tempDir, filename)
			dir := filepath.Dir(fullPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", fmt.Errorf("failed to create dir %s: %w", dir, err)
			}
			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				return "", fmt.Errorf("failed to write %s: %w", filename, err)
			}
		}
	} else if mainCode != "" {
		// Fallback: write single main file
		var fileName string
		switch language {
		case "javascript":
			fileName = "main.js"
		case "python":
			fileName = "main.py"
		case "java":
			fileName = "Main.java"
		case "go":
			fileName = "main.go"
		case "rust":
			fileName = "main.rs"
		case "php":
			fileName = "main.php"
		default:
			return "", fmt.Errorf("unsupported language: %s", language)
		}
		if err := os.WriteFile(filepath.Join(tempDir, fileName), []byte(mainCode), 0644); err != nil {
			return "", fmt.Errorf("failed to write source file: %w", err)
		}
	}

	var runCmd []string
	var compileCmd []string

	switch language {
	case "javascript":
		// Find entry file
		entry := findEntry(files, "index.js", "main.js")
		runCmd = []string{"node", filepath.Join(tempDir, entry)}
	case "python":
		entry := findEntry(files, "main.py", "app.py")
		runCmd = []string{"python3", filepath.Join(tempDir, entry)}
	case "java":
		// Compile all java files then run Main
		compileCmd = []string{"sh", "-c", fmt.Sprintf("cd %s && javac *.java", tempDir)}
		runCmd = []string{"java", "-cp", tempDir, "Main"}
	case "go":
		// Always create go.mod if not already present in the temp dir
		goModPath := filepath.Join(tempDir, "go.mod")
		if _, statErr := os.Stat(goModPath); os.IsNotExist(statErr) {
			gomod := "module project\n\ngo 1.22\n"
			_ = os.WriteFile(goModPath, []byte(gomod), 0644)
		}
		// Build the binary first (resolves all subpackages), then run it
		runCmd = []string{"sh", "-c", fmt.Sprintf("cd %s && go build -o %s/app . && %s/app", tempDir, tempDir, tempDir)}
	case "rust":
		entry := findEntry(files, "main.rs", "lib.rs")
		compileCmd = []string{"rustc", "-o", filepath.Join(tempDir, "main"), filepath.Join(tempDir, entry)}
		runCmd = []string{filepath.Join(tempDir, "main")}
	case "php":
		entry := findEntry(files, "index.php", "main.php")
		runCmd = []string{"php", filepath.Join(tempDir, entry)}
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	execCtx, cancel := context.WithTimeout(ctx, execTimeout)
	defer cancel()

	// Compile if needed
	if len(compileCmd) > 0 {
		var stderr bytes.Buffer
		cmd := exec.CommandContext(execCtx, compileCmd[0], compileCmd[1:]...)
		cmd.Stderr = &stderr
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			return truncate(stderr.String()), fmt.Errorf("compilation error: %s", stderr.String())
		}
	}

	// Run
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(execCtx, runCmd[0], runCmd[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = tempDir
	cmd.Env = append(os.Environ(), "GOPATH="+filepath.Join(tempDir, ".gopath"))

	if err := cmd.Run(); err != nil {
		errOut := stderr.String()
		if errOut == "" {
			errOut = err.Error()
		}
		return truncate(errOut), fmt.Errorf("execution error: %s", errOut)
	}

	result := stdout.String()
	if se := stderr.String(); se != "" {
		result += "\n" + se
	}
	return truncate(result), nil
}

// ExecuteCodeWithTests runs code + test file together
func ExecuteCodeWithTests(ctx context.Context, language, code, testCode string) (string, error) {
	tempDir, err := os.MkdirTemp("", "code-test-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	var mainFile, testFile string
	var runCmd []string

	switch language {
	case "javascript":
		mainFile = "main.js"
		testFile = "test.js"
		runCmd = []string{"node", filepath.Join(tempDir, testFile)}
	case "python":
		mainFile = "main.py"
		testFile = "test_main.py"
		runCmd = []string{"python3", "-m", "pytest", filepath.Join(tempDir, testFile), "-v", "--tb=short", "--no-header"}
	case "go":
		mainFile = "main.go"
		testFile = "main_test.go"
		_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module project\n\ngo 1.22\n"), 0644)
		runCmd = []string{"sh", "-c", fmt.Sprintf("cd %s && go test -v ./...", tempDir)}
	case "java":
		mainFile = "Main.java"
		testFile = "MainTest.java"
		runCmd = []string{"sh", "-c", fmt.Sprintf("cd %s && javac *.java && java MainTest", tempDir)}
	case "rust":
		mainFile = "main.rs"
		testFile = "test.rs"
		runCmd = []string{"sh", "-c", fmt.Sprintf("cd %s && rustc --test test.rs -o %s/test_bin && %s/test_bin", tempDir, tempDir, tempDir)}
	case "php":
		mainFile = "main.php"
		testFile = "test.php"
		runCmd = []string{"php", filepath.Join(tempDir, testFile)}
	default:
		return "", fmt.Errorf("unsupported language for tests: %s", language)
	}

	if err := os.WriteFile(filepath.Join(tempDir, mainFile), []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write source file: %w", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, testFile), []byte(testCode), 0644); err != nil {
		return "", fmt.Errorf("failed to write test file: %w", err)
	}

	execCtx, cancel := context.WithTimeout(ctx, execTimeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(execCtx, runCmd[0], runCmd[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = tempDir

	if err := cmd.Run(); err != nil {
		combined := stdout.String() + stderr.String()
		if combined != "" {
			return truncate(combined), nil
		}
		return "", fmt.Errorf("test execution error: %w", err)
	}

	return truncate(stdout.String() + stderr.String()), nil
}

func findEntry(files map[string]string, candidates ...string) string {
	for _, c := range candidates {
		if _, ok := files[c]; ok {
			return c
		}
	}
	// Return first candidate as default
	if len(candidates) > 0 {
		return candidates[0]
	}
	return "main"
}

func truncate(s string) string {
	if len(s) > maxOutputSize {
		return s[:maxOutputSize] + "\n... (output truncated)"
	}
	return s
}
