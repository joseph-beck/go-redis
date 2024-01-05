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
		{
			Path:        "/api/v1/users",
			Params:      "",
			Method:      routey.Post,
			HandlerFunc: s.Post(),
		},
		{
			Path:        "/api/v1/users",
			Params:      "",
			Method:      routey.Put,
			HandlerFunc: s.Put(),
		},
		{
			Path:        "/api/v1/users",
			Params:      "",
			Method:      routey.Patch,
			HandlerFunc: s.Patch(),
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
		e, err := s.cache.Contains(fmt.Sprintf("%d", m.ID))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if e {
			err := s.cache.Get(&m, fmt.Sprintf("%d", m.ID))
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

		_, err = s.cache.Set(&m, fmt.Sprintf("%d", m.ID))
		if err != nil {
			log.Printf("%v, failed to add to cache", err)
		}

		c.JSON(http.StatusOK, r)
	}
}

func (s *UserService) Post() routey.HandlerFunc {
	return func(c *routey.Context) {
		var m models.User
		err := c.ShouldBindJSON(&m)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		e, err := s.cache.Contains(fmt.Sprintf("%d", m.ID))
		if err != nil {
			log.Println(err)
		}
		if e && m.ID != 0 {
			c.Status(http.StatusBadRequest)
			return
		}

		e = s.db.Contains(&models.User{Model: models.Model{ID: m.ID}}, s.table)
		if e && m.ID != 0 {
			c.Status(http.StatusBadRequest)
			return
		}

		err = s.db.Add(&m, s.table)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		_, err = s.cache.Set(&m, fmt.Sprintf("%d", m.ID))
		if err != nil {
			log.Printf("%v, failed to add to cache", err)
		}

		c.Status(http.StatusNoContent)
	}
}

func (s *UserService) Put() routey.HandlerFunc {
	return func(c *routey.Context) {
		var m models.User
		err := c.ShouldBindJSON(&m)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		e := s.db.Contains(&models.User{Model: models.Model{ID: m.ID}}, s.table)
		if !e || m.ID == 0 {
			err = s.db.Add(&m, s.table)
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			_, err = s.cache.Set(&m, fmt.Sprintf("%d", m.ID))
			if err != nil {
				log.Printf("%v, failed to add to cache", err)
			}

			c.Status(http.StatusNoContent)
			return
		}

		err = s.db.Update(&m, s.table)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		_, err = s.cache.Set(&m, fmt.Sprintf("%d", m.ID))
		if err != nil {
			log.Printf("%v, failed to add to cache", err)
		}

		c.Status(http.StatusNoContent)
	}
}

func (s *UserService) Patch() routey.HandlerFunc {
	return func(c *routey.Context) {
		var m models.User
		err := c.ShouldBindJSON(&m)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if m.ID == 0 {
			c.Status(http.StatusBadRequest)
			return
		}

		e := s.db.Contains(&models.User{Model: models.Model{ID: m.ID}}, s.table)
		if !e {
			c.Status(http.StatusBadRequest)
			return
		}

		err = s.db.Update(&m, s.table)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		_, err = s.cache.Set(&m, fmt.Sprintf("%d", m.ID))
		if err != nil {
			log.Printf("%v, failed to add to cache", err)
		}

		c.Status(http.StatusNoContent)
	}
}

func (s *UserService) Delete() routey.HandlerFunc {
	return func(c *routey.Context) {
		i, err := c.ParamInt("user")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		e := s.db.Contains(&models.User{Model: models.Model{ID: uint(i)}}, s.table)
		if !e {
			c.Status(http.StatusBadRequest)
			return
		}

		err = s.db.Delete(&models.User{Model: models.Model{ID: uint(i)}}, s.table)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		e, err = s.cache.Contains(fmt.Sprintf("%d", i))
		if err != nil {
			log.Println(err)
		}
		if e {
			_, err = s.cache.Delete(fmt.Sprintf("%d", i))
			if err != nil {
				log.Println(err)
			}
		}

		c.Status(http.StatusNoContent)
	}
}
