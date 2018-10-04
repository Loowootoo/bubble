package vector

import (
	"math"
)

// Vector 3D vector
type Vector struct {
	X, Y, Z float64
}

// Add add 2 vector
func Add(a, b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Mult multiplexer 2 vector
func Mult(a Vector, b float64) Vector {
	return Vector{a.X * b, a.Y * b, a.Z * b}
}

// Length get vector length
func (a Vector) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

// Distance get two vector distance
func Distance(a, b Vector) float64 {
	xDiff := a.X - b.X
	yDiff := a.Y - b.Y
	zDiff := a.Z - b.Z
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff + zDiff*zDiff)
}

// DistanceSquared get two vector distance squared
func DistanceSquared(a, b Vector) float64 {
	xDiff := a.X - b.X
	yDiff := a.Y - b.Y
	zDiff := a.Z - b.Z
	return xDiff*xDiff + yDiff*yDiff + zDiff*zDiff
}

// Normalize a vector
func Normalize(a Vector) Vector {
	len := a.Length()
	return Vector{a.X / len, a.Y / len, a.Z / len}
}
