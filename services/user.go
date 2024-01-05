package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joseph-beck/go-redis/cache"
	"github.com/joseph-beck/go-redis/database"
	"github.com/joseph-beck/go-redis/models"
	routey "github.com/joseph-beck/routey/pkg/router"
)

type UserService struct {
	db    *database.Store
	cache *cache.Cache
	table string
}

func NewUserService(d *database.Store, c *cache.Cache) UserService {
	return UserService{
		db:    d,
		cache: c,
		table: "users",
	}
}

func (s UserService) Add() []routey.Route {
	return []routey.Route{
		{
			Path:        "/api/v1/users",
			Params:      "",
			Method:      routey.Get,
			HandlerFunc: s.List(),
		},
		{
			Path:        "/api/v1/users",
			Params:      "/:user",
			Method:      routey.Get,
			HandlerFunc: s.Get(),
		},
	}
}

func (s *UserService) List() routey.HandlerFunc {
	type Response struct {
		Users []models.User `json:"users"`
		Count int           `json:"user_count"`
	}

	return func(c *routey.Context) {
		var m []models.User
		err := s.db.List(&m, s.table)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		r := Response{
			Users: m,
			Count: len(m),
		}

		c.JSON(http.StatusOK, r)
	}
}

func (s *UserService) Get() routey.HandlerFunc {
	return func(c *routey.Context) {
		i, err := c.ParamInt("user")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		m := models.User{Model: models.Model{ID: uint(i)}}
		e, err := s.cache.Contains(string(rune(m.ID)))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if e {
			fmt.Println("fetched from cache")
			err := s.cache.Get(&m, string(rune(m.ID)))
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			c.JSON(http.StatusOK, m)
			return
		}

		e = s.db.Contains(&m, s.table)
		if !e {
			c.Status(http.StatusNotFound)
			return
		}

		err = s.db.Get(&m, s.table)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		r := m

		_, err = s.cache.Set(&m, string(rune(m.ID)))
		if err != nil {
			log.Printf("%v, failed to add to cache", err)
		}

		c.JSON(http.StatusOK, r)
	}
}
