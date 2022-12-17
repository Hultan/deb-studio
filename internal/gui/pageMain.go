package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) setupMainPage() {
	// Add installation Button
	btn := m.builder.GetObject("miscPage_addInstallationButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.addInstallation)
}

func (m *MainForm) addInstallation() {
	// dialog := m.builder.GetObject("addInstallationDialog").(*gtk.Dialog)
	// versionName := m.builder.GetObject("addInstallationDialog_versionNameEntry").(*gtk.Entry)
	// architectureCombo := m.builder.GetObject("addInstallationDialog_architectureCombo").(*gtk.Dialog)
	// dialog.AddButton("Add", gtk.RESPONSE_ACCEPT)
	// dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	//
	// architectureCombo.set
	//
	// // Show the dialog
	// responseId := dialog.Run()
	// if responseId == gtk.RESPONSE_ACCEPT {
	// 	// Save installation
	//
	// }
	//
	// dialog.Hide()

}
