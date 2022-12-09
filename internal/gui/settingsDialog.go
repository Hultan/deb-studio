package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) openSettingsDialog() {
	// Get the dialog window from glade
	dialog := m.builder.GetObject("settingsDialog").(*gtk.Dialog)

	dialog.SetTitle("Settings dialog")
	dialog.SetTransientFor(m.window)
	dialog.SetModal(true)

	// Show the dialog
	responseId := dialog.Run()
	if responseId == gtk.RESPONSE_ACCEPT {
		// Save settings
	}

	dialog.Hide()
	dialog = nil
}
