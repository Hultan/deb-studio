package engine

import (
	"os"
	"path"
)

type Project struct {
	Path     string
	Name     string
	Versions []*Version
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
	for i := range p.Versions {
		if p.Versions[i].Name == name {
			return p.Versions[i]
		}
	}
	return nil
}

func (p *Project) GetArchitecture(versionName, architectureName string) *Architecture {
	for i := range p.Versions {
		if p.Versions[i].Name == versionName {
			for a := range p.Versions[i].Architectures {
				if p.Versions[i].Architectures[a].Name == architectureName {
					return p.Versions[i].Architectures[a]
				}
			}
			return nil
		}
	}
	return nil
}

// TODO : Switch to another solution that is not dependant on sorting versions
// and also move to Version.IsLatestVersion().

func (p *Project) IsLatestVersion(v *Version) bool {
	if len(p.Versions) == 0 || v == nil {
		return false
	}
	return p.Versions[0].Name == v.Name
}
