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
	log           *logger.Logger
	addFileDialog *addFileDialog
	listBox       *gtk.ListBox
	infoBar       *gtk.InfoBar
	infoBarLabel  *gtk.Label
	popup         *gtk.Menu
}

var currentProject *engine.Project
var currentVersion *engine.Version
var currentArchitecture *engine.Architecture

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
	m.setupMainPage()
	m.setupInstallPage()
	m.setupControlPage()

	// Show the main window
	m.window.ShowAll()

	projectFolder := "/home/per/installs/softtube"
	projectName := "test"

	var err error
	e := engine.NewEngine(m.log)
	if e.IsProjectFolder(projectFolder) {
		currentProject, err = e.OpenProject(projectFolder)
		if err != nil {
			m.log.Error.Printf("failure during opening of '%s': %s", projectFolder, err)
			os.Exit(1)
		}
	} else {
		currentProject, err = e.SetupProject(projectFolder, projectName)
		if err != nil {
			m.log.Error.Printf("failure during setup of '%s': %s", projectFolder, err)
			os.Exit(1)
		}
	}

	m.listPackages()
	m.printTraceInfo()
	m.updateInfoBar()
	m.enableDisableStackPages(m.getInfoBarStatus() != infoBarStatusNoPackageSelected)

	// v, err := currentProject.AddVersion("testVersion1.0.0")
	// if err != nil {
	// 	m.log.Error.Printf("failed to add version")
	// 	os.Exit(1)
	// }
	// a, err := v.AddArchitecture("amd75")
	// if err != nil {
	// 	m.log.Error.Printf("failed to add architecture")
	// 	os.Exit(1)
	// }
	// fmt.Println(a.Name)
	//
	// err = a.AddFile("/home/per/temp/", "empty", "/usr/bin/", false)
	// if err != nil {
	// 	m.log.Error.Printf("failed to add file")
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
	m.log.Info.Printf(
		"Project %s (path: %s) contains %d versions:\n",
		currentProject.Name,
		currentProject.Path,
		len(currentProject.Versions),
	)

	for _, version := range currentProject.Versions {
		m.log.Info.Printf("\tVersion: %s (path: %s)\n", version.Name, version.Path)
		if len(version.Architectures) > 0 {
			for _, architecture := range version.Architectures {
				m.log.Info.Printf("\t\tArchitecture: %s (path: %s)\n", architecture.Name, architecture.Path)
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

	var log *logger.Logger
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
	m.log = log
}

// shutDown : shuts down the application
func (m *MainForm) shutDown() {
	if m.log != nil {
		m.log.Close()
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
	btn.Connect("clicked", m.new)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(newIcon, "new"))

	// Toolbar open button
	btn = m.builder.GetObject("toolbar_openButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.open)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(openIcon, "open"))

	// Toolbar save button
	btn = m.builder.GetObject("toolbar_saveButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.save)
	btn.SetIsImportant(true)
	btn.SetIconWidget(createImageFromBytes(saveIcon, "save"))

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

// addFile: Handler for the add file button clicked signal
func (m *MainForm) addFile() {
	if m.addFileDialog == nil {
		m.addFileDialog = m.newAddFileDialog()
	}
	m.addFileDialog.openForNewFile("/home/per/temp/dragon.ply")
}

// new: Handler for the new button clicked signal
func (m *MainForm) new() {
	// // TODO : new project here
	// // Open setup dialog
	// result, err := m.openSetupDialog()
	// if err != nil {
	// 	// TODO : Error handling
	// 	return
	// }
	//
	// // Create project file
	// currentProject, err := engine.SetupProject(result.path, result.name)
	// if err != nil {
	// 	// TODO : Error handling
	// 	return
	// }
	//
	// return
}

// open: Handler for the open button clicked signal
func (m *MainForm) open() {
	// TODO : open project here
}

// save: Handler for the save button clicked signal
func (m *MainForm) save() {
	// TODO : save project here
}

func isTraceMode() bool {
	return len(os.Args) >= 2 && strings.HasPrefix(os.Args[1], "-t")
}

func (m *MainForm) addFileButtonClicked() {

}

func (m *MainForm) editFileButtonClicked() {

}

func (m *MainForm) removeFileButtonClicked() {

}

func (m *MainForm) enableDisableStackPages(status bool) {
	m.enableDisableStackPage("mainWindow_controlPage", status)
	m.enableDisableStackPage("mainWindow_preinstallPage", status)
	m.enableDisableStackPage("mainWindow_installPage", status)
	m.enableDisableStackPage("mainWindow_postinstallPage", status)
	m.enableDisableStackPage("mainWindow_copyrightPage", status)
}

func (m *MainForm) enableDisableStackPage(name string, status bool) {
	// TODO : Fix this code, should be doable with *gtk.Widget
	box, ok := m.builder.GetObject(name).(*gtk.Box)
	if !ok {
		grid, ok := m.builder.GetObject(name).(*gtk.Grid)
		if !ok {
			m.log.Error.Printf("failed to retrieve stack page: %s", name)
		}
		grid.SetSensitive(status)
		return
	}
	box.SetSensitive(status)
}
