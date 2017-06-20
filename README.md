# Textify
![Gif](http://i.imgur.com/v9Yz47F.gif)
______
<!-- TOC -->

- [Textify](#textify)
- [Command Line](#command-line)
    - [View a gif](#view-a-gif)
    - [View a single image](#view-a-single-image)
    - [Write gif to file](#write-gif-to-file)
    - [Flags](#flags)
- [Options struct](#options-struct)

<!-- /TOC -->

# Command Line
In the following examples, the height parameter is not set so that it will be calculated to preserve the aspect ratio of the image. The URL can be either a local file or an http link.

## View a gif
`textify -w 236 -pg http://i.imgur.com/v9Yz47F.gif`

## View a single image
`textify -w 236 http://i.imgur.com/v9Yz47F.gif`

## Write gif to file
`textify -w 236 -g -o "gif.txt" http://i.imgur.com/v9Yz47F.gif`


## Flags

| Flag      | Description                                                                                                                                                           |
|-----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| w         | Width parameter                                                                                                                                                       |
| h         | Height parameter                                                                                                                                                      |
| r         | Resize to width and height: `Default( True)`                                                                                                                          |
| t         | Resize the image as a thumbnail. Sets the maximum values for width and height but preserves the aspect ratio                                                          |
| sw        | Compensates for the width of a character. `Default (1)`                                                                                                               |
| sh        | Compensates for the height of a character `Default (2)` Text characters are usually twice as tall as they are wide.                                                   |
| g         | Encode every frame of the supplied gif image. If not set to true, it will convert the first frame of the image                                                        |
| pg        | Play the gif to the output destination, which is stdout by default                                                                                                    |
| nl        | No loop, do not loop the gif when playing using pg                                                                                                                    |
| cl        | Crops `n` pixels from the left of the image                                                                                                                           |
| cr        | Crops `n` pixels from the right of the image                                                                                                                          |
| ct        | Crops `n` pixels from the top of the image                                                                                                                            |
| cb        | Crios `n` pixels from the bottom of the image                                                                                                                         |
| cropfirst | Crop the image before resizing `default (false)`                                                                                                                      |
| p         | Set the text palette of the image from darkest to lightest `default (" .░▒▓█")`. The default palette is reversed because terminals are usually light on dark colours. |
| o         | Set the path of the output file. It will default to Stdout if not set                                                                                                 |

# Options struct
```go
// Options contains optional parameters for converting an image to text
type Options struct {
	Palette []string // Default: PaletteNormal

	// Resize will resize the image to the supplied Width and Height dimensions when set to true
	// If one of the width or height values is left as zero, but not both, it will be calculated
	// To preserve the aspect ratio of the image.
	Resize bool // Default: false

	// Thumbnail when set to true will set a maximum resize value and if one of the bounds of the image
	// Exceeds the set value, it will be calculated to fit inside the given bounds while preserving the
	// Original image's aspect ratio
	Thumbnail bool // Default: false

	// Width and height values used by Resize and Thumbnail.
	Width  uint // Default: 0
	Height uint // Default: 0

	// StrideW and StrideH accomodate for the fact that text characters do not have to be entirely square.
	// And will allow you to compensete by setting the stride value. The default values are
	// 		StrideW: 1.
	//		StrideH: 2.
	// Stride H is defaulted to two because text characters usually take up two times the width.
	StrideW float64
	StrideH float64

	// CropFirst defines if the image should be cropped before or after resizing the image
	CropFirst bool // Default: false

	// Values for cropping the image.
	CropLeft, CropRight, CropBottom, CropTop uint // Default: 0
}
```