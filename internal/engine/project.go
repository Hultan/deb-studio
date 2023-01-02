package engine

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
	"github.com/hultan/deb-studio/internal/config/packageConfig"
	"github.com/hultan/deb-studio/internal/config/projectConfig"
	"github.com/hultan/deb-studio/internal/logger"
)

var log *logger.Logger

type Project struct {
	Config         *projectConfig.ProjectConfig
	Path           string
	CurrentPackage *Package
	Packages       []*Package
}

func OpenProject(l *logger.Logger, projectPath string) (*Project, error) {
	log = l
	if !isProjectFolder(projectPath) {
		// TODO : Fix error handling
		return nil, errors.New("missing project.json")
	}

	config, err := projectConfig.Load(path.Join(projectPath, common.ProjectJsonFileName))
	if err != nil {
		return nil, err
	}

	p := &Project{Path: projectPath, Config: config}

	err = p.scanForPackages()
	if err != nil {
		log.Error.Printf("Failed to scan project path '%s'\n", projectPath)
		return nil, err
	}

	log.Info.Printf("Successfully opened project %s...\n", p.Config.Name)

	p.SetCurrentPackage()

	return p, nil
}

func NewProject(l *logger.Logger, projectPath, projectName string) (*Project, error) {
	log = l

	// Create project and config
	p := &Project{Path: projectPath}
	p.Config = &projectConfig.ProjectConfig{Name: projectName}

	// Save config
	err := p.Config.Save(path.Join(projectPath, common.ProjectJsonFileName))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func isProjectFolder(projectPath string) bool {
	info, err := os.Stat(path.Join(projectPath, common.ProjectJsonFileName))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (p *Project) AddPackage(versionName, architectureName string) (*Package, error) {
	log.Trace.Println("Entering AddPackage...")
	defer log.Trace.Println("Exiting AddPackage...")

	packageName := fmt.Sprintf("%s_%s", versionName, architectureName)
	packagePath := path.Join(p.Path, packageName)
	err := os.MkdirAll(packagePath, 0775)
	if err != nil {
		log.Error.Printf(
			"Failed to create directory '%s' for package '%s': %s\n",
			packagePath, packageName, err,
		)
		return nil, err
	}

	config := &packageConfig.PackageConfig{
		Name:         packageName,
		Version:      versionName,
		Architecture: architectureName,
		Files:        nil,
	}

	// Add to version slice
	pkg := newPackage(packagePath, config)
	p.Packages = append(p.Packages, pkg)
	p.Config.LatestVersion = versionName

	log.Info.Printf("Created package %s...\n", packageName)

	return pkg, nil
}

func (p *Project) GetPackageListStore(checkIcon []byte) *gtk.ListStore {
	log.Trace.Println("Entering GetPackageListStore...")
	defer log.Trace.Println("Exiting GetPackageListStore...")

	// Icon, Version name, Architecture name, package name
	s, err := gtk.ListStoreNew(
		gdk.PixbufGetType(), glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING,
	)
	if err != nil {
		log.Error.Printf("failed to create new list store: %s\n", err)
		return nil
	}

	check, _ := gdk.PixbufNewFromBytesOnly(checkIcon)
	for _, pkg := range p.Packages {
		iter := s.InsertAfter(nil)
		data := []interface{}{
			nil,
			pkg.Config.Name, pkg.Config.Version, pkg.Config.Architecture,
		}
		if pkg.Config.Version == p.Config.LatestVersion {
			data[0] = check
		}
		_ = s.Set(iter, []int{0, 1, 2, 3}, data)
	}

	s.SetSortFunc(
		1, func(model *gtk.TreeModel, a, b *gtk.TreeIter) int {
			va, _ := model.GetValue(a, 2)
			vb, _ := model.GetValue(b, 2)
			vaName, _ := va.GetString()
			vbName, _ := vb.GetString()

			return strings.Compare(vbName, vaName)
		},
	)

	s.SetSortColumnId(1, gtk.SORT_ASCENDING)

	return s
}

func (p *Project) WorkingWithLatestVersion() bool {
	if p.CurrentPackage == nil {
		return false
	}
	return p.CurrentPackage.Config.Version == p.Config.LatestVersion
}

func (p *Project) SetCurrentPackage() {
	for i := range p.Packages {
		pkg := p.Packages[i]
		if pkg.Config.Name == p.Config.CurrentPackage {
			p.CurrentPackage = pkg
		}
	}
}

func (p *Project) scanForPackages() error {
	log.Trace.Println("Entering scanForVersions...")
	defer log.Trace.Println("Exiting scanForVersions...")

	// Open path to scan
	f, err := os.Open(p.Path)
	if err != nil {
		log.Error.Printf("Failed to open path '%s': %s\n", p.Path, err)
		return err
	}

	dirs, err := f.Readdirnames(-1)
	if err != nil {
		log.Error.Printf("Failed to read dir names of path '%s': %s\n", p.Path, err)
		return err
	}

	for _, dir := range dirs {
		packagePath := path.Join(p.Path, dir)
		packageConfigPath := path.Join(packagePath, common.PackageJsonFileName)

		info, err := os.Stat(packageConfigPath)
		if err != nil || info.IsDir() {
			continue
		}

		config, err := packageConfig.Load(packageConfigPath)
		if err != nil {
			// TODO : Log and error handling
		}

		// Add package
		pkg := newPackage(packagePath, config)
		p.Packages = append(p.Packages, pkg)
	}

	return nil
}

func (p *Project) GetPackageByName(name string) *Package {
	for i := range p.Packages {
		pkg := p.Packages[i]
		if pkg.Config.Name == name {
			return pkg
		}
	}
	return nil
}

func (p *Project) SetAsCurrent(name string) {
	pkg := p.GetPackageByName(name)
	if pkg == nil {
		// TODO : Error handling
		return
	}
	p.Config.CurrentPackage = pkg.Config.Name
}

func (p *Project) SetAsLatest(name string) {
	pkg := p.GetPackageByName(name)
	if pkg == nil {
		// TODO : Error handling
		return
	}
	p.Config.LatestVersion = pkg.Config.Version
}
