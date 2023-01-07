package builder

import (
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
	"github.com/hultan/deb-studio/internal/logger"
)

type Builder struct {
	builder *gtk.Builder
}

var log *logger.Logger

// NewBuilder creates a gtk.Builder, and wrap it in a Builder struct
func NewBuilder(l *logger.Logger, glade string) *Builder {
	log = l

	// Create a new builder
	b, err := gtk.BuilderNewFromString(glade)
	if err != nil {
		log.Error.Printf("failed to create builder: %d", err)
		os.Exit(common.ExitCodeGtkError)
	}
	return &Builder{b}
}

// GetObject : Gets a gtk object by name
func (b *Builder) GetObject(name string) glib.IObject {
	if b.builder == nil {
		log.Error.Printf("gtk builder is not set")
		os.Exit(common.ExitCodeGtkError)
	}
	obj, err := b.builder.GetObject(name)
	if err != nil {
		log.Error.Printf("failed to find object with name='%s': %s", name, err)
		os.Exit(common.ExitCodeGtkError)
	}

	return obj
}
