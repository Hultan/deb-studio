package projectList

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type ProjectList struct {
	treeView *gtk.TreeView
}

func NewProjectList(treeView *gtk.TreeView) *ProjectList {
	p := &ProjectList{treeView: treeView}
	p.setupColumns()
	return p
}

func (p *ProjectList) RefreshList(store *gtk.ListStore) {
	p.treeView.SetModel(store)
	p.treeView.ShowAll()
}

func (p *ProjectList) setupColumns() {
	// p.treeView.AppendColumn(p.createTextColumn("Is latest", 0, 70, 300))
	p.treeView.AppendColumn(p.createImageColumn("latest"))
	p.treeView.AppendColumn(p.createTextColumn("Version name", packageListColumnVersionName, 0, 300))
	p.treeView.AppendColumn(p.createTextColumn("Architecture name", packageListColumnArchitectureName, 0, 300))
}

// createTextColumn : Add a column to the tree view (during the initialization of the tree view)
func (p *ProjectList) createTextColumn(title string, id int, width int, weight int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}
	_ = cellRenderer.SetProperty("weight", weight)

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	// column.AddAttribute(cellRenderer, "background", 3)
	// column.AddAttribute(cellRenderer, "foreground", 4)
	if width == 0 {
		column.SetExpand(true)
	} else {
		column.SetFixedWidth(width)
	}

	return column
}

// createImageColumn : Add a column to the tree view (during the initialization of the tree view)
func (p *ProjectList) createImageColumn(title string) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		log.Fatal("Unable to create pixbuf cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", packageListColumnIsLatest)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(70)
	column.SetVisible(true)
	column.SetExpand(false)

	return column
}

func (p *ProjectList) GetSelectedVersionGuid() string {
	selection, err := p.treeView.GetSelection()
	if err != nil {
		return ""
	}
	model, iter, ok := selection.GetSelected()
	if ok {
		return ""
	}
	v, _ := model.ToTreeModel().GetValue(iter, packageListColumnVersionGuid)
	guid, _ := v.GetString()
	return guid
}

func (p *ProjectList) GetSelectedArchitectureGuid() string {
	selection, err := p.treeView.GetSelection()
	if err != nil {
		return ""
	}
	model, iter, ok := selection.GetSelected()
	if ok {
		return ""
	}
	v, _ := model.ToTreeModel().GetValue(iter, packageListColumnArchitectureGuid)
	guid, _ := v.GetString()
	return guid
}
