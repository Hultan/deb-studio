package engine

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"pault.ag/go/debian/control"

	"github.com/hultan/deb-studio/internal/common"
	"github.com/hultan/deb-studio/internal/config/packageConfig"
)

// See advanced control file example here: https://gist.github.com/citrusui/c3358f9661550e8cb849

const emptyControlFile = `Package: 
Version: 
Architecture:
Priority: optional
Depends: 
Maintainer: 
Description: `

type Package struct {
	*control.Control
	Path   string
	Config *packageConfig.PackageConfig
}

func newPackage(packagePath string, config *packageConfig.PackageConfig) (*Package, error) {
	log.Trace.Println("Entering newPackage...")
	defer log.Trace.Println("Exiting newPackage...")

	p := &Package{Config: config, Path: packagePath}

	// Load control file
	controlFilePath := path.Join(packagePath, p.GetPackageFolderName(), common.FolderNameDebian, common.FileNameControl)
	controlFile, err := control.ParseControlFile(controlFilePath)
	if os.IsNotExist(err) {
		// Control file does not exist, create an empty one
		err = createEmptyControlFile(controlFilePath)
		if err != nil {
			log.Error.Printf("failed to create empty control file: %s", err)
			return nil, err
		}
		controlFile, err = control.ParseControlFile(controlFilePath)
		if err != nil {
			log.Error.Printf("failed to parse empty control file: %s", err)
			return nil, err
		}
	} else if err != nil {
		log.Error.Printf("failed to parse control file: %s", err)
		return nil, err
	}
	p.Control = controlFile

	return p, nil
}

func (p *Package) SaveControlFile() error {
	log.Trace.Println("Entering SaveControlFile...")
	defer log.Trace.Println("Exiting SaveControlFile...")

	controlFilePath := path.Join(p.Path, p.GetPackageFolderName(), common.FolderNameDebian, common.FileNameControl)
	f, _ := os.OpenFile(controlFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	err := p.Source.WriteTo(f)
	if err != nil {
		log.Error.Printf("failed to save control file: %s", err)
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

func (p *Package) GetPackageFolderName() string {
	log.Trace.Println("Entering GetPackageFolderName...")
	defer log.Trace.Println("Exiting GetPackageFolderName...")

	return fmt.Sprintf(
		"%s-%s-%s",
		strings.ToLower(p.Config.Project), p.Config.Version, p.Config.Architecture,
	)
}

func copyFile(src, dst string) (int64, error) {
	log.Trace.Println("Entering copyFile...")
	defer log.Trace.Println("Exiting copyFile...")

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		log.Error.Printf("failed to stat source file: %s", err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		log.Warning.Printf("source file is not a regular file: %s", src)
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		log.Error.Printf("failed to open source file (%s): %s", src, err)
		return 0, err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Error.Printf("failed to open destination file (%s): %s", dst, err)
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func createEmptyControlFile(controlPath string) error {
	log.Trace.Println("Entering createEmptyControlFile...")
	defer log.Trace.Println("Exiting createEmptyControlFile...")

	f, err := os.Create(controlPath)
	if err != nil {
		log.Error.Printf("failed to create file '%s': %s", controlPath, err)
		return err
	}
	_, err = f.WriteString(emptyControlFile)
	if err != nil {
		log.Error.Printf("failed to write to file '%s': %s", controlPath, err)
		return err
	}
	return nil
}
