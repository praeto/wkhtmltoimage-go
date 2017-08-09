package wkhtmltoimage

import (
	"fmt"
	"os"
	"testing"
	"github.com/praeto/wkhtmltoimage-go/api"
)

func TestConverter_Output(t *testing.T) {
	gs := api.NewGlobalSettings()
	gs.Set("fmt", "png")
	gs.Set("encoding", "UTF-8")
	gs.Set("crop-h", "630")
	gs.Set("crop-w", "1200")
	gs.Set("crop-x", "0")
	gs.Set("crop-y", "0")
	gs.Set("height", "630")
	gs.Set("width", "1200")
	//gs.Set("quiet", "")
	gs.Set("out", "")
	//gs.Set("in", "http://google.com")
	//gs.Set("transparent", "true")

	c := gs.NewConverter("<html><head><meta charset=\"UTF-8\"></head><body><h1>Есть тут кто?</h1></body></html>")
	defer c.Destroy()

	c.ProgressChanged = func(c *api.Converter, b int) {
		fmt.Printf("Progress: %d\n", b)
	}
	c.Error = func(c *api.Converter, msg string) {
		fmt.Printf("error: %s\n", msg)
	}
	c.Warning = func(c *api.Converter, msg string) {
		fmt.Printf("error: %s\n", msg)
	}
	c.Phase = func(c *api.Converter) {
		fmt.Printf("Phase\n")
	}

	if err := c.Convert(); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Got error code: %d\n", c.ErrorCode())

	lout, outp := c.Output()

	fmt.Printf("Output %d char.s from conversion\n", lout)

	if lout > 0 {
		f, err := os.OpenFile("test.png", os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to open file: %s\n", err)
		}
		defer f.Close()
		f.Truncate(0)
		f.Write([]byte(outp))
	}
}
