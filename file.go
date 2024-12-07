package main

import (
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/widget"
)

const FILE_EXTENSION string = ".gol"

func triggerSaveImage(saveInputWidget *widget.Entry) {
	if saveInputWidget.Text == "" {
		log.Println("File needs a name")
	}
	
	triggerPause()

	f, err := os.Create(saveInputWidget.Text+FILE_EXTENSION)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, _tappableImage.canvas.Image); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func triggerLoadImage(fileName string) {
	triggerPause()

	reader, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	_tappableImage.canvas.Image = img
}

func findGOLimages() []string {
	var fileNames[]string

	filepath.WalkDir("./", func(s string, d fs.DirEntry, e error) error {
	   if e != nil { return e }
	   if filepath.Ext(d.Name()) == FILE_EXTENSION {
		fileNames = append(fileNames, s)
	   }
	   return nil
	})
	return fileNames
}