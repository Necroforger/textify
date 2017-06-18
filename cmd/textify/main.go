package main

import (
	"flag"
	"image"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Necroforger/textify"
)

// Flags
var (
	Width     = flag.Uint("w", 0, "Sets the image width parameter")
	Height    = flag.Uint("h", 0, "Sets the image height parameter")
	Resize    = flag.Bool("r", true, "Resizes the image using the width and height parameters")
	Thumbnail = flag.Bool("t", false, "Resizes the image as a thumbnail")
	StrideW   = flag.Float64("sw", 1, "Stride width parameter")
	StrideH   = flag.Float64("sh", 2, "Stride height parameter")

	CropLeft   = flag.Uint("cl", 0, "Crop left parameter")
	CropRight  = flag.Uint("cr", 0, "Crop right parameter")
	CropTop    = flag.Uint("ct", 0, "Crop top parameter")
	CropBottom = flag.Uint("cb", 0, "Crop bottom parameter")
	CropFirst  = flag.Bool("cropfirst", false, "Crop first parameter")

	Palette    = flag.String("p", strings.Join(textify.PaletteReverse, ""), "Palette parameter")
	OutputPath = flag.String("o", "", "File output path parameter, If not set, will be set to stdout")

	Dest   io.Writer
	Source io.Reader

	Options *textify.Options
)

func main() {
	flag.Parse()
	var err error
	Options = textify.NewOptions()
	Options.Width = *Width
	Options.Height = *Height
	Options.Resize = *Resize
	Options.Thumbnail = *Thumbnail
	Options.StrideW = *StrideW
	Options.StrideH = *StrideH
	Options.CropTop = *CropTop
	Options.CropBottom = *CropBottom
	Options.CropLeft = *CropLeft
	Options.CropRight = *CropRight
	Options.CropFirst = *CropFirst
	Options.Palette = strings.Split(*Palette, "")

	// Initialize input stream
	if flag.Arg(0) == "" {
		log.Println("No file path path provided.")
		return
	}

	Source, err = os.Open(flag.Arg(0))
	if err != nil {
		log.Println(err)
		return
	}

	// Decode image
	img, _, err := image.Decode(Source)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize destination stream
	if *OutputPath != "" {
		f, err := os.OpenFile(*OutputPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		Dest = f
	} else {
		Dest = os.Stdout
	}

	err = textify.NewEncoder(Dest).Encode(img, Options)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Millisecond * 300)

}
