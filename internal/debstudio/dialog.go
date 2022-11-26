package debstudio

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/gtkBuilder"
)

func openSettingsDialog(parent gtk.IWindow) {
	// Create a new softBuilder
	builder, err := gtkBuilder.Create("main.glade")
	if err != nil {
		panic(err)
	}

	// Get the dialog window from glade
	dialog := builder.GetObject("settingsDialog").(*gtk.Dialog)

	dialog.SetTitle("Settings dialog")
	dialog.SetTransientFor(parent)
	dialog.SetModal(true)

	// Show the dialog
	responseId := dialog.Run()
	if responseId == gtk.RESPONSE_ACCEPT {
		// Save settings
	}

	dialog.Destroy()
}
