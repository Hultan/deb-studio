package gui

type pageText struct {
	parent *MainWindow
}

func (m *MainWindow) setupTextPage() *pageText {
	p := &pageText{parent: m}

	return p
}

func (p *pageText) update() {

}
