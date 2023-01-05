package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type pageInstall struct {
	parent *MainWindow
}

// setupInstallPage: Set up the install page
func (m *MainWindow) setupInstallPage() *pageInstall {
	p := &pageInstall{parent: m}

	// AddFileButton
	btn := m.builder.GetObject("addFileButton").(*gtk.ToolButton)
	btn.Connect("clicked", p.addFileButtonClicked)

	return p
}

func (p *pageInstall) addFileButtonClicked() {
	if p.parent.addFileDialog == nil {
		p.parent.addFileDialog = p.parent.newAddFileDialog()
	}
	p.parent.addFileDialog.openForNewFile("/home/per/temp/dragon.ply")
}

func (p *pageInstall) editFileButtonClicked() {
	fmt.Println("Edit file clicked!")
}

func (p *pageInstall) removeFileButtonClicked() {
	fmt.Println("Remove file clicked!")
}

func (p *pageInstall) update() {

}
