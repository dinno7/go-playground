package main

import "fmt"

type Shape interface {
	getArea() float64
}

type Square struct {
	sideLength float64
}

func (s Square) getArea() float64 {
	return s.sideLength * s.sideLength
}

type Triangle struct {
	height float64
	base   float64
}

func (t Triangle) getArea() float64 {
	return 0.5 * t.base * t.height
}

func main() {
	t := Triangle{height: 12, base: 10}
	printArea(t)

	s := Square{10}
	printArea(s)
}

func printArea(shape Shape) {
	fmt.Println("💀 The area is > ", shape.getArea())
}
