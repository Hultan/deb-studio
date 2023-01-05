package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type pageInstall struct {
	parent *MainWindow
}

// setupInstallPage: Set up the install page
func (m *MainWindow) setupInstallPage() {
	m.installPage = &pageInstall{parent: m}

	// AddFileButton
	btn := m.builder.GetObject("addFileButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.addFileButtonClicked)
}

func (m *MainWindow) addFileButtonClicked() {
	if m.addFileDialog == nil {
		m.addFileDialog = m.newAddFileDialog()
	}
	m.addFileDialog.openForNewFile("/home/per/temp/dragon.ply")
}

func (m *MainWindow) editFileButtonClicked() {
	fmt.Println("Edit file clicked!")
}

func (m *MainWindow) removeFileButtonClicked() {
	fmt.Println("Remove file clicked!")
}
