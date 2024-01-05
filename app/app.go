package app

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/joseph-beck/go-redis/cache"
	"github.com/joseph-beck/go-redis/database"
	routey "github.com/joseph-beck/routey/pkg/router"
)

type App struct {
	Router routey.App
	Store  database.Store
	Cache  cache.Cache
}

var r routey.App
var c cache.Cache

func Run() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	c = cache.New()
	s, err := c.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	r = *routey.New()
	go shutdown()
	r.Run()
}

func shutdown() {
	r.Shutdown()
	c.Close()
}
