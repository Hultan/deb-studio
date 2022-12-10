package gui

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// createPixBufFromBytes: Create a *gdk.Pixbuf from a slice of bytes
func createPixBufFromBytes(bytes []byte, name string) *gdk.Pixbuf {
	pix, err := gdk.PixbufNewFromBytesOnly(bytes)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create pix buf for %s icon\n", name)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return pix
}

// createImageFromBytes: Creates a *gtk.Image from []byte
func createImageFromBytes(bytes []byte, name string) *gtk.Image {
	pix := createPixBufFromBytes(bytes, name)
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
