package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) getInfoBarStatus() infoBarStatus {
	if currentVersion == nil || currentArchitecture == nil {
		return infoBarStatusNoPackageSelected
	} else if !currentProject.IsLatestVersion(currentVersion) {
		return infoBarStatusNotLatestVersion
	}
	return infoBarStatusLatestVersion
}

func (m *MainForm) getInfoBarText() string {
	switch m.getInfoBarStatus() {
	case infoBarStatusNoPackageSelected:
		return "You need to select or add a package to edit!"
	case infoBarStatusNotLatestVersion:
		return fmt.Sprintf(
			"You are currently not editing the latest version! You are editing version %s and architecture %s.",
			currentVersion.Name, currentArchitecture.Name,
		)
	case infoBarStatusLatestVersion:
		return fmt.Sprintf(
			"You are currently editing version %s and architecture %s.",
			currentVersion.Name, currentArchitecture.Name,
		)
	default:
		log.Error.Println("Invalid infoBarStatus in getInfoBarText()")
		return ""
	}
}

func (m *MainForm) setInfoBarColor() {
	switch m.getInfoBarStatus() {
	case infoBarStatusNoPackageSelected:
		m.infoBar.SetMessageType(gtk.MESSAGE_ERROR)
	case infoBarStatusNotLatestVersion:
		m.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case infoBarStatusLatestVersion:
		m.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	default:
		log.Error.Println("Invalid infoBarStatus in SetInfoBarColor()")
	}
}

func (m *MainForm) updateInfoBar() {
	m.infoBarLabel.SetText(m.getInfoBarText())
	m.setInfoBarColor()
	m.window.QueueDraw()
}
