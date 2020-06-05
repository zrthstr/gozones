package main

import (
	"fmt"
)

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return 3.1415926 * c.radius * c.radius
}

func (c circle) perim() float64 {
	return 2 * c.radius * 3.1415926
}

type rect struct {
	x, y float64
}

func (r rect) area() float64 {
	return r.x * r.y
}

func (r rect) perim() float64 {
	return 2*r.x + 2*r.y
}

func main() {
	ball := circle{radius: 10.2}
	fmt.Println(ball, ball.area(), ball.perim())

	choco := rect{x: 12.8, y: 3.55}
	fmt.Println(choco, choco.area(), choco.perim())

}
