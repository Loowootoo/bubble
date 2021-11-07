package ui2d

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type AnimFrame struct {
	Image        *ebiten.Image
	MaxFrames    int
	CurrFrame    int
	FrameWidth   int
	FrameHeight  int
	FrameCounter *TimeCounter
	RunOnce      bool
}

func newAnimFrameFromFile(fileName string, duration int, frames int) *AnimFrame {
	var err error
	animFrame := new(AnimFrame)
	animFrame.Image, _, err = ebitenutil.NewImageFromFile(fileName)
	if err != nil {
		panic(err)
	}
	animFrame.MaxFrames = frames

	width, height := animFrame.Image.Size()
	animFrame.FrameWidth = width / animFrame.MaxFrames
	animFrame.FrameHeight = height

	animFrame.FrameCounter = NewCounter(duration)
	animFrame.RunOnce = false
	return animFrame
}

func newAnimFrameFromBytes(data []byte, duration int, frames int) *AnimFrame {
	animFrame := new(AnimFrame)
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	animFrame.Image = ebiten.NewImageFromImage(img)

	animFrame.MaxFrames = frames

	width, height := animFrame.Image.Size()
	animFrame.FrameWidth = width / animFrame.MaxFrames
	animFrame.FrameHeight = height

	animFrame.FrameCounter = NewCounter(duration)
	animFrame.RunOnce = false
	return animFrame
}

func (animFrame *AnimFrame) SetFrameDuration(duration int) {
	animFrame.FrameCounter.ResetCounter(duration)
}

type Sprite struct {
	// Animation label currently displayed
	CurrAnimFrame string
	// Array of animations
	AnimFrames         map[string]*AnimFrame
	Pos                Vec3
	Direction          Vec3
	Area               Vec3
	Speed              float64
	Alpha              float64
	Scale              float64
	AreaW              float64
	AreaH              float64
	AreaZ              float64
	Visible            bool
	Animated           bool
	CenterCoordonnates bool
}

func NewSprite(w, h, z float64) *Sprite {
	sprite := new(Sprite)
	sprite.CurrAnimFrame = "default"
	sprite.AnimFrames = make(map[string]*AnimFrame)
	sprite.Alpha = 1
	sprite.Animated = false
	sprite.CenterCoordonnates = true
	sprite.Direction = Vec3{0, 0, 0}
	sprite.Pos = Vec3{0, 0, 0}
	sprite.AreaW = w
	sprite.AreaH = h
	sprite.AreaZ = z
	sprite.Speed = 5
	sprite.Visible = true
	return sprite
}

func (sprite *Sprite) AddAnimFrameFromFile(label string, path string, duration int, steps int) {
	sprite.AnimFrames[label] = newAnimFrameFromFile(path, duration, steps)
}

func (sprite *Sprite) AddAnimFrameFromBytes(label string, data []byte, duration int, steps int) {
	sprite.AnimFrames[label] = newAnimFrameFromBytes(data, duration, steps)
}
func (sprite *Sprite) GetScale() float64 {
	return (sprite.Pos.Z/100 + 1) / 2
}

func (sprite *Sprite) Update() {
	if sprite.Pos.X > sprite.AreaW || sprite.Pos.X < 0 {
		sprite.Direction.X = -sprite.Direction.X
	}
	sprite.Pos.X += sprite.Speed * sprite.Direction.X
	if sprite.Pos.Y > sprite.AreaH || sprite.Pos.Y < 0 {
		sprite.Direction.Y = -sprite.Direction.Y
	}
	sprite.Pos.Y += sprite.Speed * sprite.Direction.Y
	if sprite.Pos.Z <= 0 {
		sprite.Direction.Z = 1
	} else if sprite.Pos.Z >= sprite.AreaZ {
		sprite.Direction.Z = -1
	}
	sprite.Pos.Z += sprite.Speed * sprite.Direction.Z
	sprite.Scale = sprite.GetScale()
}
func (sprite *Sprite) nextFrame() {
	currAnimFrame := sprite.AnimFrames[sprite.CurrAnimFrame]
	if sprite.Animated {
		if currAnimFrame.FrameCounter.TimeUp() {
			currAnimFrame.CurrFrame++
			if currAnimFrame.CurrFrame+1 > currAnimFrame.MaxFrames {
				if currAnimFrame.RunOnce {
					sprite.Stop()
				} else {
					currAnimFrame.CurrFrame = 0
				}
			}
		}
	}
}

//Draw calculates new coordonnates and draw the sprite on the screen, after drawing, go to the next step of animation
func (sprite *Sprite) Draw(surface *ebiten.Image) {
	if sprite.Visible {
		currAnimFrame := sprite.AnimFrames[sprite.CurrAnimFrame]
		options := &ebiten.DrawImageOptions{}

		if sprite.CenterCoordonnates {
			options.GeoM.Translate(-float64(currAnimFrame.FrameWidth)/2, -float64(currAnimFrame.FrameHeight)/2)
		}
		options.GeoM.Scale(sprite.Scale, sprite.Scale)
		options.GeoM.Translate(sprite.Pos.X, sprite.Pos.Y)
		options.ColorM.Scale(1, 1, 1, sprite.Alpha)
		x0 := currAnimFrame.CurrFrame * currAnimFrame.FrameWidth
		x1 := x0 + currAnimFrame.FrameWidth
		r := image.Rect(x0, 0, x1, currAnimFrame.FrameHeight)
		surface.DrawImage(currAnimFrame.Image.SubImage(r).(*ebiten.Image), options)
		sprite.nextFrame()
	}
}

func (sprite *Sprite) Stop() {
	sprite.Animated = false
	sprite.Visible = false
}

func (sprite *Sprite) Reset() {
	sprite.Animated = false
	sprite.Visible = false
	sprite.AnimFrames[sprite.CurrAnimFrame].CurrFrame = 0
}

func (sprite *Sprite) Start() {
	sprite.Animated = true
	sprite.Visible = true
}

//GetWidth returns width of the current animation displayed
func (sprite *Sprite) GetWidth() float64 {
	currAnimFrame := sprite.AnimFrames[sprite.CurrAnimFrame]
	return float64(currAnimFrame.FrameWidth)
}

//GetHeight returns height of the current animation displayed
func (sprite *Sprite) GetHeight() float64 {
	currAnimFrame := sprite.AnimFrames[sprite.CurrAnimFrame]
	return float64(currAnimFrame.FrameHeight)
}

func (sprite *Sprite) Position(arg ...float64) (float64, float64) {
	if len(arg) == 2 {
		sprite.Pos.X = arg[0]
		sprite.Pos.Y = arg[1]
	}
	return sprite.Pos.X, sprite.Pos.Y
}
