package main

import (
	"image"
	"image/color"
	"log"
	"slices"
	"time"

	"github.com/vctaragao/snake"
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
	black = color.RGBA{0x00, 0x00, 0x00, 0xaa}
	green = color.RGBA{0x00, 0x7f, 0x00, 0x7f}

	validMovementKeys = []rune{'j', 'k', 'h', 'l'}
	movKeysDirection  = map[rune]snake.Direction{'j': snake.Down, 'k': snake.Up, 'h': snake.Left, 'l': snake.Right}
)

func main() {
	game := snake.NewGame(Width, Height)

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
				time.Sleep(time.Second / 10)
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

				game.SetSnakeDirection(movKeysDirection[e.Rune])
			case paint.Event:
				w.Fill(image.Rect(0, 0, Width, Height), black, screen.Src)

				drawFood(game, w)
				drawSnake(game, w)

				if err := game.MoveSnake(); err != nil {
					log.Print(err)
					w.Send(lifecycle.Event{To: lifecycle.StageDead})
				}
			case error:
				log.Print(e)
			}

			w.Publish()
		}
	})
}
func drawFood(g snake.Game, w screen.Window) {
	g.DrawFood(func(f snake.Food) error {
		w.Fill(image.Rect(f.Min.X, f.Min.Y, f.Max.X, f.Max.Y), color.RGBA{0x7f, 0x00, 0x00, 0x7f}, screen.Over)
		return nil
	})
}

func drawSnake(g snake.Game, w screen.Window) {
	g.DrawSnake(func(s snake.Snake) error {
		for _, r := range s.Body {
			w.Fill(r, green, screen.Over)
		}
		return nil
	})
}
