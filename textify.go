package textify

import (
	"bufio"
	"bytes"
	"image"
	"io"
	"strings"

	"image/gif"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	rainbow "github.com/raphamorim/go-rainbow"
)

// ColorModes
const (
	ColorDefault = iota
	ColorTerminal
)

// Text palettes
var (
	//Palette characters in order from lightest to darkest, used in converting images.
	PaletteNormal = strings.Split("█▓▒░. ", "")

	//PaletteReverse used for light on black screens.
	PaletteReverse = strings.Split(" .░▒▓█", "")

	//PaletteAsciiNormal is an ascii palette from darkest to lightest
	PaletteASCIIReverse = strings.Split(" .:-=+*#%@", "")

	// PaletteAsciiNormal is an ascii palette from lightest to darkest
	PaletteASCIINormal = strings.Split("@%#*+=-:. ", "")
)

// Options contains optional parameters for converting an image to text
type Options struct {
	Palette []string // Default: PaletteNormal

	// Ouput options

	// OutputMode is an integer representing the output mode to use.
	ColorMode int // Default ColorDefault.

	// Resize will resize the image to the supplied Width and Height dimensions when set to true
	// If one of the width or height values is left as zero, but not both, it will be calculated
	// To preserve the aspect ratio of the image.
	Resize bool // Default: false

	// Thumbnail when set to true will set a maximum resize value and if one of the bounds of the image
	// Exceeds the set value, it will be calculated to fit inside the given bounds while preserving the
	// Original image's aspect ratio
	Thumbnail bool // Default: false

	// Width and height values used by Resize and Thumbnail.
	Width  uint // Default: 0
	Height uint // Default: 0

	// StrideW and StrideH accomodate for the fact that text characters do not have to be entirely square.
	// And will allow you to compensete by setting the stride value. The default values are
	// 		StrideW: 1.
	//		StrideH: 2.
	// Stride H is defaulted to two because text characters usually take up two times the width.
	StrideW float64
	StrideH float64

	// CropFirst defines if the image should be cropped before or after resizing the image
	CropFirst bool // Default: false

	// Values for cropping the image.
	CropLeft, CropRight, CropBottom, CropTop uint // Default: 0
}

// NewOptions Returns default option parameters
func NewOptions() *Options {
	return &Options{
		Palette:   PaletteNormal,
		ColorMode: ColorDefault,
		StrideW:   1,
		StrideH:   2,
		Resize:    false,
		Thumbnail: false,
		Width:     0,
		Height:    0,
		CropFirst: false,
	}
}

// Encode encodes an image to text and returns a string.
//		img:  Image interface to encode.
//		opts: Optional parameters. Leave nil for default.
func Encode(img image.Image, opts *Options) (string, error) {
	var out bytes.Buffer
	err := NewEncoder(&out).Encode(img, opts)
	return string(out.Bytes()), err
}

// EncodeGif returns a GifDecoder from which you can receive frames from
//		gi:   Gif to encode
//		opts: Optional parameters. Leave nil for default.
func EncodeGif(gi *gif.GIF, opts *Options) *GifDecoder {
	rd, wr := io.Pipe()
	go func() {
		NewGifEncoder(wr).Encode(gi, opts)
		wr.Close()
	}()
	reader := bufio.NewReaderSize(rd, 5012)
	return NewGifDecoder(reader)
}

// ColorToText returns a textual representation of the supplied RGB values
// By calculating its brightness and returning the corresponding index in the palette.
//		r: Red value
//		g: Green value
//		b: Blue value
//		palette: Colour palette to use in order from darkest to brightest.
func ColorToText(r, g, b uint32, palette []string) string {
	return palette[int((float32((r+g+b)/3)/65536.0)*float32(len(palette)))]
}

// ColorToColoredTerminalText returns text coloured for a terminal
//		r: Red value
//		g: Green value
//		b: Blue value
//		palette: Colour palette to use in order from darkest to brightest.
func ColorToColoredTerminalText(r, g, b uint32, palette []string) string {
	return rainbow.FromInt32(((r&0xff)<<24)|((g&0xFF)<<16)|(b&0xFF)<<8, palette[int((float32((r+g+b)/3)/65536.0)*float32(len(palette)))])
}

func cropImage(img image.Image, opts *Options) *image.NRGBA {
	var (
		bounds = img.Bounds()
		w      = bounds.Dx()
		h      = bounds.Dy()
	)
	return imaging.Crop(img, image.Rect(int(opts.CropLeft), int(opts.CropTop), w-int(opts.CropLeft), h-int(opts.CropRight)))
}

func resizeImage(img image.Image, opts *Options) image.Image {
	if opts.Resize && (opts.Width != 0 || opts.Height != 0) {
		if opts.Thumbnail {
			return resize.Thumbnail(opts.Width, opts.Height, img, resize.Lanczos3)
		}
		return resize.Resize(opts.Width, opts.Height, img, resize.Lanczos3)
	}

	return img
}
