package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

var locationEntry *gtk.Entry

type setup struct {
	path string
	name string
}

// openSetupDialog: Open the setup dialog
func (m *MainForm) openSetupDialog() (*setup, error) {
	var result *setup

	// Get the dialog window from glade
	dialog := m.builder.GetObject("setupDialog").(*gtk.Dialog)
	locationEntry = m.builder.GetObject("setupDialog_projectLocationEntry").(*gtk.Entry)
	browseButton := m.builder.GetObject("setupDialog_browseProjectLocation").(*gtk.Button)
	browseButton.Connect("clicked", m.browseLocationButtonClicked)
	nameEntry := m.builder.GetObject("setupDialog_projectNameEntry").(*gtk.Entry)

	dialog.SetTitle("Setup dialog")
	dialog.SetTransientFor(m.window)
	dialog.SetModal(true)
	dialog.AddButton("SaveControlFile", gtk.RESPONSE_ACCEPT)
	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)

	// Show the dialog
	responseId := dialog.Run()
	defer dialog.Hide()

	switch responseId {
	case gtk.RESPONSE_ACCEPT:
		// SaveControlFile setup information
		p, err := locationEntry.GetText()
		if err != nil {
			return nil, err
		}
		n, err := nameEntry.GetText()
		if err != nil {
			return nil, err
		}
		result = &setup{name: n, path: p}
		return result, nil
	default:
		return nil, nil
	}
}

func (m *MainForm) browseLocationButtonClicked() {
	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Select folder...", m.window, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		"OK", gtk.RESPONSE_ACCEPT,
		"Cancel", gtk.RESPONSE_CANCEL,
	)
	if err != nil {
		showErrorDialog("Failed to open fileChooserDialog!", err)
		log.Error.Printf("failed to create fileChooserDialog: %s\n", err)
		return
	}
	dlg.SetCurrentFolder("/home/per/installs")

	response := dlg.Run()
	dlg.Hide()
	if response == gtk.RESPONSE_ACCEPT {
		dir, err := dlg.GetCurrentFolder()
		if err != nil {

		}
		locationEntry.SetText(dir)
	}
}
