package engine

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"strings"
)

func readAllText(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to readAllText : %w", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to readAllText : %w", err)
	}
	return string(b), nil
}

// getFirstLine : gets the first line of the text, ignoring the prefix
func getFirstLine(text string, prefix string, sep string) (string, error) {
	if !strings.HasPrefix(text, prefix) {
		return "", errors.New("does not contain prefix")
	}
	truncated := text[len(prefix):]
	i := strings.IndexAny(truncated, sep)
	if i < 0 {
		return truncated, nil
	}
	return strings.Trim(truncated[:i], " \t"), nil
}

// doesDirectoryExist : Check if a directory exists
func doesDirectoryExist(workspacePath string) bool {
	folderInfo, err := os.Stat(workspacePath)
	if os.IsNotExist(err) {
		return false
	}
	return folderInfo.IsDir()
}

// getUserHomeDirectory : Get current users home directory
func getUserHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}

// createTextFile : creates a text file containing the string in the argument content
func createTextFile(filePath, fileName, content string) error {
	// Create version file
	p := path.Join(filePath, fileName)
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write program file contents
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
