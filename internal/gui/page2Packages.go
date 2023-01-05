package gui

import (
	"fmt"
	"os/exec"

	"github.com/google/uuid"
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
	infoBar      *gtk.InfoBar
	infoBarLabel *gtk.Label

	// Toolbar
	showOnlyCheckBox *gtk.CheckButton

	// Popup menu
	popup *gtk.Menu
}

func (m *MainWindow) setupPackagePage() *pagePackage {
	p := &pagePackage{parent: m}

	// General
	treeView := m.builder.GetObject("mainWindow_packageList").(*gtk.TreeView)
	treeView.Connect("row_activated", p.setPackageAsCurrentClicked)
	treeView.Connect("button-press-event", p.showPopupMenu)
	p.projectList = packageList.NewProjectList(treeView)
	p.infoBar = m.builder.GetObject("mainWindow_infoBar").(*gtk.InfoBar)
	p.infoBarLabel = m.builder.GetObject("mainWindow_infoBarLabel").(*gtk.Label)

	// Toolbar
	btn := m.builder.GetObject("mainWindow_addPackageButton").(*gtk.ToolButton)
	btn.Connect("clicked", p.addPackageClicked)
	btn = m.builder.GetObject("mainWindow_removePackageButton").(*gtk.ToolButton)
	btn.Connect("clicked", p.removePackageClicked)
	p.showOnlyCheckBox = m.builder.GetObject("mainPage_toolbarShowOnlyCurrentAndLatest").(*gtk.CheckButton)
	p.showOnlyCheckBox.Connect("toggled", p.showOnlyCurrentAndLatestToggled)

	// Popup
	p.popup = m.builder.GetObject("mainWindow_popupPackageMenu").(*gtk.Menu)
	tool := m.builder.GetObject("mainWindow_popupAddPackage").(*gtk.MenuItem)
	tool.Connect("activate", p.addPackageClicked)
	tool = m.builder.GetObject("mainWindow_popupRemovePackage").(*gtk.MenuItem)
	tool.Connect("activate", p.removePackageClicked)
	tool = m.builder.GetObject("mainWindow_popupSetAsLatest").(*gtk.MenuItem)
	tool.Connect("activate", p.setAsLatestVersionClicked)
	tool = m.builder.GetObject("mainWindow_popupSetAsCurrent").(*gtk.MenuItem)
	tool.Connect("activate", p.setPackageAsCurrentClicked)
	tool = m.builder.GetObject("mainWindow_popupOpenProject").(*gtk.MenuItem)
	tool.Connect("activate", p.openProjectFolder)
	tool = m.builder.GetObject("mainWindow_popupOpenPackage").(*gtk.MenuItem)
	tool.Connect("activate", p.openPackageFolder)

	return p
}

func (p *pagePackage) setAsLatestVersionClicked() {
	// Set version as latest
	id := p.projectList.GetSelectedPackageId()
	if id == uuid.Nil {
		return
	}
	project.SetAsLatest(id)

	// Update gui
	p.parent.pages.update()
}

func (p *pagePackage) setPackageAsCurrentClicked() {
	// Set package as current
	id := p.projectList.GetSelectedPackageId()
	if id == uuid.Nil {
		return
	}
	project.SetAsCurrent(id)

	// Update some things
	p.parent.pages.update()
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
	row.SetName(pkg.Config.GetPackageName())

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
	p.listPackages()
	p.updateInfoBar()
}

func (p *pagePackage) showOnlyCurrentAndLatestToggled(check *gtk.CheckButton) {
	checked := check.GetActive()
	project.SetShowOnlyLatestVersion(checked)
	p.parent.pages.update()
}

func (p *pagePackage) getInfoBarText() string {
	switch getProjectStatus() {
	case projectStatusNoProjectOpened:
		return "You need to open or create a new project..."
	case projectStatusNoPackageSelected:
		return "You need to select or add a package to edit!"
	case projectStatusNotLatestVersion:
		return fmt.Sprintf(
			"You are currently not editing the latest version! You are editing <b>version %s</b> and <b>architecture %s</b>.",
			project.CurrentPackage.Config.Version, project.CurrentPackage.Config.Architecture,
		)
	case projectStatusLatestVersion:
		return fmt.Sprintf(
			"You are currently editing <b>version %s</b> and <b>architecture %s</b>.",
			project.CurrentPackage.Config.Version, project.CurrentPackage.Config.Architecture,
		)
	default:
		log.Error.Println("Invalid projectStatus in getInfoBarText()")
		return ""
	}
}

func (p *pagePackage) setInfoBarColor() {
	switch getProjectStatus() {
	case projectStatusNoProjectOpened:
		p.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	case projectStatusNoPackageSelected:
		p.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case projectStatusNotLatestVersion:
		p.infoBar.SetMessageType(gtk.MESSAGE_WARNING)
	case projectStatusLatestVersion:
		p.infoBar.SetMessageType(gtk.MESSAGE_INFO)
	default:
		log.Error.Println("Invalid projectStatus in SetInfoBarColor()")
	}
}

func (p *pagePackage) updateInfoBar() {
	p.infoBarLabel.SetMarkup(p.getInfoBarText())
	p.setInfoBarColor()
	// Force a redraw to update the info bar
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
	id := p.projectList.GetSelectedPackageId()
	if id == uuid.Nil {
		return
	}
	pkg := project.GetPackageById(id)
	if pkg == nil {
		return
	}
	cmd := exec.Command("xdg-open", pkg.Path)
	cmd.Run()
}
