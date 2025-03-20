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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.POST("/transactions", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.CreateTransaction(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"transaction_nik": "1234567890123456",
		"transaction_amount": 5000000,
		"transaction_otr": 5500000,
		"transaction_admin_fee": 200000,
		"transaction_installment": 12,
		"transaction_interest": 2.5,
		"transaction_asset_name": "Motorcycle"
	}`
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	transactionUsecase.On("GetCustomerByNIK", "1234567890123456").Return(&domain.Customer{CustomerNIK: "1234567890123456"}, nil)
	transactionUsecase.On("CreateTransactionWithLimitUpdate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Transaction created successfully"`)
}

func TestCreateTransaction_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.POST("/transactions", transactionHandler.CreateTransaction)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
}

func TestCreateTransaction_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.POST("/transactions", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.CreateTransaction(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{}`
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error"`, "Response should contain error message")
}

func TestCreateTransaction_CustomerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.POST("/transactions", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.CreateTransaction(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"transaction_nik": "9999999999999999",
		"transaction_amount": 5000000,
		"transaction_otr": 5500000,
		"transaction_admin_fee": 200000,
		"transaction_installment": 12,
		"transaction_interest": 2.5,
		"transaction_asset_name": "Motorcycle"
	}`
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	transactionUsecase.On("GetCustomerByNIK", "9999999999999999").Return(nil, errors.New("customer not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestCreateTransaction_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.POST("/transactions", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.CreateTransaction(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"transaction_nik": "1234567890123456",
		"transaction_amount": 5000000,
		"transaction_otr": 5500000,
		"transaction_admin_fee": 200000,
		"transaction_installment": 12,
		"transaction_interest": 2.5,
		"transaction_asset_name": "Motorcycle"
	}`
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	transactionUsecase.On("GetCustomerByNIK", "1234567890123456").Return(&domain.Customer{CustomerNIK: "1234567890123456"}, nil)
	transactionUsecase.On("CreateTransactionWithLimitUpdate", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestGetTransaction_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.GET("/transactions", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.GetTransaction(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions?limit=10&page=1", nil)

	mockTransactions := []domain.Transaction{
		{TransactionID: 1, TransactionContractNumber: "TX123456", TransactionOTR: 5000000, TransactionInstallment: 500000},
		{TransactionID: 2, TransactionContractNumber: "TX654321", TransactionOTR: 10000000, TransactionInstallment: 800000},
	}
	transactionUsecase.On("GetAllTransactions", 10, 0).Return(mockTransactions, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"transactions"`)
	assert.Contains(t, w.Body.String(), `"TX123456"`)
	assert.Contains(t, w.Body.String(), `"TX654321"`)
}

func TestGetTransaction_InvalidLimitOrPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.GET("/transactions", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.GetTransaction(c)
	})

	testCases := []struct {
		url              string
		expectedErrorMsg string
	}{
		{"/transactions?limit=-1&page=1", `"error":"Invalid limit"`},
		{"/transactions?limit=0&page=1", `"error":"Invalid limit"`},
		{"/transactions?limit=10&page=0", `"error":"Invalid page value"`},
		{"/transactions?limit=10&page=-5", `"error":"Invalid page value"`},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", tc.url, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request for "+tc.url)
		assert.Contains(t, w.Body.String(), tc.expectedErrorMsg, "Response body does not match expected error message")
	}
}

func TestGetTransaction_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.GET("/transactions", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.GetTransaction(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions?limit=10&page=1", nil)

	transactionUsecase.On("GetAllTransactions", 10, 0).Return(nil, errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to retrieve transactions"`)
}

func TestGetTransactionByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.GET("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.GetTransactionByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions/1", nil)

	transactionUsecase.On("GetTransactionByID", uint(1)).Return(&domain.Transaction{
		TransactionID:             1,
		TransactionContractNumber: "TX123456",
		TransactionNIK:            "1234567890123456",
		NIKCustomer: domain.Customer{
			CustomerNIK:      "1234567890123456",
			CustomerFullName: "John Doe",
		},
		TransactionLimit: 1,
		IDLimit: domain.Limit{
			LimitID:     1,
			LimitAmount: 5000000,
		},
		TransactionOTR:         5000000,
		TransactionAdminFee:    250000,
		TransactionInstallment: 1250000,
		TransactionInterest:    5.5,
		TransactionAssetName:   "Motorcycle",
		TransactionDate:        time.Now(),
		TransactionCreatedBy:   1,
		CreatedByUser: domain.User{
			UserID:       1,
			UserUsername: "admin",
			UserRole:     "admin",
		},
		TransactionCreatedAt: time.Now(),
		TransactionEditedBy:  nil,
		TransactionEditedAt:  nil,
	}, nil)

	router.ServeHTTP(w, req)

	expectedResponse := `{
		"transaction_id": 1,
		"transaction_contract_number": "TX123456",
		"transaction_nik": "1234567890123456",
		"NIKCustomer": {
			"customer_nik": "1234567890123456",
			"customer_full_name": "John Doe"
		},
		"transaction_limit": 1,
		"IDLimit": {
			"limit_id": 1,
			"limit_amount": 5000000
		},
		"transaction_otr": 5000000,
		"transaction_admin_fee": 250000,
		"transaction_installment": 1250000,
		"transaction_interest": 5.5,
		"transaction_asset_name": "Motorcycle",
		"transaction_created_by": 1,
		"CreatedByUser": {
			"user_id": 1,
			"user_username": "admin",
			"user_role": "admin"
		}
	}`

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.JSONEq(t, expectedResponse, w.Body.String(), "Response JSON should match")
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.GET("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.GetTransactionByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions/999", nil)

	transactionUsecase.On("GetTransactionByID", uint(999)).Return(nil, errors.New("transaction not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Transaction not found"`)
}

func TestGetTransactionByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.GET("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.GetTransactionByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions/invalidID", nil)

	router.ServeHTTP(w, req)

	// ✅ Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error":"Invalid transaction ID"`)
}

func TestGetTransactionByID_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.GET("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		transactionHandler.GetTransactionByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions/1", nil)

	transactionUsecase.On("GetTransactionByID", uint(1)).Return(nil, errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Transaction not found"`)
}

func TestUpdateTransaction_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.PUT("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		transactionHandler.UpdateTransaction(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"transaction_nik": "1234567890123456",
		"transaction_amount": 7000000,
		"transaction_otr": 8000000,
		"transaction_admin_fee": 500000,
		"transaction_installment": 12,
		"transaction_interest": 5.0,
		"transaction_asset_name": "Motorcycle"
	}`
	req, _ := http.NewRequest("PUT", "/transactions/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	// ✅ Mock Usecase Calls
	transactionUsecase.On("GetTransactionByID", uint(1)).Return(&domain.Transaction{
		TransactionID:             1,
		TransactionContractNumber: "TX123456",
		TransactionNIK:            "1234567890123456",
	}, nil)

	transactionUsecase.On("GetCustomerByNIK", "1234567890123456").Return(&domain.Customer{
		CustomerNIK: "1234567890123456",
	}, nil)

	transactionUsecase.On("UpdateTransactionWithLimitUpdate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Transaction edited successfully"`)
}

func TestUpdateTransaction_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.PUT("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		transactionHandler.UpdateTransaction(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"transaction_nik": "1234567890123456",
		"transaction_amount": 7000000
	}`
	req, _ := http.NewRequest("PUT", "/transactions/999", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	transactionUsecase.On("GetTransactionByID", uint(999)).Return(nil, errors.New("transaction not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Transaction not found"`)
}

