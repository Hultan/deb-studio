package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

// setupControlPage: Set up the control page
func (m *MainWindow) setupControlPage() {
	m.setControlImage("imgPackage", imageTypeMandatory)
	m.setControlImage("imgSource", imageTypeOptional)
	m.setControlImage("imgVersion", imageTypeMandatory)
	m.setControlImage("imgSection", imageTypeRecommended)
	m.setControlImage("imgPriority", imageTypeRecommended)
	m.setControlImage("imgArchitecture", imageTypeMandatory)
	m.setControlImage("imgEssential", imageTypeOptional)
	m.setControlImage("imgDepends", imageTypeOptional)
	m.setControlImage("imgInstalledSize", imageTypeOptional)
	m.setControlImage("imgMaintainer", imageTypeMandatory)
	m.setControlImage("imgDescription", imageTypeMandatory)
	m.setControlImage("imgHomePage", imageTypeOptional)
	m.setControlImage("imgBuiltUsing", imageTypeOptional)
}

// setControlImage: Set a control tab image
func (m *MainWindow) setControlImage(imgName string, imgType imageType) {
	img := m.builder.GetObject(imgName).(*gtk.Image)
	bytes := m.getControlIcon(imgType)
	img.SetFromPixbuf(createPixBufFromBytes(bytes, imgName))
}

// getControlIcon: Get the icon bytes from an image type
func (m *MainWindow) getControlIcon(imgType imageType) []byte {
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
