package gui

import (
	"fmt"
	"os/exec"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
	"github.com/hultan/deb-studio/internal/engine"
	"github.com/hultan/deb-studio/internal/packageList"
)

type pagePackage struct {
	parent *MainWindow

	// General
	projectList  *packageList.PackageList
	treeView     *gtk.TreeView
	infoBar      *gtk.InfoBar
	infoBarLabel *gtk.Label

	// Toolbar
	addPackageButton    *gtk.ToolButton
	removePackageButton *gtk.ToolButton
	showOnlyCheckBox    *gtk.CheckButton

	// Popup menu
	popup              *gtk.Menu
	popupAddPackage    *gtk.MenuItem
	popupRemovePackage *gtk.MenuItem
	popupSetLatest     *gtk.MenuItem
	popupSetCurrent    *gtk.MenuItem
	popupOpenProject   *gtk.MenuItem
	popupOpenPackage   *gtk.MenuItem
}

func (m *MainWindow) setupPackagePage() {
	p := &pagePackage{}

	// General
	p.treeView = m.builder.GetObject("mainWindow_packageList").(*gtk.TreeView)
	p.projectList = packageList.NewProjectList(p.treeView)
	p.treeView.Connect("row_activated", p.setPackageAsCurrentClicked)
	p.treeView.Connect("button-press-event", p.showPopupMenu)
	p.infoBar = m.builder.GetObject("mainWindow_infoBar").(*gtk.InfoBar)
	p.infoBarLabel = m.builder.GetObject("mainWindow_infoBarLabel").(*gtk.Label)

	// Toolbar
	p.addPackageButton = m.builder.GetObject("mainWindow_addPackageButton").(*gtk.ToolButton)
	p.addPackageButton.Connect("clicked", p.addPackageClicked)
	p.removePackageButton = m.builder.GetObject("mainWindow_removePackageButton").(*gtk.ToolButton)
	p.removePackageButton.Connect("clicked", p.removePackageClicked)
	p.showOnlyCheckBox = m.builder.GetObject("mainPage_toolbarShowOnlyCurrentAndLatest").(*gtk.CheckButton)
	p.showOnlyCheckBox.Connect("toggled", p.showOnlyCurrentAndLatestToggled)

	// Popup
	p.popup = m.builder.GetObject("mainWindow_popupPackageMenu").(*gtk.Menu)
	p.popupAddPackage = m.builder.GetObject("mainWindow_popupAddPackage").(*gtk.MenuItem)
	p.popupAddPackage.Connect("activate", p.addPackageClicked)
	p.popupRemovePackage = m.builder.GetObject("mainWindow_popupRemovePackage").(*gtk.MenuItem)
	p.popupRemovePackage.Connect("activate", p.removePackageClicked)
	p.popupSetLatest = m.builder.GetObject("mainWindow_popupSetAsLatest").(*gtk.MenuItem)
	p.popupSetLatest.Connect("activate", p.setAsLatestVersionClicked)
	p.popupSetCurrent = m.builder.GetObject("mainWindow_popupSetAsCurrent").(*gtk.MenuItem)
	p.popupSetCurrent.Connect("activate", p.setPackageAsCurrentClicked)
	p.popupOpenProject = m.builder.GetObject("mainWindow_popupOpenProject").(*gtk.MenuItem)
	p.popupOpenProject.Connect("activate", p.openProjectFolder)
	p.popupOpenPackage = m.builder.GetObject("mainWindow_popupOpenPackage").(*gtk.MenuItem)
	p.popupOpenPackage.Connect("activate", p.openPackageFolder)

	p.parent = m

	m.packagePage = p
}

func (p *pagePackage) setAsLatestVersionClicked() {
	// Set version as latest
	pkgName := p.projectList.GetSelectedPackageName()
	if pkgName == "" {
		return
	}
	project.SetAsLatest(pkgName)

	// Update some things
	p.parent.projectPage.update()
	p.listPackages()
	p.updateInfoBar()
}

func (p *pagePackage) setPackageAsCurrentClicked() {
	// Set package as current
	pkgName := p.projectList.GetSelectedPackageName()
	if pkgName == "" {
		return
	}
	project.SetAsCurrent(pkgName)

	// Update some things
	p.parent.projectPage.update()
	p.listPackages()
	p.updateInfoBar()
}

