package gui

import (
	"fmt"
	"os"
	"path"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/builder"
	"github.com/hultan/deb-studio/internal/engine"
	"github.com/hultan/deb-studio/internal/logger"
)

// MainWindow : Struct for the main form
type MainWindow struct {
	builder       *builder.Builder
	window        *gtk.ApplicationWindow
	projectPage   *pageProject
	packagePage   *pagePackage
	controlPage   *pageControl
	addFileDialog *addFileDialog
}

var project *engine.Project
var log *logger.Logger

// NewMainWindow : Creates a new MainWindow object
func NewMainWindow() *MainWindow {
	mainForm := new(MainWindow)
	return mainForm
}

// Open : Opens the MainWindow window
func (m *MainWindow) Open(app *gtk.Application) {
	// TODO : Move log path to config
	m.startLogging()

	// Initialize gtk and create a builder
	gtk.Init(&os.Args)

	m.builder = builder.NewBuilder(mainGlade)

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

	// Setup pages
	m.setupProjectPage()
	m.setupPackagePage()
	m.setupInstallPage()
	m.setupControlPage()

	// Show the main window
	m.window.ShowAll()

	// Disable pages until a project has been opened
	m.enableDisableStackPages()

	// Update info bar
	m.packagePage.updateInfoBar()
}

func (m *MainWindow) startLogging() {
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

// setupMenu: Set up the menu bar
func (m *MainWindow) setupMenu() {
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

// setupStatusBar: Set up the status bar
func (m *MainWindow) setupStatusBar() {
	// Status bar
	statusBar := m.builder.GetObject("mainWindow_StatusBar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("debstudio"), getApplicationName())
}

// shutDown : shuts down the application
func (m *MainWindow) shutDown() {
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
