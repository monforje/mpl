package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func (u *User) SayHello() {
	fmt.Println("Hello,", u.Name)
}

func (u *User) Birthday() {
	u.Age = u.Age + 1
}

func (u *User) Info() {
	fmt.Println("User:", u.Name)
	fmt.Println("Age:", u.Age)
}

func Add(a int, b int) int {
	return a + b
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	user := User{Name: "Artyom", Age: 20}
	user.SayHello()
	user.Birthday()
	user.Info()

	fmt.Println("ADD:", Add(2, 3))
	fmt.Println("MAX:", Max(10, 4))
}
