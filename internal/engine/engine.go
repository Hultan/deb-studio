package engine

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hultan/deb-studio/internal/logger"
)

const projectFileName = ".project"

type engine struct {
	log *logger.Logger
}

type Project struct {
	Path     string
	Name     string
	Versions []*Version
}

type Version struct {
	Path          string
	Name          string
	Architectures []*Architecture
}

type Architecture struct {
	Path string
	Name string
}

func NewEngine(log *logger.Logger) *engine {
	return &engine{log: log}
}

func newProject(projectPath, projectName string) *Project {
	return &Project{
		Name: projectName,
		Path: projectPath,
	}
}

func (e *engine) IsProjectFolder(projectPath string) bool {
	// Check if .project file exists...
	p := path.Join(projectPath, projectFileName)
	_, err := os.Stat(p)
	if err != nil {
		// File .project does not exist, or possibly a permission error...
		return false
	}
	return true
}

func (e *engine) OpenProject(projectPath string) (*Project, error) {
	// Make sure that directory exists
	if !doesDirectoryExist(projectPath) {
		return nil, ErrorProjectFolderMissing
	}

	projectName, err := getProjectName(projectPath)
	if err != nil {
		return nil, ErrorNewProjectFolder
	}

	w := newProject(projectPath, projectName)

	// err = w.scanVersions()
	err = w.scanFolder(nil)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (e *engine) SetupProject(projectPath, projectName string) (*Project, error) {
	// Make sure that directory exists
	if !doesDirectoryExist(projectPath) {
		return nil, ErrorProjectFolderMissing
	}

	// Create project file
	filePath := path.Join(projectPath, projectFileName)
	content := fmt.Sprintf("PROJECT=%s", projectName)
	err := createTextFile(filePath, content)
	if err != nil {
		return nil, err
	}

	return newProject(projectPath, projectName), nil
}

func getProjectName(projectPath string) (string, error) {
	p := path.Join(projectPath, projectFileName)
	text, err := readAllText(p)
	if err != nil {
		return "", fmt.Errorf("failed to find project name : %w", err)
	}
	name, err := getFirstLine(text, "PROJECT=", "\t\n")
	if err != nil {
		return "", fmt.Errorf("failed to find project name : %w", err)
	}

	return name, nil
}

func (w *Project) scanFolder(version *Version) error {
	typeName := "version"
	scanPath := w.Path
	if version != nil {
		typeName = "architecture"
		scanPath = version.Path
	}

	// Open path to scan
	f, err := os.Open(scanPath)
	if err != nil {
		return fmt.Errorf("failed to find %ss : %w", typeName, err)
	}

	dirs, err := f.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("failed to find %ss : %w", typeName, err)
	}

	var text, content string

	for _, dir := range dirs {
		// Find file .version or .architecture
		p := path.Join(scanPath, dir, fmt.Sprintf(".%s", typeName))
		if _, err = os.Stat(p); err != nil {
			continue
		}
		text, err = readAllText(p)
		if err != nil {
			return fmt.Errorf("failed to find %ss : %w", typeName, err)
		}
		// Create prefix, i e "VERSION=" or "ARCHITECTURE="
		prefix := fmt.Sprintf("%s=", strings.ToUpper(typeName))
		content, err = getFirstLine(text, prefix, " \t\n")
		if err != nil {
			return fmt.Errorf("failed to find %ss : %w", typeName, err)
		}

		switch version == nil {
		case true:
			// We are scanning for versions
			v := &Version{Name: content, Path: path.Join(w.Path, dir)}
			w.Versions = append(w.Versions, v)

			if version == nil {
				// Scan architecture folders
				err = w.scanFolder(v)
				if err != nil {
					return err
				}
			}
		case false:
			// We are scanning for architectures
			a := &Architecture{Name: content, Path: p}
			version.Architectures = append(version.Architectures, a)
		}
	}

	return nil
}
