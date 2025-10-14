package main

import (
	"fmt"

	"technoBro/internal/app/configer"
)

//go:generate go run ./main.go

func main() {
	c := configer.New(
		"./../../internal/app/configer/svc",
		"./../../cmd",
	)
	if err := c.Gen(); err != nil {
		panic(err)
	}
	fmt.Println("configs generated successfully")
}
