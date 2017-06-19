package textify

import (
	"image"
	"image/draw"
	"image/gif"
	"io"
	"strconv"

	"github.com/nfnt/resize"
)

// GifEncoder ...
type GifEncoder struct {
	Dest io.Writer
}

// NewGifEncoder ...
//		dest: Destination io.writer to encode to
func NewGifEncoder(dest io.Writer) *GifEncoder {
	return &GifEncoder{
		Dest: dest,
	}
}

// Encode ...
func (e *GifEncoder) Encode(gi *gif.GIF, opts *Options) error {
	if opts == nil {
		opts = NewOptions()
	}

	width, height := getGifDimensions(gi)
	currentFrame := image.NewNRGBA(image.Rect(0, 0, width, height))

	for i, v := range gi.Image {
		draw.Draw(currentFrame, currentFrame.Bounds(), v, image.ZP, draw.Over)

		var resizedImage image.Image = currentFrame
		if opts.Resize {
			resizedImage = resizeImage(resizedImage, opts)
		}

		if opts.StrideH > 1 || opts.StrideW > 1 {
			bounds := resizedImage.Bounds()
			resizedImage = resize.Resize(uint(float64(bounds.Dx())/opts.StrideW), uint(float64(bounds.Dy())/opts.StrideH), resizedImage, resize.Lanczos3)
		}

		switch gi.Disposal[i] {
		// Clear background to transparent.
		case gif.DisposalNone:
		case gif.DisposalBackground:
			draw.Draw(currentFrame, currentFrame.Bounds(), image.Transparent, image.ZP, draw.Src)
		case gif.DisposalPrevious:
		default:
		}

		// Text conversion
		bounds := resizedImage.Bounds()
		_, err := e.Dest.Write([]byte(strconv.Itoa(gi.Delay[i]) + "\r\n"))
		if err != nil {
			return err
		}
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := resizedImage.At(x, y).RGBA()
				_, err := e.Dest.Write([]byte(ColorToText(r, g, b, opts.Palette)))
				if err != nil {
					return err
				}
			}
			_, err := e.Dest.Write([]byte("\r\n"))
			if err != nil {
				return err
			}
		}
		_, err = e.Dest.Write([]byte("\r\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}
