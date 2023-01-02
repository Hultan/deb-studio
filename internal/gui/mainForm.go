package gui

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
	"github.com/hultan/deb-studio/internal/logger"
	"github.com/hultan/deb-studio/internal/projectList"
)

// MainForm : Struct for the main form
type MainForm struct {
	builder          *Builder
	window           *gtk.ApplicationWindow
	addFileDialog    *addFileDialog
	treeView         *gtk.TreeView
	projectList      *projectList.ProjectList
	infoBar          *gtk.InfoBar
	infoBarLabel     *gtk.Label
	popup            *gtk.Menu
	showOnlyCheckBox *gtk.CheckButton
}

var project *engine.Project
var log *logger.Logger

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// Open : Opens the MainForm window
func (m *MainForm) Open(app *gtk.Application) {
	// TODO : Move log path to config
	m.startLogging()

	// Initialize gtk and create a builder
	gtk.Init(&os.Args)

	m.builder = newBuilder()

	// Main window
	m.window = m.builder.GetObject("mainWindow").(*gtk.ApplicationWindow)
	m.window.SetIcon(createPixBufFromBytes(applicationIcon, "application"))
	m.window.SetApplication(app)
	m.window.SetTitle(getApplicationName())
	m.window.SetPosition(gtk.WIN_POS_CENTER)
	m.window.Connect("destroy", m.shutDown)

	// Toolbar & menu
	m.setupStatusBar()
	m.setupToolbar()
	m.setupMenu()
	m.setupPopupMenu()

	// Setup pages
	m.setupPackagePage()
	m.setupInstallPage()
	m.setupControlPage()

	// Show the main window
	m.window.ShowAll()

	// Disable pages until a project has been opened
	m.enableDisableStackPages()

	// Update info bar
	m.updateInfoBar()
}

func (m *MainForm) startLogging() {
	logPath := "/home/per/.softteam/debstudio"
	logFile := "debstudio.log"
	fullLogPath := path.Join(logPath, logFile)

	var err error

	// Create log path if it does not exist
	if _, err = os.Stat(logPath); os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to create log path at %s\n", logPath)
			_, _ = fmt.Fprintf(os.Stderr, "Continuing without logging...\n")
			return
		}
	}

	// Create log file if it does not exist
	if _, err = os.Stat(fullLogPath); os.IsNotExist(err) {
		_, err = os.OpenFile(fullLogPath, os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to create log file at %s\n", fullLogPath)
			_, _ = fmt.Fprintf(os.Stderr, "Continuing without logging...\n")
			return
		}
	}

	if isTraceMode() {
		log, err = logger.NewDebugLogger(fullLogPath)
	} else {
		log, err = logger.NewStandardLogger(fullLogPath)
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create log at %s\n", logPath)
		_, _ = fmt.Fprintf(os.Stderr, "Continuing without logging...\n")
		return
	}
}

func isTraceMode() bool {
	return len(os.Args) >= 2 && strings.HasPrefix(os.Args[1], "-t")
}

// shutDown : shuts down the application
func (m *MainForm) shutDown() {
	if project != nil {
		project.Save()
	}
	if log != nil {
		log.Close()
	}
	if m.window != nil {
		m.window.Close()
	}
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
	// Toolbar new button
	btn := m.builder.GetObject("toolbar_newButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.newButtonClicked)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(newIcon, "new"))

	// Toolbar open button
	btn = m.builder.GetObject("toolbar_openButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.openButtonClicked)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(openIcon, "open"))

	// Toolbar save button
	btn = m.builder.GetObject("toolbar_saveButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.saveButtonClicked)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(saveIcon, "saveButtonClicked"))

	// Toolbar build button
	btn = m.builder.GetObject("toolbar_buildButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.buildButtonClicked)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(buildIcon, "build"))

	// Toolbar quit button
	btn = m.builder.GetObject("toolbar_quitButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.window.Close)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(exitIcon, "quit"))
}

// setupStatusBar: Set up the status bar
func (m *MainForm) setupStatusBar() {
	// Status bar
	statusBar := m.builder.GetObject("mainWindow_StatusBar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("debstudio"), getApplicationName())
}
