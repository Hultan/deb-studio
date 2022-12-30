package engine

import (
	"fmt"
	"os"
	"path"

	"github.com/hultan/deb-studio/internal/logger"
)

var log *logger.Logger

type Engine struct {
	Project
}

func NewEngine(l *logger.Logger) *Engine {
	log = l
	return &Engine{}
}

func (e *Engine) IsProjectFolder(projectPath string) bool {
	log.Trace.Println("Entering IsProjectFolder...")
	defer log.Trace.Println("Exiting IsProjectFolder...")

	// Check if .project file exists...
	p := path.Join(projectPath, projectFileName)
	_, err := os.Stat(p)
	if err != nil {
		// File .project does not exist, or possibly a permission error...
		log.Info.Println(err)
		return false
	}
	return true
}

func (e *Engine) OpenProject(projectPath string) (*Project, error) {
	log.Trace.Println("Entering OpenProject...")
	defer log.Trace.Println("Exiting OpenProject...")

	// Make sure that directory exists
	if !doesDirectoryExist(projectPath) {
		log.Error.Printf("Path %s is missing!", projectPath)
		return nil, ErrorProjectFolderMissing
	}

	log.Info.Printf("Opening project : %s\n", projectPath)

	projectName, err := getProjectName(projectPath)
	if err != nil {
		log.Error.Printf("Failed to get project name from path '%s': %s\n", projectPath, err)
		return nil, err
	}

	w := &Project{
		Name: projectName,
		Path: projectPath,
	}

	err = w.scanProjectPath(nil)
	if err != nil {
		log.Error.Printf("Failed to scan project path '%s'\n", projectPath)
		return nil, err
	}

	log.Info.Printf("Successfully opened project path %s...\n", projectPath)
	return w, nil
}

func (e *Engine) SetupProject(projectPath, projectName string) (*Project, error) {
	log.Trace.Println("Entering SetupProject...")
	defer log.Trace.Println("Exiting SetupProject...")

	// Make sure that directory exists
	if !doesDirectoryExist(projectPath) {
		log.Error.Printf("Path %s is missing!", projectPath)
		return nil, ErrorProjectFolderMissing
	}

	log.Info.Printf("Setting up new project at : %s\n", projectPath)

	// Create project file
	filePath := path.Join(projectPath, projectFileName)
	content := fmt.Sprintf("PROJECT=%s", projectName)
	err := createTextFile(filePath, content)
	if err != nil {
		log.Error.Printf("Failed to write .project file to path '%s': %s\n", projectPath, err)
		return nil, err
	}

	log.Info.Printf("Successfully setup project path %s...\n", projectPath)
	return &Project{Name: projectName, Path: projectPath}, nil
}
