package gui

// Application constants

const applicationTitle = "Deb Studio"
const applicationVersion = "v 0.1.1"
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

// Info bar

type infoBarStatus int

const (
	infoBarStatusNoPackageSelected infoBarStatus = iota
	infoBarStatusLatestVersion
	infoBarStatusNotLatestVersion
)

// Misc constants

const RightMouseButton = 3
