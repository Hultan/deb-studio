package gui

type pageScript struct {
	parent *MainWindow
}

func (m *MainWindow) setupScriptPage() *pageScript {
	p := &pageScript{parent: m}

	return p
}

func (p *pageScript) update() {

}
