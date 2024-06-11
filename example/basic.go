package main

import "github.com/insanXYZ/mapes"

func main() {
	m := mapes.New()
	m.Get("/", func(ctx *mapes.Context) {
		ctx.W.Write([]byte("hello world"))
	})
	err := m.Start("1323")
	if err != nil {
		panic(err.Error())
	}
}
