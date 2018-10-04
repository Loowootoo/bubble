package main

import (
	"Loowootoo/bubble/ui2d"
	"image/color"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
)

var ui *ui2d.UI2d

var (
	sampleText = `
君不見黃河之水天上來，奔流到海不復回？
君不見高堂明鏡悲白髮，朝如青絲暮成雪？
人生得意須盡歡，莫使金樽空對月。
天生我材必有用，千金散盡還復來。
烹羊宰牛且為樂，會須一飲三百杯。
岑夫子，丹丘生，
將進酒，君莫停。
與君歌一曲，請君為我側耳聽。
鐘鼓饌玉不足貴，但願長醉不願醒。
古來聖賢皆寂寞，惟有飲者留其名。
陳王昔時宴平樂，斗酒十千恣歡謔。
主人何為言少錢，徑須沽取對君酌。
五花馬，千金裘，
呼兒將出換美酒，與爾同銷萬古愁！
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
	for i, line := range strings.Split(sampleText, "\n") {
		ui.TextOut(screen, line, 10, 40+i*30, color.White)
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
