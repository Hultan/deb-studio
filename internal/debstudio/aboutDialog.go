package debstudio

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/gtkBuilder"
)

func (m *MainForm) openAboutDialog() {
	about := m.builder.GetObject("aboutDialog").(*gtk.AboutDialog)
	about.SetDestroyWithParent(true)
	about.SetTransientFor(m.window)
	about.SetProgramName(applicationTitle)
	about.SetComments("Create deb packages for debian based linux distributions.")
	about.SetVersion(applicationVersion)
	about.SetCopyright(applicationCopyRight)
	image, err := gdk.PixbufNewFromFile(gtkBuilder.GetResourcePath("debstudio_256.png"))
	if err == nil {
		about.SetLogo(image)
	}
	about.SetModal(true)
	about.SetPosition(gtk.WIN_POS_CENTER)

	responseId := about.Run()

	if responseId == gtk.RESPONSE_DELETE_EVENT {
		about.Destroy()
	}
}
