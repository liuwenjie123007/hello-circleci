package main

import (
	"fmt"
	"hello-circleci/internal/config"
	"hello-circleci/internal/repository"
	"hello-circleci/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ca = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	c, err := config.GetAPIConfig()
	if err != nil {
		panic(err)
	}

	d, err := db.NewDB(&db.Config{
		URI:      fmt.Sprintf("(%s:%s)/circle", c.MySQL.Host, c.MySQL.Port),
		User:     c.MySQL.User,
		Password: c.MySQL.Password,
	})

	if err != nil {
		panic(err)
	}
	ur := repository.NewUserRepository(d)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		user, err := ur.Get(name)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"id": user.Id, "name": user.Name})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	r.POST("/user/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")

		err := ur.Add(name)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "no"})
		}
	})

	r.DELETE("/user/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		err := ur.Delete(name)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "no"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
