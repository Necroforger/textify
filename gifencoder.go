package textify

import (
	"image"
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

// Encode encodes a gif into text
//		gi: The gif to be encoded.
func (e GifEncoder) Encode(gi *gif.GIF, opts *Options) error {
	if opts == nil {
		opts = NewOptions()
	}

	var lastFrame *image.Paletted
	var lastNoDispose *image.Paletted
	for i, v := range gi.Image {
		_, err := e.Dest.Write([]byte(strconv.Itoa(gi.Delay[i]) + "\r\n"))
		if err != nil {
			return err
		}

		// Initialize lastFrame
		if lastFrame == nil {
			lastFrame = v
		} else {
			// Overwrite lastFrame with changed data
			for y := v.Rect.Min.Y; y < v.Rect.Max.Y; y++ {
				for x := v.Rect.Min.X; x < v.Rect.Max.X; x++ {
					clr := v.At(x, y)
					if _, _, _, a := clr.RGBA(); a != 0 {
						lastFrame.Set(x, y, v.At(x, y))
					}
				}
			}
		}

		// Deal with frame disposal
		switch gi.Disposal[i] {

		case gif.DisposalNone:
			lastNoDispose = v

		case gif.DisposalBackground:
			lastFrame = image.NewPaletted(lastFrame.Rect, lastFrame.Palette)
			for y := lastFrame.Rect.Min.Y; y < lastFrame.Rect.Max.Y; y++ {
				for x := lastFrame.Rect.Min.X; x < lastFrame.Rect.Max.X; x++ {
					lastFrame.Set(x, y, lastFrame.Palette[gi.BackgroundIndex])
				}
			}

		case gif.DisposalPrevious:
			if lastNoDispose != nil {
				lastFrame = lastNoDispose
			}

		default:
			//Clear frame for next image
			if i+1 < len(gi.Image) {
				bnds := gi.Image[i+1].Bounds()
				lastFrame = image.NewPaletted(bnds, gi.Image[i+1].Palette)
			}
		}

		var resizedImage image.Image = lastFrame
		if opts.Resize {
			resizedImage = resizeImage(resizedImage, opts)
		}

		if opts.StrideH > 1 || opts.StrideW > 1 {
			bounds := resizedImage.Bounds()
			resizedImage = resize.Resize(uint(float64(bounds.Dx())/opts.StrideW), uint(float64(bounds.Dy())/opts.StrideH), resizedImage, resize.Lanczos3)
		}

		bounds := resizedImage.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := resizedImage.At(x, y).RGBA()
				_, err := e.Dest.Write([]byte(ColorToText(r, g, b, opts.Palette)))
				if err != nil {
					return err
				}
			}
			e.Dest.Write([]byte("\r\n"))
		}
		e.Dest.Write([]byte("\r\n"))
	}

	return nil
}

//func (e GifEncoder) Encode(gi *gif.GIF, opt *Options) error {
// 	if opt == nil {
// 		opt = NewOptions()
// 	}
// 	tmpmin, tmpmax := gi.Image[0].Bounds().Min, gi.Image[0].Bounds().Max
// 	lastFrame := image.NewPaletted(image.Rect(tmpmin.X, tmpmin.Y, tmpmax.X, tmpmax.Y), gi.Image[0].Palette)
// 	lastNoDispose := image.NewPaletted(image.Rect(tmpmin.X, tmpmin.Y, tmpmax.X, tmpmax.Y), gi.Image[0].Palette)

// 	for i, img := range gi.Image {
// 		_, err := e.Dest.Write([]byte(strconv.Itoa(gi.Delay[i]) + "\r\n"))
// 		if err != nil {
// 			return err
// 		}

// 		//Obtain image bounds
// 		min, max := img.Bounds().Min, img.Bounds().Max

// 		//Overwrite frame
// 		for y := min.Y; y < max.Y; y++ {
// 			for x := min.X; x < max.X; x++ {
// 				if _, _, _, a := img.At(x, y).RGBA(); a != 0 {
// 					lastFrame.Set(x, y, img.At(x, y))
// 				}
// 			}
// 		}
// 		img = lastFrame

// 		//Disposal
// 		switch gi.Disposal[i] {

// 		//Do not remove pixels from last frame
// 		case gif.DisposalNone:
// 			lastNoDispose = img

// 		//Dispose to background colour.
// 		case gif.DisposalBackground:
// 			if i+1 < len(gi.Image) {
// 				bnds := lastFrame.Bounds()
// 				min, max := bnds.Min, bnds.Max
// 				lastFrame = image.NewPaletted(image.Rect(min.X, min.Y, max.X, max.Y), lastFrame.Palette)
// 				//Set to background colour.

// 				for y := min.Y; y < max.Y; y++ {
// 					for x := min.X; x < max.X; x++ {
// 						lastFrame.Set(x, y, lastFrame.Palette[gi.BackgroundIndex])
// 					}
// 				}
// 			}

// 		//Restore to last non-disposed frame
// 		case gif.DisposalPrevious:
// 			lastFrame = lastNoDispose

// 		//Clear image for next frame
// 		default:
// 			if i+1 < len(gi.Image) {
// 				bnds := gi.Image[i+1].Bounds()
// 				min := bnds.Min
// 				max := bnds.Max
// 				lastFrame = image.NewPaletted(image.Rect(min.X, min.Y, max.X, max.Y), gi.Image[i+1].Palette)
// 			}
// 		}

// 		for y := lastFrame.Rect.Min.Y; y < lastFrame.Rect.Max.Y; y++ {
// 			for x := lastFrame.Rect.Min.X; x < lastFrame.Rect.Max.X; x++ {
// 				r, g, b, _ := lastFrame.At(x, y).RGBA()
// 				_, err := e.Dest.Write([]byte(ColorToText(r, g, b, opt.Palette)))
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}

// 	}

// 	return nil
// }
