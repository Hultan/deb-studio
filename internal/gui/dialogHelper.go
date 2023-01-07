package gui

import (
	"github.com/hultan/dialog"
)

func showErrorDialog(msg string, err error) {
	dlg := dialog.
		Title("An error occurred...").
		Width(300).
		Text(msg).
		ErrorIcon().
		OkButton()

	if err != nil {
		dlg = dlg.Extra(err.Error())
	}

	_, internalErr := dlg.Show()

	if internalErr != nil {
		log.Error.Println("failed to show information dialog: %s\n", internalErr)
	}
}

func showInformationDialog(title, msg string) {
	_, err := dialog.
		Title(title).
		Text(msg).
		InfoIcon().
		OkButton().
		Show()

	if err != nil {
		log.Error.Println("failed to show information dialog: %s\n", err)
	}
}
