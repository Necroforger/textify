package textify

import (
	"image"
	"image/draw"
	"image/gif"
	"io"
	"strconv"
)

// GifEncoder ...
type GifEncoder struct {
	Dest        io.Writer
	textencoder *Encoder
}

// NewGifEncoder ...
//		dest: Destination io.writer to encode to
func NewGifEncoder(dest io.Writer) *GifEncoder {
	return &GifEncoder{
		Dest:        dest,
		textencoder: NewEncoder(dest),
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

		err := e.EncodeFrame(currentFrame, gi.Delay[i], opts)
		if err != nil {
			return err
		}

		switch gi.Disposal[i] {
		// Clear background to transparent.
		case gif.DisposalNone:
		case gif.DisposalBackground:
			draw.Draw(currentFrame, currentFrame.Bounds(), image.Transparent, image.ZP, draw.Src)
		case gif.DisposalPrevious:
		default:
		}
	}

	return nil
}

// EncodeFrame writes a frame to the destination
func (e *GifEncoder) EncodeFrame(img image.Image, delay int, opts *Options) error {
	if opts == nil {
		opts = NewOptions()
	}

	_, err := e.Dest.Write([]byte(strconv.Itoa(delay) + "\r\n"))
	if err != nil {
		return err
	}

	err = e.textencoder.Encode(img, opts)
	if err != nil {
		return err
	}

	_, err = e.Dest.Write([]byte("\r\n"))
	if err != nil {
		return err
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
