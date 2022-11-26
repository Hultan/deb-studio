package gui

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"

	builder "github.com/hultan/deb-studio/internal/gtk"
)

const applicationTitle = "Deb Studio"
const applicationVersion = "v 0.01"
const applicationCopyRight = "©SoftTeam AB, 2022"

type MainForm struct {
	window  *gtk.ApplicationWindow
	builder *builder.Builder
}

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// OpenMainForm : Opens the MainForm window
func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new softBuilder
	builder, err := builder.Create("main.glade")
	if err != nil {
		panic(err)
	}
	m.builder = builder

	// Get the main window from the glade file
	m.window = m.builder.GetObject("mainWindow").(*gtk.ApplicationWindow)

	// Set up main window
	m.window.SetApplication(app)
	m.window.SetTitle(m.getApplicationString())
	m.window.SetPosition(gtk.WIN_POS_CENTER)

	// Hook up the destroy event
	m.window.Connect("destroy", m.window.Close)

	// Toolbar Quit button
	button := m.builder.GetObject("toolbar_QuitButton").(*gtk.ToolButton)
	button.Connect("clicked", m.window.Close)

	// Status bar
	statusBar := m.builder.GetObject("mainWindow_StatusBar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("debstudio"), m.getApplicationString())

	// Open form button
	openFormButton := m.builder.GetObject("mainWindow_OpenExtraFormButton").(*gtk.Button)
	openFormButton.Connect(
		"clicked", func() {
			openExtraForm()
		},
	)

	// Open dialog button
	openDialogButton := m.builder.GetObject("mainWindow_OpenSettingsDialogButton").(*gtk.Button)
	openDialogButton.Connect(
		"clicked", func() {
			openSettingsDialog(m.window)
		},
	)

	// Menu
	m.setupMenu()

	// Show the main window
	m.window.ShowAll()
}

func (m *MainForm) getApplicationString() string {
	return fmt.Sprintf("%s %s - %s", applicationTitle, applicationVersion, applicationCopyRight)
}

func (m *MainForm) setupMenu() {
	menuQuit := m.builder.GetObject("menu_FileQuit").(*gtk.MenuItem)
	menuQuit.Connect("activate", m.window.Close)

	menuEditPreferences := m.builder.GetObject("menu_EditPreferences").(*gtk.MenuItem)
	menuEditPreferences.Connect(
		"activate", func() {
			openSettingsDialog(m.window)
		},
	)

	menuHelpAbout := m.builder.GetObject("menu_HelpAbout").(*gtk.MenuItem)
	menuHelpAbout.Connect(
		"activate", func() {
			m.openAboutDialog()
		},
	)
}
