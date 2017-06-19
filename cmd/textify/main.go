package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"io"
	"log"
	"net/http"
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

	IsGif   = flag.Bool("g", false, "Encode the image as a gif")
	PlayGif = flag.Bool("pg", false, "Play the supplied gif image")

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

	path := flag.Arg(0)
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		Source = f
	} else {
		resp, err := http.Get(path)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		Source = resp.Body
	}

	var (
		img    image.Image
		gifimg *gif.GIF
	)

	// Decode as image
	if !*IsGif && !*PlayGif {
		img, _, err = image.Decode(Source)
		if err != nil {
			log.Println(err)
			return
		}

		// Decode as gif
	} else {
		gifimg, err = gif.DecodeAll(Source)
		if err != nil {
			log.Println(err)
			return
		}
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

	if *PlayGif {
		type Frame struct {
			Delay int
			Text  string
		}
		frames := []Frame{}
		decoder := textify.EncodeGif(gifimg, Options)
		var (
			frame string
			err   error
			delay int
		)

		lastWrite := time.Now()
		for frame, delay, err = decoder.NextFrame(); err == nil; frame, delay, err = decoder.NextFrame() {
			// Do not allow the gif to display too fast while rendering.
			if d, u := (time.Millisecond * 10 * time.Duration(delay)), time.Now().Sub(lastWrite); u < d {
				time.Sleep(d - u)
			}
			frames = append(frames, Frame{Text: frame, Delay: delay})
			fmt.Fprintln(Dest, frame)
			lastWrite = time.Now()
		}
		for {
			for _, f := range frames {
				fmt.Fprintln(Dest, f.Text)
				time.Sleep(time.Millisecond * 10 * time.Duration(f.Delay))
			}
		}
	}

	// Using a buffered writer helps to display the image to the console smoother.
	bufwriter := bufio.NewWriterSize(Dest, 2048)

	// Encode gif
	if *IsGif {
		err = textify.NewGifEncoder(bufwriter).Encode(gifimg, Options)
		// Encode image
	} else {
		err = textify.NewEncoder(bufwriter).Encode(img, Options)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Flush the writer to append any data still remaining in the buffer
	err = bufwriter.Flush()
	if err != nil {
		log.Println(err)
		return
	}

}
