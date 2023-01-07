package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
)

// setupToolbar: Set up the toolbar
func (m *MainWindow) setupToolbar() {
	log.Trace.Println("Entering setupToolbar...")
	defer log.Trace.Println("Exiting setupToolbar...")

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

// newButtonClicked: Handler for the newButtonClicked button clicked signal
func (m *MainWindow) newButtonClicked() {
	log.Trace.Println("Entering newButtonClicked...")
	defer log.Trace.Println("Exiting newButtonClicked...")

	defer func() {
		m.pages.update()
	}()

	// Open setup dialog
	result, err := m.openSetupDialog()
	if err != nil {
		// TODO : Error handling
		return
	}
	if result == nil {
		// User cancelled out of setup dialog
		return
	}

	// Create project file
	project, err = engine.NewProject(log, result.path, result.name)
	if err != nil {
		// TODO : Handle error
	}
}

// openButtonClicked: Handler for the openButtonClicked button clicked signal
func (m *MainWindow) openButtonClicked() {
	log.Trace.Println("Entering openButtonClicked...")
	defer log.Trace.Println("Exiting openButtonClicked...")

	// TODO : Handle if a project is already open

	var err error
	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Select folder...", m.window, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		"OK", gtk.RESPONSE_ACCEPT,
		"Cancel", gtk.RESPONSE_CANCEL,
	)
	if err != nil {
		log.Error.Printf("failed to create fileChooserDialog: %s\n", err)
		showErrorDialog("Failed to open fileChooserDialog!", err)
		return
	}
	dlg.SetCurrentFolder("/home/per/installs/softtube")

	response := dlg.Run()
	dlg.Hide()
	if response == gtk.RESPONSE_ACCEPT {
		projectFolder, err := dlg.GetCurrentFolder()
		if err != nil {
			log.Error.Printf("failed to get folder from fileChooserDialog: %s\n", err)
			msg := "Failed to get folder from fileChooserDialog!"
			showErrorDialog(msg, err)
			return
		}
		project, err = engine.OpenProject(log, projectFolder)
		// TODO: Handle if project.json does not exist
		if err != nil {
			log.Error.Printf("failure during opening of '%s': %s", projectFolder, err)
			msg := fmt.Sprintf("failed to open project folder: %s", err)
			showErrorDialog(msg, err)
			return
		}

		// Update gui
		m.pages.controlPage.init()
		m.pages.scriptPage.init()
		m.pages.textPage.init()
		if project.Config.ShowOnlyLatestVersion {
			m.pages.packagePage.showOnlyCheckBox.SetActive(true)
		}
		m.pages.update()
	}
}

// buildButtonClicked: Handler for the build button clicked signal
func (m *MainWindow) buildButtonClicked() {
	log.Trace.Println("Entering buildButtonClicked...")
	defer log.Trace.Println("Exiting buildButtonClicked...")

	fmt.Println("Build clicked")
}
