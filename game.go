package snake

import (
	"image"
)

type Game struct {
	width, height int

	snake Snake
	food  Food
}

type (
	DrawFoodFn  func(f Food) error
	DrawSnakeFn func(s Snake) error
)

func NewGame(width, height int) Game {
	x0, y0 := width/2, height/2

	return Game{
		snake:  NewSnake(x0, y0),
		food:   NewFood(width, height),
		width:  width,
		height: height,
	}
}

func (g *Game) DrawFood(fn DrawFoodFn) error {
	return fn(g.food)
}

func (g *Game) DrawSnake(fn DrawSnakeFn) error {
	return fn(g.snake)
}

func (g *Game) MoveSnake() error {
	g.snake.Move()

	if g.bounderyCheck(g.snake.head()) {
		return GameOverErr
	}

	if g.chkFoodColision() {
		g.snake.grow()
		g.food = NewFood(g.width, g.height)
	}

	return nil
}

func (g *Game) bounderyCheck(head image.Rectangle) bool {
	return head.Max.X > g.width ||
		head.Max.Y > g.height ||
		head.Min.X < 0 ||
		head.Min.Y < 0
}

func (g *Game) chkFoodColision() bool {
	return g.snake.head().Eq(image.Rectangle(g.food))
}

func (g *Game) isSnakeGoing(d Direction) bool {
	return g.snake.Dir == d
}

func (g *Game) SetSnakeDirection(d Direction) {
	if g.isInvalidNewDirection(d) {
		return
	}

	g.snake.Dir = d
}

func (g *Game) isInvalidNewDirection(d Direction) bool {
	return (d == Up && g.isSnakeGoing(Down)) ||
		(d == Down && g.isSnakeGoing(Up)) ||
		(d == Left && g.isSnakeGoing(Right)) ||
		(d == Right && g.isSnakeGoing(Left))
}
