package controller

import (
	"calls-service/rest-service/internal/controller/apierrors"
	"calls-service/rest-service/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CallsHandler) register(c *gin.Context) {
	var req entity.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid request format"})
	}

	err := h.u.RegisterUser(c.Request.Context(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid request format"})
			case codes.AlreadyExists:
				c.JSON(http.StatusConflict, apierrors.Response{Error: "user already exists"})
			default:
				c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "internal error"})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "unknown error"})
		return
	}
	h.l.Info().Str("user", req.Username).Msg("User save successfully")

	c.Status(http.StatusCreated)
}

func (h *CallsHandler) login(c *gin.Context) {
	var req entity.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid request format"})
		return
	}

	token, err := h.u.LoginUser(c.Request.Context(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid request format"})
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, apierrors.Response{Error: "Invalid username or password"})
			default:
				c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "internal error"})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "unknown error"})
		return
	}

	h.l.Info().Str("user", req.Username).Str("token", token).Msg("User logged in successfully")

	c.JSON(http.StatusOK, gin.H{"token": token})
}
