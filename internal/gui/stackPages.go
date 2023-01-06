package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

type stackPages struct {
	parent *MainWindow

	stack *gtk.Stack

	projectPage *pageProject
	packagePage *pagePackage
	controlPage *pageControl
	scriptPage  *pageScript
	installPage *pageInstall
	textPage    *pageText

	infoBar *infoBar
}

func (m *MainWindow) setupStackPages() {
	p := &stackPages{parent: m}
	m.pages = p

	p.projectPage = m.setupProjectPage()
	p.packagePage = m.setupPackagePage()
	p.controlPage = m.setupControlPage()
	p.scriptPage = m.setupScriptPage()
	p.installPage = m.setupInstallPage()
	p.textPage = m.setupTextPage()

	p.infoBar = m.setupInfoBar()
}

func (s *stackPages) update() {
	s.enableDisablePages()
	s.infoBar.update()
	s.projectPage.update()

	if project == nil {
		return
	}

	switch getProjectStatus() {
	case projectStatusNoProjectOpened:
	case projectStatusNoPackageSelected:
		s.projectPage.update()
	default:
		s.packagePage.update()
		s.controlPage.update()
		s.scriptPage.update()
		s.installPage.update()
		s.textPage.update()
	}
}

func (s *stackPages) enableDisablePages() {
	status := getProjectStatus()
	haveOpenedProject := status > projectStatusNoProjectOpened
	haveChosenPackage := status > projectStatusNoPackageSelected
	s.enableDisablePage("mainWindow_packagePage", haveOpenedProject)
	s.enableDisablePage("mainWindow_controlPage", haveChosenPackage)
	s.enableDisablePage("mainWindow_scriptPage", haveChosenPackage)
	s.enableDisablePage("mainWindow_installPage", haveChosenPackage)
	s.enableDisablePage("mainWindow_textPage", haveChosenPackage)
}

func (s *stackPages) enableDisablePage(name string, status bool) {
	w := s.parent.builder.GetObject(name)
	switch item := w.(type) {
	case *gtk.Box:
		item.SetSensitive(status)
	case *gtk.Grid:
		item.SetSensitive(status)
	}
}
