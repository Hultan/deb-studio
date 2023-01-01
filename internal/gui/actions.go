package gui

import (
	"fmt"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
)

// TODO : All actions needs to check currentVersion and currentArchitecture

func (m *MainForm) setAsLatestVersionClicked() {
	fmt.Println("Set as latest version!")
}

func (m *MainForm) setPackageAsCurrentClickedPopup() {
	row := m.listBox.GetSelectedRow()
	m.setPackageAsCurrent(row)
}

func (m *MainForm) setPackageAsCurrentClicked(l *gtk.ListBox, row *gtk.ListBoxRow) {
	m.setPackageAsCurrent(row)
}

func (m *MainForm) setPackageAsCurrent(row *gtk.ListBoxRow) {
	defer func() {
		m.updateInfoBar()
		m.enableDisableStackPages()
	}()

	name, err := row.GetName()
	if err != nil {
		log.Error.Printf("failed to set current package")
		return
	}
	list := strings.Split(name, separator)
	versionName := list[0]
	architectureName := list[1]

	project.SetCurrent(versionName, architectureName)
	if project.CurrentVersion == nil || project.CurrentArchitecture == nil {
		log.Error.Printf(
			"failed to set current version '%s' and architecture '%s'",
			versionName, architectureName,
		)
	}
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
	e := engine.NewEngine(log)
	if e.IsProjectFolder(result.path) {
		showErrorDialog("This folder is already a project folder for another project.", err)
		log.Warning.Printf("folder is already a %s project folder : %s", applicationTitle, result.path)
		return
	} else {
		project, err = e.SetupProject(result.path, result.name)
		if err != nil {
			log.Error.Printf("failure during setup of '%s': %s", result.path, err)
			os.Exit(1)
		}
	}
}

// openButtonClicked: Handler for the openButtonClicked button clicked signal
func (m *MainForm) openButtonClicked() {
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
		e := engine.NewEngine(log)
		if !e.IsProjectFolder(projectFolder) {
			// TODO : Ask if the user wants to create a new project here
			title := "Invalid folder"
			msg := fmt.Sprintf("This is not a %s project folder.", applicationTitle)
			showInformationDialog(title, msg)
			log.Trace.Printf("opened a non project folder: %s\n", projectFolder)
			return
		}

		project, err = e.OpenProject(projectFolder)
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
