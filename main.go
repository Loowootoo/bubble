package main

import (
	//	"image/color"
	"math/rand"
	//	"strings"
	"time"

	"github.com/Loowootoo/bubble/ui2d"

	"github.com/hajimehoshi/ebiten/v2"
)

var ui *ui2d.UI2d

type Game struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}
func (g *Game) Update() error {
	ui.Update()
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	ui.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 480
}

func main() {
	ui = ui2d.NewUI2d()
	ebiten.SetWindowSize(800, 480)
	ebiten.SetWindowTitle("Bubble")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
