package gui

type pageScript struct {
	parent *MainWindow
}

func (m *MainWindow) setupScriptPage() {
	m.scriptPage = &pageScript{parent: m}
}
