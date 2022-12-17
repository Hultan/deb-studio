package gui

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
	logger2 "github.com/hultan/deb-studio/internal/logger"
)

// MainForm : Struct for the main form
type MainForm struct {
	builder       *Builder
	window        *gtk.ApplicationWindow
	logger        *logger2.Logger
	addFileDialog *addFileDialog
}

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

	// Setup pages
	m.setupMainPage()
	m.setupInstallPage()
	m.setupControlPage()

	// Show the main window
	m.window.ShowAll()

	w, err := engine.Open("/home/per/installs/softtube")
	if err != nil {
		switch {
		case errors.Is(err, engine.ErrorNewWorkspaceFolder):
			w, err = m.setupWorkspaceFolder("/home/per/installs/softtube")
			if err != nil {
				m.logger.Error.Println("failure during setup")
				m.logger.Error.Println(err)
				os.Exit(1)
			}
		default:
			m.logger.Error.Println("failure during setup")
			m.logger.Error.Println(err)
			os.Exit(1)
		}
	}

	fmt.Printf("Program %s contains %d versions:\n", w.ProgramName, len(w.Versions))
	for _, version := range w.Versions {
		architectures := ""
		if len(version.Architectures) > 0 {
			architectures = version.Architectures[0].Name
			for i := 1; i < len(version.Architectures); i++ {
				architectures += "," + version.Architectures[i].Name
			}
		}
		fmt.Printf("    %s (for architectures: %s)\n", version.Name, architectures)
	}

	// //
	// // Save
	// //
	//
	// configPath := getConfigPath()
	//
	// av := &installationConfig.InstallationConfig{}
	// av.Version = "1.0.0"
	// av.Architecture = "amd64"
	//
	// av.Control.Package = "debStudio"
	// av.Control.Source = "source"
	// av.Control.Version = "1.0.0"
	// av.Control.Section = "section"
	// av.Control.Priority = "high"
	// av.Control.Architecture = "amd64"
	// av.Control.Essential = true
	// av.Control.Depends = "dpkg"
	// av.Control.InstalledSize = "1024"
	// av.Control.Maintainer = "Per Hultqvist"
	// av.Control.Description = "A deb file creator"
	// av.Control.Homepage = "www.softteam.se"
	// av.Control.BuiltUsing = "debStudio"
	//
	// file := installationConfig.FileSection{}
	// file.FilePath = "/home/per/temp/dragon.ply"
	// file.InstallPath = "/usr/bin"
	// file.Static = false
	// file.RunScript = true
	// file.Script = "go build /home/per/code"
	//
	// av.Files = append(av.Files, file)
	//
	// err = av.Save(configPath)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// //
	// // Load
	// //
	//
	// av, err = installationConfig.Load(configPath)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", av)
	//
	// // Modify the indent level of the ConfigState only.  The global
	// // configuration is not modified.
	// scs := spew.ConfigState{
	// 	Indent:                  "\t",
	// 	DisableCapacities:       true,
	// 	DisableMethods:          true,
	// 	DisablePointerMethods:   true,
	// 	DisablePointerAddresses: true,
	// }
	// scs.Dump(av)
}

func (m *MainForm) startLogging() {
	logPath := "/home/per/.softteam/debstudio"
	logFile := "debstudio.log"
	fullLogPath := path.Join(logPath, logFile)

	// Create log path if it does not exist
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to create log path at %s\n", logPath)
			_, _ = fmt.Fprintf(os.Stderr, "Continuing without logging...\n")
			return
		}
	}

	// Create log file if it does not exist
	if _, err := os.Stat(fullLogPath); os.IsNotExist(err) {
		_, err = os.OpenFile(fullLogPath, os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to create log file at %s\n", fullLogPath)
			_, _ = fmt.Fprintf(os.Stderr, "Continuing without logging...\n")
			return
		}
	}

	logger, err := logger2.NewStandardLogger(fullLogPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create log at %s\n", logPath)
		_, _ = fmt.Fprintf(os.Stderr, "Continuing without logging...\n")
		return
	}
	m.logger = logger
}

// shutDown : shuts down the application
func (m *MainForm) shutDown() {
	if m.logger != nil {
		m.logger.Close()
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

func (m *MainForm) setupWorkspaceFolder(s string) (*engine.Workspace, error) {
	// Open setup dialog
	result, err := m.openSetupDialog()
	if err != nil {
		return nil, err
	}

	// Create program file
	w, err := engine.SetupWorkspaceFolder(s, result.name)
	if err != nil {
		return nil, err
	}

	return w, nil
}
