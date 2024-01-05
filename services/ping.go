package services

import (
	"net/http"

	routey "github.com/joseph-beck/routey/pkg/router"
)

type PingService struct {
}

func NewPingService() PingService {
	return PingService{}
}

func (s PingService) Add() []routey.Route {
	return []routey.Route{
		{
			Path:        "/api/v1/ping",
			Params:      "",
			Method:      routey.Get,
			HandlerFunc: s.Get(),
		},
	}
}

func (s *PingService) Get() routey.HandlerFunc {
	return func(c *routey.Context) {
		c.Status(http.StatusOK)
	}
}
