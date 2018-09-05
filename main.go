package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gui-moreira/golang-rest-introduction/user"
)

var repo user.Repo

func main() {
	repo = user.InmemUserRepo{Users: make(map[int]user.User)}
	r := configureRoutes()
	r.Run(":8080")
}

func configureRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/users", func(c *gin.Context) {
		u := user.User{}
		err := c.ShouldBindJSON(&u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if err := repo.Save(&u); err != nil {
			if err == user.ErrValidateUser {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": u.ID})
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		u, err := repo.Get(id)
		if err != nil {
			if err == user.ErrUserNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, u)
	})

	return r
}
