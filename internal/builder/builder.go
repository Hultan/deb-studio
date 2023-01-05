package builder

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
)

type Builder struct {
	builder *gtk.Builder
}

// NewBuilder: create a gtk builder, and wrap it in a Builder
func NewBuilder(glade string) *Builder {
	// Create a new builder
	b, err := gtk.BuilderNewFromString(glade)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(common.ExitCodeGtkError)
	}
	return &Builder{b}
}

// GetObject : Gets a gtk object by name
func (b *Builder) GetObject(name string) glib.IObject {
	if b.builder == nil {
		_, _ = fmt.Fprintln(os.Stderr, "gtk builder is not set")
		os.Exit(common.ExitCodeGtkError)
	}
	obj, err := b.builder.GetObject(name)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to find object with name='%s'\n", name)
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(common.ExitCodeGtkError)
	}

	return obj
}
