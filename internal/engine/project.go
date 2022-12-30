package engine

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type Project struct {
	Path     string
	Name     string
	Versions []*Version
}

func (p *Project) AddVersion(versionName string) error {
	log.Trace.Println("Entering AddVersion...")
	defer log.Trace.Println("Exiting AddVersion...")

	// Create version directory
	versionPath := path.Join(p.Path, versionName)
	err := os.MkdirAll(versionPath, 0775)
	if err != nil {
		log.Error.Printf("Failed to create directory at path '%s': %s\n", versionPath, err)
		return err
	}

	// Create .version file
	filePath := path.Join(versionPath, versionFileName)
	content := fmt.Sprintf("VERSION=%s", versionName)
	err = createTextFile(filePath, content)
	if err != nil {
		log.Error.Printf("Failed to write .version file to path '%s': %s\n", filePath, err)
		return err
	}

	// Add to version slice
	v := &Version{Name: versionName, Path: versionPath}
	p.Versions = append(p.Versions, v)

	log.Info.Printf("Created version %s...\n", versionName)
	return nil
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

func (p *Project) scanProjectPath(version *Version) error {
	log.Trace.Println("Entering scanProjectPath...")
	defer log.Trace.Println("Exiting scanProjectPath...")

	typeName := "version"
	scanPath := p.Path
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
		scanDir := path.Join(scanPath, dir)

		// Find file .version or .architecture
		filePath := path.Join(scanDir, fmt.Sprintf(".%s", typeName))
		if _, err = os.Stat(filePath); err != nil {
			continue
		}
		text, err = readAllText(filePath)
		if err != nil {
			log.Error.Printf("Failed to read all text of path '%s': %s\n", filePath, err)
			return err
		}
		// Create prefix, i e "VERSION=" or "ARCHITECTURE="
		prefix := fmt.Sprintf("%s=", strings.ToUpper(typeName))
		content, err = getFirstLine(text, prefix, " \t\n")
		if err != nil {
			log.Error.Printf("Failed to get first line of path '%s': %s\n", filePath, err)
			return err
		}

		switch version == nil {
		case true:
			// We are scanning for versions
			v := &Version{Name: content, Path: scanDir}
			p.Versions = append(p.Versions, v)

			if version == nil {
				// Scan architecture folders
				err = p.scanProjectPath(v)
				if err != nil {
					log.Error.Printf("Failed to scan path '%s': %s\n", filePath, err)
					return err
				}
			}
		case false:
			// We are scanning for architectures
			a := &Architecture{Name: content, Path: scanDir}
			version.Architectures = append(version.Architectures, a)
		}
	}

	return nil
}
