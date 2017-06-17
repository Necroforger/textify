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
func (e Encoder) Encode(img image.Image, opts *Options) {
	var (
		bounds = img.Bounds()
		w      = bounds.Dx()
		h      = bounds.Dy()
	)
	if opts == nil {
		opts = NewOptions()
	}

	if opts.Resize && (opts.Width != 0 && opts.Height != 0) {
		if opts.Thumbnail {
			img = resize.Thumbnail(opts.Width/opts.StrideW, opts.Height/opts.StrideH, img, resize.Lanczos3)
		} else {
			img = resize.Resize(opts.Width/opts.StrideW, opts.Height/opts.StrideH, img, resize.Lanczos3)
		}
	} else {
		if opts.Thumbnail {
			img = resize.Thumbnail(uint(w)/opts.StrideW, uint(h)/opts.StrideH, img, resize.Lanczos3)
		} else {
			img = resize.Resize(uint(w)/opts.StrideW, uint(h)/opts.StrideH, img, resize.Lanczos3)
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			e.Dest.Write([]byte(ColorToText(r, g, b, opts.Palette)))
		}
		e.Dest.Write([]byte("\r\n"))
	}

}

// ColorToText returns a textual representation of the supplied RGB values
// By calculating its brightness and returning the corresponding index in the palette.
//		r: Red value
//		g: Green value
//		b: Blue value
//		palette: Colour palette to use in order from darkest to brightest.
func ColorToText(r, g, b uint32, palette []string) string {
	return palette[int((float32((r+g+b)/3)/65535.0)*float32(len(palette)-1))]
}
