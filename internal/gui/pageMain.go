package gui

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
	"github.com/hultan/deb-studio/internal/projectList"
)

func (m *MainForm) setupPackagePage() {
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
}

func (m *MainForm) listPackages() {
	if project == nil {
		return
	}

	store := project.GetPackageListStore(checkIcon)
	m.projectList.RefreshList(store)
	//
	// // Sort packages by newest first
	// m.listBox.SetSortFunc(
	// 	func(row1 *gtk.ListBoxRow, row2 *gtk.ListBoxRow) int {
	// 		name1, _ := row1.GetName()
	// 		name2, _ := row2.GetName()
	// 		if name1 < name2 {
	// 			return 1
	// 		} else if name1 == name2 {
	// 			return 0
	// 		}
	// 		return -1
	// 	},
	// )

	m.treeView.ShowAll()
}

func (m *MainForm) createPackageListRow(p *engine.Package) (*gtk.ListBoxRow, error) {
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
	row.SetName(p.Config.Name)

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
