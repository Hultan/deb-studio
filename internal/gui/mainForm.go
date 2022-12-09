package gui

import (
	_ "embed"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/applicationVersionConfig"
)

const applicationTitle = "Deb Studio"
const applicationVersion = "v 0.01"
const applicationCopyRight = "©SoftTeam AB, 2022"

const exitCodeSetupError = 1
const exitCodeGtkError = 2

type imageType int
const (
	imageTypeMandatory imageType = iota
	imageTypeRecommended
	imageTypeOptional
)

//go:embed assets/main.glade
var mainGlade string

//go:embed assets/debstudio_256.png
var applicationIcon []byte

//go:embed assets/mandatory.png
var mandatoryIcon []byte

//go:embed assets/recommended.png
var recommendedIcon []byte

//go:embed assets/optional.png
var optionalIcon []byte

//go:embed assets/save.png
var saveIcon []byte

//go:embed assets/exit.png
var exitIcon []byte

type MainForm struct {
	builder *Builder
	window  *gtk.ApplicationWindow
	addFileDialog *addFileDialog
}

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// OpenMainForm : Opens the MainForm window
func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk and create a builder
	gtk.Init(&os.Args)
	m.builder = createBuilder()

	// Main window
	m.window = m.builder.GetObject("mainWindow").(*gtk.ApplicationWindow)
	m.window.SetIcon(createPixBufFromBytes(applicationIcon))
	m.window.SetApplication(app)
	m.window.SetTitle(getApplicationName())
	m.window.SetPosition(gtk.WIN_POS_CENTER)
	m.window.Connect("destroy", m.window.Close)

	// Status bar
	statusBar := m.builder.GetObject("mainWindow_StatusBar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("debstudio"), getApplicationName())

	// AddFileButton
	btn := m.builder.GetObject("addFileButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.addFile)

	// Toolbar & menu
	m.setupToolbar()
	m.setupMenu()

	// Set images on control tab
	setControlImages(m.builder)

	// Show the main window
	m.window.ShowAll()

	load := true
	configPath := getConfigPath()

	if load {
		av, err := applicationVersionConfig.Load(configPath)
		if err != nil {
			_,_ = fmt.Fprintln(os.Stderr,err)
			os.Exit(exitCodeSetupError)
		}
		fmt.Printf("%+v\n", av)

		// Modify the indent level of the ConfigState only.  The global
		// configuration is not modified.
		scs := spew.ConfigState{
			Indent:                  "\t",
			DisableCapacities:       true,
			DisableMethods:          true,
			DisablePointerMethods:   true,
			DisablePointerAddresses: true,
		}
		scs.Dump(av)
	} else {
		av := &applicationVersionConfig.ApplicationVersion{}
		av.Control.Package = "debStudio"
		av.Control.Source = "source"
		av.Control.Version = "1.0.0"
		av.Control.Section = "section"
		av.Control.Priority = "high"
		av.Control.Architecture = "amd64"
		av.Control.Essential = true
		av.Control.Depends = "dpkg"
		av.Control.InstalledSize = "1024"
		av.Control.Maintainer = "Per Hultqvist"
		av.Control.Description = "A deb file creator"
		av.Control.Homepage = "www.softteam.se"
		av.Control.BuiltUsing = "debStudio"
		err := av.Save(configPath)
		if err != nil {
			_,_ = fmt.Fprintln(os.Stderr,err)
			os.Exit(exitCodeSetupError)
		}
	}
}

func (m *MainForm) setupMenu() {
	menuQuit := m.builder.GetObject("menu_FileQuit").(*gtk.MenuItem)
	menuQuit.Connect("activate", m.window.Close)

	menuEditPreferences := m.builder.GetObject("menu_EditPreferences").(*gtk.MenuItem)
	menuEditPreferences.Connect(
		"activate", func() {
			m.openSettingsDialog()
		},
	)

	menuHelpAbout := m.builder.GetObject("menu_HelpAbout").(*gtk.MenuItem)
	menuHelpAbout.Connect(
		"activate", func() {
			m.openAboutDialog()
		},
	)
}

// getConfigPath: Get path to the config file
func getConfigPath() string {
	home := getHomeDirectory()

	return path.Join(home, "code/deb-studio/test.json")
}

// getHomeDirectory: Get current users home directory
func getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		_,_ = fmt.Fprintf(os.Stderr,"Failed to get user home directory : %s\n", err)
		os.Exit(exitCodeSetupError)
	}
	return u.HomeDir
}

