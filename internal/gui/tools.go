package gui

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// setControlImages: Set control tab images
func setControlImages(b *Builder) {
	setControlImage(b, "imgPackage", imageTypeMandatory)
	setControlImage(b, "imgSource", imageTypeOptional)
	setControlImage(b, "imgVersion", imageTypeMandatory)
	setControlImage(b, "imgSection", imageTypeRecommended)
	setControlImage(b, "imgPriority", imageTypeRecommended)
	setControlImage(b, "imgArchitecture", imageTypeMandatory)
	setControlImage(b, "imgEssential", imageTypeOptional)
	setControlImage(b, "imgDepends", imageTypeOptional)
	setControlImage(b, "imgInstalledSize", imageTypeOptional)
	setControlImage(b, "imgMaintainer", imageTypeMandatory)
	setControlImage(b, "imgDescription", imageTypeMandatory)
	setControlImage(b, "imgHomePage", imageTypeOptional)
	setControlImage(b, "imgBuiltUsing", imageTypeOptional)
}

// setControlImage: Set a control tab image
func setControlImage(b *Builder, imgName string, imgType imageType) {
	bytes := getControlIcon(imgType)
	img := b.GetObject(imgName).(*gtk.Image)
	img.SetFromPixbuf(createPixBufFromBytes(bytes))
}

// getControlIcon: Get the icon bytes from an image type
func getControlIcon(imgType imageType) []byte {
	var bytes []byte
	switch imgType {
	case imageTypeMandatory:
		bytes = mandatoryIcon
	case imageTypeRecommended:
		bytes = recommendedIcon
	case imageTypeOptional:
		bytes = optionalIcon
	}
	return bytes
}

// createBuilder: create a gtk builder, and wrap it in a Builder
func createBuilder() *Builder {
	// Create a new builder
	b, err := gtk.BuilderNewFromString(mainGlade)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCodeSetupError)
	}
	return &Builder{b}
}

// createPixBufFromBytes: Create a gdk.pixbuf from a slice of bytes
func createPixBufFromBytes(bytes []byte) *gdk.Pixbuf {
	pix, err := gdk.PixbufNewFromBytesOnly(bytes)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return pix
}

// createImageFromBytes: Creates a *gtk.Image from []byte
func createImageFromBytes(bytes []byte, name string) *gtk.Image {
	pix := createPixBufFromBytes(bytes)
	img, err := gtk.ImageNewFromPixbuf(pix)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create image for %s icon\n", name)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return img
}

// getApplicationName: Get the application name, version and copyright
func getApplicationName() string {
	return fmt.Sprintf("%s %s - %s", applicationTitle, applicationVersion, applicationCopyRight)
}

// getConfigPath: Get path to the config file
func getConfigPath() string {
	home := getHomeDirectory()

	return path.Join(home, "code/deb-studio/test.json")
}

// getHomeDirectory: Get current users home directory
func getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to get user home directory : %s\n", err)
		os.Exit(exitCodeSetupError)
	}
	return u.HomeDir
}
