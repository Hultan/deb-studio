package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

// setupInstallPage: Set up the install page
func (m *MainForm) setupInstallPage() {
	// AddFileButton
	btn := m.builder.GetObject("addFileButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.addFileButtonClicked)
}
