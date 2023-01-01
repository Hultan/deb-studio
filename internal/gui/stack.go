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
	w := m.builder.GetObject(name)
	switch item := w.(type) {
	case *gtk.Box:
		item.SetSensitive(status)
	case *gtk.Grid:
		item.SetSensitive(status)
	}
}
