package main

import (
	"os"

	debStudio "github.com/hultan/deb-studio/internal/debstudio"

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
		panic("Failed to create GTK Application : " + err.Error())
	}

	mainForm := debStudio.NewMainForm()
	// Hook up the activate event handler
	application.Connect("activate", mainForm.OpenMainForm)
	if err != nil {
		panic("Failed to connect Application.Activate event : " + err.Error())
	}

	// Start the application (and exit when it is done)
	os.Exit(application.Run(nil))
}
