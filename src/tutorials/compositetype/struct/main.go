package main

import "time"

type Employee struct {
	ID        int
	Name      string
	Address   string
	DoB       time.Time
	Position  string
	Salary    int
	ManagerID int
}

// 结构体字面值
type Point struct{ X, Y int }

type Circle struct {
	Point
	Radius int
}

type Wheel struct {
	Circle
	Spokes int
}

// // 结合匿名成员，简化访问
// type Circle2 struct {
// 	Point
// 	Radius int
// }

// type Wheel2 struct {
// 	Circle
// 	Spokes int
// }

func main() {
	var dilbert Employee
	dilbert.Salary -= 5000
	position := &dilbert.Position
	*position = "Senior " + *position

	var employeeOfTheMonth *Employee = &dilbert
	employeeOfTheMonth.Position += " (proactive team player)"
	// 等价于
	(*employeeOfTheMonth).Position += " (proactive team player)"

	_ = Point{1, 2}

	// var w Wheel
	// w.Circle.Center.X = 8 // 这种访问机制很麻烦
	// w.Circle.Center.Y = 8
	// w.Circle.Radius = 5
	// w.Spokes = 20

	var w2 Wheel
	w2.X = 8
	w2.Y = 8
	w2.Radius = 5
	w2.Spokes = 20
}