func TestUpdateTransaction_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.PUT("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		transactionHandler.UpdateTransaction(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"transaction_nik": "",
		"transaction_amount": -500
	}`
	req, _ := http.NewRequest("PUT", "/transactions/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	transactionUsecase.On("GetTransactionByID", uint(1)).Return(&domain.Transaction{
		TransactionID:             1,
		TransactionContractNumber: "TX123456",
		TransactionNIK:            "1234567890123456",
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestUpdateTransaction_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.PUT("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		transactionHandler.UpdateTransaction(c)
	})

	w := httptest.NewRecorder()
	reqBody := `{
		"transaction_nik": "1234567890123456",
		"transaction_amount": 5000000,
		"transaction_otr": 8000000,
		"transaction_admin_fee": 500000,
		"transaction_installment": 12,
		"transaction_interest": 5.0,
		"transaction_asset_name": "Motorcycle"
	}`
	req, _ := http.NewRequest("PUT", "/transactions/1", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	transactionUsecase.On("GetTransactionByID", uint(1)).Return(&domain.Transaction{
		TransactionID:             1,
		TransactionContractNumber: "TX123456",
		TransactionNIK:            "1234567890123456",
	}, nil)

	transactionUsecase.On("GetCustomerByNIK", "1234567890123456").Return(&domain.Customer{
		CustomerNIK: "1234567890123456",
	}, nil)

	transactionUsecase.On("UpdateTransactionWithLimitUpdate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to update transaction"`)
}

func TestDeleteTransaction_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.DELETE("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		transactionHandler.DeleteTransaction(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/transactions/1", nil)

	transactionUsecase.On("GetTransactionByID", uint(1)).Return(&domain.Transaction{
		TransactionContractNumber: "TX123456",
	}, nil)
	transactionUsecase.On("DeleteTransactionWithLimitUpdate", mock.Anything, mock.Anything).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Transaction deleted successfully"`)
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.DELETE("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		transactionHandler.DeleteTransaction(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/transactions/999", nil)

	transactionUsecase.On("GetTransactionByID", uint(999)).Return(nil, errors.New("transaction not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Transaction not found"`)
}

func TestDeleteTransaction_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUsecase := new(mocks.TransactionUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	router.DELETE("/transactions/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1, UserRole: "admin"})
		transactionHandler.DeleteTransaction(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/transactions/1", nil)

	transactionUsecase.On("GetTransactionByID", uint(1)).Return(&domain.Transaction{
		TransactionContractNumber: "TX123456",
	}, nil)
	transactionUsecase.On("DeleteTransactionWithLimitUpdate", mock.Anything, mock.Anything).Return(errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to delete transaction"`)
}
