package gui

import _ "embed"

// TODO : Switch buildIcon to package (see aboutDialog)

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

//go:embed assets/new.png
var newIcon []byte

//go:embed assets/open.png
var openIcon []byte

//go:embed assets/save.png
var saveIcon []byte

//go:embed assets/build.png
var buildIcon []byte

//go:embed assets/exit.png
var exitIcon []byte

//go:embed assets/addFile.png
var addFileIcon []byte

//go:embed assets/editFile.png
var editFileIcon []byte

//go:embed assets/removeFile.png
var removeFileIcon []byte

//go:embed assets/check.png
var checkIcon []byte

//go:embed assets/edit.png
var editIcon []byte
