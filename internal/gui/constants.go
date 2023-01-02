package gui

// Application constants

const applicationTitle = "Deb Studio"
const applicationVersion = "v 0.1.4"
const applicationCopyRight = "©SoftTeam AB, 2022"

// Exit codes

const exitCodeSetupError = 1
const exitCodeGtkError = 2

// Image type constants

type imageType int

const (
	imageTypeMandatory imageType = iota
	imageTypeRecommended
	imageTypeOptional
)

// Package list columns

const (
	packageListColumnIsLatest = iota
	packageListColumnVersionName
	packageListColumnArchitectureName
	packageListColumnVersionGuid
	packageListColumnArchitectureGuid
)

// Info bar

type infoBarStatus int

const (
	infoBarStatusNoProjectOpened infoBarStatus = iota
	infoBarStatusNoPackageSelected
	infoBarStatusLatestVersion
	infoBarStatusNotLatestVersion
)

// Misc constants

const RightMouseButton = 3
const separator = "$$$"
