package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomePage(c *gin.Context) {
	c.String(http.StatusOK, "BallRoom Home Page")
}
