package main

import (
	"image"

	"golang.org/x/image/draw"
)

func editImage(img *image.RGBA, point image.Point) {
	if img.At(point.X, point.Y) == DEAD {
		img.Set(point.X, point.Y, ALIVE)
	} else {
		img.Set(point.X, point.Y, DEAD)
	}
}

func runImage(img *image.RGBA) {
	nextImage := copyImage(img)
	for y := 0; y < IMAGE_HEIGHT; y++ {
		for x := 0; x < IMAGE_WIDTH; x++ {
			point := image.Point{X: x, Y: y}
			nbNeighbors := countNeighbors(img, point)
			if img.At(x, y) == ALIVE && (nbNeighbors < 2 || nbNeighbors > 3) {
				killCell(nextImage, point)
			} else if img.At(x, y) == DEAD && (nbNeighbors == 3) {
				aliveCell(nextImage, point)
			}
		}
	}

	*img = *nextImage
}

func countNeighbors(img *image.RGBA, point image.Point) int {
	var count int
	for y := -1; y < 2; y++ {
		for x := -1; x < 2; x++ {
			if x == 0 && y == 0 { // self compare
				continue
			}

			if img.At(x + point.X, y + point.Y) == ALIVE {
				count++
			}
		}
	}

	return count
}

func aliveCell(img *image.RGBA, point image.Point) {
	img.Set(point.X, point.Y, ALIVE)
}

func killCell(img *image.RGBA, point image.Point) {
	img.Set(point.X, point.Y, DEAD)
}

func copyImage(src *image.RGBA) *image.RGBA {
	rectangle := image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT)
	dst := image.NewRGBA(rectangle)
	draw.Draw(dst, rectangle, src, image.Point{0, 0}, draw.Src)

	return dst
}