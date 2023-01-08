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
	projectList *packageList.PackageList

	// Toolbar
	showOnlyCheckBox *gtk.CheckButton

	// Popup menu
	popup *gtk.Menu
}

func (m *MainWindow) setupPackagePage() *pagePackage {
	log.Trace.Println("Entering setupPackagePage...")
	defer log.Trace.Println("Exiting setupPackagePage...")

	p := &pagePackage{parent: m}

	// General
	treeView := m.builder.GetObject("mainWindow_packageList").(*gtk.TreeView)
	treeView.Connect("row_activated", p.setPackageAsCurrentClicked)
	treeView.Connect("button-press-event", p.showPopupMenu)
	p.projectList = packageList.NewProjectList(treeView)

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

func (p *pagePackage) update() {
	log.Trace.Println("Entering update...")
	defer log.Trace.Println("Exiting update...")

	p.listPackages()
}

func (p *pagePackage) setAsLatestVersionClicked() {
	log.Trace.Println("Entering setAsLatestVersionClicked...")
	defer log.Trace.Println("Exiting setAsLatestVersionClicked...")

	// Set version as latest
	id := p.projectList.GetSelectedPackageId()
	if id == "" {
		log.Warning.Println("no selected package found")
		return
	}
	project.SetAsLatest(id)

	// Update gui
	p.parent.pages.update()
}

func (p *pagePackage) setPackageAsCurrentClicked() {
	log.Trace.Println("Entering setPackageAsCurrentClicked...")
	defer log.Trace.Println("Exiting setPackageAsCurrentClicked...")

	// Set package as current
	id := p.projectList.GetSelectedPackageId()
	if id == "" {
		log.Warning.Println("no selected package found")
		return
	}
	project.SetAsCurrent(id)

	// Update some things
	p.parent.pages.update()
}

func (p *pagePackage) addPackageClicked() {
	log.Trace.Println("Entering addPackageClicked...")
	defer log.Trace.Println("Exiting addPackageClicked...")

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
	log.Trace.Println("Entering removePackageClicked...")
	defer log.Trace.Println("Exiting removePackageClicked...")

	fmt.Println("Remove package clicked!")
}

func (p *pagePackage) listPackages() {
	log.Trace.Println("Entering listPackages...")
	defer log.Trace.Println("Exiting listPackages...")

	if project == nil {
		return
	}

	store := project.GetPackageListStore()
	p.projectList.RefreshList(store)
}

func (p *pagePackage) createPackageListRow(pkg *engine.Package) (*gtk.ListBoxRow, error) {
	log.Trace.Println("Entering createPackageListRow...")
	defer log.Trace.Println("Exiting createPackageListRow...")

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

func (p *pagePackage) showOnlyCurrentAndLatestToggled(check *gtk.CheckButton) {
	log.Trace.Println("Entering showOnlyCurrentAndLatestToggled...")
	defer log.Trace.Println("Exiting showOnlyCurrentAndLatestToggled...")

	checked := check.GetActive()
	project.SetShowOnlyLatestVersion(checked)
	p.parent.pages.update()
}

func (p *pagePackage) showPopupMenu(_ *gtk.ListBox, e *gdk.Event) {
	log.Trace.Println("Entering showPopupMenu...")
	defer log.Trace.Println("Exiting showPopupMenu...")

	ev := gdk.EventButtonNewFromEvent(e)
	if ev.Button() == common.RightMouseButton {
		p.popup.PopupAtPointer(e)
	}
}

// openProjectFolder: Handler for the open project folder button clicked signal
func (p *pagePackage) openProjectFolder() {
	log.Trace.Println("Entering openProjectFolder...")
	defer log.Trace.Println("Exiting openProjectFolder...")

	cmd := exec.Command("xdg-open", project.Path)
	err := cmd.Run()
	if err != nil {
		log.Error.Printf("failed to open project folder: %s\n", err)
	}
}

// openPackageFolder: Handler for the open package folder button clicked signal
func (p *pagePackage) openPackageFolder() {
	log.Trace.Println("Entering openPackageFolder...")
	defer log.Trace.Println("Exiting openPackageFolder...")

	// Set version as latest
	id := p.projectList.GetSelectedPackageId()
	if id == "" {
		return
	}
	pkg := project.GetPackageById(id)
	if pkg == nil {
		return
	}
	cmd := exec.Command("xdg-open", pkg.Path)
	err := cmd.Run()
	if err != nil {
		log.Error.Printf("failed to open project folder: %s\n", err)
	}
}
