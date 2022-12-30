package engine

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const projectFileName = ".project"
const versionFileName = ".version"
const architectureFileName = ".architecture"

var ErrorProjectFolderMissing = errors.New("missing project folder")

func readAllText(path string) (string, error) {
	log.Trace.Println("Entering readAllText...")
	defer log.Trace.Println("Exiting readAllText...")

	f, err := os.Open(path)
	if err != nil {
		log.Error.Printf("Failed to open path '%s': %s\n", path, err)
		return "", err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Error.Printf("Failed to read all text of path '%s': %s\n", path, err)
		return "", err
	}
	return string(b), nil
}

// getFirstLine : gets the first line of the text, ignoring the prefix
func getFirstLine(text string, prefix string, sep string) (string, error) {
	log.Trace.Println("Entering getFirstLine...")
	defer log.Trace.Println("Exiting getFirstLine...")

	if !strings.HasPrefix(text, prefix) {
		log.Error.Printf("Text does not contain the prefix: %s\n", prefix)
		return "", errors.New("does not contain prefix '" + prefix + "'")
	}
	truncated := text[len(prefix):]
	i := strings.IndexAny(truncated, sep)
	if i < 0 {
		return truncated, nil
	}
	return strings.Trim(truncated[:i], " \t"), nil
}

// doesDirectoryExist : Check if a directory exists
func doesDirectoryExist(path string) bool {
	log.Trace.Println("Entering doesDirectoryExist...")
	defer log.Trace.Println("Exiting doesDirectoryExist...")

	folderInfo, err := os.Stat(path)
	if err != nil {
		// Directory does not exist, or user does not
		// have permissions, or ...
		return false
	}
	return folderInfo.IsDir()
}

// createTextFile : creates a text file containing the string in the argument content
func createTextFile(filePath, content string) error {
	log.Trace.Println("Entering createTextFile...")
	defer log.Trace.Println("Exiting createTextFile...")

	// Create file
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Error.Printf("Failed to get absolute path of path '%s': %s\n", filePath, err)
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		log.Error.Printf("Failed to create file '%s': %s\n", filePath, err)
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Error.Printf("Failed to create file '%s': %s\n", filePath, err)
		}
	}(f)

	// Write file contents
	_, err = f.WriteString(content)
	if err != nil {
		log.Error.Printf("Failed to write content to file '%s': %s\n", filePath, err)
		log.Error.Printf("Content '%s'\n", content)
		return err
	}

	return nil
}
