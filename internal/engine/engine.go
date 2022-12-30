package engine

import (
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
	if !haveDescriptor(projectDescriptor, projectPath) {
		return false
	}
	return true
}

func (e *Engine) OpenProject(projectPath string) (*Project, error) {
	log.Trace.Println("Entering OpenProject...")
	defer log.Trace.Println("Exiting OpenProject...")

	// Read .project file
	projectName, err := readDescriptor(projectDescriptor, projectPath)
	if err != nil {
		log.Error.Printf("Failed to get project name from path '%s': %s\n", projectPath, err)
		return nil, err
	}

	// Create project
	p := &Project{
		Name: projectName,
		Path: projectPath,
	}

	// Scan for versions
	err = p.scanForVersions(nil)
	if err != nil {
		log.Error.Printf("Failed to scan project path '%s'\n", projectPath)
		return nil, err
	}

	log.Info.Printf("Successfully opened project path %s...\n", projectPath)
	return p, nil
}

func (e *Engine) SetupProject(projectPath, projectName string) (*Project, error) {
	log.Trace.Println("Entering SetupProject...")
	defer log.Trace.Println("Exiting SetupProject...")

	log.Info.Printf("Setting up new project at : %s\n", projectPath)

	// Create .project file
	err := writeDescriptor(projectDescriptor, projectPath, projectName)
	if err != nil {
		log.Error.Printf("Failed to write .project file to path '%s': %s\n", projectPath, err)
		return nil, err
	}

	log.Info.Printf("Successfully setup project path %s...\n", projectPath)
	return &Project{Name: projectName, Path: projectPath}, nil
}
