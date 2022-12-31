package gui

//
// type setup struct {
// 	path string
// 	name string
// }
//
// // openSetupDialog: Open the setup dialog
// func (m *MainForm) openSetupDialog() (*setup, error) {
// 	var result *setup
//
// 	// Get the dialog window from glade
// 	dialog := m.builder.GetObject("setupDialog").(*gtk.Dialog)
// 	locationFileChooser := m.builder.GetObject("setupDialog_locationFileChooser").(*gtk.FileChooserButton)
// 	projectNameEntry := m.builder.GetObject("setupDialog_programNameEntry").(*gtk.Entry)
// 	dialog.SetTitle("Setup dialog")
// 	dialog.SetTransientFor(m.window)
// 	dialog.SetModal(true)
// 	dialog.AddButton("Save", gtk.RESPONSE_ACCEPT)
// 	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
//
// 	// Show the dialog
// 	responseId := dialog.Run()
// 	defer dialog.Hide()
//
// 	switch responseId {
// 	case gtk.RESPONSE_ACCEPT:
// 		// Save setup information
// 		p := locationFileChooser.GetFilename()
// 		n, err := projectNameEntry.GetText()
// 		if err != nil {
// 			return nil, err
// 		}
// 		result = &setup{name: n, path: p}
// 		return result, nil
// 	default:
// 		return nil, userCancelError
// 	}
// }
