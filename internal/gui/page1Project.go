package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *MainWindow) updateProjectPage() {
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
