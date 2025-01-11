package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

func ExecuteCode(ctx context.Context, language, code string) (string, error) {
	var imageName string
	var compileCommand []string
	var runCommand []string
	tempDir, err := os.MkdirTemp("", "code-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	var fileName string
	switch language {
	case "javascript":
		imageName = "node:latest"
		fileName = "main.js"
		runCommand = []string{"node", fmt.Sprintf("/code/%s", fileName)}
	case "python":
		imageName = "python:latest"
		fileName = "main.py"
		runCommand = []string{"python3", fmt.Sprintf("/code/%s", fileName)}
	case "java":
		imageName = "openjdk:latest"
		fileName = "Main.java"
		compileCommand = []string{"javac", fmt.Sprintf("/code/%s", fileName)}
		runCommand = []string{"java", "-cp", "/code", "Main"}
	case "go":
		imageName = "golang:latest"
		fileName = "main.go"
		runCommand = []string{"go", "run", fmt.Sprintf("/code/%s", fileName)}
	case "rust":
		imageName = "rust:latest"
		fileName = "main.rs"
		compileCommand = []string{"rustc", "-o", "/code/main", "/code/main.rs"}
		runCommand = []string{"/code/main"}
	case "php":
		imageName = "php:latest"
		fileName = "main.php"
		runCommand = []string{"php", fmt.Sprintf("/code/%s", fileName)}
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	sourceFilePath := fmt.Sprintf("%s/%s", tempDir, fileName)
	if err := os.WriteFile(sourceFilePath, []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write source file: %w", err)
	}

	dockerRunBase := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/code", tempDir), imageName}

	var output bytes.Buffer
	var stderr bytes.Buffer

	if len(compileCommand) > 0 {
		cmd := exec.CommandContext(ctx, "docker", append(dockerRunBase, compileCommand...)...)
		cmd.Stdout = &output
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return stderr.String(), fmt.Errorf("compilation error: %w", err)
		}
	}

	cmd := exec.CommandContext(ctx, "docker", append(dockerRunBase, runCommand...)...)
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stderr.String(), fmt.Errorf("execution error: %w", err)
	}

	return output.String(), nil
}
