package gui

// Image type constants

type imageType int

const (
	imageTypeMandatory imageType = iota
	imageTypeRecommended
	imageTypeOptional
)

// Info bar

type infoBarStatus int

const (
	infoBarStatusNoProjectOpened infoBarStatus = iota
	infoBarStatusNoPackageSelected
	infoBarStatusLatestVersion
	infoBarStatusNotLatestVersion
)
