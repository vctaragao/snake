package snake

import (
	"image"
	"math/rand"
)

type Food image.Rectangle

func NewFood(width, height int) Food {
	x0, y0 := rand.Intn(width), rand.Intn(height)
	x0 = x0 - x0%Size
	y0 = y0 - y0%Size

	return Food(image.Rect(x0, y0, x0+Size, y0+Size))
}
