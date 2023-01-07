package gui

type pageText struct {
	parent *MainWindow
}

func (m *MainWindow) setupTextPage() *pageText {
	log.Trace.Println("Entering setupTextPage...")
	defer log.Trace.Println("Exiting setupTextPage...")

	p := &pageText{parent: m}

	return p
}

func (p *pageText) update() {
	log.Trace.Println("Entering update...")
	defer log.Trace.Println("Exiting update...")
}
