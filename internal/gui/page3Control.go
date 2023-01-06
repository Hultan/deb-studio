package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type pageControl struct {
	parent *MainWindow

	packageEntry         *gtk.Entry
	versionEntry         *gtk.Entry
	architectureEntry    *gtk.Entry
	sourceEntry          *gtk.Entry
	sectionEntry         *gtk.Entry
	priorityEntry        *gtk.Entry
	essentialCheckButton *gtk.CheckButton
	dependsEntry         *gtk.Entry
	installSizeEntry     *gtk.Entry
	installSizeButton    *gtk.Button
	maintainerEntry      *gtk.Entry
	descriptionTextView  *gtk.TextView
	homePageEntry        *gtk.Entry
	builtUsingEntry      *gtk.Entry
}

// setupControlPage: Set up the control page
func (m *MainWindow) setupControlPage() *pageControl {
	p := &pageControl{parent: m}

	p.packageEntry = m.builder.GetObject("mainWindow_packageEntry").(*gtk.Entry)
	p.packageEntry.Connect("focus-out-event", p.focusOut)
	p.versionEntry = m.builder.GetObject("mainWindow_versionEntry").(*gtk.Entry)
	p.versionEntry.Connect("focus-out-event", p.focusOut)
	p.architectureEntry = m.builder.GetObject("mainWindow_architectureEntry").(*gtk.Entry)
	p.architectureEntry.Connect("focus-out-event", p.focusOut)
	p.sourceEntry = m.builder.GetObject("mainWindow_sourceEntry").(*gtk.Entry)
	p.sourceEntry.Connect("focus-out-event", p.focusOut)
	p.sectionEntry = m.builder.GetObject("mainWindow_sectionEntry").(*gtk.Entry)
	p.sectionEntry.Connect("focus-out-event", p.focusOut)
	p.priorityEntry = m.builder.GetObject("mainWindow_priorityEntry").(*gtk.Entry)
	p.priorityEntry.Connect("focus-out-event", p.focusOut)
	p.essentialCheckButton = m.builder.GetObject("mainWindow_essentialCheckButton").(*gtk.CheckButton)
	p.essentialCheckButton.Connect("focus-out-event", p.focusOut)
	p.dependsEntry = m.builder.GetObject("mainWindow_dependsEntry").(*gtk.Entry)
	p.dependsEntry.Connect("focus-out-event", p.focusOut)
	p.installSizeEntry = m.builder.GetObject("mainWindow_installedSizeEntry").(*gtk.Entry)
	p.installSizeEntry.Connect("focus-out-event", p.focusOut)
	p.installSizeButton = m.builder.GetObject("mainWindow_installedSizeButton").(*gtk.Button)
	p.maintainerEntry = m.builder.GetObject("mainWindow_maintainerEntry").(*gtk.Entry)
	p.maintainerEntry.Connect("focus-out-event", p.focusOut)
	p.descriptionTextView = m.builder.GetObject("mainWindow_descriptionTextView").(*gtk.TextView)
	p.descriptionTextView.Connect("focus-out-event", p.focusOut)
	p.homePageEntry = m.builder.GetObject("mainWindow_homePageEntry").(*gtk.Entry)
	p.homePageEntry.Connect("focus-out-event", p.focusOut)
	p.builtUsingEntry = m.builder.GetObject("mainWindow_builtUsingEntry").(*gtk.Entry)
	p.builtUsingEntry.Connect("focus-out-event", p.focusOut)

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

	return p
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

func (p *pageControl) update() {

}

func (p *pageControl) init() {
	p.packageEntry.SetText(p.GetControlField("Package"))
	p.versionEntry.SetText(p.GetControlField("Version"))
	p.architectureEntry.SetText(p.GetControlField("Architecture"))
	p.sourceEntry.SetText(p.GetControlField("Source"))
	p.sectionEntry.SetText(p.GetControlField("Section"))
	p.priorityEntry.SetText(p.GetControlField("Priority"))
	active := p.GetControlField("Essential") == "yes"
	p.essentialCheckButton.SetActive(active)
	p.dependsEntry.SetText(p.GetControlField("Depends"))
	p.installSizeEntry.SetText(p.GetControlField("Installed-Size"))
	p.maintainerEntry.SetText(p.GetControlField("Maintainer"))
	p.setDescriptionText(p.GetControlField("Description"))
	p.homePageEntry.SetText(p.GetControlField("Homepage"))
	p.builtUsingEntry.SetText(p.GetControlField("Built-using"))
}

func (p *pageControl) focusOut(obj glib.IObject) {
	switch widget := obj.(type) {
	case *gtk.Entry:
		name := p.getEntryName(widget)
		value := p.getEntryText(widget)
		p.SetControlField(name, value)
	case *gtk.TextView:
		// TODO : Replace description with an entry (short description) and a textview (long description)
		value := p.getDescriptionText()
		p.SetControlField("Description", value)
	case *gtk.CheckButton:
		value := widget.GetActive()
		text := "no"
		if value {
			text = "yes"
		}
		p.SetControlField("Essential", text)
	}
}

func (p *pageControl) getEntryText(e *gtk.Entry) string {
	text, err := e.GetText()
	if err != nil {
		// TODO : Log error
		return ""
	}
	return text
}

func (p *pageControl) getEntryName(e *gtk.Entry) string {
	name, err := e.GetName()
	if err != nil {
		// TODO : Log error
		return ""
	}
	return name
}

func (p *pageControl) getDescriptionText() string {
	buffer, err := p.descriptionTextView.GetBuffer()
	if err != nil {
		// TODO : Log error
		return ""
	}
	text, err := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if err != nil {
		// TODO : Log error
		return ""
	}
	return text
}

func (p *pageControl) SetControlField(name, value string) {
	project.CurrentPackage.Source.Set(name, value)
	err := project.CurrentPackage.SaveControlFile()
	if err != nil {
		// TODO : Log error
		return
	}
	// fmt.Printf("Set '%s' to '%s'\n", name, value)
}

func (p *pageControl) GetControlField(name string) string {
	value := project.CurrentPackage.Source.Get(name)
	fmt.Printf("Set '%s' to '%s'\n", name, value)
	return value
}

func (p *pageControl) setDescriptionText(text string) {
	buffer, err := p.descriptionTextView.GetBuffer()
	if err != nil {
		// TODO : Log error
		return
	}
	buffer.SetText(text)
}
