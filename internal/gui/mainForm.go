package gui

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/builder"
	"github.com/hultan/deb-studio/internal/common"
	"github.com/hultan/deb-studio/internal/engine"
	"github.com/hultan/deb-studio/internal/logger"
)

// MainWindow : Struct for the main form
type MainWindow struct {
	builder *builder.Builder
	window  *gtk.ApplicationWindow

	// Pages
	pages *stackPages

	// dialogs
	addFileDialog *addFileDialog
}

var project *engine.Project
var log *logger.Logger

// NewMainWindow : Creates a new MainWindow object
func NewMainWindow() *MainWindow {
	return &MainWindow{}
}

// Open : Opens the MainWindow window
func (m *MainWindow) Open(app *gtk.Application) {
	m.startLogging()

	log.Trace.Println("Entering Open...")
	defer log.Trace.Println("Exiting Open...")

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
	m.setupStackPages()

	// Show the main window
	m.window.ShowAll()

	// Update gui
	m.pages.update()
}

func (m *MainWindow) startLogging() {
	// TODO : Move log file to debStudio config
	fullLogPath := path.Join(common.FolderNameLog, common.FileNameLog)

	var err error

	// Create log path if it does not exist
	if _, err = os.Stat(common.FolderNameLog); os.IsNotExist(err) {
		err = os.MkdirAll(common.FolderNameLog, 0755)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to create log path at %s\n", common.FolderNameLog)
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
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create log at %s\n", common.FolderNameLog)
		_, _ = fmt.Fprintf(os.Stderr, "Continuing without logging...\n")
		return
	}
}

func isTraceMode() bool {
	return len(os.Args) >= 2 && strings.Trim(os.Args[1], " \t") == "--trace"
}

// setupMenu: Set up the menu bar
func (m *MainWindow) setupMenu() {
	log.Trace.Println("Entering setupMenu...")
	defer log.Trace.Println("Exiting setupMenu...")

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
	log.Trace.Println("Entering setupStatusBar...")
	defer log.Trace.Println("Exiting setupStatusBar...")

	// Status bar
	statusBar := m.builder.GetObject("mainWindow_StatusBar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("debstudio"), getApplicationName())
}

// shutDown : shuts down the application
func (m *MainWindow) shutDown() {
	log.Trace.Println("Entering shutdown...")
	defer log.Trace.Println("Exiting shutdown...")

	if project != nil {
		err := project.Save()
		if err != nil {
			log.Error.Printf("failed to save project : %s\n", err)
		}
	}
	if log != nil {
		log.Close()
	}
	if m.window != nil {
		m.window.Close()
	}
}
