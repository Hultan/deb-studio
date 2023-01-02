package gui

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
)

// TODO : All actions needs to check project.IsPackageSelected()

func (m *MainForm) setAsLatestVersionClicked() {
}

func (m *MainForm) setPackageAsCurrentClickedPopup() {
}

func (m *MainForm) setPackageAsCurrentClicked() {
}

func (m *MainForm) addPackageClicked() {
	fmt.Println("Add package clicked!")

	dialog := m.builder.GetObject("addPackageDialog").(*gtk.Dialog)
	// versionEntry := m.builder.GetObject("addInstallationDialog_versionNameEntry").(*gtk.Entry)
	// architectureCombo := m.builder.GetObject("addInstallationDialog_architectureCombo").(*gtk.Dialog)
	dialog.AddButton("Add", gtk.RESPONSE_ACCEPT)
	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)

	// Show the dialog
	responseId := dialog.Run()
	if responseId == gtk.RESPONSE_ACCEPT {
		// Save installation

	}

	dialog.Hide()

}

func (m *MainForm) removePackageClicked() {
	fmt.Println("Remove package clicked!")
}

func (m *MainForm) addFileButtonClicked() {
	if m.addFileDialog == nil {
		m.addFileDialog = m.newAddFileDialog()
	}
	m.addFileDialog.openForNewFile("/home/per/temp/dragon.ply")
}

func (m *MainForm) editFileButtonClicked() {

}

func (m *MainForm) removeFileButtonClicked() {

}

// newButtonClicked: Handler for the newButtonClicked button clicked signal
func (m *MainForm) newButtonClicked() {
	defer func() {
		m.listPackages()
		m.printTraceInfo()
		m.updateInfoBar()
		m.enableDisableStackPages()
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
func (m *MainForm) openButtonClicked() {
	var err error
	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Select folder...", m.window, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		"OK", gtk.RESPONSE_ACCEPT,
		"Cancel", gtk.RESPONSE_CANCEL,
	)
	if err != nil {
		showErrorDialog("Failed to open fileChooserDialog!", err)
		log.Error.Printf("failed to create fileChooserDialog: %s\n", err)
		return
	}
	dlg.SetCurrentFolder("/home/per/installs/softtube")

	response := dlg.Run()
	dlg.Hide()
	if response == gtk.RESPONSE_ACCEPT {
		projectFolder, err := dlg.GetCurrentFolder()
		if err != nil {
			msg := "Failed to get folder from fileChooserDialog!"
			showErrorDialog(msg, err)
			log.Error.Printf("failed to get folder from fileChooserDialog: %s\n", err)
			return
		}
		project, err = engine.OpenProject(log, projectFolder)
		// TODO: Handle if project.json does not exist
		if err != nil {
			msg := fmt.Sprintf("failed to open project folder: %s", err)
			showErrorDialog(msg, err)
			log.Error.Printf("failure during opening of '%s': %s", projectFolder, err)
			os.Exit(1)
		}

		m.listPackages()
		m.printTraceInfo()
		m.updateInfoBar()
		m.enableDisableStackPages()
	}
}

// saveButtonClicked: Handler for the saveButtonClicked button clicked signal
// TODO : Do we even need a save button?
func (m *MainForm) saveButtonClicked() {
	// TODO : saveButtonClicked project here
	fmt.Println("Save clicked")
}

// buildButtonClicked: Handler for the build button clicked signal
func (m *MainForm) buildButtonClicked() {
	// TODO : build project here
	fmt.Println("Build clicked")
}
