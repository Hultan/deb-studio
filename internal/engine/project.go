package engine

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
	"github.com/hultan/deb-studio/internal/config/packageConfig"
	"github.com/hultan/deb-studio/internal/config/projectConfig"
	"github.com/hultan/deb-studio/internal/iconFactory"
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

	log.Trace.Println("Entering OpenProject...")
	defer log.Trace.Println("Exiting OpenProject...")

	if !isProjectFolder(projectPath) {
		log.Warning.Println("missing project.json")
		return nil, errors.New("missing project.json")
	}

	config, err := projectConfig.Load(path.Join(projectPath, common.FileNameProjectJson))
	if err != nil {
		log.Error.Printf("failed to load project.json")
		return nil, err
	}

	p := &Project{Path: projectPath, Config: config}

	err = p.scanForPackages()
	if err != nil {
		log.Error.Printf("Failed to scan project path '%s'\n", projectPath)
		return nil, err
	}

	log.Info.Printf("Successfully opened project %s...\n", p.Config.Name)

	p.CurrentPackage = p.GetPackageById(p.Config.CurrentPackageId)

	return p, nil
}

func NewProject(l *logger.Logger, projectPath, projectName string) (*Project, error) {
	log = l

	log.Trace.Println("Entering NewProject...")
	defer log.Trace.Println("Exiting NewProject...")

	// Create project and config
	p := &Project{Path: projectPath}
	p.Config = &projectConfig.ProjectConfig{Name: projectName}

	// SaveControlFile config
	err := p.Config.Save(path.Join(projectPath, common.FileNameProjectJson))
	if err != nil {
		log.Error.Printf("failed to save control file: %s\n", err)
		return nil, err
	}

	return p, nil
}

func (p *Project) AddPackage(versionName, architectureName string) (*Package, error) {
	log.Trace.Println("Entering AddPackage...")
	defer log.Trace.Println("Exiting AddPackage...")

	config := &packageConfig.PackageConfig{
		Id:           uuid.New().String(),
		Project:      p.Config.Name,
		Version:      versionName,
		Architecture: architectureName,
		Files:        nil,
	}

	packageFolderName := fmt.Sprintf("%s-%s-%s", config.Project, config.Version, config.Architecture)
	packagePath := path.Join(p.Path, packageFolderName)
	err := os.MkdirAll(packagePath, 0775)
	if err != nil {
		log.Error.Printf(
			"Failed to create directory '%s' for package '%s': %s\n",
			packagePath, packageFolderName, err,
		)
		return nil, err
	}

	// Add to version slice
	pkg, err := newPackage(packagePath, config)
	if err != nil {
		log.Error.Printf("failed to create new package : %s\n", err)
		return nil, err
	}
	p.Packages = append(p.Packages, pkg)
	p.Config.LatestVersion = versionName

	log.Info.Printf("Created package %s...\n", packageFolderName)

	return pkg, nil
}

func (p *Project) IsWorkingWithLatestVersion() bool {
	log.Trace.Println("Entering IsWorkingWithLatestVersion...")
	defer log.Trace.Println("Exiting IsWorkingWithLatestVersion...")

	if p.CurrentPackage == nil {
		return false
	}
	return p.CurrentPackage.Config.Version == p.Config.LatestVersion
}

func (p *Project) GetPackageById(id string) *Package {
	log.Trace.Println("Entering GetPackageById...")
	defer log.Trace.Println("Exiting GetPackageById...")

	for i := range p.Packages {
		pkg := p.Packages[i]
		if pkg.Config.Id == id {
			return pkg
		}
	}
	return nil
}

func (p *Project) SetAsCurrent(id string) {
	log.Trace.Println("Entering SetAsCurrent...")
	defer log.Trace.Println("Exiting SetAsCurrent...")

	pkg := p.GetPackageById(id)
	if pkg == nil {
		log.Warning.Printf("no package with id '%s' found", id)
		return
	}
	p.Config.CurrentPackageId = pkg.Config.Id
	p.CurrentPackage = pkg
}

func (p *Project) SetAsLatest(id string) error {
	log.Trace.Println("Entering SetAsLatest...")
	defer log.Trace.Println("Exiting SetAsLatest...")

	pkg := p.GetPackageById(id)
	if pkg == nil {
		log.Error.Printf("failed to get package id\n")
		return errors.New("failed to get package id")
	}
	p.Config.LatestVersion = pkg.Config.Version

	return nil
}

