package controller

import (
	"calls-service/rest-service/internal/controller/middleware"
	"calls-service/rest-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CallsHandler struct {
	u usecase.UseCase
	l zerolog.Logger
}

func New(u *usecase.UseCase, l zerolog.Logger) *CallsHandler {
	return &CallsHandler{u: *u, l: l}
}

func NewCallsRoutes(router *gin.Engine, h *CallsHandler) {

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", h.register)
		authGroup.POST("/login", h.login)
	}

	callsGroup := router.Group("/calls")

	callsGroup.Use(middleware.Auth())
	{
		callsGroup.POST("", h.SaveCall)
		callsGroup.GET("", h.GetUserCalls)
		callsGroup.GET("/:id", h.GetUserCallByID)
		callsGroup.PATCH("/:id/status", h.UpdateCallStatus)
		callsGroup.DELETE("/:id", h.DeleteCall)
	}
}
