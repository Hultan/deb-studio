package gui

import (
	"os/exec"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
)

func (m *MainWindow) setupPopupMenu() {
	m.popup = m.builder.GetObject("mainWindow_popupPackageMenu").(*gtk.Menu)
	menuItem := m.builder.GetObject("mainWindow_popupAddPackage").(*gtk.MenuItem)
	menuItem.Connect("activate", m.addPackageClicked)
	menuItem = m.builder.GetObject("mainWindow_popupRemovePackage").(*gtk.MenuItem)
	menuItem.Connect("activate", m.removePackageClicked)
	menuItem = m.builder.GetObject("mainWindow_popupSetAsLatest").(*gtk.MenuItem)
	menuItem.Connect("activate", m.setAsLatestVersionClicked)
	menuItem = m.builder.GetObject("mainWindow_popupSetAsCurrent").(*gtk.MenuItem)
	menuItem.Connect("activate", m.setPackageAsCurrentClicked)
	menuItem = m.builder.GetObject("mainWindow_popupOpenProject").(*gtk.MenuItem)
	menuItem.Connect("activate", m.openProjectFolder)
	menuItem = m.builder.GetObject("mainWindow_popupOpenPackage").(*gtk.MenuItem)
	menuItem.Connect("activate", m.openPackageFolder)
}

func (m *MainWindow) showPopupMenu(_ *gtk.ListBox, e *gdk.Event) {
	ev := gdk.EventButtonNewFromEvent(e)
	if ev.Button() == common.RightMouseButton {
		m.popup.PopupAtPointer(e)
	}
}

// openProjectFolder: Handler for the open project folder button clicked signal
func (m *MainWindow) openProjectFolder() {
	cmd := exec.Command("xdg-open", project.Path)
	cmd.Run()
}

// openPackageFolder: Handler for the open package folder button clicked signal
func (m *MainWindow) openPackageFolder() {
	// Set version as latest
	pkgName := m.projectList.GetSelectedPackageName()
	if pkgName == "" {
		return
	}
	pkg := project.GetPackageByName(pkgName)
	if pkg == nil {
		return
	}
	cmd := exec.Command("xdg-open", pkg.Path)
	cmd.Run()
}
