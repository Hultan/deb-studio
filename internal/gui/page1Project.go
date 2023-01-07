package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type pageProject struct {
	parent *MainWindow

	headerLabel         *gtk.Label
	subHeaderLabel      *gtk.Label
	projectNameLabel    *gtk.Label
	projectPathLabel    *gtk.Label
	latestVersionLabel  *gtk.Label
	currentVersionLabel *gtk.Label
}

func (m *MainWindow) setupProjectPage() *pageProject {
	log.Trace.Println("Entering setupProjectPage...")
	defer log.Trace.Println("Exiting setupProjectPage...")

	p := &pageProject{parent: m}

	p.headerLabel = m.builder.GetObject("mainWindow_projectHeaderLabel").(*gtk.Label)
	p.subHeaderLabel = m.builder.GetObject("mainWindow_projectSubheaderLabel").(*gtk.Label)
	p.projectNameLabel = m.builder.GetObject("mainWindow_projectNameLabel").(*gtk.Label)
	p.projectPathLabel = m.builder.GetObject("mainWindow_projectPathLabel").(*gtk.Label)
	p.latestVersionLabel = m.builder.GetObject("mainWindow_latestVersionLabel").(*gtk.Label)
	p.currentVersionLabel = m.builder.GetObject("mainWindow_currentPackageLabel").(*gtk.Label)

	return p
}

func (p *pageProject) update() {
	log.Trace.Println("Entering update...")
	defer log.Trace.Println("Exiting update...")

	if project == nil {
		p.headerLabel.SetText("Please open or create a project...")
		return
	}

	p.headerLabel.SetText(project.Config.Name)
	p.subHeaderLabel.SetText("Project information")
	p.projectNameLabel.SetMarkup("Project name: <b>" + project.Config.Name + "</b>")
	p.projectPathLabel.SetMarkup("Project path: <b>" + project.Path + "</b>")
	p.latestVersionLabel.SetMarkup("Latest version: <b>" + project.Config.LatestVersion + "</b>")
	current := ""
	if project.CurrentPackage != nil {
		current = fmt.Sprintf(
			"%s (%s)",
			project.CurrentPackage.Path,
			project.Config.CurrentPackageId,
		)
	}
	p.currentVersionLabel.SetMarkup("Current package: <b>" + current + "</b>")
}
