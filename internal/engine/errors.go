package engine

import (
	"errors"
)

var ErrorProjectFolderMissing = errors.New("missing project folder")
var ErrorNewProjectFolder = errors.New("missing .program file")
