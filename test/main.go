package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Necroforger/textify"
)

func handleErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	// res, err := http.Get("https://avatars3.githubusercontent.com/u/16108486?v=3&s=460")
	res, err := http.Get("https://68.media.tumblr.com/d91fd6043b15751bf4bcefee41fe59cf/tumblr_nk82xzFkaW1rr5vcmo1_500.gif")
	handleErr(err)
	defer res.Body.Close()

	img, err := gif.DecodeAll(res.Body)
	handleErr(err)

	// out, err := os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	// handleErr(err)
	// defer out.Close()

	options := textify.NewOptions()
	options.Resize = true
	options.Width = 100
	options.Palette = textify.PaletteReverse
	rd, wr := io.Pipe()
	go func() {
		textify.NewGifEncoder(wr).Encode(img, options)
		wr.Close()
	}()

	frames, _, err := textify.NewGifDecoder(rd).DecodeAll()
	handleErr(err)

	for _, v := range frames {
		fmt.Println(v)
	}
}
