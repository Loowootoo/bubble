package ui2d

import (
	"github.com/Loowootoo/bubble/assets/fonts"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	//	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const WinWidth, WinHeight, WinDepth int32 = 800, 480, 100

func loadFromFile(fileName string) *ebiten.Image {
	inFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	img, err := png.Decode(inFile)
	if err != nil {
		panic(err)
	}
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}
	tex := ebiten.NewImage(w, h)
	if err != nil {
		panic(err)
	}
	tex.ReplacePixels(pixels)
	return tex
}

type UI2d struct {
	Bubbles         []*Sprite
	Bubbleexp       *Sprite
	Background      *ebiten.Image
	normalFont      font.Face
	bigFont         font.Face
	bubbleExploded  bool
	bubbleExploding bool
}

func NewUI2d() *UI2d {
	Bubbles := loadBubbles(20)
	Bubbleexp := loadBubbleexp()
	Background := loadFromFile("assets/moon.png")
	tt, err := truetype.Parse(fonts.Water_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	normalFont := truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	bigFont := truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	return &UI2d{Bubbles, Bubbleexp, Background, normalFont, bigFont, false, false}
}

func loadBubbleexp() *Sprite {
	Spr := NewSprite(float64(WinWidth), float64(WinHeight), float64(WinDepth))
	Spr.AddAnimFrameFromFile("default", "assets/explosion1.png", 1000, 5)
	Spr.CenterCoordonnates = true
	Spr.Animated = true
	Spr.Start()
	return Spr
}

func (ui *UI2d) Draw(screen *ebiten.Image) {
	ui.DrawBackground(screen)
	for i, line := range strings.Split(waterText, "\n") {
		ui.DrawTextWithShadowCenter(screen, line, 10, 40+i*30, 1, color.White, int(WinWidth))
	}
	for i := 0; i < len(ui.Bubbles); i++ {
		ui.Bubbles[i].Draw(screen)
	}
}

func loadBubbles(numBubbles int) []*Sprite {
	bubbleStrs := []string{"assets/mm_blue.png", "assets/mm_brown.png", "assets/mm_green.png", "assets/mm_orange.png", "assets/mm_purple.png", "assets/mm_red.png", "assets/mm_teal.png", "assets/mm_yellow.png"}
	var bstr string
	bubbles := make([]*Sprite, numBubbles)
	for i := 0; i < len(bubbles); i++ {
		bstr = bubbleStrs[i%8]
		bubbles[i] = NewSprite(float64(WinWidth), float64(WinHeight), float64(WinDepth))
		bubbles[i].AddAnimFrameFromFile("default", bstr, 1, 1)
		bubbles[i].CenterCoordonnates = true
		bubbles[i].Pos = Vec3{X: rand.Float64() * float64(WinWidth), Y: rand.Float64() * float64(WinHeight), Z: rand.Float64() * float64(WinDepth)}
		bubbles[i].Direction = Vec3{X: rand.Float64()*.5 - .25, Y: rand.Float64()*.5 - .25, Z: rand.Float64() * .25}
		bubbles[i].Start()
	}
	return bubbles
}
func (ui *UI2d) DrawBackground(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, 1)
	op.GeoM.Translate(0, 0)
	screen.DrawImage(ui.Background, op)
}

func (ui *UI2d) TextOut(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	// Draw the sample text
	text.Draw(screen, str, ui.bigFont, x, y, clr)
}
func (ui *UI2d) Update() {
	for i := len(ui.Bubbles) - 1; i >= 0; i-- {
		ui.Bubbles[i].Update()
	}
}

func (ui *UI2d) textWidth(str string) int {
	b, _ := font.BoundString(ui.normalFont, str)
	return (b.Max.X - b.Min.X).Ceil()
}

var (
	shadowColor  = color.NRGBA{0, 0, 0, 0x80}
	fontBaseSize = 16
)

func (ui *UI2d) DrawTextWithShadow(rt *ebiten.Image, str string, x, y, scale int, clr color.Color) {
	offsetY := fontBaseSize * scale
	for _, line := range strings.Split(str, "\n") {
		y += offsetY
		text.Draw(rt, line, ui.normalFont, x+2, y+2, shadowColor)
		text.Draw(rt, line, ui.normalFont, x, y, clr)
	}
}

func (ui *UI2d) DrawTextWithShadowCenter(rt *ebiten.Image, str string, x, y, scale int, clr color.Color, width int) {
	w := ui.textWidth(str) * scale
	x += (width - w) / 2
	ui.DrawTextWithShadow(rt, str, x, y, scale, clr)
}

func (ui *UI2d) DrawTextWithShadowRight(rt *ebiten.Image, str string, x, y, scale int, clr color.Color, width int) {
	w := ui.textWidth(str) * scale
	x += width - w
	ui.DrawTextWithShadow(rt, str, x, y, scale, clr)
}

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
