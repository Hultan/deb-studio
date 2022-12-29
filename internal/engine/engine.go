package engine

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hultan/deb-studio/internal/logger"
)

const projectFileName = ".project"

var log *logger.Logger

type Engine struct {
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

func NewEngine(l *logger.Logger) *Engine {
	log = l
	return &Engine{}
}

func newProject(projectPath, projectName string) *Project {
	return &Project{
		Name: projectName,
		Path: projectPath,
	}
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

	w := newProject(projectPath, projectName)

	err = w.scanProjectPath(nil)
	if err != nil {
		log.Error.Printf("Failed to scan project path '%s'\n", projectPath)
		return nil, err
	}

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

	return newProject(projectPath, projectName), nil
}

func getProjectName(projectPath string) (string, error) {
	log.Trace.Println("Entering getProjectName...")
	defer log.Trace.Println("Exiting getProjectName...")

	p := path.Join(projectPath, projectFileName)
	text, err := readAllText(p)
	if err != nil {
		log.Error.Printf("Failed to read .project file '%s': %s\n", p, err)
		return "", err
	}
	name, err := getFirstLine(text, "PROJECT=", "\t\n")
	if err != nil {
		log.Error.Printf("Failed to get first line of file '%s': %s\n", p, err)
		return "", err
	}

	return name, nil
}

func (w *Project) scanProjectPath(version *Version) error {
	log.Trace.Println("Entering scanProjectPath...")
	defer log.Trace.Println("Exiting scanProjectPath...")

	typeName := "version"
	scanPath := w.Path
	if version != nil {
		typeName = "architecture"
		scanPath = version.Path
	}

	// Open path to scan
	f, err := os.Open(scanPath)
	if err != nil {
		log.Error.Printf("Failed to open path '%s': %s\n", scanPath, err)
		return err
	}

	dirs, err := f.Readdirnames(-1)
	if err != nil {
		log.Error.Printf("Failed to read dir names of path '%s': %s\n", scanPath, err)
		return err
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
			log.Error.Printf("Failed to read all text of path '%s': %s\n", p, err)
			return err
		}
		// Create prefix, i e "VERSION=" or "ARCHITECTURE="
		prefix := fmt.Sprintf("%s=", strings.ToUpper(typeName))
		content, err = getFirstLine(text, prefix, " \t\n")
		if err != nil {
			log.Error.Printf("Failed to get first line of path '%s': %s\n", p, err)
			return err
		}

		switch version == nil {
		case true:
			// We are scanning for versions
			v := &Version{Name: content, Path: path.Join(w.Path, dir)}
			w.Versions = append(w.Versions, v)

			if version == nil {
				// Scan architecture folders
				err = w.scanProjectPath(v)
				if err != nil {
					log.Error.Printf("Failed to scan path '%s': %s\n", p, err)
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
