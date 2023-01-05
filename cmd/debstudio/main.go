package main

import (
	"fmt"
	"os"

	debStudio "github.com/hultan/deb-studio/internal/gui"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	ApplicationId    = "se.softteam.debstudio"
	ApplicationFlags = glib.APPLICATION_FLAGS_NONE
)

func main() {
	// Create a new application
	application, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to create GTK Application")
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create the main form
	mainForm := debStudio.NewMainWindow()
	application.Connect("activate", mainForm.Open)

	// Start the application (and exit when it is done)
	os.Exit(application.Run(nil))
}
