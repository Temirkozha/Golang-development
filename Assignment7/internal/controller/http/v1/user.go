package v1

import (
	"net/http"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/utils"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	t usecase.UserInterface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface) {
	r := &userRoutes{t}
	h := handler.Group("/users")
	h.Use(utils.RateLimiterMiddleware())
	{
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		{
			protected.GET("/me", r.GetMe)
			adminOnly := protected.Group("/")
			adminOnly.Use(utils.RoleMiddleware("admin"))
			{
				adminOnly.PATCH("/promote/:id", r.PromoteUser)
			}
		}
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var dto entity.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := utils.HashPassword(dto.Password)
	user := &entity.User{Username: dto.Username, Email: dto.Email, Password: hash, Role: "user"}
	u, s, err := r.t.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": u, "session_id": s})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var dto entity.LoginUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := r.t.LoginUser(&dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	id := c.GetString("userID")
	user, err := r.t.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
	id := c.Param("id")
	if err := r.t.PromoteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user promoted"})
}