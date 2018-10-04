package ui2d

import (
	"Loowootoo/bubble/assets/fonts"
	"Loowootoo/bubble/vector"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

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
	pos               vector.Vector
	dir               vector.Vector
	w, h              int
	exploding         bool
	exploded          bool
	explosionCount    float64
	explosionInterval float64
	explosionTexture  *ebiten.Image
}

func newBubble(tex *ebiten.Image, pos, dir vector.Vector, explosionTexture *ebiten.Image) *bubble {
	w, h := tex.Size()
	return &bubble{tex, pos, dir, w, h, false, false, 0, 50.0, explosionTexture}
}

func (bubble *bubble) getScale() float64 {
	return (bubble.pos.Z/200 + 1) / 2
}

func (bubble *bubble) getCircle() (x, y, r float64) {
	x = bubble.pos.X
	y = bubble.pos.Y - 30*bubble.getScale()
	r = float64(bubble.w) / 2 * bubble.getScale()
	return x, y, r
}

func (bubble *bubble) Draw(screen *ebiten.Image) {
	scale := float64(bubble.getScale())
	newW := int32(float64(bubble.w) * scale)
	newH := int32(float64(bubble.h) * scale)
	x := float64(bubble.pos.X - float64(newW)/2)
	y := float64(bubble.pos.Y - float64(newH)/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
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
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(x, y)
		screen.DrawImage(bubble.explosionTexture, op)
	}
}

func loadBubbles(numBubbles int) []*bubble {
	explosionTexture := loadFromFile("assets/explosion.png")
	bubbleStrs := []string{"assets/balloon_blue.png", "assets/balloon_green.png", "assets/balloon_red.png"}
	bubbleTextures := make([]*ebiten.Image, len(bubbleStrs))

	for i, bstr := range bubbleStrs {
		bubbleTextures[i] = loadFromFile(bstr)
	}
	bubbles := make([]*bubble, numBubbles)
	for i := range bubbles {
		tex := bubbleTextures[i%3]
		pos := vector.Vector{X: rand.Float64() * float64(WinWidth), Y: rand.Float64() * float64(WinHeight), Z: rand.Float64() * float64(WinDepth)}
		dir := vector.Vector{X: rand.Float64()*.5 - .25, Y: rand.Float64()*.5 - .25, Z: rand.Float64() * .25}
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
	text.Draw(screen, str, ui.normalFont, x, y, clr)
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
			xDiff := float64(mouseX) - x
			yDiff := float64(mouseY) - y
			dist := math.Sqrt(float64(xDiff*xDiff + yDiff*yDiff))
			if dist < r {
				bubbleClicked = true
				bubble.exploding = true
				bubble.explosionCount = 0
			}
		}
		p := vector.Add(bubble.pos, vector.Mult(bubble.dir, elapsedTime))
		if p.X < 0 || p.X > float64(WinWidth) {
			bubble.dir.X = -bubble.dir.X
		}
		if p.Y < 0 || p.Y > float64(WinHeight) {
			bubble.dir.Y = -bubble.dir.Y
		}
		if p.Z < 0 || p.Z > float64(WinDepth) {
			bubble.dir.Z = -bubble.dir.Z
		}
		bubble.pos = vector.Add(bubble.pos, vector.Mult(bubble.dir, elapsedTime))
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
