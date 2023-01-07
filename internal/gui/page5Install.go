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
	log.Trace.Println("Entering setupInstallPage...")
	defer log.Trace.Println("Exiting setupInstallPage...")

	p := &pageInstall{parent: m}

	// AddFileButton
	btn := m.builder.GetObject("addFileButton").(*gtk.ToolButton)
	btn.Connect("clicked", p.addFileButtonClicked)

	return p
}

func (p *pageInstall) update() {
	log.Trace.Println("Entering update...")
	defer log.Trace.Println("Exiting update...")

}

func (p *pageInstall) addFileButtonClicked() {
	log.Trace.Println("Entering addFileButtonClicked...")
	defer log.Trace.Println("Exiting addFileButtonClicked...")

	if p.parent.addFileDialog == nil {
		p.parent.addFileDialog = p.parent.newAddFileDialog()
	}
	p.parent.addFileDialog.openForNewFile("/home/per/temp/dragon.ply")
}

func (p *pageInstall) editFileButtonClicked() {
	log.Trace.Println("Entering editFileButtonClicked...")
	defer log.Trace.Println("Exiting editFileButtonClicked...")

	fmt.Println("Edit file clicked!")
}

func (p *pageInstall) removeFileButtonClicked() {
	log.Trace.Println("Entering removeFileButtonClicked...")
	defer log.Trace.Println("Exiting removeFileButtonClicked...")

	fmt.Println("Remove file clicked!")
}
