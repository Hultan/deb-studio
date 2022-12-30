package engine

type Version struct {
	Path          string
	Name          string
	Architectures []*Architecture
}

func (v *Version) AddArchitecture(architectureName string) {

}
