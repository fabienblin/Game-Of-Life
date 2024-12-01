package main

import (
	"image"
	"time"
)

// refreshImage modifies the underlying RGBA image buffer and updates the canvas image.
func refreshImageRoutine() {
	go func() {
		// image := _tappableImage.canvas.Image.(*image.RGBA)
		for {
			if _ihm.mode == EDIT {
				// edit mode logic
				select {
				case point := <-_tappableImage.tappChan:
					editImage(_tappableImage.canvas.Image.(*image.RGBA), point)
				default:
				}
			} else if _ihm.mode == RUN {
				// run mode logic
				runImage(_tappableImage.canvas.Image.(*image.RGBA))
				time.Sleep(time.Second / time.Duration(_ihm.speed))
			} else {
				continue
			}

			_tappableImage.Refresh()
		}
	}()
}