// setControlImages: Set control tab images
func setControlImages(b *Builder) {
	setControlImage(b, "imgPackage", imageTypeMandatory)
	setControlImage(b, "imgSource", imageTypeOptional)
	setControlImage(b, "imgVersion", imageTypeMandatory)
	setControlImage(b, "imgSection", imageTypeRecommended)
	setControlImage(b, "imgPriority", imageTypeRecommended)
	setControlImage(b, "imgArchitecture", imageTypeMandatory)
	setControlImage(b, "imgEssential", imageTypeOptional)
	setControlImage(b, "imgDepends", imageTypeOptional)
	setControlImage(b, "imgInstalledSize", imageTypeOptional)
	setControlImage(b, "imgMaintainer", imageTypeMandatory)
	setControlImage(b, "imgDescription", imageTypeMandatory)
	setControlImage(b, "imgHomePage", imageTypeOptional)
	setControlImage(b, "imgBuiltUsing", imageTypeOptional)
}

// setControlImage: Set a control tab image
func setControlImage(b *Builder, imgName string, imgType imageType) {
	bytes := getControlIcon(imgType)
	img := b.GetObject(imgName).(*gtk.Image)
	img.SetFromPixbuf(createPixBufFromBytes(bytes))
}

// getControlIcon: Get the icon bytes from an image type
func getControlIcon(imgType imageType) []byte {
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

// createBuilder: create a gtk builder, and wrap it in a Builder
func createBuilder() *Builder {
	// Create a new builder
	b, err := gtk.BuilderNewFromString(mainGlade)
	if err != nil {
		_,_ = fmt.Fprintln(os.Stderr,err)
		os.Exit(exitCodeSetupError)
	}
	return &Builder{b}
}

// createPixBufFromBytes: Create a gdk.pixbuf from a slice of bytes
func createPixBufFromBytes(bytes []byte) *gdk.Pixbuf {
	pix, err := gdk.PixbufNewFromBytesOnly(bytes)
	if err != nil {
		_,_ = fmt.Fprintln(os.Stderr,err)
		os.Exit(exitCodeSetupError)
	}
	return pix
}

// getApplicationName: Get the application name, version and copyright
func getApplicationName() string {
	return fmt.Sprintf("%s %s - %s", applicationTitle, applicationVersion, applicationCopyRight)
}

func (m *MainForm) addFile() {
	if m.addFileDialog == nil {
		m.addFileDialog = m.newAddFileDialog()
	}
	m.addFileDialog.openForNewFile("/home/per/temp/dragon.ply")
}

func (m *MainForm) setupToolbar() {
	// Toolbar quit button
	btn := m.builder.GetObject("toolbar_quitButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.window.Close)
	btn.SetIsImportant(true)
	btn.SetIconWidget(getImage(exitIcon,"quit"))

	// Toolbar save button
	btn = m.builder.GetObject("toolbar_saveButton").(*gtk.ToolButton)
	btn.Connect("clicked", m.Save)
	btn.SetIsImportant(true)
	btn.SetIconWidget(getImage(saveIcon,"save"))
}

func getImage(bytes []byte, name string) *gtk.Image {
	pix, err := gdk.PixbufNewFromBytesOnly(bytes)
	if err != nil {
		_,_ = fmt.Fprintf(os.Stderr,"failed to create pixbuf for %s icon\n", name)
		_,_ = fmt.Fprintln(os.Stderr,err)
		return nil
	}
	img, err := gtk.ImageNewFromPixbuf(pix)
	if err != nil {
		_,_ = fmt.Fprintf(os.Stderr,"failed to create image for %s icon\n", name)
		_,_ = fmt.Fprintln(os.Stderr,err)
		return nil
	}
	return img
}

func (m *MainForm) Save() {
	// TODO : Save here
}
