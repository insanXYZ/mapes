package main

import "github.com/insanXYZ/mapes"

func main() {
	m := mapes.New()
	m.Get("/", func(ctx *mapes.Context) error {
		return ctx.String(200, "Hello world")
	})
	err := m.Start("1323")
	if err != nil {
		panic(err.Error())
	}
}
