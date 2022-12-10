package gui

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/applicationVersionConfig"
)

// MainForm : Struct for the main form
type MainForm struct {
	builder       *Builder
	window        *gtk.ApplicationWindow
	addFileDialog *addFileDialog
}

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// Open : Opens the MainForm window
func (m *MainForm) Open(app *gtk.Application) {
	// Initialize gtk and create a builder
	gtk.Init(&os.Args)
	m.builder = newBuilder()

	// Main window
	m.window = m.builder.GetObject("mainWindow").(*gtk.ApplicationWindow)
	m.window.SetIcon(createPixBufFromBytes(applicationIcon, "application"))
	m.window.SetApplication(app)
	m.window.SetTitle(getApplicationName())
	m.window.SetPosition(gtk.WIN_POS_CENTER)
	m.window.Connect("destroy", m.window.Close)

	// Toolbar & menu
	m.setupStatusBar()
	m.setupToolbar()
	m.setupMenu()

	// Setup pages
	m.setupInstallPage()
	m.setupControlPage()

	// Show the main window
	m.window.ShowAll()

	//
	// Save
	//

	configPath := getConfigPath()

	av := &applicationVersionConfig.ApplicationVersion{}
	av.Version = "1.0.0"
	av.Architecture="amd64"

	av.Control.Package = "debStudio"
	av.Control.Source = "source"
	av.Control.Version = "1.0.0"
	av.Control.Section = "section"
	av.Control.Priority = "high"
	av.Control.Architecture = "amd64"
	av.Control.Essential = true
	av.Control.Depends = "dpkg"
	av.Control.InstalledSize = "1024"
	av.Control.Maintainer = "Per Hultqvist"
	av.Control.Description = "A deb file creator"
	av.Control.Homepage = "www.softteam.se"
	av.Control.BuiltUsing = "debStudio"

	file := applicationVersionConfig.File{}
	file.FilePath="/home/per/temp/dragon.ply"
	file.InstallPath="/usr/bin"
	file.Static=false
	file.RunScript=true
	file.Script="go build /home/per/code"

	av.Files = append(av.Files, file)

	err := av.Save(configPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCodeSetupError)
	}

	//
	// Load
	//

	av, err = applicationVersionConfig.Load(configPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCodeSetupError)
	}
	fmt.Printf("%+v\n", av)

	// Modify the indent level of the ConfigState only.  The global
	// configuration is not modified.
	scs := spew.ConfigState{
		Indent:                  "\t",
		DisableCapacities:       true,
		DisableMethods:          true,
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
	}
	scs.Dump(av)
}

// setupMenu: Set up the menu bar
func (m *MainForm) setupMenu() {
	menuQuit := m.builder.GetObject("menu_FileQuit").(*gtk.MenuItem)
	menuQuit.Connect("activate", m.window.Close)

	menuEditPreferences := m.builder.GetObject("menu_EditPreferences").(*gtk.MenuItem)
	menuEditPreferences.Connect(
		"activate", func() {
			m.openSettingsDialog()
		},
	)

	menuHelpAbout := m.builder.GetObject("menu_HelpAbout").(*gtk.MenuItem)
	menuHelpAbout.Connect(
		"activate", func() {
			m.openAboutDialog()
		},
	)
}

// setupToolbar: Set up the toolbar
func (m *MainForm) setupToolbar() {
	// Toolbar quit button
	btn := m.builder.GetObject("toolbar_quitButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.window.Close)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(exitIcon, "quit"))

	// Toolbar save button
	btn = m.builder.GetObject("toolbar_saveButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.save)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(saveIcon, "save"))
}

// setupStatusBar: Set up the status bar
func (m *MainForm) setupStatusBar() {
	// Status bar
	statusBar := m.builder.GetObject("mainWindow_StatusBar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("debstudio"), getApplicationName())
}

// addFile: Handler for the add file button clicked signal
func (m *MainForm) addFile() {
	if m.addFileDialog == nil {
		m.addFileDialog = m.newAddFileDialog()
	}
	m.addFileDialog.openForNewFile("/home/per/temp/dragon.ply")
}

// save: Handler for the save button clicked signal
func (m *MainForm) save() {
	// TODO : save here
}
