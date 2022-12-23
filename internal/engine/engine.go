package engine

import (
	"fmt"
	"os"
	"path"
)

type Workspace struct {
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

func Open(workspacePath string) (*Workspace, error) {
	if !doesDirectoryExist(workspacePath) {
		return nil, ErrorWorkspaceFolderMissing
	}

	name, err := getProgramName(workspacePath)
	if err != nil {
		return nil, ErrorNewWorkspaceFolder
	}

	w := &Workspace{
		ProgramName: name,
		Path:        workspacePath,
	}

	err = w.scanVersions()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func SetupWorkspaceFolder(workspacePath, programName string) (*Workspace, error) {
	// Create program file
	content := fmt.Sprintf("PROGRAM=%s", programName)
	err := createTextFile(workspacePath, ".program", content)
	if err != nil {
		return nil, err
	}

	return &Workspace{ProgramName: programName}, nil
}

func (w *Workspace) scanVersions() error {
	// Open workspace path
	f, err := os.Open(w.Path)
	if err != nil {
		return fmt.Errorf("failed to find versions : %w", err)
	}

	dirs, err := f.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("failed to find versions : %w", err)
	}

	for _, dir := range dirs {
		p := path.Join(w.Path, dir, ".version")
		if _, err = os.Stat(p); err != nil {
			continue
		}
		text, err := readAllText(p)
		if err != nil {
			return fmt.Errorf("failed to find versions : %w", err)
		}
		version, err := getFirstLine(text, "VERSION=", " \t\n")
		if err != nil {
			return fmt.Errorf("failed to find versions : %w", err)
		}
		// We found a version folder
		v := &Version{Name: version, Path: path.Join(w.Path, dir)}
		w.Versions = append(w.Versions, v)

		err = w.scanArchitectures(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Workspace) scanArchitectures(v *Version) error {
	// Open version path
	f, err := os.Open(v.Path)
	if err != nil {
		return fmt.Errorf("failed to find architectures : %w", err)
	}

	// Get names of all files and folders
	// Use readdir instead?
	dirs, err := f.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("failed to find architectures : %w", err)
	}

	for _, dir := range dirs {
		p := path.Join(v.Path, dir, ".architecture")
		if _, err = os.Stat(p); err != nil {
			continue
		}
		text, err := readAllText(p)
		if err != nil {
			return fmt.Errorf("failed to find architectures : %w", err)
		}
		version, err := getFirstLine(text, "ARCHITECTURE=", " \t\n")
		if err != nil {
			return fmt.Errorf("failed to find architectures : %w", err)
		}
		// We found an architecture folder
		a := &Architecture{Name: version, Path: p}
		v.Architectures = append(v.Architectures, a)
	}

	return nil
}

func getProgramName(workspacePath string) (string, error) {
	p := path.Join(workspacePath, ".program")
	if _, err := os.Stat(p); err != nil {
		return "", err
	}
	text, err := readAllText(p)
	if err != nil {
		return "", fmt.Errorf("failed to find program name : %w", err)
	}
	name, err := getFirstLine(text, "PROGRAM=", "\n")
	if err != nil {
		return "", fmt.Errorf("failed to find program name : %w", err)
	}

	return name, nil
}
