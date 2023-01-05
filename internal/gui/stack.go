package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *MainWindow) enableDisableStackPages() {
	status := m.getInfoBarStatus()
	haveOpenedProject := status > infoBarStatusNoProjectOpened
	haveChosenPackage := status > infoBarStatusNoPackageSelected
	m.enableDisableStackPage("mainWindow_packagePage", haveOpenedProject)
	m.enableDisableStackPage("mainWindow_controlPage", haveChosenPackage)
	m.enableDisableStackPage("mainWindow_scriptPage", haveChosenPackage)
	m.enableDisableStackPage("mainWindow_installPage", haveChosenPackage)
	m.enableDisableStackPage("mainWindow_textPage", haveChosenPackage)
}

func (m *MainWindow) enableDisableStackPage(name string, status bool) {
	w := m.builder.GetObject(name)
	switch item := w.(type) {
	case *gtk.Box:
		item.SetSensitive(status)
	case *gtk.Grid:
		item.SetSensitive(status)
	}
}
