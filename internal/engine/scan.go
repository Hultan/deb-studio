package engine

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type scanFolderType int

const (
	scanHolderTypeVersion scanFolderType = iota
	scanFolderTypeArchitecture
)

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
		p := path.Join(w.Path, dir, fmt.Sprintf(".%s", typeName))
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
