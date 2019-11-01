package main

import "fmt"

type Figure interface {
	SetWidth(width rune)
	SetHeight(height rune)
	GetWidth() rune
	GetHeight() rune
}

type Rectangle struct {
	width, height rune
}

func (r *Rectangle) SetWidth(w rune) {
	r.width = w
}

func (r *Rectangle) SetHeight(h rune) {
	r.height = h
}

func (r *Rectangle) GetWidth() rune {
	return r.width
}

func (r *Rectangle) GetHeight() rune {
	return r. height
}

type Square struct {
	*Rectangle
	height, width rune
}

func (s *Square) GetWidth() rune {
	return s.width
}

func (s *Square) GetHeight() rune {
	return s.height
}

func Area(f Figure, height, width rune) rune {
	f.SetHeight(height)
	f.SetWidth(width)
	return f.GetHeight() * f.GetWidth()
}

func main() {
	rectangle := &Rectangle{width: 2, height: 2}
	square := &Square{Rectangle: rectangle}
	fmt.Println("square:", square.GetWidth(), square.GetHeight())
	square.SetWidth(3)
	square.SetHeight(4)
	fmt.Println("square:", square.GetWidth(), square.GetHeight())
	fmt.Println("rectangle:", rectangle.GetWidth(), rectangle.GetHeight())
}

