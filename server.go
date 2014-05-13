package main

import (
	"encoding/json"
	"flag"
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
	Port     int64  `json:"port"`
	Database string `json:"database"`
	Password string `json:"password"`
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadRedis() *redis.Client {
	configPath := flag.String("config", "", "config file for redis connection")
	flag.Parse()

	var err error

	config := redisConfig{
		Host:     os.Getenv("HOST"),
		Password: os.Getenv("PASSWORD"),
	}

	if port := os.Getenv("PORT"); port != "" {
		config.Port, err = strconv.ParseInt(port, 0, 64)
		panicIfErr(err)
	}

	if *configPath != "" {
		fmt.Println("Loading config from", *configPath)
		configFile, err := os.Open(*configPath)
		defer configFile.Close()
		panicIfErr(err)

		jsonParser := json.NewDecoder(configFile)
		err = jsonParser.Decode(&config)
		panicIfErr(err)
	}
	fmt.Printf("%#v\n", config)

	return redis.NewTCPClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.FormatInt(config.Port, 10),
		Password: config.Password,
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
		_ = client.Set("name", form.Name)
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
