package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/iconFactory"
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
	log.Trace.Println("Entering setupControlPage...")
	defer log.Trace.Println("Exiting setupControlPage...")

	p := &pageControl{parent: m}

	p.packageEntry = m.builder.GetObject("mainWindow_packageEntry").(*gtk.Entry)
	// secondaryIcon, _ := gdk.PixbufNewFromBytesOnly(optionalIcon)
	// p.packageEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_SECONDARY, secondaryIcon)
	// p.packageEntry.SetIconTooltipText(gtk.ENTRY_ICON_SECONDARY, "Detta är ett test")
	// p.packageEntry.Connect("icon-press", test)
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

	p.setControlImage("imgPackage", iconFactory.ImageMandatory)
	p.setControlImage("imgSource", iconFactory.ImageOptional)
	p.setControlImage("imgVersion", iconFactory.ImageMandatory)
	p.setControlImage("imgSection", iconFactory.ImageRecommended)
	p.setControlImage("imgPriority", iconFactory.ImageRecommended)
	p.setControlImage("imgArchitecture", iconFactory.ImageMandatory)
	p.setControlImage("imgEssential", iconFactory.ImageOptional)
	p.setControlImage("imgDepends", iconFactory.ImageOptional)
	p.setControlImage("imgInstalledSize", iconFactory.ImageOptional)
	p.setControlImage("imgMaintainer", iconFactory.ImageMandatory)
	p.setControlImage("imgDescription", iconFactory.ImageMandatory)
	p.setControlImage("imgHomePage", iconFactory.ImageOptional)
	p.setControlImage("imgBuiltUsing", iconFactory.ImageOptional)

	return p
}

func (p *pageControl) update() {
	log.Trace.Println("Entering update...")
	defer log.Trace.Println("Exiting update...")
}

func (p *pageControl) init() {
	log.Trace.Println("Entering init...")
	defer log.Trace.Println("Exiting init...")

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
	setTextViewText(p.descriptionTextView, p.GetControlField("Description"))
	p.homePageEntry.SetText(p.GetControlField("Homepage"))
	p.builtUsingEntry.SetText(p.GetControlField("Built-using"))
}

func (p *pageControl) focusOut(obj glib.IObject) {
	log.Trace.Println("Entering focusOut...")
	defer log.Trace.Println("Exiting focusOut...")

	switch widget := obj.(type) {
	case *gtk.Entry:
		name, err := getEntryName(widget)
		if err != nil {
			log.Error.Printf("failed to get name from entry (field: %s): %s\n", name, err)
			msg := fmt.Sprintf("Failed to get name from entry (field: %s)", name)
			showErrorDialog(msg, err)
			return
		}
		value, err := getEntryText(widget)
		if err != nil {
			log.Error.Printf("failed to get text from entry (field: %s): %s\n", name, err)
			msg := fmt.Sprintf("Failed to get text from entry (field: %s)", name)
			showErrorDialog(msg, err)
			return
		}
		err = p.SetControlField(name, value)
		if err != nil {
			log.Error.Printf("failed to save to control file (field: %s): %s\n", name, err)
			msg := fmt.Sprintf("Failed to save to control file (field: %s)", name)
			showErrorDialog(msg, err)
			return
		}
	case *gtk.TextView:
		// TODO : Replace description with an entry (short description) and a textview (long description)
		value, err := getTextViewText(p.descriptionTextView)
		if err != nil {
			log.Error.Printf("failed to get text from textview (field: Description): %s\n", err)
			showErrorDialog("Failed to get text from textview (field: Description)", err)
			return
		}
		err = p.SetControlField("Description", value)
		if err != nil {
			log.Error.Printf("failed to save to control file (field: Description): %s\n", err)
			showErrorDialog("Failed to save to control file (field: Description)", err)
			return
		}
	case *gtk.CheckButton:
		value := widget.GetActive()
		text := "no"
		if value {
			text = "yes"
		}
		err := p.SetControlField("Essential", text)
		if err != nil {
			log.Error.Printf("failed to save to control file (field: Essential): %s\n", err)
			showErrorDialog("Failed to save to control file (field: Essential)", err)
			return
		}
	}
}

func (p *pageControl) SetControlField(name, value string) error {
	log.Trace.Println("Entering SetControlField...")
	defer log.Trace.Println("Exiting SetControlField...")

	project.CurrentPackage.Source.Set(name, value)
	err := project.CurrentPackage.SaveControlFile()
	if err != nil {
		log.Error.Printf("failed to set control field (field: %s, value: %s): %s\n", name, value, err)
		return err
	}
	return nil
}

func (p *pageControl) GetControlField(name string) string {
	log.Trace.Println("Entering GetControlField...")
	defer log.Trace.Println("Exiting GetControlField...")

	return project.CurrentPackage.Source.Get(name)
}

// setControlImage: Set a control tab image
func (p *pageControl) setControlImage(imgName string, imgType iconFactory.Image) {
	log.Trace.Println("Entering setControlImage...")
	defer log.Trace.Println("Exiting setControlImage...")

	img := p.parent.builder.GetObject(imgName).(*gtk.Image)
	pix := p.parent.image.GetPixBuf(imgType)
	img.SetFromPixbuf(pix)
}
