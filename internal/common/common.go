package common

// Application constants

const ApplicationTitle = "Deb Studio"
const ApplicationVersion = "v 0.1.4"
const ApplicationCopyRight = "©SoftTeam AB, 2022"

// Exit codes

const ExitCodeSetupError = 1
const ExitCodeGtkError = 2

// File names

const PackageJsonFileName = "package.json"
const ProjectJsonFileName = "project.json"

// Package list columns

const (
	PackageListColumnIsLatest = iota
	PackageListColumnPackageName
	PackageListColumnVersionName
	PackageListColumnArchitectureName
)

// Misc constants

const RightMouseButton = 3