func (p *pagePackage) addPackageClicked() {
	fmt.Println("Add package clicked!")

	// dialog := p.builder.GetObject("addPackageDialog").(*gtk.Dialog)
	// // versionEntry := m.builder.GetObject("addInstallationDialog_versionNameEntry").(*gtk.Entry)
	// // architectureCombo := m.builder.GetObject("addInstallationDialog_architectureCombo").(*gtk.Dialog)
	// _, err := dialog.AddButton("Add", gtk.RESPONSE_ACCEPT)
	// if err != nil {
	// 	return
	// }
	// _, err = dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	// if err != nil {
	// 	return
	// }
	//
	// // Show the dialog
	// responseId := dialog.Run()
	// if responseId == gtk.RESPONSE_ACCEPT {
	// 	// Add package
	// }
	//
	// dialog.Hide()
}

func (p *pagePackage) removePackageClicked() {
	fmt.Println("Remove package clicked!")
}

func (p *pagePackage) listPackages() {
	if project == nil {
		return
	}

	store := project.GetPackageListStore(checkIcon, editIcon)
	p.projectList.RefreshList(store)
}

func (p *pagePackage) createPackageListRow(pkg *engine.Package) (*gtk.ListBoxRow, error) {
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
	row.SetName(pkg.Config.GetFolderName())

	// Add version label
	label, err := gtk.LabelNew(pkg.Config.Version)
	label.SetHAlign(gtk.ALIGN_START)
	if err != nil {
		log.Error.Printf("failed to create package version label")
		return nil, err
	}
	box.PackStart(label, false, true, 20)

	// Add architecture label
	label, err = gtk.LabelNew(pkg.Config.Architecture)
	label.SetHAlign(gtk.ALIGN_START)
	if err != nil {
		log.Error.Printf("failed to create package architecture label")
		return nil, err
	}
	box.PackStart(label, false, true, 20)
	return row, nil
}

func (p *pagePackage) update() {
	if project.Config.ShowOnlyLatestVersion {
		p.showOnlyCheckBox.SetActive(true)
	}
	p.listPackages()
	p.updateInfoBar()
}

func (p *pagePackage) showOnlyCurrentAndLatestToggled(check *gtk.CheckButton) {
	checked := check.GetActive()
	project.SetShowOnlyCurrentAndLatest(checked)
	p.listPackages()
}

func (p *pagePackage) getInfoBarStatus() infoBarStatus {
	if project == nil {
		return infoBarStatusNoProjectOpened
	} else if project.CurrentPackage == nil {
		return infoBarStatusNoPackageSelected
	} else if !project.WorkingWithLatestVersion() {
		return infoBarStatusNotLatestVersion
	}
	return infoBarStatusLatestVersion
}

func (p *pagePackage) getInfoBarText() string {
	switch p.getInfoBarStatus() {
	case infoBarStatusNoProjectOpened:
		return "You need to open or create a new project..."
	case infoBarStatusNoPackageSelected:
		return "You need to select or add a package to edit!"
	case infoBarStatusNotLatestVersion:
		return fmt.Sprintf(
			"You are currently not editing the latest version! You are editing <b>version %s</b> and <b>architecture %s</b>.",
			project.CurrentPackage.Config.Version, project.CurrentPackage.Config.Architecture,
		)
	case infoBarStatusLatestVersion:
		return fmt.Sprintf(
			"You are currently editing <b>version %s</b> and <b>architecture %s</b>.",
			project.CurrentPackage.Config.Version, project.CurrentPackage.Config.Architecture,
		)
	default:
		log.Error.Println("Invalid infoBarStatus in getInfoBarText()")
		return ""
	}
}

func (p *pagePackage) setInfoBarColor() {
	switch p.getInfoBarStatus() {
	case infoBarStatusNoProjectOpened:
		p.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	case infoBarStatusNoPackageSelected:
		p.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case infoBarStatusNotLatestVersion:
		p.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case infoBarStatusLatestVersion:
		p.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	default:
		log.Error.Println("Invalid infoBarStatus in SetInfoBarColor()")
	}
}

func (p *pagePackage) updateInfoBar() {
	p.infoBarLabel.SetMarkup(p.getInfoBarText())
	p.setInfoBarColor()
	p.parent.window.QueueDraw()
}

func (p *pagePackage) showPopupMenu(_ *gtk.ListBox, e *gdk.Event) {
	ev := gdk.EventButtonNewFromEvent(e)
	if ev.Button() == common.RightMouseButton {
		p.popup.PopupAtPointer(e)
	}
}

// openProjectFolder: Handler for the open project folder button clicked signal
func (p *pagePackage) openProjectFolder() {
	cmd := exec.Command("xdg-open", project.Path)
	cmd.Run()
}

// openPackageFolder: Handler for the open package folder button clicked signal
func (p *pagePackage) openPackageFolder() {
	// Set version as latest
	pkgName := p.projectList.GetSelectedPackageName()
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
