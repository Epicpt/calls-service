package controller

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"calls-service/rest-service/internal/controller/apierrors"
	"calls-service/rest-service/internal/entity"
	"calls-service/rest-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

const statusOpen = "открыта"

// SaveCall handles the creation of a new call record.
//
// @Summary Create a new call
// @Description Saves a new call with client name, phone number, and description
// @Tags calls
// @Accept json
// @Produce json
// @Param input body entity.CallDTO true "Call data"
// @Success 201 {string} string "Created"
// @Failure 400 {object} apierrors.Response "Invalid input"
// @Failure 401 {object} apierrors.Response "Unauthorized"
// @Failure 500 {object} apierrors.Response "Internal server error"
// @Router /calls [post]
func (h *CallsHandler) SaveCall(c *gin.Context) {
	var input entity.CallDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, apierrors.Response{Error: "Invalid request format"})
		return
	}

	if !ValidatePhoneNumber(input.PhoneNumber) {
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

	h.l.Info().Interface("call", newCall).Msg("Call success save")

	c.Status(http.StatusCreated)
	c.Writer.Write([]byte{})
}

// GetUserCalls returns all calls for the authenticated user.
//
// @Summary Get user calls
// @Description Retrieves a list of calls belonging to the authenticated user
// @Tags calls
// @Produce json
// @Success 200 {array} entity.CallResponse "List of calls"
// @Failure 401 {object} apierrors.Response "Unauthorized"
// @Failure 500 {object} apierrors.Response "Internal server error"
// @Router /calls [get]
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

	c.JSON(http.StatusOK, calls)
}

// GetUserCallByID returns a specific call by ID for the authenticated user.
//
// @Summary Get user call by ID
// @Description Retrieves details of a specific call belonging to the authenticated user
// @Tags calls
// @Produce json
// @Param id path int true "Call ID"
// @Success 200 {object} entity.CallResponse "Call details"
// @Failure 400 {object} apierrors.Response "Invalid call ID"
// @Failure 401 {object} apierrors.Response "Unauthorized"
// @Failure 404 {object} apierrors.Response "Call not found"
// @Failure 500 {object} apierrors.Response "Internal server error"
// @Router /calls/{id} [get]
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
		return
	}

	h.l.Info().Interface("call", call).Msg("Call success get")

	c.JSON(http.StatusOK, call)
}

// UpdateCallStatus updates the status of a specific call for the authenticated user.
//
// @Summary Update call status
// @Description Updates the status (open or closed) of a specific user call
// @Tags calls
// @Accept json
// @Produce json
// @Param id path int true "Call ID"
// @Param input body entity.UpdateCallStatusDTO true "New status"
// @Success 204 "No Content"
// @Failure 400 {object} apierrors.Response "Invalid input or status value"
// @Failure 401 {object} apierrors.Response "Unauthorized"
// @Failure 404 {object} apierrors.Response "Call not found or does not belong to user"
// @Failure 500 {object} apierrors.Response "Internal server error"
// @Router /calls/{id}/status [put]
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

	h.l.Info().Int64("callID", callID).Msg("Call success update")

	c.JSON(http.StatusNoContent, nil)
}

// DeleteCall deletes a specific call for the authenticated user.
//
// @Summary Delete call
// @Description Deletes a call belonging to the authenticated user by its ID
// @Tags calls
// @Produce json
// @Param id path int true "Call ID"
// @Success 204 "No Content"
// @Failure 400 {object} apierrors.Response "Invalid call ID"
// @Failure 401 {object} apierrors.Response "Unauthorized"
// @Failure 404 {object} apierrors.Response "Call not found or does not belong to user"
// @Failure 500 {object} apierrors.Response "Internal server error"
// @Router /calls/{id} [delete]
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

	h.l.Info().Int64("callID", callID).Msg("Call success deleted")

	c.JSON(http.StatusNoContent, nil)
}

func ValidatePhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^(\+?\d{1,3}|\d)?[\d\-]{7,15}$`)
	return re.MatchString(phone)
}
