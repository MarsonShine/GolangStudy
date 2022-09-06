package main

import (
	"github.com/gin-gonic/gin"
)

func ginMiddleware() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Run("localhost:5000")
}
