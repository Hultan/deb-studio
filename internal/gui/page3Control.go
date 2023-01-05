package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

type pageControl struct {
	parent *MainWindow
}

// setupControlPage: Set up the control page
func (m *MainWindow) setupControlPage() {
	p := &pageControl{parent: m}
	m.controlPage = p

	p.setControlImage("imgPackage", imageTypeMandatory)
	p.setControlImage("imgSource", imageTypeOptional)
	p.setControlImage("imgVersion", imageTypeMandatory)
	p.setControlImage("imgSection", imageTypeRecommended)
	p.setControlImage("imgPriority", imageTypeRecommended)
	p.setControlImage("imgArchitecture", imageTypeMandatory)
	p.setControlImage("imgEssential", imageTypeOptional)
	p.setControlImage("imgDepends", imageTypeOptional)
	p.setControlImage("imgInstalledSize", imageTypeOptional)
	p.setControlImage("imgMaintainer", imageTypeMandatory)
	p.setControlImage("imgDescription", imageTypeMandatory)
	p.setControlImage("imgHomePage", imageTypeOptional)
	p.setControlImage("imgBuiltUsing", imageTypeOptional)

}

// setControlImage: Set a control tab image
func (p *pageControl) setControlImage(imgName string, imgType imageType) {
	img := p.parent.builder.GetObject(imgName).(*gtk.Image)
	bytes := p.getControlIcon(imgType)
	img.SetFromPixbuf(createPixBufFromBytes(bytes, imgName))
}

// getControlIcon: Get the icon bytes from an image type
func (p *pageControl) getControlIcon(imgType imageType) []byte {
	var bytes []byte
	switch imgType {
	case imageTypeMandatory:
		bytes = mandatoryIcon
	case imageTypeRecommended:
		bytes = recommendedIcon
	case imageTypeOptional:
		bytes = optionalIcon
	}
	return bytes
}
