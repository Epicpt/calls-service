package controller

import (
	"calls-service/rest-service/internal/controller/apierrors"
	"calls-service/rest-service/internal/entity"
	"calls-service/rest-service/internal/usecase"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

const statusOpen = "открыта"

func (h *CallsHandler) SaveCall(c *gin.Context) {
	var input entity.CallDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid request format"})
		return
	}

	if !validatePhoneNumber(input.PhoneNumber) {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid phone number format"})
		return
	}

	userIDAny, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apierrors.Response{Error: "Unauthorized"})
		return
	}

	userID, ok := userIDAny.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Invalid user ID in context"})
		return
	}

	newCall := entity.Call{
		ClientName:  input.ClientName,
		PhoneNumber: input.PhoneNumber,
		Description: input.Description,
		Status:      statusOpen,
		UserID:      userID,
	}

	err := h.u.SaveCall(c.Request.Context(), newCall)
	if err != nil {
		h.l.Error().Err(err).Msg("Failed to save call")
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Failed to save call"})
		return
	}

	//h.l.Info().

	c.Status(http.StatusCreated)
}

func validatePhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^(\+?\d{1,3}|\d)?[\d\-]{7,15}$`)
	return re.MatchString(phone)
}

func (h *CallsHandler) GetUserCalls(c *gin.Context) {
	userIDAny, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apierrors.Response{Error: "Unauthorized"})
		return
	}

	userID, ok := userIDAny.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Invalid user ID in context"})
		return
	}

	calls, err := h.u.GetUserCalls(c.Request.Context(), userID)
	if err != nil {
		h.l.Error().Err(err).Msg("Failed to get user calls")
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Failed to get user calls"})
		return
	}

	responses := make([]entity.CallResponse, len(calls))

	for i, c := range calls {
		responses[i] = entity.CallResponse{
			ID:          c.ID,
			ClientName:  c.ClientName,
			PhoneNumber: c.PhoneNumber,
			Description: c.Description,
			Status:      c.Status,
			CreatedAt:   c.CreatedAt,
		}
	}

	//h.l.Info().

	c.JSON(http.StatusOK, responses)
}

func (h *CallsHandler) GetUserCallByID(c *gin.Context) {
	userIDAny, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apierrors.Response{Error: "Unauthorized"})
		return
	}

	userID, ok := userIDAny.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Invalid user ID in context"})
		return
	}

	callIDStr := c.Param("id")
	callID, err := strconv.ParseInt(callIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid call ID"})
		return
	}

	call, err := h.u.GetUserCallByID(c.Request.Context(), callID, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrCallNotFound) {
			c.JSON(http.StatusNotFound, apierrors.Response{Error: "Call not found"})
			return
		}
		h.l.Error().Err(err).Msg("Failed to get user call by ID")
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Failed to get user call"})
	}
	c.JSON(http.StatusOK, call)
}

func (h *CallsHandler) UpdateCallStatus(c *gin.Context) {
	userIDAny, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apierrors.Response{Error: "Unauthorized"})
		return
	}

	userID, ok := userIDAny.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Invalid user ID in context"})
		return
	}

	callIDStr := c.Param("id")
	callID, err := strconv.ParseInt(callIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid call ID"})
		return
	}

	var input entity.UpdateCallStatusDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid request format"})
		return
	}

	if input.Status != "открыта" && input.Status != "закрыта" {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid status value"})
		return
	}

	if err := h.u.UpdateCallStatus(c.Request.Context(), callID, userID, input.Status); err != nil {
		if errors.Is(err, usecase.ErrCallNotFound) {
			c.JSON(http.StatusNotFound, apierrors.Response{Error: "Call not found or does not belong to user"})
			return
		}
		h.l.Error().Err(err).Msg("Failed to update call status")
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Failed to update call status"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CallsHandler) DeleteCall(c *gin.Context) {
	userIDAny, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apierrors.Response{Error: "Unauthorized"})
		return
	}

	userID, ok := userIDAny.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Invalid user ID in context"})
		return
	}

	callIDStr := c.Param("id")
	callID, err := strconv.ParseInt(callIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid call ID"})
		return
	}

	if err := h.u.DeleteCall(c.Request.Context(), callID, userID); err != nil {
		if errors.Is(err, usecase.ErrCallNotFound) {
			c.JSON(http.StatusNotFound, apierrors.Response{Error: "Call not found or does not belong to user"})
			return
		}
		h.l.Error().Err(err).Msg("Failed to delete call")
		c.JSON(http.StatusInternalServerError, apierrors.Response{Error: "Failed to delete call"})
		return
	}

	c.Status(http.StatusNoContent)
}
