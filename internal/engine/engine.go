package engine

import (
	"fmt"
	"os"
	"path"

	"github.com/hultan/deb-studio/internal/logger"
)

type engine struct {
	log *logger.Logger
}

type Project struct {
	Path        string
	ProgramName string
	Versions    []*Version
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

func newProject(projectPath, programName string) *Project {
	return &Project{
		ProgramName: programName,
		Path:        projectPath,
	}
}

func (e *engine) IsProjectFolder(projectPath string) bool {
	// Check if .program file exists...
	p := path.Join(projectPath, ".program")
	_, err := os.Stat(p)
	if err != nil {
		// File .program does not exist, or possibly a permission error...
		return false
	}
	return true
}

func (e *engine) OpenProject(projectPath string) (*Project, error) {
	// Make sure that directory exists
	if !doesDirectoryExist(projectPath) {
		return nil, ErrorProjectFolderMissing
	}

	programName, err := getProgramName(projectPath)
	if err != nil {
		return nil, ErrorNewProjectFolder
	}

	w := newProject(projectPath, programName)

	// err = w.scanVersions()
	err = w.scanFolder(nil)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (e *engine) SetupProject(projectPath, programName string) (*Project, error) {
	// Make sure that directory exists
	if !doesDirectoryExist(projectPath) {
		return nil, ErrorProjectFolderMissing
	}

	// Create program file
	filePath := path.Join(projectPath, ".program")
	content := fmt.Sprintf("PROGRAM=%s", programName)
	err := createTextFile(filePath, content)
	if err != nil {
		return nil, err
	}

	return newProject(projectPath, programName), nil
}

//
// func (w *Project) scanVersions() error {
// 	// Open project path
// 	f, err := os.Open(w.Path)
// 	if err != nil {
// 		return fmt.Errorf("failed to find versions : %w", err)
// 	}
//
// 	dirs, err := f.Readdirnames(-1)
// 	if err != nil {
// 		return fmt.Errorf("failed to find versions : %w", err)
// 	}
//
// 	for _, dir := range dirs {
// 		p := path.Join(w.Path, dir, ".version")
// 		if _, err = os.Stat(p); err != nil {
// 			continue
// 		}
// 		text, err := readAllText(p)
// 		if err != nil {
// 			return fmt.Errorf("failed to find versions : %w", err)
// 		}
// 		version, err := getFirstLine(text, "VERSION=", " \t\n")
// 		if err != nil {
// 			return fmt.Errorf("failed to find versions : %w", err)
// 		}
// 		// We found a version folder
// 		v := &Version{Name: version, Path: path.Join(w.Path, dir)}
// 		w.Versions = append(w.Versions, v)
//
// 		err = w.scanArchitectures(v)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// func (w *Project) scanArchitectures(v *Version) error {
// 	// Open version path
// 	f, err := os.Open(v.Path)
// 	if err != nil {
// 		return fmt.Errorf("failed to find architectures : %w", err)
// 	}
//
// 	// Get names of all files and folders
// 	// Use readdir instead?
// 	dirs, err := f.Readdirnames(-1)
// 	if err != nil {
// 		return fmt.Errorf("failed to find architectures : %w", err)
// 	}
//
// 	for _, dir := range dirs {
// 		p := path.Join(v.Path, dir, ".architecture")
// 		if _, err = os.Stat(p); err != nil {
// 			continue
// 		}
// 		text, err := readAllText(p)
// 		if err != nil {
// 			return fmt.Errorf("failed to find architectures : %w", err)
// 		}
// 		version, err := getFirstLine(text, "ARCHITECTURE=", " \t\n")
// 		if err != nil {
// 			return fmt.Errorf("failed to find architectures : %w", err)
// 		}
// 		// We found an architecture folder
// 		a := &Architecture{Name: version, Path: p}
// 		v.Architectures = append(v.Architectures, a)
// 	}
//
// 	return nil
// }

func getProgramName(projectPath string) (string, error) {
	p := path.Join(projectPath, ".program")
	text, err := readAllText(p)
	if err != nil {
		return "", fmt.Errorf("failed to find program name : %w", err)
	}
	name, err := getFirstLine(text, "PROGRAM=", "\t\n")
	if err != nil {
		return "", fmt.Errorf("failed to find program name : %w", err)
	}

	return name, nil
}
