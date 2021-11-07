package ui2d

import (
	"fmt"
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

//向量常數值
var (
	Zero     = Vec3{0, 0, 0}
	Up       = Vec3{0, 1, 0}
	Down     = Vec3{0, -1, 0}
	Left     = Vec3{-1, 0, 0}
	Right    = Vec3{1, 0, 0}
	Forward  = Vec3{0, 0, 1}
	Backward = Vec3{0, 0, -1}
	One      = Vec3{1, 1, 1}
	MinusOne = Vec3{-1, -1, -1}
)

func Roundf(val float64, places int) float64 {
	if places < 0 {
		panic("places should be >= 0")
	}

	factor := math.Pow10(places)
	val = val * factor
	tmp := float64(int(val))
	return tmp / factor
}

func Lerpf(from, to float64, t float64) float64 {
	return from + ((to - from) * t)
}

func LerpAngle(from, to float64, t float64) float64 {
	for to-from > 180 {
		from += 360
	}
	for from-to > 180 {
		to += 360
	}
	return from + ((to - from) * t)
}

func (v *Vec3) String() string {
	return fmt.Sprintf("(%f,%f,%f)", v.X, v.Y, v.Z)
}

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{x, y, z}
}

func (v *Vec3) Add(vect Vec3) Vec3 {
	return Vec3{v.X + vect.X, v.Y + vect.Y, v.Z + vect.Z}
}

func (v *Vec3) Sub(vect Vec3) Vec3 {
	return Vec3{v.X - vect.X, v.Y - vect.Y, v.Z - vect.Z}
}

func (v *Vec3) Mul(vect Vec3) Vec3 {
	return Vec3{v.X * vect.X, v.Y * vect.Y, v.Z * vect.Z}
}

func (v *Vec3) Mul2(vect float64) Vec3 {
	return Vec3{v.X * vect, v.Y * vect, v.Z * vect}
}

func (v *Vec3) Distance(vect Vec3) float64 {
	x := v.X - vect.X
	y := v.Y - vect.Y
	return math.Sqrt(float64(x*x + y*y))
}

func (v *Vec3) Div(vect Vec3) Vec3 {
	return Vec3{v.X / vect.X, v.Y / vect.Y, v.Z / vect.Z}
}

func (v *Vec3) fixAngle() {
	for v.X >= 360 {
		v.X -= 360
	}
	for v.X <= -360 {
		v.X += 360
	}

	for v.Y >= 360 {
		v.Y -= 360
	}
	for v.Y <= -360 {
		v.Y += 360
	}

	for v.Z >= 360 {
		v.Z -= 360
	}
	for v.Z <= -360 {
		v.Z += 360
	}
}

func (v *Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func Lerp(from, to Vec3, t float64) Vec3 {
	return NewVec3(from.X+((to.X-from.X)*t), from.Y+((to.Y-from.Y)*t), from.Z+((to.Z-from.Z)*t))
}

func (v *Vec3) Normalize() {
	l := v.Length()
	v.X /= l
	v.Y /= l
	v.Z /= l
}

func (v *Vec3) Normalized() Vec3 {
	l := v.Length()
	if l == 0 {
		return NewVec3(0, 0, 0)
	}
	return NewVec3(v.X/l, v.Y/l, v.Z/l)
}
