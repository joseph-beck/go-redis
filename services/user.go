package services

import (
	"net/http"

	"github.com/joseph-beck/go-redis/database"
	"github.com/joseph-beck/go-redis/models"
	routey "github.com/joseph-beck/routey/pkg/router"
)

type UserService struct {
	db    *database.Store
	table string
}

func NewUserService(db *database.Store) UserService {
	return UserService{
		db:    db,
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
	return func(c *routey.Context) {

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
		e := s.db.Contains(&m, s.table)
		if !e {
			c.Status(http.StatusNotFound)
			return
		}

		err = s.db.Get(&m, s.table)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, m)
	}
}
