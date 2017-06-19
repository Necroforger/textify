package main

import (
	"image"
	"log"
	"net/http"
	"os"

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
	res, err := http.Get("https://avatars3.githubusercontent.com/u/16108486?v=3&s=460")
	handleErr(err)
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	handleErr(err)

	options := textify.NewOptions()
	options.Resize = true
	options.Width = 180
	options.Palette = textify.PaletteReverse

	textify.NewEncoder(os.Stdout).Encode(img, options)
}
