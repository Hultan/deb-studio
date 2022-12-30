package engine

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func Test_copyFile(t *testing.T) {
	testPath := "./../../test"
	testPath, err := filepath.Abs(testPath)
	if err != nil {
		t.Errorf("failed to get absolute path to test folder")
	}
	_, err = copyFile(path.Join(testPath, "fromFile"), path.Join(testPath, "toFile"))
	if err != nil {
		t.Errorf("CopyFile() failed to copy file to non-existing file: %s\n", err)
	}
	_, err = copyFile(path.Join(testPath, "fromFile"), path.Join(testPath, "existingFile"))
	if err != nil {
		t.Errorf("CopyFile() failed to copy file to existing file: %s\n", err)
	}
	os.Remove(path.Join(testPath, "toFile"))
}
