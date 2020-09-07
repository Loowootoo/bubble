package ui2d

import (
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"

	"bubble/assets/fonts"
	"bubble/vec3"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"math"
	"math/rand"

	"github.com/Loowootoo/go-sprite"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
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
	tex, err := ebiten.NewImage(w, h, ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
	tex.ReplacePixels(pixels)
	return tex
}

type UI2d struct {
	Bubbles    []*bubble
	Background *ebiten.Image
	normalFont font.Face
	bigFont    font.Face
}

func NewUI2d() *UI2d {
	Bubbles := loadBubbles(20)
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
	return &UI2d{Bubbles, Background, normalFont, bigFont}
}

type bubble struct {
	bubbleSpr  *sprite.Sprite
	pos        vec3.Vec2
	dir        vec3.Vec2
	w, h       float64
	exploding  bool
	exploded   bool
	explodeSpr *sprite.Sprite
}

func newBubble(bubbleSpr *sprite.Sprite, pos, dir vec3.Vec2, explodeSpr *sprite.Sprite) *bubble {
	w := bubbleSpr.GetWidth()
	h := bubbleSpr.GetHeight()
	return &bubble{bubbleSpr, pos, dir, w, h, false, false, explodeSpr}
}

func (bubble *bubble) getScale() float64 {
	return (bubble.pos.Z/200 + 1) / 2
}

func (bubble *bubble) getCircle() (x, y, r float64) {
	x = bubble.pos.X
	y = bubble.pos.Y
	r = bubble.w / 2 * bubble.getScale()
	return x, y, r
}

func (bubble *bubble) Draw(screen *ebiten.Image) {
	scale := float64(bubble.getScale())
	bubble.bubbleSpr.Zoom(scale)
	bubble.bubbleSpr.Position(float64(bubble.pos.X), float64(bubble.pos.Y))
	bubble.bubbleSpr.Draw(screen)
	if bubble.exploding {
		bubble.explodeSpr.Position(float64(bubble.pos.X), float64(bubble.pos.Y))
		bubble.explodeSpr.Draw(screen)
	}
}

func loadBubbles(numBubbles int) []*bubble {
	explodeSpr := sprite.NewSprite()
	explodeSpr.AddAnimation("default", "assets/explosion1.png", 1000, 5, ebiten.FilterDefault)
	explodeSpr.CenterCoordonnates = true
	explodeSpr.Animated = true
	explodeSpr.Start()
	bubbleStrs := []string{"assets/mm_blue.png", "assets/mm_brown.png", "assets/mm_green.png", "assets/mm_orange.png", "assets/mm_purple.png", "assets/mm_red.png", "assets/mm_teal.png", "assets/mm_yellow.png"}
	bubblesprites := make([]*sprite.Sprite, len(bubbleStrs))

	for i, bstr := range bubbleStrs {
		bubblesprites[i] = sprite.NewSprite()
		bubblesprites[i].AddAnimation("default", bstr, 1, 1, ebiten.FilterDefault)
		bubblesprites[i].CenterCoordonnates = true
	}
	bubbles := make([]*bubble, numBubbles)
	for i := range bubbles {
		tex := bubblesprites[i%8]
		pos := vec3.Vec2{X: rand.Float64() * float64(WinWidth), Y: rand.Float64() * float64(WinHeight), Z: rand.Float64() * float64(WinDepth)}
		dir := vec3.Vec2{X: rand.Float64()*.5 - .25, Y: rand.Float64()*.5 - .25, Z: rand.Float64() * .25}
		bubbles[i] = newBubble(tex, pos, dir, explodeSpr)
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
func (ui *UI2d) UpdateBubbles(elapsedTime float64) {

	bubbleClicked := false
	bubbleExploded := false
	for i := len(ui.Bubbles) - 1; i >= 0; i-- {
		bubble := ui.Bubbles[i]
		if bubble.exploding {
			currentAnim := bubble.explodeSpr.Animations[bubble.explodeSpr.CurrentAnimation]
			if currentAnim.CurrentStep+1 >= currentAnim.Steps {
				bubble.exploding = false
				bubble.exploded = true
				bubbleExploded = true
			}
		}
		if !bubbleClicked && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y, r := bubble.getCircle()
			mouseX, mouseY := ebiten.CursorPosition()
			xDiff := float64(mouseX) - x
			yDiff := float64(mouseY) - y
			dist := math.Sqrt(xDiff*xDiff + yDiff*yDiff)
			if dist < r {
				bubbleClicked = true
				bubble.exploding = true
			}
		}
		p := bubble.pos.Add(bubble.dir.Mul2(elapsedTime))
		if p.X < 0 || p.X > float64(WinWidth) {
			bubble.dir.X = -bubble.dir.X
		}
		if p.Y < 0 || p.Y > float64(WinHeight) {
			bubble.dir.Y = -bubble.dir.Y
		}
		if p.Z < 0 || p.Z > float64(WinDepth) {
			bubble.dir.Z = -bubble.dir.Z
		}
		bubble.pos = bubble.pos.Add(bubble.dir.Mul2(elapsedTime))
	}
	if bubbleExploded {
		filteredBubbles := ui.Bubbles[0:0]
		for _, bubble := range ui.Bubbles {
			if !bubble.exploded {
				filteredBubbles = append(filteredBubbles, bubble)
			}
		}
		ui.Bubbles = filteredBubbles
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
