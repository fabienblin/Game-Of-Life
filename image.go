package main

// "image"

// refreshImage modifies the underlying RGBA image buffer and updates the canvas image.
func refreshImageRoutine(imageCanvas *tappableImageWidget) {
	go func() {
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
