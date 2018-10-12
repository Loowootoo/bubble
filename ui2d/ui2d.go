package ui2d

import (
	"Loowootoo/bubble/assets/fonts"
	"Loowootoo/bubble/vec3"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

const WinWidth, WinHeight, WinDepth int32 = 800, 600, 100

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
	tex               *ebiten.Image
	pos               vec3.Vector
	dir               vec3.Vector
	w, h              int
	exploding         bool
	exploded          bool
	explosionCount    float64
	explosionInterval float64
	explosionTexture  *ebiten.Image
}

func newBubble(tex *ebiten.Image, pos, dir vec3.Vector, explosionTexture *ebiten.Image) *bubble {
	w, h := tex.Size()
	return &bubble{tex, pos, dir, w, h, false, false, 0, 50.0, explosionTexture}
}

func (bubble *bubble) getScale() float32 {
	return (bubble.pos.Z/200 + 1) / 2
}

func (bubble *bubble) getCircle() (x, y, r float32) {
	x = bubble.pos.X
	y = bubble.pos.Y
	r = float32(bubble.w) / 2 * bubble.getScale()
	return x, y, r
}

func (bubble *bubble) Draw(screen *ebiten.Image) {
	scale := bubble.getScale()
	newW := int32(float32(bubble.w) * scale)
	newH := int32(float32(bubble.h) * scale)
	x := float64(bubble.pos.X - float32(newW)/2)
	y := float64(bubble.pos.Y - float32(newH)/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(scale), float64(scale))
	op.GeoM.Translate(x, y)
	screen.DrawImage(bubble.tex, op)

	if bubble.exploding {
		numAnimations := 16
		animationIndex := numAnimations - 1 - int(bubble.explosionCount/bubble.explosionInterval)
		animationX := animationIndex % 4
		animationY := 64 * ((animationIndex - animationX) / 4)
		animationX *= 64
		rect := image.Rect(int(animationX), int(animationY), animationX+64, animationY+64)
		op.GeoM.Reset()
		op.SourceRect = &rect
		op.GeoM.Scale(float64(scale), float64(scale))
		op.GeoM.Translate(x, y)
		screen.DrawImage(bubble.explosionTexture, op)
	}
}

func loadBubbles(numBubbles int) []*bubble {
	explosionTexture := loadFromFile("assets/explosion.png")
	bubbleStrs := []string{"assets/mm_blue.png", "assets/mm_brown.png", "assets/mm_green.png", "assets/mm_orange.png", "assets/mm_purple.png", "assets/mm_red.png", "assets/mm_teal.png", "assets/mm_yellow.png"}
	bubbleTextures := make([]*ebiten.Image, len(bubbleStrs))

	for i, bstr := range bubbleStrs {
		bubbleTextures[i] = loadFromFile(bstr)
	}
	bubbles := make([]*bubble, numBubbles)
	for i := range bubbles {
		tex := bubbleTextures[i%8]
		pos := vec3.Vector{X: rand.Float32() * float32(WinWidth), Y: rand.Float32() * float32(WinHeight), Z: rand.Float32() * float32(WinDepth)}
		dir := vec3.Vector{X: rand.Float32()*.5 - .25, Y: rand.Float32()*.5 - .25, Z: rand.Float32() * .25}
		bubbles[i] = newBubble(tex, pos, dir, explosionTexture)
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

	numAnimations := 16
	bubbleClicked := false
	bubbleExploded := false
	for i := len(ui.Bubbles) - 1; i >= 0; i-- {
		bubble := ui.Bubbles[i]
		if bubble.exploding {
			bubble.explosionCount += float64(elapsedTime)
			animationIndex := numAnimations - 1 - int(bubble.explosionCount/bubble.explosionInterval)
			if animationIndex < 0 {
				bubble.exploding = false
				bubble.exploded = true
				bubbleExploded = true
			}
		}
		if !bubbleClicked && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y, r := bubble.getCircle()
			mouseX, mouseY := ebiten.CursorPosition()
			xDiff := float32(mouseX) - x
			yDiff := float32(mouseY) - y
			dist := float32(math.Sqrt(float64(xDiff*xDiff + yDiff*yDiff)))
			if dist < r {
				bubbleClicked = true
				bubble.exploding = true
				bubble.explosionCount = 0
			}
		}
		p := bubble.pos.Add(bubble.dir.Mul2(float32(elapsedTime)))
		if p.X < 0 || p.X > float32(WinWidth) {
			bubble.dir.X = -bubble.dir.X
		}
		if p.Y < 0 || p.Y > float32(WinHeight) {
			bubble.dir.Y = -bubble.dir.Y
		}
		if p.Z < 0 || p.Z > float32(WinDepth) {
			bubble.dir.Z = -bubble.dir.Z
		}
		bubble.pos = bubble.pos.Add(bubble.dir.Mul2(float32(elapsedTime)))
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
		text.Draw(rt, line, ui.normalFont, x+1, y+1, shadowColor)
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
