package v1

import (
	"github.com/gin-gonic/gin"
	"practice-7/internal/usecase"
)

func NewRouter(handler *gin.Engine, t usecase.UserInterface) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, t)
	}
}