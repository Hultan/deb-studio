package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

type stackPages struct {
	parent *MainWindow
}

func (m *MainWindow) setupStackPages() {
	m.stackPages = &stackPages{parent: m}
}

func (s *stackPages) enableDisable() {
	status := s.parent.packagePage.getInfoBarStatus()
	haveOpenedProject := status > infoBarStatusNoProjectOpened
	haveChosenPackage := status > infoBarStatusNoPackageSelected
	s.enableDisableStackPage("mainWindow_packagePage", haveOpenedProject)
	s.enableDisableStackPage("mainWindow_controlPage", haveChosenPackage)
	s.enableDisableStackPage("mainWindow_scriptPage", haveChosenPackage)
	s.enableDisableStackPage("mainWindow_installPage", haveChosenPackage)
	s.enableDisableStackPage("mainWindow_textPage", haveChosenPackage)
}

func (s *stackPages) enableDisableStackPage(name string, status bool) {
	w := s.parent.builder.GetObject(name)
	switch item := w.(type) {
	case *gtk.Box:
		item.SetSensitive(status)
	case *gtk.Grid:
		item.SetSensitive(status)
	}
}
