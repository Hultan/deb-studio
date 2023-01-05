package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
	"github.com/hultan/deb-studio/internal/projectList"
)

func (m *MainWindow) setupPackagePage() {
	m.treeView = m.builder.GetObject("mainWindow_packageList").(*gtk.TreeView)
	m.projectList = projectList.NewProjectList(m.treeView)
	m.treeView.Connect("row_activated", m.setPackageAsCurrentClicked)
	m.treeView.Connect("button-press-event", m.showPopupMenu)

	m.infoBar = m.builder.GetObject("mainWindow_infoBar").(*gtk.InfoBar)
	m.infoBarLabel = m.builder.GetObject("mainWindow_infoBarLabel").(*gtk.Label)

	// Installation Button
	btn := m.builder.GetObject("mainWindow_addPackageButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.addPackageClicked)

	btn = m.builder.GetObject("mainWindow_removePackageButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.removePackageClicked)

	m.showOnlyCheckBox = m.builder.GetObject("mainPage_toolbarShowOnlyCurrentAndLatest").(*gtk.CheckButton)
	m.showOnlyCheckBox.Connect("toggled", m.showOnlyCurrentAndLatestToggled)
}

func (m *MainWindow) setAsLatestVersionClicked() {
	// Set version as latest
	pkgName := m.projectList.GetSelectedPackageName()
	if pkgName == "" {
		return
	}
	project.SetAsLatest(pkgName)

	// Update some things
	m.updateProjectPage()
	m.listPackages()
	m.updateInfoBar()
}

func (m *MainWindow) setPackageAsCurrentClicked() {
	// Set package as current
	pkgName := m.projectList.GetSelectedPackageName()
	if pkgName == "" {
		return
	}
	project.SetAsCurrent(pkgName)

	// Update some things
	m.updateProjectPage()
	m.listPackages()
	m.updateInfoBar()
}

func (m *MainWindow) addPackageClicked() {
	fmt.Println("Add package clicked!")

	dialog := m.builder.GetObject("addPackageDialog").(*gtk.Dialog)
	// versionEntry := m.builder.GetObject("addInstallationDialog_versionNameEntry").(*gtk.Entry)
	// architectureCombo := m.builder.GetObject("addInstallationDialog_architectureCombo").(*gtk.Dialog)
	_, err := dialog.AddButton("Add", gtk.RESPONSE_ACCEPT)
	if err != nil {
		return
	}
	_, err = dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	if err != nil {
		return
	}

	// Show the dialog
	responseId := dialog.Run()
	if responseId == gtk.RESPONSE_ACCEPT {
		// Add package
	}

	dialog.Hide()
}

func (m *MainWindow) removePackageClicked() {
	fmt.Println("Remove package clicked!")
}

func (m *MainWindow) listPackages() {
	if project == nil {
		return
	}

	store := project.GetPackageListStore(checkIcon, editIcon)
	m.projectList.RefreshList(store)
}

func (m *MainWindow) createPackageListRow(p *engine.Package) (*gtk.ListBoxRow, error) {
	row, err := gtk.ListBoxRowNew()
	if err != nil {
		log.Error.Printf("failed to create package list row")
		return nil, err
	}
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 20)
	box.SetHomogeneous(true)
	if err != nil {
		log.Error.Printf("failed to create package box")
		return nil, err
	}
	row.Add(box)
	// TODO : Change to a map instead?
	row.SetName(p.Config.GetFolderName())

	// Add version label
	label, err := gtk.LabelNew(p.Config.Version)
	label.SetHAlign(gtk.ALIGN_START)
	if err != nil {
		log.Error.Printf("failed to create package version label")
		return nil, err
	}
	box.PackStart(label, false, true, 20)

	// Add architecture label
	label, err = gtk.LabelNew(p.Config.Architecture)
	label.SetHAlign(gtk.ALIGN_START)
	if err != nil {
		log.Error.Printf("failed to create package architecture label")
		return nil, err
	}
	box.PackStart(label, false, true, 20)
	return row, nil
}

func (m *MainWindow) updatePackagePage() {
	if project.Config.ShowOnlyLatestVersion {
		m.showOnlyCheckBox.SetActive(true)
	}
}

func (m *MainWindow) showOnlyCurrentAndLatestToggled(check *gtk.CheckButton) {
	checked := check.GetActive()
	project.SetShowOnlyCurrentAndLatest(checked)
	m.listPackages()
}
