package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) enableDisableStackPages(status bool) {
	m.enableDisableStackPage("mainWindow_controlPage", status)
	m.enableDisableStackPage("mainWindow_preinstallPage", status)
	m.enableDisableStackPage("mainWindow_installPage", status)
	m.enableDisableStackPage("mainWindow_postinstallPage", status)
	m.enableDisableStackPage("mainWindow_copyrightPage", status)
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
