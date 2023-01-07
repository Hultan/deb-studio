package gui

import (
	"fmt"
	"path"
	"strconv"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
)

type pageScript struct {
	parent *MainWindow

	scriptTextView *gtk.TextView
	currentScript  scriptType
}

type scriptType int

const (
	scriptTypeNone scriptType = iota
	scriptTypePreInstall
	scriptTypePostInstall
	scriptTypePreRemove
	scriptTypePostRemove
)

func (m *MainWindow) setupScriptPage() *pageScript {
	log.Trace.Println("Entering setupScriptPage...")
	defer log.Trace.Println("Exiting setupScriptPage...")

	p := &pageScript{parent: m}

	cmb := m.builder.GetObject("mainWindow_scriptCombo").(*gtk.ComboBoxText)
	cmb.Connect("changed", p.scriptChanged)

	p.scriptTextView = m.builder.GetObject("mainWindow_scriptTextView").(*gtk.TextView)
	p.currentScript = scriptTypeNone

	return p
}

func (p *pageScript) init() {
	log.Trace.Println("Entering init...")
	defer log.Trace.Println("Exiting init...")

	text, err := readTextFile(p.getScriptPath(scriptTypePreInstall))
	if err != nil {
		log.Error.Printf("failed to read text from text file '%s': %s\n", p.getScriptPath(scriptTypePreInstall), err)
		msg := fmt.Sprintf("failed to read text from text file '%s'", p.getScriptPath(scriptTypePreInstall))
		showErrorDialog(msg, err)
		return
	}
	p.currentScript = scriptTypePreInstall
	setTextViewText(p.scriptTextView, text)
}

func (p *pageScript) update() {
	log.Trace.Println("Entering update...")
	defer log.Trace.Println("Exiting update...")
}

func (p *pageScript) scriptChanged(cmb *gtk.ComboBoxText) {
	log.Trace.Println("Entering scriptChanged...")
	defer log.Trace.Println("Exiting scriptChanged...")

	var err error

	// Save previous script
	content, err := getTextViewText(p.scriptTextView)
	if err != nil {
		log.Error.Printf("failed to get text from textview: %s\n", err)
		showErrorDialog("Failed to get text from textview", err)
		return

	}
	err = writeTextFile(p.getScriptPath(p.currentScript), content)
	if err != nil {
		log.Error.Printf("failed to write text to file '%s': %s\n", p.getScriptPath(p.currentScript), err)
		msg := fmt.Sprintf("failed to write text to file %s", p.getScriptPath(p.currentScript))
		showErrorDialog(msg, err)
		return
	}

	// Load next script
	newIndexStr := cmb.GetActiveID()
	newIndex, err := strconv.Atoi(newIndexStr)
	if err != nil {
		log.Error.Printf("failed to get active id: %s\n", err)
		showErrorDialog("failed to get active id", err)
		return
	}
	text, err := readTextFile(p.getScriptPath(scriptType(newIndex)))
	if err != nil {
		log.Error.Printf("failed to read text from text file '%s': %s\n", p.getScriptPath(scriptType(newIndex)), err)
		msg := fmt.Sprintf("failed to read text from text file '%s'", p.getScriptPath(scriptType(newIndex)))
		showErrorDialog(msg, err)
		return
	}
	p.currentScript = scriptType(newIndex)
	setTextViewText(p.scriptTextView, text)
}

func (p *pageScript) getScriptPath(script scriptType) string {
	log.Trace.Println("Entering getScriptPath...")
	defer log.Trace.Println("Exiting getScriptPath...")

	fileName := ""

	switch script {
	case scriptTypePreInstall:
		fileName = "preinst"
	case scriptTypePostInstall:
		fileName = "postinst"
	case scriptTypePreRemove:
		fileName = "prerm"
	case scriptTypePostRemove:
		fileName = "postrm"
	default:
		log.Error.Printf("invalid script type in getScriptPath(): %d\n", int(script))
		return ""
	}

	return path.Join(
		project.CurrentPackage.Path,
		project.CurrentPackage.GetPackageFolderName(),
		common.DebianFolderName,
		fileName,
	)
}
