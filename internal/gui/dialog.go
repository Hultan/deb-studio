package gui

import (
	"github.com/hultan/dialog"
)

func showErrorDialog(msg string, err error) {
	dlg := dialog.
		Title("An error occurred...").
		Text(msg).
		ErrorIcon().
		OkButton()

	if err != nil {
		dlg = dlg.Extra(err.Error())
	}

	dlg.Show()
}

func showInformationDialog(title, msg string) {
	dialog.
		Title(title).
		Text(msg).
		InfoIcon().
		OkButton().
		Show()
}
