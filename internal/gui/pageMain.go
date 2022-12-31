package gui

import (
	"os"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/engine"
)

func (m *MainForm) setupMainPage() {
	m.listBox = m.builder.GetObject("mainWindow_packageListBox").(*gtk.ListBox)
	m.listBox.Connect("row-activated", m.setPackageAsCurrent)
	m.listBox.Connect("button-press-event", m.showPopupMenu)

	m.infoBar = m.builder.GetObject("mainWindow_infoBar").(*gtk.InfoBar)
	m.infoBarLabel = m.builder.GetObject("mainWindow_infoBarLabel").(*gtk.Label)

	// Installation Button
	btn := m.builder.GetObject("mainWindow_addPackageButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.addPackageClicked)

	btn = m.builder.GetObject("mainWindow_removePackageButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.removePackageClicked)
}

func (m *MainForm) listPackages() {
	for _, version := range currentProject.Versions {
		for _, architecture := range version.Architectures {
			box, err := m.createPackageListRow(version, architecture)
			if err != nil {
				m.log.Error.Println(err)
				os.Exit(1)
			}
			m.listBox.Add(box)
		}
	}

	// Sort packages by newest first
	m.listBox.SetSortFunc(
		func(row1 *gtk.ListBoxRow, row2 *gtk.ListBoxRow) int {
			name1, _ := row1.GetName()
			name2, _ := row2.GetName()
			if name1 < name2 {
				return 1
			} else if name1 == name2 {
				return 0
			}
			return -1
		},
	)
	m.listBox.ShowAll()
}

func (m *MainForm) createPackageListRow(v *engine.Version, a *engine.Architecture) (*gtk.ListBoxRow, error) {
	row, err := gtk.ListBoxRowNew()
	if err != nil {
		m.log.Error.Printf("failed to create package list row")
		return nil, err
	}
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 20)
	box.SetHomogeneous(true)
	if err != nil {
		m.log.Error.Printf("failed to create package box")
		return nil, err
	}
	row.Add(box)
	row.SetName(v.Name + "$$$" + a.Name)

	// Add version label
	label, err := gtk.LabelNew(v.Name)
	label.SetHAlign(gtk.ALIGN_START)
	if err != nil {
		m.log.Error.Printf("failed to create package version label")
		return nil, err
	}
	box.PackStart(label, false, true, 20)

	// Add architecture label
	label, err = gtk.LabelNew(a.Name)
	label.SetHAlign(gtk.ALIGN_START)
	if err != nil {
		m.log.Error.Printf("failed to create package architecture label")
		return nil, err
	}
	box.PackStart(label, false, true, 20)
	return row, nil
}
