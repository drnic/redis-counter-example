package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"

	redis "gopkg.in/redis.v1"
)

type redisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Password string `json:"password"`
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadRedis() *redis.Client {
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	panicIfErr(err)

	config := redisConfig{}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	panicIfErr(err)

	return redis.NewTCPClient(&redis.Options{
		Addr: config.Host + ":" + strconv.Itoa(config.Port),
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
	m.Map(loadRedis())
	m.Use(render.Renderer(render.Options{}))

	m.Get("/", func(ren render.Render, client *redis.Client) {
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

	m.Post("/name", binding.Form(nameForm{}), func(form nameForm, err binding.Errors, ren render.Render, client *redis.Client) {
		set := client.Set("name", form.Name)
		fmt.Printf("%#v\n", set)
		fmt.Printf("%#v\n", set.Err())

		ren.Redirect("/")
	})

	m.Post("/clear", func(ren render.Render, client *redis.Client) {
		_, err := client.Del("count").Result()
		if err != nil {
			ren.HTML(200, "index", &nameFormRender{Error: err.Error()})
			return
		}

		_, err = client.Del("name").Result()
		if err != nil {
			ren.HTML(200, "index", &nameFormRender{Error: err.Error()})
			return
		}

		ren.Redirect("/")
	})

	m.Run()
}
