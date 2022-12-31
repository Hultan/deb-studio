package gui

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func (m *MainForm) setPackageAsCurrent(l *gtk.ListBox, row *gtk.ListBoxRow) {
	name, err := row.GetName()
	if err != nil {
		m.log.Error.Printf("failed to set current package")
		return
	}
	list := strings.Split(name, "$$$")
	versionName := list[0]
	architectureName := list[1]

	currentVersion = currentProject.GetVersion(versionName)
	if currentVersion == nil {
		m.log.Error.Printf("failed to find version '%s'", currentVersion.Name)
		return
	}
	currentArchitecture = currentVersion.GetArchitecture(architectureName)
	if currentArchitecture == nil {
		m.log.Error.Printf(
			"failed to find architecture '%s' in version '%s'",
			architectureName, versionName,
		)
		return
	}

	m.updateInfoBar()
	m.enableDisableStackPages(m.getInfoBarStatus() != infoBarStatusNoPackageSelected)
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

func (m *MainForm) setAsLatestVersionClicked() {
	fmt.Println("Set as latest version!")
}

func (m *MainForm) setAsCurrentPackageClicked() {
	row := m.listBox.GetSelectedRow()
	m.setPackageAsCurrent(m.listBox, row)
}
