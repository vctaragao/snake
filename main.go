package main

import (
	"errors"
	"image"
	"image/color"
	"log"
	"math/rand"
	"slices"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

const (
	Width  = 600
	Height = 400
)

var (
	GameOverErr = errors.New("Game Over")

	black = color.RGBA{0x00, 0x00, 0x00, 0xaa}
	green = color.RGBA{0x00, 0x7f, 0x00, 0x7f}

	validMovementKeys = []rune{'j', 'k', 'h', 'l'}
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Food image.Rectangle

func NewFood() Food {
	size := 10
	x0, y0 := rand.Intn(Width), rand.Intn(Height)
	x0 = x0 - x0%size
	y0 = y0 - y0%size

	return Food(image.Rect(x0, y0, x0+size, y0+size))
}

func (f Food) Draw(w screen.Window) {
	w.Fill(image.Rect(f.Min.X, f.Min.Y, f.Max.X, f.Max.Y), color.RGBA{0x7f, 0x00, 0x00, 0x7f}, screen.Over)
}

type Snake struct {
	body []image.Rectangle
	dir  Direction
}

func (s *Snake) draw(w screen.Window) {
	for _, r := range s.body {
		w.Fill(r, green, screen.Over)
	}
}

func (s *Snake) move() error {
	for i := len(s.body) - 1; i > 0; i-- {
		s.body[i] = s.body[i-1]
	}

	p := image.Point{10, 0}

	switch s.dir {
	case Up:
		p = image.Point{0, -10}
	case Down:
		p = image.Point{0, 10}
	case Left:
		p = image.Point{-10, 0}
	}

	s.body[0] = s.body[0].Add(p)

	if s.head().Max.X > Width || s.head().Max.Y > Height || s.head().Min.X < 0 || s.head().Min.Y < 0 {
		return GameOverErr
	}

	return nil
}

func (s *Snake) chkFoodColision(food Food) bool {
	return s.head().Eq(image.Rectangle(food))
}

func (s *Snake) head() image.Rectangle {
	return s.body[0]
}

func (s *Snake) grow() {
	s.body = append(s.body, s.body[len(s.body)-1])
}

func main() {
	snake := Snake{
		body: []image.Rectangle{
			image.Rect(20, 0, 30, 10),
			image.Rect(10, 0, 20, 10),
			image.Rect(0, 0, 10, 10),
		},
		dir: Right,
	}

	food := NewFood()

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "Snake Game",
			Width:  Width,
			Height: Height,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		go func() {
			for {
				time.Sleep(time.Second / 8)
				w.Send(paint.Event{})
			}
		}()

		for {
			e := w.NextEvent()

			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			case key.Event:
				if e.Code == key.CodeEscape {
					w.Send(lifecycle.Event{To: lifecycle.StageDead})
					return
				}

				if !slices.Contains(validMovementKeys, e.Rune) {
					break
				}

				snake.handleMovement(e.Rune)

			case paint.Event:
				log.Print("paint event")
				w.Fill(image.Rect(0, 0, Width, Height), black, screen.Over)

				food.Draw(w)
				snake.draw(w)

				if err := snake.move(); err != nil {
					log.Print(err)
					w.Send(lifecycle.Event{To: lifecycle.StageDead})
				}

				if snake.chkFoodColision(food) {
					snake.grow()
					food = NewFood()
				}

			case error:
				log.Print(e)
			}

			w.Publish()
		}
	})
}

func (s *Snake) handleMovement(key rune) {
	switch key {
	case 'j':
		if s.dir == Up {
			return
		}
		s.dir = Down
	case 'k':
		if s.dir == Down {
			return
		}
		s.dir = Up
	case 'h':
		if s.dir == Right {
			return
		}
		s.dir = Left
	case 'l':
		if s.dir == Left {
			return
		}

		s.dir = Right
	}
}
