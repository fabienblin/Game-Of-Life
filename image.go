package main

import (
	"log"
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
					_tappableImage.canvas.Image = GenerateRandomImage(IMAGE_WIDTH, IMAGE_HEIGHT)
					log.Println(point)
				default:
				}
			} else if _ihm.mode == RUN {
				// run mode logic
				_tappableImage.canvas.Image = GenerateRandomImage(IMAGE_WIDTH, IMAGE_HEIGHT)
				time.Sleep(time.Second / time.Duration(_ihm.speed))
			} else {
				continue
			}

			_tappableImage.Refresh()
		}
	}()
}
