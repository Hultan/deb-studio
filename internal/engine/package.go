package engine

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/hultan/deb-studio/internal/config/packageConfig"
)

type Package struct {
	Path   string
	Config *packageConfig.PackageConfig
}

func newPackage(path string, config *packageConfig.PackageConfig) *Package {
	return &Package{
		Path:   path,
		Config: config,
	}
}

func (p *Package) AddFile(fromPath, fileName, userToPath string, copy bool) error {
	log.Trace.Println("Entering AddFile...")
	defer log.Trace.Println("Exiting AddFile...")

	// Create localToPath by joining architecture path with userToPath
	// and create it locally
	localToPath := path.Join(p.Path, userToPath)
	err := os.MkdirAll(localToPath, 0755)
	if err != nil {
		log.Error.Printf("failed to create directory %s : %s", localToPath, err)
		return err
	}

	fromFile := path.Join(fromPath, fileName)
	toFile := path.Join(localToPath, fileName)

	_, err = copyFile(fromFile, toFile)
	if err != nil {
		log.Error.Printf("failed to copy file from %s to %s : %s", fromFile, toFile, err)
		return err
	}

	log.Info.Printf("successfully added file %s to %s...\n", fromFile, toFile)
	return nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}