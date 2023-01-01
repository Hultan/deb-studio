package engine

import (
	"os"
	"path"
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
	v := &Version{Name: versionName, Path: versionPath}
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
		v := &Version{Name: versionName, Path: versionPath}
		p.Versions = append(p.Versions, v)

		// Scan architecture folders
		err = v.scanForArchitectures()
		if err != nil {
			log.Error.Printf("Failed to scan path '%s': %s\n", versionPath, err)
			return err
		}
	}

	return nil
}

func (p *Project) GetVersion(name string) *Version {
	log.Trace.Println("Entering GetVersion...")
	defer log.Trace.Println("Exiting GetVersion...")

	for i := range p.Versions {
		if p.Versions[i].Name == name {
			return p.Versions[i]
		}
	}
	log.Trace.Printf("failed to find version '%s'", name)
	return nil
}

func (p *Project) GetArchitecture(versionName, architectureName string) *Architecture {
	log.Trace.Println("Entering GetArchitecture...")
	defer log.Trace.Println("Exiting GetArchitecture...")

	for i := range p.Versions {
		if p.Versions[i].Name == versionName {
			for a := range p.Versions[i].Architectures {
				if p.Versions[i].Architectures[a].Name == architectureName {
					return p.Versions[i].Architectures[a]
				}
			}
			log.Trace.Printf("failed to find architecture '%s' in version '%s'", architectureName, versionName)
			return nil
		}
	}
	log.Trace.Printf("failed to find version '%s'", versionName)
	return nil
}

// TODO : Switch to another solution that is not dependant on sorting versions
// and also move to Version.IsLatestVersion().

func (p *Project) IsLatestVersion(v *Version) bool {
	log.Trace.Println("Entering IsLatestVersion...")
	defer log.Trace.Println("Exiting IsLatestVersion...")

	if len(p.Versions) == 0 || v == nil {
		return false
	}
	return p.Versions[0].Name == v.Name
}

func (p *Project) SetCurrent(versionName, architectureName string) {
	log.Trace.Println("Entering SetCurrent...")
	defer log.Trace.Println("Exiting SetCurrent...")

	p.CurrentVersion = p.GetVersion(versionName)
	if p.CurrentVersion == nil {
		log.Error.Printf("failed to find version '%s'\n", versionName)
		return
	}
	p.CurrentArchitecture = p.CurrentVersion.GetArchitecture(architectureName)
	if p.CurrentArchitecture == nil {
		log.Error.Printf("failed to find architecture '%s' in version '%s'\n", architectureName, versionName)
		p.CurrentVersion = nil
		return
	}
}
