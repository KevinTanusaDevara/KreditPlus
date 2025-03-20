package handler_test

import (
	"bytes"
	"errors"
	"kreditplus/internal/domain"
	"kreditplus/internal/handler"
	"kreditplus/internal/usecase/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLimit_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.POST("/limits", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.CreateLimit(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"limit_nik": "1234567890123456",
		"limit_tenor": 12,
		"limit_amount": 50000000
	}`
	req, _ := http.NewRequest("POST", "/limits", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	limitUsecase.On("GetCustomerByNIK", "1234567890123456").Return(&domain.Customer{CustomerNIK: "1234567890123456"}, nil)
	limitUsecase.On("CreateLimit", mock.Anything).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Limit created successfully"`)
}

func TestCreateLimit_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.POST("/limits", limitHandler.CreateLimit)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/limits", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
}

func TestCreateLimit_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.POST("/limits", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.CreateLimit(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{}`
	req, _ := http.NewRequest("POST", "/limits", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestCreateLimit_CustomerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.POST("/limits", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.CreateLimit(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"limit_nik": "9999999999999999",
		"limit_tenor": 12,
		"limit_amount": 50000000
	}`
	req, _ := http.NewRequest("POST", "/limits", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	limitUsecase.On("GetCustomerByNIK", "9999999999999999").Return(nil, errors.New("customer not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Customer NIK not found"`)
}

func TestCreateLimit_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.POST("/limits", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.CreateLimit(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"limit_nik": "1234567890123456",
		"limit_tenor": 12,
		"limit_amount": 50000000
	}`
	req, _ := http.NewRequest("POST", "/limits", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	limitUsecase.On("GetCustomerByNIK", "1234567890123456").Return(&domain.Customer{CustomerNIK: "1234567890123456"}, nil)
	limitUsecase.On("CreateLimit", mock.Anything).Return(errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to create limit"`)
}

func TestGetLimit_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.GET("/limits", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.GetLimit(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/limits?limit=10&page=1", nil)

	limitUsecase.On("GetAllLimits", 10, 0).Return([]domain.Limit{
		{LimitID: 1, LimitNIK: "1234567890123456", LimitAmount: 50000000},
		{LimitID: 2, LimitNIK: "9876543210987654", LimitAmount: 75000000},
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"limits"`)
	assert.Contains(t, w.Body.String(), `"1234567890123456"`)
	assert.Contains(t, w.Body.String(), `"9876543210987654"`)
}

func TestGetLimit_InvalidLimitOrPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.GET("/limits", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.GetLimit(c)
	})

	testCases := []struct {
		url      string
		errorMsg string
	}{
		{"/limits?limit=-1&page=1", `"error":"Invalid limit"`},
		{"/limits?limit=0&page=1", `"error":"Invalid limit"`},
		{"/limits?limit=10&page=0", `"error":"Invalid page value"`},
		{"/limits?limit=10&page=-5", `"error":"Invalid page value"`},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", tc.url, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
		assert.Contains(t, w.Body.String(), tc.errorMsg)
	}
}

func TestGetLimit_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.GET("/limits", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.GetLimit(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/limits?limit=10&page=1", nil)

	limitUsecase.On("GetAllLimits", 10, 0).Return(nil, errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to retrieve limits"`)
}

func TestGetLimitByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.GET("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.GetLimitByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/limits/1", nil)

	limitUsecase.On("GetLimitByID", uint(1)).Return(&domain.Limit{
		LimitID:     1,
		LimitNIK:    "1234567890123456",
		LimitTenor:  12,
		LimitAmount: 50000000,
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"limit_nik":"1234567890123456"`)
}

func TestGetLimitByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.GET("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.GetLimitByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/limits/999", nil)

	limitUsecase.On("GetLimitByID", uint(999)).Return(nil, errors.New("limit not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Limit not found"`)
}

func TestGetLimitByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.GET("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.GetLimitByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/limits/invalidID", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error":"Invalid limit ID"`)
}

func TestGetLimitByID_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.GET("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.GetLimitByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/limits/1", nil)

	limitUsecase.On("GetLimitByID", uint(1)).Return(nil, errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Limit not found"`)
}

func TestUpdateLimit_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.PUT("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		limitHandler.UpdateLimit(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
        "limit_nik": "1234567890123456",
        "limit_tenor": 24,
        "limit_amount": 75000000,
        "limit_used_amount": 5000000,
        "limit_remaining_amount": 70000000
    }`
	req, _ := http.NewRequest("PUT", "/limits/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	limitUsecase.On("GetLimitByID", uint(1)).Return(&domain.Limit{LimitNIK: "1234567890123456"}, nil)
	limitUsecase.On("UpdateLimit", mock.Anything).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Limit updated successfully"`)
}

func TestUpdateLimit_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.PUT("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		limitHandler.UpdateLimit(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"limit_nik": "1234567890123456",
		"limit_tenor": 24,
		"limit_amount": 75000000
	}`
	req, _ := http.NewRequest("PUT", "/limits/999", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	limitUsecase.On("GetLimitByID", uint(999)).Return(nil, errors.New("limit not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Limit not found"`)
}

func TestUpdateLimit_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.PUT("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		limitHandler.UpdateLimit(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{}`
	req, _ := http.NewRequest("PUT", "/limits/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	limitUsecase.On("GetLimitByID", uint(1)).Return(&domain.Limit{LimitNIK: "1234567890123456"}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")

	responseBody := w.Body.String()

	expectedErrors := []string{
		"Key: 'EditLimitInput.LimitNIK' Error:Field validation for 'LimitNIK' failed on the 'required' tag",
		"Key: 'EditLimitInput.LimitTenor' Error:Field validation for 'LimitTenor' failed on the 'required' tag",
		"Key: 'EditLimitInput.LimitAmount' Error:Field validation for 'LimitAmount' failed on the 'required' tag",
		"Key: 'EditLimitInput.LimitUsedAmount' Error:Field validation for 'LimitUsedAmount' failed on the 'required' tag",
		"Key: 'EditLimitInput.LimitRemainingAmount' Error:Field validation for 'LimitRemainingAmount' failed on the 'required' tag",
	}

	for _, errMsg := range expectedErrors {
		assert.Contains(t, responseBody, errMsg, "Validation error message missing")
	}
}

func TestUpdateLimit_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.PUT("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		limitHandler.UpdateLimit(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"limit_nik": "1234567890123456",
		"limit_tenor": 24,
		"limit_amount": 75000000,
		"limit_used_amount": 5000000,
		"limit_remaining_amount": 70000000
	}`
	req, _ := http.NewRequest("PUT", "/limits/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	limitUsecase.On("GetLimitByID", uint(1)).Return(&domain.Limit{LimitNIK: "1234567890123456"}, nil)
	limitUsecase.On("UpdateLimit", mock.Anything).Return(errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"database error"`)
}

func TestDeleteLimit_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.DELETE("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.DeleteLimit(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/limits/1", nil)

	limitUsecase.On("GetLimitByID", uint(1)).Return(&domain.Limit{
		LimitNIK: "1234567890123456",
	}, nil)
	limitUsecase.On("DeleteLimit", uint(1)).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Limit deleted successfully"`)
}

func TestDeleteLimit_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.DELETE("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.DeleteLimit(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/limits/999", nil)

	limitUsecase.On("GetLimitByID", uint(999)).Return(nil, errors.New("limit not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Limit not found"`)
}

func TestDeleteLimit_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	limitUsecase := new(mocks.LimitUsecase)
	limitHandler := handler.NewLimitHandler(limitUsecase)

	router.DELETE("/limits/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		limitHandler.DeleteLimit(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/limits/1", nil)

	limitUsecase.On("GetLimitByID", uint(1)).Return(&domain.Limit{
		LimitNIK: "1234567890123456",
	}, nil)
	limitUsecase.On("DeleteLimit", uint(1)).Return(errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to delete limit"`)
}

