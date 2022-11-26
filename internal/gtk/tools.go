package gtk

import (
	"os"
	"path/filepath"
)

// fileExists : Checks if the specified file exists
func fileExists(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return !info.IsDir()
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

// directoryExists : Checks if the specified directory exists
func directoryExists(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return info.IsDir()
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

// getExecutableFile : Returns the path of the executable file
func getExecutableFile() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return ex
}

// getExecutablePath : Returns the directory of the executable file
func getExecutablePath() string {
	return filepath.Dir(getExecutableFile())
}
