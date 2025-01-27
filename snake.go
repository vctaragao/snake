package snake

import "image"

type (
	Snake struct {
		Body []image.Rectangle
		Dir  Direction
	}
)

func NewSnake(x0, y0 int) Snake {
	return Snake{
		Body: []image.Rectangle{
			image.Rect(x0, y0, x0+Size, y0+Size),
			image.Rect(x0-Size, y0, x0, y0+Size),
			image.Rect(x0-2*Size, y0, x0-Size, y0+Size),
		},
		Dir: Right,
	}
}

func (s *Snake) Move() {
	for i := len(s.Body) - 1; i > 0; i-- {
		s.Body[i] = s.Body[i-1]
	}

	p := image.Point{10, 0}

	switch s.Dir {
	case Up:
		p = image.Point{0, -10}
	case Down:
		p = image.Point{0, 10}
	case Left:
		p = image.Point{-10, 0}
	}

	s.Body[0] = s.Body[0].Add(p)
}

func (s *Snake) head() image.Rectangle {
	return s.Body[0]
}

func (s *Snake) grow() {
	s.Body = append(s.Body, s.Body[len(s.Body)-1])
}
