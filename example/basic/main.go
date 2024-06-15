package main

import (
	"fmt"
	"github.com/insanXYZ/mapes"
)

type FormUser struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

func main() {
	m := mapes.New()
	m.Get("/", func(ctx *mapes.Context) error {
		return ctx.String(200, "Hello world")
	})
	m.Get("/hello/:name/from/:address", func(ctx *mapes.Context) error {
		return ctx.String(200, fmt.Sprintf("Hello, my name is %s from %s", ctx.Param("name"), ctx.Param("address")))
	})
	m.Get("/queryParams", func(ctx *mapes.Context) error {
		return ctx.String(200, ctx.Query("last"))
	})
	m.Post("/form", func(ctx *mapes.Context) error {
		form := FormUser{}
		if err := ctx.Bind(&form); err != nil {
			panic(err.Error())
		}
		return ctx.Json(200, form)
	})
	err := m.Start("1323")
	if err != nil {
		panic(err.Error())
	}
}
