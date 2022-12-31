package gui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) setupPopupMenu() {
	m.popup = m.builder.GetObject("mainWindow_popupPackageMenu").(*gtk.Menu)
	menuItem := m.builder.GetObject("mainWindow_popupAddPackage").(*gtk.MenuItem)
	menuItem.Connect("activate", m.addPackageClicked)
	menuItem = m.builder.GetObject("mainWindow_popupRemovePackage").(*gtk.MenuItem)
	menuItem.Connect("activate", m.removePackageClicked)
	menuItem = m.builder.GetObject("mainWindow_popupSetAsLatest").(*gtk.MenuItem)
	menuItem.Connect("activate", m.setAsLatestVersionClicked)
	menuItem = m.builder.GetObject("mainWindow_popupSetAsCurrent").(*gtk.MenuItem)
	menuItem.Connect("activate", m.setAsCurrentPackageClicked)
}

func (m *MainForm) showPopupMenu(_ *gtk.ListBox, e *gdk.Event) {
	ev := gdk.EventButtonNewFromEvent(e)
	if ev.Button() == RightMouseButton {
		m.popup.PopupAtPointer(e)
	}
}
