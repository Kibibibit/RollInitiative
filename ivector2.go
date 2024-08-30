package main

import "fmt"

type IVector2 struct {
	x int
	y int
}

func (v IVector2) String() string {
	return fmt.Sprintf("IVector2:{x:%d, y:%d}", v.x, v.y)
}
