package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"calls-service/rest-service/internal/controller"
	"calls-service/rest-service/internal/controller/apierrors"
	"calls-service/rest-service/internal/entity"
	"calls-service/rest-service/internal/mocks"
	"calls-service/rest-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid with plus and country code", "+71234567890", true},
		{"Valid without plus", "81234567890", true},
		{"Valid with dashes", "812-345-6789", true},
		{"Too short", "12345", false},
		{"Too long", "12345678901234567890", false},
		{"Letters inside", "123ABC7890", false},
		{"Empty string", "", false},
		{"Plus only", "+", false},
		{"Country code only", "+7", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := controller.ValidatePhoneNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSaveCall(t *testing.T) {
	tests := []struct {
		name             string
		input            entity.CallDTO
		mockSaveCallErr  error
		expectedStatus   int
		expectedResponse apierrors.Response
		setupContext     func(c *gin.Context)
		shouldCallMock   bool
	}{
		{
			name: "Successful save",
			input: entity.CallDTO{
				ClientName:  "John Doe",
				PhoneNumber: "+79876543211",
				Description: "Test call",
			},
			mockSaveCallErr: nil,
			expectedStatus:  http.StatusCreated,
			setupContext: func(c *gin.Context) {
				c.Set("id", int64(123))
			},
			shouldCallMock: true,
		},
		{
			name: "Invalid phone number format",
			input: entity.CallDTO{
				ClientName:  "John Doe",
				PhoneNumber: "invalid-phone",
				Description: "Test call",
			},
			mockSaveCallErr: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedResponse: apierrors.Response{
				Error: "Invalid phone number format",
			},
			setupContext: func(c *gin.Context) {
				c.Set("id", int64(123))
			},
			shouldCallMock: false,
		},
		{
			name: "Unauthorized (missing user ID)",
			input: entity.CallDTO{
				ClientName:  "John Doe",
				PhoneNumber: "+79876543211",
				Description: "Test call",
			},
			mockSaveCallErr: nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedResponse: apierrors.Response{
				Error: "Unauthorized",
			},
			setupContext: func(c *gin.Context) { // No user ID
			},
			shouldCallMock: false,
		},
		{
			name: "Failed to save call",
			input: entity.CallDTO{
				ClientName:  "John Doe",
				PhoneNumber: "+79876543211",
				Description: "Test call",
			},
			mockSaveCallErr: errors.New("database error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedResponse: apierrors.Response{
				Error: "Failed to save call",
			},
			setupContext: func(c *gin.Context) {
				c.Set("id", int64(123))
			},
			shouldCallMock: true,
		},
		{
			name: "Empty client name",
			input: entity.CallDTO{
				ClientName:  "",
				PhoneNumber: "+79876543211",
				Description: "Test call",
			},
			mockSaveCallErr: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedResponse: apierrors.Response{
				Error: "Invalid request format",
			},
			setupContext: func(c *gin.Context) {
				c.Set("id", int64(123))
			},
			shouldCallMock: false,
		},
		{
			name: "Empty description",
			input: entity.CallDTO{
				ClientName:  "John Doe",
				PhoneNumber: "+79876543211",
				Description: "",
			},
			mockSaveCallErr: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedResponse: apierrors.Response{
				Error: "Invalid request format",
			},
			setupContext: func(c *gin.Context) {
				c.Set("id", int64(123))
			},
			shouldCallMock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mocks.NewMockUseCase(t)
			if tt.shouldCallMock {
				mockUseCase.On("SaveCall", mock.Anything, mock.AnythingOfType("entity.Call")).
					Return(tt.mockSaveCallErr)
			}

			requestBody, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			tt.setupContext(c)
			c.Request = httptest.NewRequest("POST", "/calls", bytes.NewBuffer(requestBody))

			handler := controller.New(mockUseCase, zerolog.Nop())

			handler.SaveCall(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			responseBody := w.Body.Bytes()
			if w.Code != http.StatusCreated {
				var response apierrors.Response
				err := json.Unmarshal(responseBody, &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			} else {
				assert.Empty(t, responseBody)
			}

			if tt.shouldCallMock {
				mockUseCase.AssertExpectations(t)
			} else {
				mockUseCase.AssertNotCalled(t, "SaveCall")
			}
		})
	}
}

func TestGetUserCalls(t *testing.T) {
	tests := []struct {
		name             string
		mockGetCallsRes  []entity.CallResponse
		mockGetCallsErr  error
		expectedStatus   int
		expectedResponse any
		setupContext     func(c *gin.Context)
		shouldCallMock   bool
	}{
		{
			name: "Successful retrieval",
			mockGetCallsRes: []entity.CallResponse{
				{
					ID:          1,
					ClientName:  "John Doe",
					PhoneNumber: "+79876543211",
					Description: "Test call",
					Status:      "completed",
					CreatedAt:   time.Now(),
				},
			},
			mockGetCallsErr: nil,
			expectedStatus:  http.StatusOK,
			expectedResponse: []entity.CallResponse{
				{
					ID:          1,
					ClientName:  "John Doe",
					PhoneNumber: "+79876543211",
					Description: "Test call",
					Status:      "completed",
					CreatedAt:   time.Now(),
				},
			},
			setupContext: func(c *gin.Context) {
				c.Set("id", int64(123))
			},
			shouldCallMock: true,
		},
		{
			name:            "Unauthorized (missing user ID)",
			mockGetCallsErr: nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedResponse: apierrors.Response{
				Error: "Unauthorized",
			},
			setupContext:   func(c *gin.Context) {}, // no id
			shouldCallMock: false,
		},
		{
			name:            "Invalid user ID type",
			mockGetCallsErr: nil,
			expectedStatus:  http.StatusInternalServerError,
			expectedResponse: apierrors.Response{
				Error: "Invalid user ID in context",
			},
			setupContext: func(c *gin.Context) {
				c.Set("id", "not-an-int64")
			},
			shouldCallMock: false,
		},
		{
			name:            "Failed to get calls (internal error)",
			mockGetCallsErr: errors.New("db failure"),
			expectedStatus:  http.StatusInternalServerError,
			expectedResponse: apierrors.Response{
				Error: "Failed to get user calls",
			},
			setupContext: func(c *gin.Context) {
				c.Set("id", int64(123))
			},
			shouldCallMock: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mocks.NewMockUseCase(t)
			if tt.shouldCallMock {
				mockUseCase.On("GetUserCalls", mock.Anything, int64(123)).
					Return(tt.mockGetCallsRes, tt.mockGetCallsErr)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			tt.setupContext(c)
			c.Request = httptest.NewRequest("GET", "/calls", nil)

			handler := controller.New(mockUseCase, zerolog.Nop())

			handler.GetUserCalls(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody []byte = w.Body.Bytes()
			if tt.expectedStatus == http.StatusOK {
				var responses []entity.CallResponse
				err := json.Unmarshal(responseBody, &responses)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockGetCallsRes), len(responses))
			} else {
				var response apierrors.Response
				err := json.Unmarshal(responseBody, &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}

			if tt.shouldCallMock {
				mockUseCase.AssertExpectations(t)
			} else {
				mockUseCase.AssertNotCalled(t, "GetUserCalls")
			}
		})
	}
}

func TestUpdateCallStatus(t *testing.T) {
	tests := []struct {
		name               string
		callIDParam        string
		inputBody          entity.UpdateCallStatusDTO
		mockUpdateErr      error
		expectedStatus     int
		expectedResponse   apierrors.Response
		setupContext       func(c *gin.Context)
		shouldCallMock     bool
		expectedCallID     int64
		expectedStatusText string
	}{
		{
			name:        "Successful update",
			callIDParam: "1",
			inputBody: entity.UpdateCallStatusDTO{
				Status: "открыта",
			},
			mockUpdateErr:      nil,
			expectedStatus:     http.StatusNoContent,
			setupContext:       func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:     true,
			expectedCallID:     1,
			expectedStatusText: "открыта",
		},
		{
			name:             "Unauthorized (missing user ID)",
			callIDParam:      "1",
			inputBody:        entity.UpdateCallStatusDTO{Status: "открыта"},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: apierrors.Response{Error: "Unauthorized"},
			setupContext:     func(c *gin.Context) {}, // no id
			shouldCallMock:   false,
		},
		{
			name:             "Invalid user ID type",
			callIDParam:      "1",
			inputBody:        entity.UpdateCallStatusDTO{Status: "открыта"},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: apierrors.Response{Error: "Invalid user ID in context"},
			setupContext:     func(c *gin.Context) { c.Set("id", "not-an-int64") },
			shouldCallMock:   false,
		},
		{
			name:             "Invalid call ID param",
			callIDParam:      "abc",
			inputBody:        entity.UpdateCallStatusDTO{Status: "открыта"},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: apierrors.Response{Error: "Invalid call ID"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   false,
		},
		{
			name:             "Invalid request format (bad JSON)",
			callIDParam:      "1",
			inputBody:        entity.UpdateCallStatusDTO{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: apierrors.Response{Error: "Invalid request format"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   false,
		},
		{
			name:             "Invalid status value",
			callIDParam:      "1",
			inputBody:        entity.UpdateCallStatusDTO{Status: "приоткрыта"},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: apierrors.Response{Error: "Invalid status value"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   false,
		},
		{
			name:               "Call not found",
			callIDParam:        "1",
			inputBody:          entity.UpdateCallStatusDTO{Status: "открыта"},
			mockUpdateErr:      usecase.ErrCallNotFound,
			expectedStatus:     http.StatusNotFound,
			expectedResponse:   apierrors.Response{Error: "Call not found or does not belong to user"},
			setupContext:       func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:     true,
			expectedCallID:     1,
			expectedStatusText: "открыта",
		},
		{
			name:               "Internal server error",
			callIDParam:        "1",
			inputBody:          entity.UpdateCallStatusDTO{Status: "открыта"},
			mockUpdateErr:      errors.New("db error"),
			expectedStatus:     http.StatusInternalServerError,
			expectedResponse:   apierrors.Response{Error: "Failed to update call status"},
			setupContext:       func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:     true,
			expectedCallID:     1,
			expectedStatusText: "открыта",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mocks.NewMockUseCase(t)
			if tt.shouldCallMock {
				mockUseCase.On("UpdateCallStatus", mock.Anything, tt.expectedCallID, int64(123), tt.expectedStatusText).
					Return(tt.mockUpdateErr)
			}

			var requestBody *bytes.Buffer
			if tt.name == "Invalid request format (bad JSON)" {
				requestBody = bytes.NewBuffer([]byte("{invalid-json}"))
			} else {
				body, err := json.Marshal(tt.inputBody)
				assert.NoError(t, err)
				requestBody = bytes.NewBuffer(body)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			tt.setupContext(c)
			c.Params = []gin.Param{{Key: "id", Value: tt.callIDParam}}
			c.Request = httptest.NewRequest("PATCH", "/calls/"+tt.callIDParam, requestBody)

			handler := controller.New(mockUseCase, zerolog.Nop())

			handler.UpdateCallStatus(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusNoContent {
				var response apierrors.Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			} else {
				assert.Empty(t, w.Body.Bytes())
			}

			if tt.shouldCallMock {
				mockUseCase.AssertExpectations(t)
			} else {
				mockUseCase.AssertNotCalled(t, "UpdateCallStatus")
			}
		})
	}
}

func TestGetUserCallByID(t *testing.T) {
	tests := []struct {
		name             string
		callIDParam      string
		mockGetCall      *entity.CallResponse
		mockGetCallErr   error
		expectedStatus   int
		expectedResponse any
		setupContext     func(c *gin.Context)
		shouldCallMock   bool
		expectedCallID   int64
	}{
		{
			name:        "Successful fetch",
			callIDParam: "1",
			mockGetCall: &entity.CallResponse{
				ID:     1,
				Status: "открыта",
			},
			mockGetCallErr:   nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: entity.CallResponse{ID: 1, Status: "открыта"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   true,
			expectedCallID:   1,
		},
		{
			name:             "Unauthorized (missing user ID)",
			callIDParam:      "1",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: apierrors.Response{Error: "Unauthorized"},
			setupContext:     func(c *gin.Context) {},
			shouldCallMock:   false,
		},
		{
			name:             "Invalid user ID type",
			callIDParam:      "1",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: apierrors.Response{Error: "Invalid user ID in context"},
			setupContext:     func(c *gin.Context) { c.Set("id", "not-an-int64") },
			shouldCallMock:   false,
		},
		{
			name:             "Invalid call ID param",
			callIDParam:      "qwerty",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: apierrors.Response{Error: "Invalid call ID"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   false,
		},
		{
			name:             "Call not found",
			callIDParam:      "1",
			mockGetCallErr:   usecase.ErrCallNotFound,
			expectedStatus:   http.StatusNotFound,
			expectedResponse: apierrors.Response{Error: "Call not found"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   true,
			expectedCallID:   1,
		},
		{
			name:             "Internal server error",
			callIDParam:      "1",
			mockGetCallErr:   errors.New("db error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: apierrors.Response{Error: "Failed to get user call"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   true,
			expectedCallID:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mocks.NewMockUseCase(t)
			if tt.shouldCallMock {
				mockUseCase.On("GetUserCallByID", mock.Anything, tt.expectedCallID, int64(123)).
					Return(tt.mockGetCall, tt.mockGetCallErr)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			tt.setupContext(c)
			c.Params = []gin.Param{{Key: "id", Value: tt.callIDParam}}
			c.Request = httptest.NewRequest("GET", "/calls/"+tt.callIDParam, nil)

			handler := controller.New(mockUseCase, zerolog.Nop())

			handler.GetUserCallByID(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entity.CallResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, *tt.mockGetCall, response)
			} else {
				var response apierrors.Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}

			if tt.shouldCallMock {
				mockUseCase.AssertExpectations(t)
			} else {
				mockUseCase.AssertNotCalled(t, "GetUserCallByID")
			}
		})
	}
}

func TestDeleteCall(t *testing.T) {
	tests := []struct {
		name             string
		callIDParam      string
		mockDeleteErr    error
		expectedStatus   int
		expectedResponse any
		setupContext     func(c *gin.Context)
		shouldCallMock   bool
		expectedCallID   int64
	}{
		{
			name:           "Successful delete",
			callIDParam:    "1",
			mockDeleteErr:  nil,
			expectedStatus: http.StatusNoContent,
			setupContext:   func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock: true,
			expectedCallID: 1,
		},
		{
			name:             "Unauthorized (missing user ID)",
			callIDParam:      "1",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: apierrors.Response{Error: "Unauthorized"},
			setupContext:     func(c *gin.Context) {},
			shouldCallMock:   false,
		},
		{
			name:             "Invalid user ID type",
			callIDParam:      "1",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: apierrors.Response{Error: "Invalid user ID in context"},
			setupContext:     func(c *gin.Context) { c.Set("id", "not-an-int64") },
			shouldCallMock:   false,
		},
		{
			name:             "Invalid call ID param",
			callIDParam:      "abc",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: apierrors.Response{Error: "Invalid call ID"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   false,
		},
		{
			name:             "Call not found",
			callIDParam:      "1",
			mockDeleteErr:    usecase.ErrCallNotFound,
			expectedStatus:   http.StatusNotFound,
			expectedResponse: apierrors.Response{Error: "Call not found or does not belong to user"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   true,
			expectedCallID:   1,
		},
		{
			name:             "Internal server error",
			callIDParam:      "1",
			mockDeleteErr:    errors.New("db error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: apierrors.Response{Error: "Failed to delete call"},
			setupContext:     func(c *gin.Context) { c.Set("id", int64(123)) },
			shouldCallMock:   true,
			expectedCallID:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mocks.NewMockUseCase(t)
			if tt.shouldCallMock {
				mockUseCase.On("DeleteCall", mock.Anything, tt.expectedCallID, int64(123)).
					Return(tt.mockDeleteErr)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			tt.setupContext(c)
			c.Params = []gin.Param{{Key: "id", Value: tt.callIDParam}}
			c.Request = httptest.NewRequest("DELETE", "/calls/"+tt.callIDParam, nil)

			handler := controller.New(mockUseCase, zerolog.Nop())

			handler.DeleteCall(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusNoContent {
				var response apierrors.Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			} else {
				assert.Empty(t, w.Body.Bytes())
			}

			if tt.shouldCallMock {
				mockUseCase.AssertExpectations(t)
			} else {
				mockUseCase.AssertNotCalled(t, "DeleteCall")
			}
		})
	}
}
