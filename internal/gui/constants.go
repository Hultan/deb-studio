package gui

// Project status

type projectStatus int

const (
	projectStatusNoProjectOpened projectStatus = iota
	projectStatusNoPackageSelected
	projectStatusLatestVersion
	projectStatusNotLatestVersion
)
