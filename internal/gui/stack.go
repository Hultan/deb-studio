package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) enableDisableStackPages() {
	status := m.getInfoBarStatus()
	haveOpenedProject := status > infoBarStatusNoProjectOpened
	haveChosenPackage := status > infoBarStatusNoPackageSelected
	m.enableDisableStackPage("mainWindow_packagePage", haveOpenedProject)
	m.enableDisableStackPage("mainWindow_controlPage", haveChosenPackage)
	m.enableDisableStackPage("mainWindow_preinstallPage", haveChosenPackage)
	m.enableDisableStackPage("mainWindow_installPage", haveChosenPackage)
	m.enableDisableStackPage("mainWindow_postinstallPage", haveChosenPackage)
	m.enableDisableStackPage("mainWindow_copyrightPage", haveChosenPackage)
}

func (m *MainForm) enableDisableStackPage(name string, status bool) {
	// TODO : Fix this code, should be doable with *gtk.Widget
	box, ok := m.builder.GetObject(name).(*gtk.Box)
	if !ok {
		grid, ok := m.builder.GetObject(name).(*gtk.Grid)
		if !ok {
			log.Error.Printf("failed to retrieve stack page: %s", name)
		}
		grid.SetSensitive(status)
		return
	}
	box.SetSensitive(status)
}
