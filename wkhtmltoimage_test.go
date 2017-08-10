package wkhtmltoimage

import (
	"fmt"
	"testing"
	"os"
)

func TestFromString(t *testing.T) {
	var (
		err error
		res []byte
	)
	ts := "<html><head><meta charset=\"UTF-8\"></head><body><h1>Есть тут кто?</h1></body></html>"
	config := Config{
		Format: "png",
		Width: 1200,
		Height: 630,
		DisableSmartWidth: true,
		Encoding: "UTF-8",
		//Quiet: true,
	}
	if res, err = FromString(ts, "testString.png", &config); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	if len(res) > 0 {
		fmt.Println("success")
	}
}

func TestFromUrl(t *testing.T) {
	var (
		err error
		res []byte
	)
	config := Config{
		Format: "png",
		Width: 1200,
		Height: 630,
		DisableSmartWidth: true,
		Encoding: "UTF-8",
		//Quiet: true,
	}
	if res, err = FromUrl("http://google.com", "testUrl.png", &config); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	if len(res) > 0 {
		fmt.Println("success")
	}
}

func TestFromFile(t *testing.T) {
	var (
		err error
		res []byte
	)
	curDir, _ := os.Getwd()
	config := Config{
		Format: "png",
		Width: 1200,
		Height: 630,
		DisableSmartWidth: true,
		Encoding: "UTF-8",
		//Quiet: true,
	}
	if res, err = FromFile(curDir + "/testfiles/html.html", "testFile.png", &config); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	if len(res) > 0 {
		fmt.Println("success")
	}
}