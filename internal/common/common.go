package common

// Application constants

const ApplicationTitle = "Deb Studio"
const ApplicationVersion = "v 0.2.4"
const ApplicationCopyRight = "©SoftTeam AB, 2022"

// Exit codes

const ExitCodeGtkError = 1

// Package list columns

const (
	PackageListColumnFilter = iota
	PackageListColumnIsCurrent
	PackageListColumnIsLatest
	PackageListColumnVersionName
	PackageListColumnArchitectureName
	PackageListColumnPackagePath
	PackageListColumnPackageId
)

// Folder names

const (
	FolderNameDebian = "DEBIAN"
	FolderNameLog    = "/home/per/.softteam/debstudio"
)

// File names

const (
	FileNameControl     = "control"
	FileNamePackageJson = "package.json"
	FileNameProjectJson = "project.json"
	FileNamePreInstall  = "preinst"
	FileNamePostInstall = "postinst"
	FileNamePreRemove   = "prerm"
	FileNamePostRemove  = "postrm"
	FileNameLog         = "debstudio.log"
	FileNameCopyRight   = "copyright"
	FileNameChangeLog   = "changelog"
	FileNameReadme      = "README.debian"
)

// Misc constants

const RightMouseButton = 3
