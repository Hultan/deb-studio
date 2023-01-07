package gui

import (
	"fmt"
	"path"
	"strconv"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
)

type pageText struct {
	parent *MainWindow

	textView    *gtk.TextView
	currentText textType
}

type textType int

const (
	textTypeNone textType = iota
	textTypeCopyRight
	textTypeChangeLog
	textTypeReadMe
)

func (m *MainWindow) setupTextPage() *pageText {
	log.Trace.Println("Entering setupTextPage...")
	defer log.Trace.Println("Exiting setupTextPage...")

	p := &pageText{parent: m}

	cmb := m.builder.GetObject("mainWindow_textCombo").(*gtk.ComboBoxText)
	cmb.Connect("changed", p.textChanged)

	p.textView = m.builder.GetObject("mainWindow_textTextView").(*gtk.TextView)
	p.currentText = textTypeNone

	return p
}

func (p *pageText) update() {
	log.Trace.Println("Entering update...")
	defer log.Trace.Println("Exiting update...")
}

func (p *pageText) init() {
	log.Trace.Println("Entering init...")
	defer log.Trace.Println("Exiting init...")

	textPath := p.getTextPath(textTypeCopyRight)
	text, err := readTextFile(textPath)
	if err != nil {
		log.Warning.Printf("failed to read text from text file '%s': %s\n", textPath, err)
		log.Warning.Printf("creating empty file...\n")
	}
	p.currentText = textTypeCopyRight
	setTextViewText(p.textView, text)
}

func (p *pageText) textChanged(cmb *gtk.ComboBoxText) {
	log.Trace.Println("Entering textChanged...")
	defer log.Trace.Println("Exiting textChanged...")

	var err error

	// Save previous text
	content, err := getTextViewText(p.textView)
	if err != nil {
		log.Error.Printf("failed to get text from textview: %s\n", err)
		showErrorDialog("Failed to get text from textview", err)
		return

	}
	textPath := p.getTextPath(p.currentText)
	err = writeTextFile(textPath, content)
	if err != nil {
		log.Error.Printf("failed to write text to file '%s': %s\n", textPath, err)
		msg := fmt.Sprintf("failed to write text to file %s", textPath)
		showErrorDialog(msg, err)
		return
	}

	// Load next text
	newIndexStr := cmb.GetActiveID()
	newIndex, err := strconv.Atoi(newIndexStr)
	if err != nil {
		log.Error.Printf("failed to get active id: %s\n", err)
		showErrorDialog("failed to get active id", err)
		return
	}
	textPath = p.getTextPath(textType(newIndex))
	text, err := readTextFile(textPath)
	if err != nil {
		log.Warning.Printf("failed to read text from text file '%s': %s\n", textPath, err)
		log.Warning.Printf("creating empty file...\n")
	}
	p.currentText = textType(newIndex)
	setTextViewText(p.textView, text)
}

func (p *pageText) getTextPath(t textType) string {
	log.Trace.Println("Entering getTextPath...")
	defer log.Trace.Println("Exiting getTextPath...")

	fileName := ""

	switch t {
	case textTypeCopyRight:
		fileName = common.FileNameCopyRight
	case textTypeChangeLog:
		fileName = common.FileNameChangeLog
	case textTypeReadMe:
		fileName = common.FileNameReadme
	default:
		log.Error.Printf("invalid text type in getTextPath(): %d\n", int(t))
		return ""
	}

	return path.Join(
		project.CurrentPackage.Path,
		project.CurrentPackage.GetPackageFolderName(),
		common.FolderNameDebian,
		fileName,
	)
}
