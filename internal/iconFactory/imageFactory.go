package iconFactory

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type Image int

const (
	ImageAddFile = iota
	ImageApplication
	ImageBuild
	ImageCheck
	ImageDebStudio
	ImageEdit
	ImageEditFile
	ImageExit
	ImageMandatory
	ImageNew
	ImageOpen
	ImageOptional
	ImageRecommended
	ImageRemoveFile
	ImageSave
)

type ImageFactory struct {
}

func NewImageFactory() *ImageFactory {
	return &ImageFactory{}
}

func (i *ImageFactory) GetImage(image Image) *gtk.Image {
	bytes := i.GetBytes(image)
	pix, err := gdk.PixbufNewFromBytesOnly(bytes)
	if err != nil {
		return nil
	}
	img, err := gtk.ImageNewFromPixbuf(pix)
	if err != nil {
		return nil
	}
	return img
}

func (i *ImageFactory) GetPixBuf(image Image) *gdk.Pixbuf {
	bytes := i.GetBytes(image)
	pix, err := gdk.PixbufNewFromBytesOnly(bytes)
	if err != nil {
		return nil
	}
	return pix
}

func (i *ImageFactory) GetBytes(image Image) []byte {
	switch image {
	case ImageAddFile:
		return addFileIcon
	case ImageApplication:
		return applicationIcon
	case ImageBuild:
		return buildIcon
	case ImageCheck:
		return checkIcon
	case ImageDebStudio:
		return debStudioIcon
	case ImageEdit:
		return editIcon
	case ImageEditFile:
		return editFileIcon
	case ImageExit:
		return exitIcon
	case ImageMandatory:
		return mandatoryIcon
	case ImageNew:
		return newIcon
	case ImageOpen:
		return openIcon
	case ImageOptional:
		return optionalIcon
	case ImageRecommended:
		return recommendedIcon
	case ImageRemoveFile:
		return removeFileIcon
	case ImageSave:
		return saveIcon
	}
	return nil
}
