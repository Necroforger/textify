# Textify
![Gif](http://i.imgur.com/v9Yz47F.gif)

# Flags

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
| cl        | Crops `n` pixels from the left of the image                                                                                                                           |
| cr        | Crops `n` pixels from the right of the image                                                                                                                          |
| ct        | Crops `n` pixels from the top of the image                                                                                                                            |
| cb        | Crios `n` pixels from the bottom of the image                                                                                                                         |
| cropfirst | Crop the image before resizing `default (false)`                                                                                                                      |
| p         | Set the text palette of the image from darkest to lightest `default (" .░▒▓█")`. The default palette is reversed because terminals are usually light on dark colours. |
| o         | Set the path of the output file. It will default to Stdout if not set                                                                                                 |
