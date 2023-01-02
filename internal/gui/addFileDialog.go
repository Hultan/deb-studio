package gui

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
)

// addFileDialog: Struct for the add file dialog
type addFileDialog struct {
	builder            *Builder
	dialog             *gtk.Dialog
	filePath           *gtk.Entry
	installPath        *gtk.Entry
	staticRadioButton  *gtk.RadioButton
	dynamicRadioButton *gtk.RadioButton
	runScriptCheckbox  *gtk.CheckButton
	runScriptTextView  *gtk.TextView
}

// newAddFileDialog: Constructor for the add file dialog
func (m *MainForm) newAddFileDialog() *addFileDialog {
	d := &addFileDialog{builder: m.builder}

	// Get the dialog window from glade
	dialog := m.builder.GetObject("addFileDialog").(*gtk.Dialog)

	dialog.SetTitle("Add file dialog...")
	dialog.SetTransientFor(m.window)
	dialog.SetModal(true)

	d.dialog = dialog
	d.setupAddFileDialog()

	return d
}

// setupAddFileDialog: Set up the add file dialog
func (a *addFileDialog) setupAddFileDialog() {
	btn := a.builder.GetObject("addFile_filePathButton").(*gtk.Button)
	btn.Connect("clicked", a.filePathButtonClicked)

	btn = a.builder.GetObject("addFile_installPathButton").(*gtk.Button)
	btn.Connect("clicked", a.installPathButtonClicked)

	a.filePath = a.builder.GetObject("addFile_filePathEntry").(*gtk.Entry)
	a.installPath = a.builder.GetObject("addFile_installPathEntry").(*gtk.Entry)
	a.staticRadioButton = a.builder.GetObject("addFile_staticRadioButton").(*gtk.RadioButton)
	a.dynamicRadioButton = a.builder.GetObject("addFile_dynamicRadioButton").(*gtk.RadioButton)
	a.runScriptCheckbox = a.builder.GetObject("addFile_runScriptCheckbox").(*gtk.CheckButton)
	a.runScriptTextView = a.builder.GetObject("addFile_scriptEntry").(*gtk.TextView)

	a.dynamicRadioButton.Connect("toggled", a.radioDynamicToggled)
}

// openForNewFile: Opens the add file dialog for a new file
func (a *addFileDialog) openForNewFile(filePath string) {
	a.setupForNewFile(filePath)

	// Show the dialog
	responseId := a.dialog.Run()
	if responseId == gtk.RESPONSE_ACCEPT {
		// Save settings
	}

	a.dialog.Hide()
}

// setupForNewFile: Set up the add file dialog for a new file
func (a *addFileDialog) setupForNewFile(filePath string) {
	a.filePath.SetText(filePath)
	a.installPath.SetText("")
	a.installPath.GrabFocus()
	buf, err := a.runScriptTextView.GetBuffer()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to get textview buffer")
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(common.ExitCodeGtkError)
	}
	buf.SetText("")
	a.staticRadioButton.SetActive(true)
	a.runScriptCheckbox.SetActive(false)
}

// filePathButtonClicked: Handles the file path button clicked signal
func (a *addFileDialog) filePathButtonClicked() {
	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Choose file to install...",
		a.dialog,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"Ok",
		gtk.RESPONSE_OK,
		"Cancel",
		gtk.RESPONSE_CANCEL,
	)

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to create file path dialog")
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(common.ExitCodeGtkError)
	}

	// Show the dialog
	responseId := dlg.Run()
	if responseId == gtk.RESPONSE_OK {
		a.filePath.SetText(dlg.GetFilename())
	}

	dlg.Hide()
}

// installPathButtonClicked: Handles the install path button clicked signal
func (a *addFileDialog) installPathButtonClicked() {
	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Choose file to install...",
		a.dialog,
		gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		"Ok",
		gtk.RESPONSE_OK,
		"Cancel",
		gtk.RESPONSE_CANCEL,
	)

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to create install path dialog")
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(common.ExitCodeGtkError)
	}

	// Show the dialog
	responseId := dlg.Run()
	if responseId == gtk.RESPONSE_OK {
		a.installPath.SetText(dlg.GetFilename())
	}

	dlg.Hide()
}

// radioDynamicToggled: Handles the toggled signal of the dynamic radio button
func (a *addFileDialog) radioDynamicToggled(btn *gtk.RadioButton) {
	// Enable/disable the run script checkbox and textview
	a.runScriptCheckbox.SetSensitive(btn.GetActive())
	a.runScriptTextView.SetEditable(btn.GetActive())
}
