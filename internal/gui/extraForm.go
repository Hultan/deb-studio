package gui

import (
	"github.com/gotk3/gotk3/gtk"

	builder "github.com/hultan/deb-studio/internal/gtk"
)

func openExtraForm() {
	// Create a new gtk helper
	builder, err := builder.Create("main.glade")
	if err != nil {
		panic(err)
	}
	// Get the extra window from glade
	extraWindow := builder.GetObject("extraWindow").(*gtk.Window)

	// Set up the extra window
	extraWindow.SetTitle("extra form")
	extraWindow.HideOnDelete()
	extraWindow.SetModal(true)
	extraWindow.SetKeepAbove(true)
	extraWindow.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

	// Hook up the destroy event
	extraWindow.Connect("destroy", extraWindow.Destroy)

	// Close button
	button := builder.GetObject("extraWindow_CloseButton").(*gtk.Button)
	button.Connect("clicked", extraWindow.Destroy)

	// Show the window
	extraWindow.ShowAll()
}
