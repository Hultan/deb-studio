package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type infoBar struct {
	parent       *MainWindow
	infoBar      *gtk.InfoBar
	infoBarLabel *gtk.Label
}

func (m *MainWindow) setupInfoBar() *infoBar {
	i := &infoBar{parent: m}

	i.infoBar = m.builder.GetObject("mainWindow_infoBar").(*gtk.InfoBar)
	i.infoBarLabel = m.builder.GetObject("mainWindow_infoBarLabel").(*gtk.Label)

	return i
}

func (i *infoBar) update() {
	i.infoBarLabel.SetMarkup(i.getInfoBarText())
	i.setInfoBarColor()

	// Force a redraw to update the info bar
	i.parent.window.QueueDraw()
}

func (i *infoBar) getInfoBarText() string {
	switch getProjectStatus() {
	case projectStatusNoProjectOpened:
		return "You need to open or create a new project..."
	case projectStatusNoPackageSelected:
		return "You need to select or add a package to edit!"
	case projectStatusNotLatestVersion:
		return fmt.Sprintf(
			"You are currently not editing the latest version! You are editing <b>version %s</b> and <b>architecture %s</b>.",
			project.CurrentPackage.Config.Version, project.CurrentPackage.Config.Architecture,
		)
	case projectStatusLatestVersion:
		return fmt.Sprintf(
			"You are currently editing <b>version %s</b> and <b>architecture %s</b>.",
			project.CurrentPackage.Config.Version, project.CurrentPackage.Config.Architecture,
		)
	default:
		log.Error.Println("Invalid projectStatus in getInfoBarText()")
		return ""
	}
}

func (i *infoBar) setInfoBarColor() {
	switch getProjectStatus() {
	case projectStatusNoProjectOpened:
		i.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	case projectStatusNoPackageSelected:
		i.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case projectStatusNotLatestVersion:
		i.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case projectStatusLatestVersion:
		i.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	default:
		log.Error.Println("Invalid projectStatus in SetInfoBarColor()")
	}
}
