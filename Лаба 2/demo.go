package main

import "fmt"

type Counter struct {
	Value int
}

func (c Counter) Inc(step int) Counter {
	c.Value = c.Value + step
	return c
}

func divmod(a int, b int) (int, int) {
	q := a / b
	r := a % b
	return q, r
}

func main() {
	c := Counter{Value: 0}

	for i := 0; i < 5; i = i + 1 {
		c = c.Inc(2)
	}

	q, r := divmod(c.Value, 3)
	if r == 0 {
		fmt.Println(q)
	} else {
		fmt.Println(q, r)
	}
}
