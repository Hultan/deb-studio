package gui

import 	_ "embed"

//go:embed assets/main.glade
var mainGlade string

//go:embed assets/debstudio_256.png
var applicationIcon []byte

//go:embed assets/mandatory.png
var mandatoryIcon []byte

//go:embed assets/recommended.png
var recommendedIcon []byte

//go:embed assets/optional.png
var optionalIcon []byte

//go:embed assets/save.png
var saveIcon []byte

//go:embed assets/exit.png
var exitIcon []byte

