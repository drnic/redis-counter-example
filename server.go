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
	DB       int64  `json:"database"`
	Password string `json:"password"`
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadRedis(options *redis.Options) *redis.Client {
	return redis.NewTCPClient(options)
}

func loadRedisOptions() *redis.Options {
	configPath := flag.String("config", "", "config file for redis connection")
	flag.Parse()

	var err error

	config := redisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	if port := os.Getenv("REDIS_PORT"); port != "" {
		config.Port, err = strconv.ParseInt(port, 0, 64)
		panicIfErr(err)
	}
	if port := os.Getenv("REDIS_DATABASE"); port != "" {
		config.DB, err = strconv.ParseInt(port, 0, 64)
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

	return &redis.Options{
		Addr:     config.Host + ":" + strconv.FormatInt(config.Port, 10),
		Password: config.Password,
		DB:       config.DB,
	}
}

type nameForm struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type nameFormRender struct {
	Count string
	Name  string
	Error string
	Redis *redis.Options
}

func main() {
	redisOptions := loadRedisOptions()
	m := martini.Classic()
	m.Map(loadRedis(redisOptions))
	m.Use(render.Renderer(render.Options{}))

	m.Get("/", func(ren render.Render, client *redis.Client) {
		_, err := client.Incr("count").Result()
		if err != nil {
			ren.HTML(200, "index", &nameFormRender{Error: err.Error(), Redis: redisOptions})
			return
		}

		count, err := client.Get("count").Result()
		if err != nil {
			ren.HTML(200, "index", &nameFormRender{Error: err.Error(), Redis: redisOptions})
			return
		}

		name, err := client.Get("name").Result()
		if err != nil {
			ren.HTML(200, "index", &nameFormRender{Error: "No one has said hello yet", Count: count, Redis: redisOptions})
			return
		}
		ren.HTML(200, "index", &nameFormRender{Name: name, Count: count, Redis: redisOptions})
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
