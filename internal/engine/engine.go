package engine

import (
	"sort"

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
	err = p.scanForVersions()
	if err != nil {
		log.Error.Printf("Failed to scan project path '%s'\n", projectPath)
		return nil, err
	}

	// Sort versions based on name
	// TODO : Every version needs an IsLatest flag
	// Relying on sorting like this is probably not good enough
	sort.Slice(
		p.Versions, func(i, j int) bool {
			return p.Versions[i].Name > p.Versions[j].Name
		},
	)

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
