package main

import (
	"image/color"
	"math/rand"
	"strings"
	"time"

	"bubble/ui2d"

	"github.com/hajimehoshi/ebiten"
)

var ui *ui2d.UI2d

var (
	waterText = `
明月幾時有，把酒問青天。
不知天上宮闕，今夕是何年。
我欲乘風歸去，又恐瓊樓玉宇，高處不勝寒。
起舞弄清影，何似在人間。
轉朱閣，低綺戶，照無眠。
不應有恨，何事長向別時圓。
人有悲歡離合，月有陰晴圓缺，此事古難全。
但願人長久，千里共嬋娟。	
`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func update(screen *ebiten.Image) error {
	ui.UpdateBubbles(ebiten.CurrentTPS() / 2)

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	ui.DrawBackground(screen)
	for i, line := range strings.Split(waterText, "\n") {
		ui.DrawTextWithShadowCenter(screen, line, 10, 40+i*30, 1, color.White, int(ui2d.WinWidth))
	}
	for _, bubble := range ui.Bubbles {
		bubble.Draw(screen)
	}
	return nil
}

func main() {
	ui = ui2d.NewUI2d()
	err := ebiten.Run(update, int(ui2d.WinWidth), int(ui2d.WinHeight), 1, "Bubble !!!")
	if err != nil {
		panic(err)
	}
}
