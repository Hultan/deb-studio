package gui

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/config"
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

	// Menu
	m.setupMenu()

	// Show the main window
	m.window.ShowAll()

	path := getConfigPath()

	c, err := config.Load(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", c)

	// Modify the indent level of the ConfigState only.  The global
	// configuration is not modified.
	scs := spew.ConfigState{
		Indent:                  "\t",
		DisableCapacities:       true,
		DisableMethods:          true,
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
	}
	scs.Dump(c)

	// c := &config.Config{}
	// c.Control.Package = "debStudio"
	// c.Control.Source = "source"
	// c.Control.Version = "1.0.0"
	// c.Control.Section = "section"
	// c.Control.Priority = "high"
	// c.Control.Architecture = "amd64"
	// c.Control.Essential = true
	// c.Control.Depends = "dpkg"
	// c.Control.InstalledSize = "1024"
	// c.Control.Maintainer = "Per Hultqvist"
	// c.Control.Description = "A deb file creator"
	// c.Control.Homepage = "www.softteam.se"
	// c.Control.BuiltUsing = "debStudio"
	// err = c.Save(path)
	// if err != nil {
	// 	fmt.Println(err)
	// }
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

// Get path to the config file
// Mode = "test" returns test config path
// otherwise returns normal config path
func getConfigPath() string {
	home := getHomeDirectory()

	return path.Join(home, "deb-studio/test.json")
}

// Get current users home directory
func getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
