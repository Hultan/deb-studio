package gui

import (
	"io"
	"os"
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
	p := &pageScript{parent: m}

	cmb := m.builder.GetObject("mainWindow_scriptCombo").(*gtk.ComboBoxText)
	cmb.Connect("changed", p.scriptChanged)

	p.scriptTextView = m.builder.GetObject("mainWindow_scriptTextView").(*gtk.TextView)
	p.currentScript = scriptTypeNone

	return p
}

func (p *pageScript) init() {
	text, err := p.loadFile(p.getScriptPath(scriptTypePreInstall))
	if err != nil {
		// TODO : Log error
		return
	}
	p.currentScript = scriptTypePreInstall
	p.setScriptText(text)
}

func (p *pageScript) update() {}

func (p *pageScript) scriptChanged(cmb *gtk.ComboBoxText) {
	var err error

	// Save previous script
	content, err := p.getScriptText()
	if err != nil {
		// TODO : Failed to retrieve text, notify user
		return

	}
	err = p.saveFile(p.getScriptPath(p.currentScript), content)
	if err != nil {
		// TODO : Log error, notify user
		return
	}

	// Load next script
	newIndexStr := cmb.GetActiveID()
	newIndex, err := strconv.Atoi(newIndexStr)
	if err != nil {
		// TODO : Log error
		return
	}
	text, err := p.loadFile(p.getScriptPath(scriptType(newIndex)))
	if err != nil {
		// TODO : Log error, notify user
		return
	}
	p.currentScript = scriptType(newIndex)
	p.setScriptText(text)
}

func (p *pageScript) loadFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		// TODO : Log error
		return "", err
	}
	text, err := io.ReadAll(f)
	if err != nil {
		// TODO : Log error
		return "", err
	}

	return string(text), nil
}

func (p *pageScript) saveFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		// TODO : Log error
		return err
	}
	n, err := io.WriteString(f, content)
	if err != nil {
		// TODO : Log error
		return err
	}
	if n != len(content) {
		// TODO : Log error
		return err
	}

	return nil
}

func (p *pageScript) getScriptPath(script scriptType) string {
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
		// TODO : Error
		return ""
	}

	return path.Join(
		project.CurrentPackage.Path,
		project.CurrentPackage.GetPackageFolderName(),
		common.DebianFolderName,
		fileName,
	)
}

func (p *pageScript) setScriptText(text string) {
	buffer, err := p.scriptTextView.GetBuffer()
	if err != nil {
		// TODO : Log error
		return
	}
	buffer.SetText(text)
}

func (p *pageScript) getScriptText() (string, error) {
	buffer, err := p.scriptTextView.GetBuffer()
	if err != nil {
		// TODO : Log error
		return "", err
	}
	text, err := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if err != nil {
		// TODO : Log error
		return "", err
	}
	return text, nil
}
