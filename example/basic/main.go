package main

import (
	"fmt"
	"github.com/insanXYZ/mapes"
)

func main() {
	m := mapes.New()
	m.Get("/", func(ctx *mapes.Context) error {
		return ctx.String(200, "Hello world")
	})
	m.Get("/hello/:name/from/:address", func(ctx *mapes.Context) error {
		return ctx.String(200, fmt.Sprintf("Hello, my name is %s from %s", ctx.Param("name"), ctx.Param("address")))
	})
	err := m.Start("1323")
	if err != nil {
		panic(err.Error())
	}
}
