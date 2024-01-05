package app

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/joseph-beck/go-redis/cache"
	"github.com/joseph-beck/go-redis/database"
	"github.com/joseph-beck/go-redis/services"
	routey "github.com/joseph-beck/routey/pkg/router"
)

type App struct {
	Router *routey.App
	Store  *database.Store
	Cache  *cache.Cache
}

func New() *App {
	fmt.Println("loading env")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println("making app")
	r := routey.New()
	s := database.New()
	c := cache.New()

	return &App{
		Router: r,
		Store:  s,
		Cache:  c,
	}
}

func (a *App) Run() {
	fmt.Println("migrating db")
	err := a.Store.AutoMigrate()
	if err != nil {
		panic(err)
	}

	fmt.Println("pinging cache")
	s, err := a.Cache.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("loading cache, ", s)

	fmt.Println("starting router")
	a.Router.Service(services.NewUserService(a.Store))
	a.Router.Service(services.NewPingService())
	go a.Router.Shutdown(a.shutdown())
	a.Router.Run()
}

func (a *App) shutdown() routey.ShutdownFunc {
	return func() {
		err := a.Cache.Close()
		if err != nil {
			log.Println(err)
		}

		err = a.Store.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
