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
	"os/exec"
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
	IsVideo = flag.Bool("v", false, "Encode from video")
	UseYTDL = flag.Bool("yt", false, "Use youtube-dl to download the video")
	FPS     = flag.Float64("fps", 10, "Video fps")

	PlayVideo = flag.Bool("pv", false, "Play the supplied video")
	PlayGif   = flag.Bool("pg", false, "Play the supplied gif image")
	PlayAudio = flag.Bool("pa", false, "Play audio using ffplay")
	NoLoop    = flag.Bool("nl", false, "Will not loop gifs when playing them")

	CropLeft   = flag.Uint("cl", 0, "Crop left parameter")
	CropRight  = flag.Uint("cr", 0, "Crop right parameter")
	CropTop    = flag.Uint("ct", 0, "Crop top parameter")
	CropBottom = flag.Uint("cb", 0, "Crop bottom parameter")
	CropFirst  = flag.Bool("cropfirst", false, "Crop first parameter")

	Palette    = flag.String("p", strings.Join(textify.PaletteReverse, ""), "Palette parameter")
	OutputPath = flag.String("o", "", "File output path parameter, If not set, will be set to stdout")

	Dest        io.Writer
	Source      io.Reader
	AudioSource io.Reader

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
	///////////////////
	//	Video
	//////////////////
	if *IsVideo || *PlayVideo {

		input := path
		if *UseYTDL {
			input = "pipe:0"
		}

		ff := exec.Command("ffmpeg", "-i", input, "-vf", "fps="+fmt.Sprint(*FPS), "-c:v", "bmp", "-f", "rawvideo", "pipe:1")
		out, err := ff.StdoutPipe()
		if err != nil {
			log.Println(err)
			return
		}

		Source = out

		if *UseYTDL {
			yt := exec.Command("youtube-dl", "-f", "best", "-o", "-", path)
			ytout, err := yt.StdoutPipe()
			if err != nil {
				log.Println(err)
				return
			}

			var videoStream io.Reader
			if *PlayAudio {
				sourceR, sourceW := io.Pipe()
				audioR, audioW := io.Pipe()

				go func() {
					writer := bufio.NewWriterSize(io.MultiWriter(sourceW, audioW), 102400)
					defer writer.Flush()
					_, err := io.Copy(writer, ytout)
					if err != nil {
						log.Println(err)
						os.Exit(1)
					}
				}()

				AudioSource = audioR
				videoStream = sourceR
			} else {
				videoStream = ytout
			}

			ff.Stdin = videoStream
			err = yt.Start()
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			if *PlayAudio {
				aout, err := os.Open(path)
				if err != nil {
					log.Println(err)
					return
				}
				AudioSource = aout
			}
		}

		err = ff.Start()
		if err != nil {
			log.Println(err)
			return
		}

		/////////////////////////
		// HTTP
		/////////////////////////
	} else if strings.HasPrefix(path, "http://") && strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		Source = resp.Body

		/////////////////////////
		// File
		////////////////////////
	} else {
		f, err := os.Open(path)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		Source = f
	}

	//////////////////////////////////////////////
	//          Decode image
	//////////////////////////////////////////////
	var (
		img    image.Image
		gifimg *gif.GIF
	)

	if !*PlayVideo && !*IsVideo {
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
	}

	////////////////////////////////////////
	//       Create Destination writer
	///////////////////////////////////////
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

	//////////////////////////////////////////////
	//           Play Gif
	//////////////////////////////////////////////

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
			// Display frames at correct time.
			if d, u := (time.Millisecond * 10 * time.Duration(delay)), time.Now().Sub(lastWrite); u < d {
				time.Sleep(d - u)
			}
			if !*NoLoop {
				frames = append(frames, Frame{Text: frame, Delay: delay})
			}
			fmt.Fprintln(Dest, frame)
			lastWrite = time.Now()
		}

		if *NoLoop {
			return
		}

		for {
			for _, f := range frames {
				fmt.Fprintln(Dest, f.Text)
				time.Sleep(time.Millisecond * 10 * time.Duration(f.Delay))
			}
		}
	}

	///////////////////////////////////////////
	//            Play video
	///////////////////////////////////////////
	if *PlayVideo {
		if *PlayAudio {
			ffplay := exec.Command("ffplay", "-nodisp", "-i", "-")
			ffplay.Stdin = AudioSource
			err := ffplay.Start()
			if err != nil {
				log.Println(err)
				return
			}
		}
		for {
			decodeStart := time.Now()
			img, _, err := image.Decode(Source)
			if err != nil {
				return
			}
			text, _ := textify.Encode(img, Options)
			fmt.Fprintln(Dest, text)
			if d, u := time.Now().Sub(decodeStart), (time.Millisecond * time.Duration((1.0 / *FPS)*1000)); d < u {
				time.Sleep(u - d)
			}
		}
	}

	////////////////////////////////////////
	//    Write to destination
	////////////////////////////////////////
	// Using a buffered writer helps to display the image to the console smoother.
	bufwriter := bufio.NewWriterSize(Dest, 2048)

	if *IsGif {
		err = textify.NewGifEncoder(bufwriter).Encode(gifimg, Options)
		if err != nil {
			log.Println(err)
			return
		}
	} else if *IsVideo {
		for {
			img, _, err := image.Decode(Source)
			if err != nil {
				if err == image.ErrFormat {
					break
				}
				continue
			}
			textify.NewEncoder(bufwriter).Encode(img, Options)
		}
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
