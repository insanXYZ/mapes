package main

import (
	"fmt"

	"github.com/insanXYZ/mapes"
)

type CreateUser struct {
	Name    string `form:"name"`
	Address string `param:"address"`
	Email   string `query:"email"`
}

func main() {
	m := mapes.New()

	m.Get("/", func(ctx *mapes.Context) error {
		return ctx.String(200, "Hello world")
	})

	m.Static("/static", "resources")

	m.Get("/hello/:name/from/:address", func(ctx *mapes.Context) error {
		return ctx.String(200, fmt.Sprintf("Hello, my name is %s from %s", ctx.Param("name"), ctx.Param("address")))
	})

	m.Get("/queryParams", func(ctx *mapes.Context) error {
		return ctx.String(200, ctx.Query("last"))
	})

	m.Post("/create-user/:address", func(ctx *mapes.Context) error {
		user := new(CreateUser)
		err := ctx.Bind(user)
		if err != nil {
			return err
		}

		return ctx.Json(200, user)
	})

	group := m.Group("/api")

	group.Get("/contacts", func(ctx *mapes.Context) error {
		return ctx.Json(200, map[string]string{
			"name":  "Jhon Doe",
			"email": "jhondoe@example.com",
		})
	})

	err := m.Start(":1323")
	if err != nil {
		panic(err.Error())
	}
}
