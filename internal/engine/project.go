package engine

import (
	"os"
	"path"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Project struct {
	Path                string
	Name                string
	Versions            []*Version
	CurrentVersion      *Version
	CurrentArchitecture *Architecture
}

func (p *Project) AddVersion(versionName string) (*Version, error) {
	log.Trace.Println("Entering AddVersion...")
	defer log.Trace.Println("Exiting AddVersion...")

	versionPath := path.Join(p.Path, versionName)
	err := os.MkdirAll(versionPath, 0775)
	if err != nil {
		log.Error.Printf("Failed to create directory at path '%s': %s\n", versionPath, err)
		return nil, err
	}

	err = writeDescriptor(versionDescriptor, versionPath, versionName)
	if err != nil {
		log.Error.Printf("Failed to create .version file at path '%s': %s\n", versionPath, err)
		return nil, err
	}

	// Add to version slice
	v := newVersion(versionName, versionPath)
	p.SetAsLatest(v)
	p.Versions = append(p.Versions, v)

	log.Info.Printf("Created version %s...\n", versionName)
	return v, nil
}

func (p *Project) scanForVersions() error {
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
		versionPath := path.Join(p.Path, dir)

		// Check if .version file exists
		if !haveDescriptor(versionDescriptor, versionPath) {
			continue
		}

		versionName, err := readDescriptor(versionDescriptor, versionPath)
		if err != nil {
			log.Error.Printf("Failed to get descriptor in path '%s': %s\n", versionPath, err)
			return err
		}

		// Add version
		v := newVersion(versionName, versionPath)
		p.Versions = append(p.Versions, v)

		// Scan architecture folders
		err = v.scanForArchitectures()
		if err != nil {
			log.Error.Printf("Failed to scan path '%s': %s\n", versionPath, err)
			return err
		}
	}

	p.setLatestAsLatest()

	return nil
}

func (p *Project) GetVersion(name string) *Version {
	for i := range p.Versions {
		if p.Versions[i].Name == name {
			return p.Versions[i]
		}
	}
	return nil
}

func (p *Project) GetVersionByGuid(guid string) *Version {
	for i := range p.Versions {
		if p.Versions[i].Guid.String() == guid {
			return p.Versions[i]
		}
	}
	return nil
}

func (p *Project) GetArchitectureByGuid(guid string) *Architecture {
	for i := range p.Versions {
		for j := range p.Versions[i].Architectures {
			if p.Versions[i].Architectures[j].Guid.String() == guid {
				return p.Versions[i].Architectures[j]
			}
		}
	}
	return nil
}

func (p *Project) SetAsCurrent(a *Architecture) {
	log.Trace.Println("Entering SetAsCurrent...")
	defer log.Trace.Println("Exiting SetAsCurrent...")

	p.CurrentVersion = p.GetVersion(a.Version.Name)
	if p.CurrentVersion == nil {
		log.Error.Printf("failed to find version '%s'\n", a.Version.Name)
		return
	}
	p.CurrentArchitecture = p.CurrentVersion.GetArchitecture(a.Name)
	if p.CurrentArchitecture == nil {
		log.Error.Printf("failed to find architecture '%s' in version '%s'\n", a.Name, a.Version.Name)
		p.CurrentVersion = nil
		return
	}
}

//
// func (p *Project) SetAsCurrent(versionName, architectureName string) {
// 	log.Trace.Println("Entering SetAsCurrent...")
// 	defer log.Trace.Println("Exiting SetAsCurrent...")
//
// 	p.CurrentVersion = p.GetVersion(versionName)
// 	if p.CurrentVersion == nil {
// 		log.Error.Printf("failed to find version '%s'\n", versionName)
// 		return
// 	}
// 	p.CurrentArchitecture = p.CurrentVersion.GetArchitecture(architectureName)
// 	if p.CurrentArchitecture == nil {
// 		log.Error.Printf("failed to find architecture '%s' in version '%s'\n", architectureName, versionName)
// 		p.CurrentVersion = nil
// 		return
// 	}
// }

func (p *Project) SetAsLatest(v *Version) {
	for i := range p.Versions {
		p.Versions[i].IsLatest = false
	}
	v.IsLatest = true
}

// TEMPORARY FUNCTION! REMOVE!
func (p *Project) setLatestAsLatest() {
	latest := 0

	for i := 1; i < len(p.Versions); i++ {
		if p.Versions[i].Name > p.Versions[latest].Name {
			latest = i
		}
	}

	p.SetAsLatest(p.Versions[latest])
}

func (p *Project) GetPackageListStore(checkIcon []byte) *gtk.ListStore {
	log.Trace.Println("Entering GetPackageListStore...")
	defer log.Trace.Println("Exiting GetPackageListStore...")

	// Icon, Version name, Architecture name, Version guid, Architecture guid
	s, err := gtk.ListStoreNew(
		gdk.PixbufGetType(), glib.TYPE_STRING, glib.TYPE_STRING,
		glib.TYPE_STRING, glib.TYPE_STRING,
	)
	if err != nil {
		log.Error.Printf("failed to create new list store: %s\n", err)
		return nil
	}

	check, _ := gdk.PixbufNewFromBytesOnly(checkIcon)
	for _, version := range p.Versions {
		for _, architecture := range version.Architectures {
			iter := s.InsertAfter(nil)
			data := []interface{}{
				nil,
				version.Name, architecture.Name,
				version.Guid, architecture.Guid,
			}
			if version.IsLatest {
				data[0] = check
			}
			_ = s.Set(iter, []int{0, 1, 2, 3, 4}, data)
		}
	}

	s.SetSortFunc(
		1, func(model *gtk.TreeModel, a, b *gtk.TreeIter) int {
			va, _ := model.GetValue(a, 1)
			vb, _ := model.GetValue(b, 1)
			vaName, _ := va.GetString()
			vbName, _ := vb.GetString()

			return strings.Compare(vbName, vaName)
		},
	)

	s.SetSortColumnId(1, gtk.SORT_ASCENDING)

	return s
}
