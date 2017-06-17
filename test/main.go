package main

import (
	"image"
	"log"
	"net/http"
	"os"

	_ "image/gif"
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
	log.Println("started")
	res, err := http.Get("https://avatars3.githubusercontent.com/u/16108486?v=3&s=460")
	handleErr(err)
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	handleErr(err)

	out, err := os.Create("output.txt")
	handleErr(err)
	defer out.Close()

	textify.NewEncoder(out).Encode(img, nil)
}
