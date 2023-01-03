package engine

import (
	"fmt"
	"io"
	"os"
	"path"

	"pault.ag/go/debian/control"

	"github.com/hultan/deb-studio/internal/config/packageConfig"
)

// See advanced control file example here: https://gist.github.com/citrusui/c3358f9661550e8cb849

const emptyControlFile = `Section: 
Package: 
License: 
Vendor: 
Version: 
Architecture: 
Essential: no
Priority: 
Depends: 
Installed-Size: 
Maintainer: 
Description: 
HomePage: `

type Package struct {
	*control.Control
	Path   string
	Config *packageConfig.PackageConfig
}

func newPackage(packagePath string, config *packageConfig.PackageConfig) (*Package, error) {
	controlFile, err := control.ParseControlFile(config.GetControlFilePath(packagePath))
	if os.IsNotExist(err) {
		// Control file does not exist, create an empty one
		err = createEmptyControlFile(config.GetControlFilePath(packagePath))
		if err != nil {
			return nil, err
		}
		controlFile, err = control.ParseControlFile(config.GetControlFilePath(packagePath))
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &Package{
		Control: controlFile,
		Path:    packagePath,
		Config:  config,
	}, nil
}

func createEmptyControlFile(controlPath string) error {
	f, err := os.Create(controlPath)
	if err != nil {
		return err
	}
	_, err = f.WriteString(emptyControlFile)
	if err != nil {
		return err
	}
	return nil
}

func (p *Package) SaveControlFile() error {
	f, _ := os.OpenFile(p.Config.GetControlFilePath(p.Path), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	err := p.Source.WriteTo(f)
	if err != nil {
		return err
	}
	return nil
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
