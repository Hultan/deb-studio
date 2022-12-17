package engine

import (
	"errors"
)

var ErrorWorkspaceFolderMissing = errors.New("missing workspace folder")
var ErrorNewWorkspaceFolder = errors.New("missing .program file")
