package gtk

import (
	"errors"
	"path"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Builder struct {
	Builder *gtk.Builder
}

// Create : Creates a Builder
func Create(fileNameOrPath string) (*Builder, error) {
	if !fileExists(fileNameOrPath) {
		fileNameOrPath = GetResourcePath(fileNameOrPath)
	}

	builder, err := gtk.BuilderNewFromFile(fileNameOrPath)
	if err != nil {
		return nil, err
	}

	return &Builder{builder}, nil
}

// GetObject : Gets a gtk object by name
func (g *Builder) GetObject(name string) glib.IObject {
	if g.Builder == nil {
		panic(errors.New("need to manually set gtk.Builder or call gtkBuilder.Create() first"))
	}
	obj, err := g.Builder.GetObject(name)
	if err != nil {
		// We panic here since the glade file is invalid and
		// the application cannot work without a valid glade file
		panic(err)
	}

	return obj
}

// GetResourcePath : Gets the path for a single resource file
func GetResourcePath(fileName string) string {
	resourcesPath := getResourcesPath()
	resourcePath := path.Join(resourcesPath, fileName)
	return resourcePath
}

// getResourcesPath : Returns the resources path
func getResourcesPath() string {
	executablePath := getExecutablePath()

	var pathsToCheck []string
	pathsToCheck = append(pathsToCheck, path.Join(executablePath, "../assets"))
	pathsToCheck = append(pathsToCheck, path.Join(executablePath, "assets"))

	dir, err := checkPathsExists(pathsToCheck)
	if err != nil {
		return executablePath
	}
	return dir
}

// checkPathsExists : Returns the first path that exists
func checkPathsExists(pathsToCheck []string) (string, error) {
	for _, pathToCheck := range pathsToCheck {
		if directoryExists(pathToCheck) {
			return pathToCheck, nil
		}
	}
	return "", errors.New("paths do not exist")
}
