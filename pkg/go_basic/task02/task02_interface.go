package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

// 注意： 这里用 r *Rectangle 时，在 var s Shape 进行赋值时，要用 Rectangle的指针进行赋值
//
//	使用 r Rectangle 时，在 var s Shape 进行赋值时，要用 Rectangle的值进行赋值
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	p          Person
	EmployeeID string
}

func (e *Employee) PrintInfo() {
	fmt.Printf("Name: %s, Age: %d, EmployeeID: %s\n", e.p.Name, e.p.Age, e.EmployeeID)
}

func main02() {
	r := &Rectangle{Width: 5, Height: 10}
	c := &Circle{Radius: 7}
	// r1 := Rectangle{Width: 5, Height: 10}
	// var s Shape = r1
	// fmt.Println("Rectangle Area:", s.Area())
	fmt.Println("Rectangle Area:", r.Area())
	fmt.Println("Rectangle Perimeter:", r.Perimeter())
	fmt.Println("Circle Area:", c.Area())
	fmt.Println("Circle Perimeter:", c.Perimeter())

	e := &Employee{
		p:          Person{Name: "Alice", Age: 30},
		EmployeeID: "E12345",
	}
	e.PrintInfo()
}
