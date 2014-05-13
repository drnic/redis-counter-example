package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	redis "gopkg.in/redis.v1"
)

var client *redis.Client

func init() {
	client = redis.NewTCPClient(&redis.Options{
		Addr: ":6379",
	})
}

type nameForm struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type nameFormRender struct {
	Count string
	Name  string
	Error string
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{}))

	m.Get("/", func(ren render.Render) {
		_, err := client.Incr("count").Result()
		if err != nil {
			ren.HTML(200, "index", &nameFormRender{Error: err.Error()})
			return
		}

		count, err := client.Get("count").Result()
		if err != nil {
			ren.HTML(200, "index", &nameFormRender{Error: err.Error()})
			return
		}

		name, err := client.Get("name").Result()
		if err != nil {
			ren.HTML(200, "index", &nameFormRender{Error: "No one has said hello yet", Count: count})
			return
		}
		ren.HTML(200, "index", &nameFormRender{Name: name, Count: count})
		return
	})

	m.Post("/name", binding.Form(nameForm{}), func(form nameForm, err binding.Errors, ren render.Render, r *http.Request) {
		set := client.Set("name", form.Name)
		fmt.Printf("%#v\n", set)
		fmt.Printf("%#v\n", set.Err())

		ren.Redirect("/")
	})

	m.Run()
}
