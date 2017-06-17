package textify

import (
	"bytes"
	"image"
	"strings"
)

// Text palettes
var (
	//Palette characters in order from lightest to darkest, used in converting images.
	NormalPalette = strings.Split("█▓▒░. ", "")

	//ReversePalette used for light on black screens.
	ReversePalette = strings.Split(" .░▒▓█", "")
)

// Options contains optional parameters for converting an image to text
type Options struct {
	Palette []string

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
	StrideW uint
	StrideH uint
}

// NewOptions Returns default option parameters
func NewOptions() *Options {
	return &Options{
		Palette: NormalPalette,
		StrideW: 1,
		StrideH: 2,
		Resize:  false,
	}
}

// Encode encodes an image to text and returns.
//		opts: Optional parameters. Leave nil for default.
func Encode(img image.Image, opts *Options) string {
	var out bytes.Buffer
	NewEncoder(&out).Encode(img, opts)
	return string(out.Bytes())
}
