package engine

import (
	"os"
	"path"
)

type Version struct {
	Path          string
	Name          string
	Architectures []*Architecture
}

func (v *Version) AddArchitecture(architectureName string) (*Architecture, error) {
	log.Trace.Println("Entering AddArchitecture...")
	defer log.Trace.Println("Exiting AddArchitecture...")

	architecturePath := path.Join(v.Path, architectureName)
	err := os.MkdirAll(architecturePath, 0775)
	if err != nil {
		log.Error.Printf("Failed to create directory at path '%s': %s\n", architecturePath, err)
		return nil, err
	}

	err = writeDescriptor(architectureDescriptor, architecturePath, architectureName)
	if err != nil {
		log.Error.Printf("Failed to create .architecture file at path '%s': %s\n", architecturePath, err)
		return nil, err
	}

	// Add to version slice
	a := &Architecture{Name: architectureName, Path: architecturePath}
	v.Architectures = append(v.Architectures, a)

	log.Info.Printf("Created architecture %s...\n", architectureName)
	return a, nil

}

func (v *Version) scanForArchitectures() error {
	log.Trace.Println("Entering scanForArchitectures...")
	defer log.Trace.Println("Exiting scanForArchitectures...")

	// Open path to scan
	f, err := os.Open(v.Path)
	if err != nil {
		log.Error.Printf("Failed to open path '%s': %s\n", v.Path, err)
		return err
	}

	dirs, err := f.Readdirnames(-1)
	if err != nil {
		log.Error.Printf("Failed to read dir names of path '%s': %s\n", v.Path, err)
		return err
	}

	for _, dir := range dirs {
		architecturePath := path.Join(v.Path, dir)

		// Check if .architecture file exists
		if !haveDescriptor(architectureDescriptor, architecturePath) {
			continue
		}

		architectureName, err := readDescriptor(architectureDescriptor, architecturePath)
		if err != nil {
			log.Error.Printf("Failed to get descriptor in path '%s': %s\n", architecturePath, err)
			return err
		}

		// Add architecture
		a := &Architecture{Name: architectureName, Path: architecturePath}
		v.Architectures = append(v.Architectures, a)
	}

	return nil
}