func (p *Project) Save() error {
	log.Trace.Println("Entering Save...")
	defer log.Trace.Println("Exiting Save...")

	configPath := path.Join(p.Path, common.FileNameProjectJson)
	err := p.Config.Save(configPath)
	if err != nil {
		log.Error.Printf("failed to save project : %s\n", err)
		return err
	}

	return nil
}

func (p *Project) SetShowOnlyLatestVersion(checked bool) {
	log.Trace.Println("Entering SetShowOnlyLatestVersion...")
	defer log.Trace.Println("Exiting SetShowOnlyLatestVersion...")

	p.Config.ShowOnlyLatestVersion = checked
}

func (p *Project) GetPackageListStore() *gtk.TreeModelFilter {
	log.Trace.Println("Entering GetPackageListStore...")
	defer log.Trace.Println("Exiting GetPackageListStore...")

	// Icon, Version name, Architecture name, package name
	s, err := gtk.ListStoreNew(
		glib.TYPE_BOOLEAN, gdk.PixbufGetType(), gdk.PixbufGetType(),
		glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING,
	)
	if err != nil {
		log.Error.Printf("failed to create new list store: %s\n", err)
		return nil
	}

	// Get images
	fac := iconFactory.NewImageFactory()
	check := fac.GetPixBuf(iconFactory.ImageCheck)
	edit := fac.GetPixBuf(iconFactory.ImageEdit)
	for _, pkg := range p.Packages {
		iter := s.InsertAfter(nil)
		data := []interface{}{
			false, nil, nil,
			pkg.Config.Version, pkg.Config.Architecture,
			pkg.Path, pkg.Config.Id,
		}
		if pkg.Config.Version == p.Config.LatestVersion {
			data[common.PackageListColumnFilter] = true
			data[common.PackageListColumnIsLatest] = check
		}
		if pkg.Config.Id == p.Config.CurrentPackageId {
			data[common.PackageListColumnFilter] = true
			data[common.PackageListColumnIsCurrent] = edit
		}
		_ = s.Set(iter, []int{0, 1, 2, 3, 4, 5, 6}, data)
	}

	// Sorting
	s.SetSortFunc(
		1, func(model *gtk.TreeModel, a, b *gtk.TreeIter) int {
			va, _ := model.GetValue(a, common.PackageListColumnVersionName)
			vb, _ := model.GetValue(b, common.PackageListColumnVersionName)
			vaName, _ := va.GetString()
			vbName, _ := vb.GetString()

			return strings.Compare(vbName, vaName)
		},
	)
	s.SetSortColumnId(1, gtk.SORT_ASCENDING)

	// Filtering
	filter, err := s.FilterNew(&gtk.TreePath{})
	if err != nil {
		return nil
	}
	filter.SetVisibleFunc(p.filterFunc)

	return filter
}

func (p *Project) filterFunc(model *gtk.TreeModel, iter *gtk.TreeIter) bool {
	log.Trace.Println("Entering filterFunc...")
	defer log.Trace.Println("Exiting filterFunc...")

	if !p.Config.ShowOnlyLatestVersion {
		return true
	}

	value, err := model.GetValue(iter, common.PackageListColumnFilter)
	if err != nil {
		return true
	}
	goValue, err := value.GoValue()
	if err != nil {
		return true
	}
	filter := goValue.(bool)
	if err != nil {
		return true
	}

	return filter
}

func (p *Project) scanForPackages() error {
	log.Trace.Println("Entering scanForPackages...")
	defer log.Trace.Println("Exiting scanForPackages...")

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
		packageConfigPath := path.Join(packagePath, common.FileNamePackageJson)

		info, err := os.Stat(packageConfigPath)
		if err != nil || info.IsDir() {
			continue
		}

		config, err := packageConfig.Load(packageConfigPath)
		if err != nil {
			log.Error.Printf("failed to load package.json at '%s': %s", packageConfigPath, err)
			continue
		}

		// Add package
		pkg, err := newPackage(packagePath, config)
		if err != nil {
			log.Error.Printf("failed to create new package: %s", err)
			return err
		}
		p.Packages = append(p.Packages, pkg)
	}

	return nil
}

func isProjectFolder(projectPath string) bool {
	log.Trace.Println("Entering isProjectFolder...")
	defer log.Trace.Println("Exiting isProjectFolder...")

	info, err := os.Stat(path.Join(projectPath, common.FileNameProjectJson))
	if err != nil {
		return false
	}
	return !info.IsDir()
}
