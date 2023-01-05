package gui

// Image type constants

type imageType int

const (
	imageTypeMandatory imageType = iota
	imageTypeRecommended
	imageTypeOptional
)

// Project status

type projectStatus int

const (
	projectStatusNoProjectOpened projectStatus = iota
	projectStatusNoPackageSelected
	projectStatusLatestVersion
	projectStatusNotLatestVersion
)
