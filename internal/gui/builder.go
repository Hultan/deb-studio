package gui

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Builder struct {
	Builder *gtk.Builder
}

// GetObject : Gets a gtk object by name
func (g *Builder) GetObject(name string) glib.IObject {
	if g.Builder == nil {
		_,_ = fmt.Fprintln(os.Stderr,"builder must be set")
		os.Exit(exitCodeSetupError)
	}
	obj, err := g.Builder.GetObject(name)
	if err != nil {
		_,_ = fmt.Fprintln(os.Stderr,err)
		os.Exit(exitCodeSetupError)
	}

	return obj
}
