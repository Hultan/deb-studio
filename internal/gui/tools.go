package gui

import (
	"fmt"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
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
	return fmt.Sprintf(
		"%s %s - %s",
		common.ApplicationTitle,
		common.ApplicationVersion,
		common.ApplicationCopyRight,
	)
}

func isTraceMode() bool {
	// TODO : Fix trace mode
	return len(os.Args) >= 2 && strings.HasPrefix(os.Args[1], "-t")
}
