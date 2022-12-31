package gui

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
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
		m.enableDisableStackPages(m.getInfoBarStatus() != infoBarStatusNoPackageSelected)
	}()

	name, err := row.GetName()
	if err != nil {
		log.Error.Printf("failed to set current package")
		return
	}
	list := strings.Split(name, "$$$")
	versionName := list[0]
	architectureName := list[1]

	currentVersion = currentProject.GetVersion(versionName)
	if currentVersion == nil {
		log.Error.Printf("failed to find version '%s'", versionName)
		return
	}
	currentArchitecture = currentVersion.GetArchitecture(architectureName)
	if currentArchitecture == nil {
		currentVersion = nil
		log.Error.Printf(
			"failed to find architecture '%s' in version '%s'",
			architectureName, versionName,
		)
		return
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

// openButtonClicked: Handler for the openButtonClicked button clicked signal
func (m *MainForm) openButtonClicked() {
	// TODO : open project here
}

// saveButtonClicked: Handler for the saveButtonClicked button clicked signal
// TODO : Do we need a save button?
func (m *MainForm) saveButtonClicked() {
	// TODO : saveButtonClicked project here
}

// buildButtonClicked: Handler for the build button clicked signal
func (m *MainForm) buildButtonClicked() {
	// TODO : build project here
}
