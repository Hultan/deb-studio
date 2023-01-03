package gui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
)

func (m *MainForm) setAsLatestVersionClicked() {
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

func (m *MainForm) setPackageAsCurrentClicked() {
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

func (m *MainForm) addPackageClicked() {
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
	fmt.Println("Edit file clicked!")
}

func (m *MainForm) removeFileButtonClicked() {
	fmt.Println("Remove file clicked!")
}

// newButtonClicked: Handler for the newButtonClicked button clicked signal
func (m *MainForm) newButtonClicked() {
	defer func() {
		m.listPackages()
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
	// TODO : Handle if a project is already open
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

		m.updateProjectPage()
		m.updatePackagePage()
		m.listPackages()
		m.updateInfoBar()
		m.enableDisableStackPages()
	}
}

func (m *MainForm) updatePackagePage() {
	if project.Config.ShowOnlyLatestVersion {
		m.showOnlyCheckBox.SetActive(true)
	}
}

func (m *MainForm) updateProjectPage() {
	entry := m.builder.GetObject("mainWindow_projectHeaderLabel").(*gtk.Label)
	entry.SetText(project.Config.Name)
	entry = m.builder.GetObject("mainWindow_projectSubheaderLabel").(*gtk.Label)
	entry.SetText("Project information")
	entry = m.builder.GetObject("mainWindow_projectNameLabel").(*gtk.Label)
	entry.SetMarkup("Project name: <b>" + project.Config.Name + "</b>")
	entry = m.builder.GetObject("mainWindow_projectPathLabel").(*gtk.Label)
	entry.SetMarkup("Project path: <b>" + project.Path + "</b>")
	entry = m.builder.GetObject("mainWindow_latestVersionLabel").(*gtk.Label)
	entry.SetMarkup("Latest version: <b>" + project.Config.LatestVersion + "</b>")
	entry = m.builder.GetObject("mainWindow_currentPackageLabel").(*gtk.Label)
	entry.SetMarkup("Current package: <b>" + project.Config.CurrentPackage + "</b>")
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

// openProjectFolder: Handler for the open project folder button clicked signal
func (m *MainForm) openProjectFolder() {
	cmd := exec.Command("xdg-open", project.Path)
	cmd.Run()
}

// openPackageFolder: Handler for the open package folder button clicked signal
func (m *MainForm) openPackageFolder() {
	// Set version as latest
	pkgName := m.projectList.GetSelectedPackageName()
	if pkgName == "" {
		return
	}
	pkg := project.GetPackageByName(pkgName)
	if pkg == nil {
		return
	}
	cmd := exec.Command("xdg-open", pkg.Path)
	cmd.Run()
}

func (m *MainForm) showOnlyCurrentAndLatestToggled(check *gtk.CheckButton) {
	checked := check.GetActive()
	project.SetShowOnlyCurrentAndLatest(checked)
	m.listPackages()
}
