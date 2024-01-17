package main

import "github.com/joseph-beck/go-redis/app"

// @title go redis
// @version 1.0
// @description testing out redis caching in golang.

// @license.name MIT

// @BasePath /
func main() {
	a := app.New()
	a.Run()
}
