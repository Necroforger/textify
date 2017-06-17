package textify

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// GifDecoder ...
type GifDecoder struct {
	Source io.Reader
}

// NewGifDecoder ...
//		source: Reader source to decode image from
func NewGifDecoder(source io.Reader) *GifDecoder {
	return &GifDecoder{Source: source}
}

// DecodeAll decodes encoded gifs from their textual representations into an
// Array of frames and an array of delays for their corresponding frames.
func (d *GifDecoder) DecodeAll() (frames []string, delays []int, err error) {
	for frame, delay, err := d.NextFrame(); err == nil; frame, delay, err = d.NextFrame() {
		frames = append(frames, frame)
		delays = append(delays, delay)
	}
	return
}

// NextFrame retrieves the next frame and delay from the source reader.
func (d *GifDecoder) NextFrame() (frame string, delay int, err error) {
	reader := bufio.NewReader(d.Source)

	var line string
	line, err = reader.ReadString('\n')
	if err != nil {
		return
	}

	if delay, err = strconv.Atoi(strings.TrimRight(line, "\r\n")); err != nil {
		return
	}

	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			return
		}
		if line == "\r\n" {
			break
		}
		frame += line
	}

	return
}
