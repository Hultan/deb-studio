package gui

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
	"github.com/hultan/deb-studio/internal/logger"
)

// MainForm : Struct for the main form
type MainForm struct {
	builder       *Builder
	window        *gtk.ApplicationWindow
	addFileDialog *addFileDialog
	listBox       *gtk.ListBox
	infoBar       *gtk.InfoBar
	infoBarLabel  *gtk.Label
	popup         *gtk.Menu
}

var currentProject *engine.Project
var currentVersion *engine.Version
var currentArchitecture *engine.Architecture
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

	projectFolder := "/home/per/installs/softtube"
	projectName := "test"

	var err error
	e := engine.NewEngine(log)
	if e.IsProjectFolder(projectFolder) {
		currentProject, err = e.OpenProject(projectFolder)
		if err != nil {
			log.Error.Printf("failure during opening of '%s': %s", projectFolder, err)
			os.Exit(1)
		}
	} else {
		currentProject, err = e.SetupProject(projectFolder, projectName)
		if err != nil {
			log.Error.Printf("failure during setup of '%s': %s", projectFolder, err)
			os.Exit(1)
		}
	}

	m.listPackages()
	m.printTraceInfo()
	m.updateInfoBar()
	m.enableDisableStackPages(m.getInfoBarStatus() != infoBarStatusNoPackageSelected)

	// v, err := currentProject.AddVersion("testVersion1.0.0")
	// if err != nil {
	// 	log.Error.Printf("failed to add version")
	// 	os.Exit(1)
	// }
	// a, err := v.AddArchitecture("amd75")
	// if err != nil {
	// 	log.Error.Printf("failed to add architecture")
	// 	os.Exit(1)
	// }
	// fmt.Println(a.Name)
	//
	// err = a.AddFile("/home/per/temp/", "empty", "/usr/bin/", false)
	// if err != nil {
	// 	log.Error.Printf("failed to add file")
	// 	os.Exit(1)
	// }
	// fmt.Println("file added successfully!")

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

func (m *MainForm) printTraceInfo() {
	log.Info.Printf(
		"Project %s (path: %s) contains %d versions:\n",
		currentProject.Name,
		currentProject.Path,
		len(currentProject.Versions),
	)

	for _, version := range currentProject.Versions {
		log.Info.Printf("\tVersion: %s (path: %s)\n", version.Name, version.Path)
		if len(version.Architectures) > 0 {
			for _, architecture := range version.Architectures {
				log.Info.Printf("\t\tArchitecture: %s (path: %s)\n", architecture.Name, architecture.Path)
			}
		}
	}
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

	// Toolbar add file button
	btn = m.builder.GetObject("toolbar_addFileButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.addFileButtonClicked)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(addFileIcon, "addFile"))

	// Toolbar edit file button
	btn = m.builder.GetObject("toolbar_editFileButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.editFileButtonClicked)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(editFileIcon, "editFile"))

	// Toolbar remove file button
	btn = m.builder.GetObject("toolbar_removeFileButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.removeFileButtonClicked)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(removeFileIcon, "removeFile"))
}

// setupStatusBar: Set up the status bar
func (m *MainForm) setupStatusBar() {
	// Status bar
	statusBar := m.builder.GetObject("mainWindow_StatusBar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("debstudio"), getApplicationName())
}
