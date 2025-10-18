package main

import (
	"fmt"
	"math"
)

// 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
// 然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
// 考察点 ：接口的定义与实现、面向对象编程风格。
// 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
// 再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
// 考察点 ：组合的使用、方法接收者。

type Shape interface {
	Area() float64
	Perimeter() float64
}
type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radios float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}
func (r Circle) Area() float64 {
	return math.Pi * r.Radios * r.Radios
}
func (r Circle) Perimeter() float64 {
	return 2 * math.Pi * r.Radios
}

type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("Name:%s ", e.Name)
	fmt.Printf("Age:%d ", e.Age)
	fmt.Printf("EmployeeID:%s", e.EmployeeID)
}

func main() {
	var s Shape
	s = Rectangle{Width: 3, Height: 4}
	fmt.Println(s.Area())
	fmt.Println(s.Perimeter())

	s = Circle{Radios: 8}
	fmt.Println(s.Area())
	fmt.Println(s.Perimeter())

	emp := Employee{
		Person: Person{
			Name: "Mike",
			Age:  18,
		},
		EmployeeID: "88888888",
	}
	emp.PrintInfo()
}
