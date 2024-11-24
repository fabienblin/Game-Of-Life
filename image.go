package main

import (
	// "image"

	"fyne.io/fyne/v2/canvas"
)

// refreshImage modifies the underlying RGBA image buffer and updates the canvas image.
func refreshImageRoutine(imageCanvas *canvas.Image) {
	go func(){
		// image := imageCanvas.Image.(*image.RGBA)

		for {
			if _mode == EDIT {
				// edit mode logic
			} else if _mode == RUN {
				// run mode logic
			} else {
				continue
			}
			imageCanvas.Refresh()
		}
	}()
}
