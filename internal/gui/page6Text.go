package gui

type pageText struct {
	parent *MainWindow
}

func (m *MainWindow) setupTextPage() {
	m.textPage = &pageText{parent: m}
}
