package wkhtmltoimage

import (
	"github.com/praeto/wkhtmltoimage-go/api"
	"log"
	"strconv"
	"fmt"
)

const (
	typeURL    = iota
	typeFile   = iota
	typeString = iota
)

type Config struct {
	// Format is the type of image to generate
	// jpg, png, svg, bmp supported. Defaults to local wkhtmltoimage default
	Format string
	// Height is the height of the screen used to render in pixels.
	// Default is calculated from page content. Default 0 (renders entire page top to bottom)
	Height int
	// Set screen width, note that this is used
	// Only as a guide line. Use disable-smart-width to make it strict.
	// Default 1024
	Width int
	// Set height for cropping
	CropH int
	// Set width for cropping
	CropW int
	// Set x coordinate for cropping
	CropX int
	// Set y coordinate for cropping
	CropY int
	// Quality determines the final image quality.
	// Values supported between 1 and 100. Default is 94
	Quality int
	// Make the background transparent in pngs
	// By default it is false
	Transparent bool
	// Be less verbose
	// By default it is false
	Quiet bool
	// true - Use the specified width even if it is not large enough for the content
	// false - Extend width to fit unbreakable content
	// By default it is false
	DisableSmartWidth bool
	// Set the text encoding, for input
	// By default is UTF-8
	Encoding string
}

// Convert URL file to IMG file
//
// params:
// url: URL to be saved
// output: path to output image file
// config: instance of wkhtmltoimage.Config()
//
// return:
// []byte:
// 		1. output name if given
//		2. bytes of result image if output param is empty string
// error: nil when success
func FromUrl(url, output string, config *Config) ([]byte, error) {
	return convert(url, output, typeURL, config)
}

// Convert HTML file to IMG file
//
// params:
// filename: path of HTML file with paths or file
// output: path to output image file
// config: instance of wkhtmltoimage.Config()
//
// return:
// []byte:
// 		1. output name if given
//		2. bytes of result image if output param is empty string
// error: nil when success
func FromFile(filename, output string, config *Config) ([]byte, error) {

	return convert(filename, output, typeFile, config)
}

// Convert given string file to IMG file
//
// params:
// string: string to be converted
// output: path to output image file
// config: instance of wkhtmltoimage.Config()
//
// return:
// []byte:
// 		1. output name if given
//		2. bytes of result image if output param is empty string
// error: nil when success
func FromString(string, output string, config *Config) ([]byte, error) {
	return convert(string, output, typeString, config)
}

func convert(input, output string, dataType int, config *Config) ([]byte, error) {
	var html string
	gs := api.NewGlobalSettings()

	if dataType < typeString {
		gs.Set("in", input)
	} else {
		html = input
	}

	if output != "" {
		gs.Set("out", output)
	}

	if config.Format != "" {
		gs.Set("fmt", config.Format)
	}

	if config.Encoding != "" {
		gs.Set("encoding", config.Encoding)
	} else {
		gs.Set("encoding", "UTF-8")
	}

	if config.Height > 0 {
		gs.Set("screenHeight", strconv.Itoa(config.Height))
	}

	if config.Width > 0 {
		gs.Set("screenWidth", strconv.Itoa(config.Width))
	}

	if config.CropH > 0 {
		gs.Set("crop-h", strconv.Itoa(config.CropH))
	}

	if config.CropW > 0 {
		gs.Set("crop-w", strconv.Itoa(config.CropW))
	}

	if config.CropX > 0 {
		gs.Set("crop-x", strconv.Itoa(config.CropX))
	}

	if config.CropY > 0 {
		gs.Set("crop-y", strconv.Itoa(config.CropY))
	}

	if config.Quality > 0 {
		gs.Set("quality", strconv.Itoa(config.Quality))
	}

	if config.Transparent {
		gs.Set("transparent", "true")
	}

	if config.DisableSmartWidth {
		gs.Set("smartWidth", "false")
	}

	c := gs.NewConverter(html, config.Quiet)
	defer c.Destroy()

	c.ProgressChanged = conversionProgressChanged
	c.Error = conversionError
	c.Warning = conversionWarning
	c.Phase = conversionPhase
	c.Finished = conversionFinished

	if errCode := c.Convert(); errCode > 0 {
		return nil, fmt.Errorf("Conversion failed with status; %d", errCode)
	}

	if lout, outp := c.Output(); output != "" {
		return []byte(output), nil
	} else {
		if !config.Quiet {
			log.Printf("Output %d char's from conversion\n", lout)
		}
		return []byte(outp), nil
	}
}

func conversionProgressChanged(c *api.Converter, b int) {
	log.Printf("Progress: %d\n", b)
}

func conversionError(c *api.Converter, msg string) {
	log.Printf("error: %s\n", msg)
}

func conversionWarning(c *api.Converter, msg string) {
	log.Printf("warning: %s\n", msg)
}

func conversionPhase(c *api.Converter) {
	log.Printf("Phase\n")
}

func conversionFinished(c *api.Converter, b int) {
	log.Printf("Finished: %d\n", b)
}
