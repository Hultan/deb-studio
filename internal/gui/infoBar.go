package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) getInfoBarStatus() infoBarStatus {
	if project == nil {
		return infoBarStatusNoProjectOpened
	} else if project.CurrentPackage == nil {
		return infoBarStatusNoPackageSelected
	} else if !project.WorkingWithLatestVersion() {
		return infoBarStatusNotLatestVersion
	}
	return infoBarStatusLatestVersion
}

func (m *MainForm) getInfoBarText() string {
	switch m.getInfoBarStatus() {
	case infoBarStatusNoProjectOpened:
		return "You need to open or create a new project..."
	case infoBarStatusNoPackageSelected:
		return "You need to select or add a package to edit!"
	case infoBarStatusNotLatestVersion:
		return fmt.Sprintf(
			"You are currently not editing the latest version! You are editing <b>version %s</b> and <b>architecture %s</b>.",
			project.CurrentPackage.Config.Version, project.CurrentPackage.Config.Architecture,
		)
	case infoBarStatusLatestVersion:
		return fmt.Sprintf(
			"You are currently editing <b>version %s</b> and <b>architecture %s</b>.",
			project.CurrentPackage.Config.Version, project.CurrentPackage.Config.Architecture,
		)
	default:
		log.Error.Println("Invalid infoBarStatus in getInfoBarText()")
		return ""
	}
}

func (m *MainForm) setInfoBarColor() {
	switch m.getInfoBarStatus() {
	case infoBarStatusNoProjectOpened:
		m.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	case infoBarStatusNoPackageSelected:
		m.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case infoBarStatusNotLatestVersion:
		m.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case infoBarStatusLatestVersion:
		m.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	default:
		log.Error.Println("Invalid infoBarStatus in SetInfoBarColor()")
	}
}

func (m *MainForm) updateInfoBar() {
	m.infoBarLabel.SetMarkup(m.getInfoBarText())
	m.setInfoBarColor()
	m.window.QueueDraw()
}
