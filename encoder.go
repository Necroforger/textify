package textify

import (
	"image"
	"io"

	"github.com/nfnt/resize"
)

// Encoder encodes an image as a stream
type Encoder struct {
	Dest io.Writer
}

// NewEncoder returns an encoder for the given destination
//		dest: writer to write the data to.
func NewEncoder(dest io.Writer) *Encoder {
	return &Encoder{Dest: dest}
}

// Encode writes an image as text the the encoders destination writer.
//		img:  Image to encode to text.
//		Opts: Optional encoding parameters. Leave nil to use default encoding options.
func (e Encoder) Encode(img image.Image, opts *Options) error {
	var (
		bounds = img.Bounds()
		w      = bounds.Dx()
		h      = bounds.Dy()
	)
	if opts == nil {
		opts = NewOptions()
	}

	// Crop and resize images in the requested order.
	if opts.CropFirst {
		img = cropImage(img, opts)
		if opts.Resize {
			img = resizeImage(img, opts)
		}
	} else {
		if opts.Resize {
			img = resizeImage(img, opts)
		}
		img = cropImage(img, opts)
	}

	// Resize the image to accomodate for character size.
	if opts.StrideH > 1 || opts.StrideW > 1 {
		bounds = img.Bounds()
		img = resize.Resize(uint(float64(bounds.Dx())/opts.StrideW), uint(float64(bounds.Dy())/opts.StrideH), img, resize.Lanczos3)
	}

	bounds = img.Bounds()
	w, h = bounds.Dx(), bounds.Dy()

	var (
		err     error
		r, g, b uint32
	)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ = img.At(x, y).RGBA()
			switch opts.ColorMode {
			case ColorTerminal:
				_, err = e.Dest.Write([]byte(ColorToColoredTerminalText(r, g, b, opts.Palette)))
			default:
				_, err = e.Dest.Write([]byte(ColorToText(r, g, b, opts.Palette)))
			}
			if err != nil {
				return err
			}
		}
		_, err := e.Dest.Write([]byte("\r\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
