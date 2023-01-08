package gui

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
	"github.com/hultan/deb-studio/internal/iconFactory"
)

// openAboutDialog: Opens the help/about dialog
func (m *MainWindow) openAboutDialog() {
	about := m.builder.GetObject("aboutDialog").(*gtk.AboutDialog)
	about.SetDestroyWithParent(true)
	about.SetTransientFor(m.window)
	about.SetProgramName(common.ApplicationTitle)
	about.SetComments("Create deb packages for debian based linux distributions.")
	about.SetVersion(common.ApplicationVersion)
	about.SetCopyright(common.ApplicationCopyRight)
	image := m.image.GetPixBuf(iconFactory.ImageDebStudio)
	about.SetLogo(image)
	about.SetModal(true)
	about.SetPosition(gtk.WIN_POS_CENTER)

	responseId := about.Run()

	if responseId == gtk.RESPONSE_DELETE_EVENT {
		about.Destroy()
	}
}
