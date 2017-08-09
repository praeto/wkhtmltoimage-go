package wkhtmltoimage

import "github.com/praeto/wkhtmltoimage-go/api"

const (
	TypeURL = iota
	TypeFile = iota
	TypeString = iota
)

type Config struct {
	// Format is the type of image to generate
	// jpg, png, svg, bmp supported. Defaults to local wkhtmltoimage default
	Format string
	// Height is the height of the screen used to render in pixels.
	// Default is calculated from page content. Default 0 (renders entire page top to bottom)
	Height int
	// Width is the width of the screen used to render in pixels.
	// Note that this is used only as a guide line. Default 1024
	Width int
	// Quality determines the final image quality.
	// Values supported between 1 and 100. Default is 94
	Quality int
	// Transparent determines image background.
	// By default it is false
	Transparent bool
	//// Transparent determines image background.
	//// By default it is false
	//Quiet string
}

func FromUrl(url, output string, config *Config) ([]byte, error) {
	return convert(url, output, TypeURL, config)
}

func FromFile(filename, output string, config *Config) ([]byte, error) {
	return convert(filename, output, TypeFile, config)
}

func FromString(string, output string, config *Config) ([]byte, error) {
	return convert(string, output, TypeString, config)
}

func convert(input, output string, dataType int, config *Config) ([]byte, error) {
	var html string
	gs := api.NewGlobalSettings()
	if dataType < TypeString {
		gs.Set("in", input)
	} else {
		html = input
	}

	return nil, nil
}
