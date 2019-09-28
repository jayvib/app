package shapes

import "math"

// - Write the test first
// - Write the minimum amount of code for the
// 	 test to run and check the failing test output
// - Write enough code to make it pass

type Areaer interface {
	Area() float64
}

type Perimeterer interface {
	Perimeter() float64
}

type Shape interface {
	Areaer
	Perimeterer
}

type Triangle struct {
	Base   float64
	Height float64
	Shape
}

func (t Triangle) Area() float64 {
	return (t.Base * t.Height) / 2
}


// Rectangle represents a rectangle.
type Rectangle struct {
	Width  float64
	Height float64
	Shape
}

// Area calculates the area of the rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}

// Circle represents a circle.
type Circle struct {
	Radius float64
	Shape
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Area calculates the  area of circle.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func Perimeter(r Rectangle) float64 {
	return 2 * (r.Width + r.Height)
}

func Area(r Rectangle) float64 {
	return r.Width * r.Height
}
