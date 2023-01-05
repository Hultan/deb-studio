package packageList

import (
	"log"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
)

type PackageList struct {
	treeView *gtk.TreeView
}

func NewProjectList(treeView *gtk.TreeView) *PackageList {
	p := &PackageList{treeView: treeView}
	p.setupColumns()
	return p
}

func (p *PackageList) RefreshList(store *gtk.TreeModelFilter) {
	p.treeView.SetModel(store)
	p.treeView.ShowAll()
}

func (p *PackageList) setupColumns() {
	// p.treeView.AppendColumn(p.createTextColumn("Is latest", 0, 70, 300))
	p.treeView.AppendColumn(p.createImageColumn("Current", common.PackageListColumnIsCurrent))
	p.treeView.AppendColumn(p.createImageColumn("Latest", common.PackageListColumnIsLatest))
	p.treeView.AppendColumn(p.createTextColumn("Version name", common.PackageListColumnVersionName, 0, 300))
	p.treeView.AppendColumn(p.createTextColumn("Architecture name", common.PackageListColumnArchitectureName, 0, 300))
	p.treeView.AppendColumn(p.createTextColumn("Package path", common.PackageListColumnPackagePath, 0, 300))
}

// createTextColumn : Add a column to the tree view (during the initialization of the tree view)
func (p *PackageList) createTextColumn(title string, id int, width int, weight int) *gtk.TreeViewColumn {
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
func (p *PackageList) createImageColumn(title string, colIndex int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		log.Fatal("Unable to create pixbuf cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", colIndex)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(70)
	column.SetVisible(true)
	column.SetExpand(false)

	return column
}

func (p *PackageList) GetSelectedPackageName() string {
	selection, err := p.treeView.GetSelection()
	if err != nil {
		return ""
	}
	model, iter, ok := selection.GetSelected()
	if !ok {
		return ""
	}
	v, err := model.ToTreeModel().GetValue(iter, common.PackageListColumnPackageName)
	if err != nil {
		return ""
	}

	name, err := v.GetString()
	if err != nil {
		return ""
	}

	return name
}
