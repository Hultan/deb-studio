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
		_,_ = fmt.Fprintln(os.Stderr,"failed to create GTK Application")
		_,_ = fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}

	mainForm := debStudio.NewMainForm()
	// Hook up the activate event handler
	application.Connect("activate", mainForm.Open)
	if err != nil {
		_,_ = fmt.Fprintln(os.Stderr,"failed to connect Application.Activate event")
		_,_ = fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}

	// Start the application (and exit when it is done)
	os.Exit(application.Run(nil))
}
