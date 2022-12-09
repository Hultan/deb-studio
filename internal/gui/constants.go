package gui

const applicationTitle = "Deb Studio"
const applicationVersion = "v 0.01"
const applicationCopyRight = "©SoftTeam AB, 2022"

const exitCodeSetupError = 1
const exitCodeGtkError = 2

type imageType int

const (
	imageTypeMandatory imageType = iota
	imageTypeRecommended
	imageTypeOptional
)

