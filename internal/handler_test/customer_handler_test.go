package handler_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"kreditplus/internal/domain"
	"kreditplus/internal/handler"
	"kreditplus/internal/usecase/mocks"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func addRealFile(writer *multipart.Writer, fieldName, filePath string) {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		panic("Failed to open image file: " + err.Error())
	}
	defer sourceFile.Close()

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filepath.Base(filePath)))
	header.Set("Content-Type", "image/jpeg")

	filePart, err := writer.CreatePart(header)
	if err != nil {
		panic("Failed to create file part: " + err.Error())
	}

	_, err = io.Copy(filePart, sourceFile)
	if err != nil {
		panic("Failed to write file data: " + err.Error())
	}
}

func createTestJPEG() string {
	sourceFilePath := "C:/Users/KevinTanusaDevara/Documents/go/KreditPlus/internal/handler_test/testdata/real_image.jpg"

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		panic("Failed to open source image file: " + err.Error())
	}
	defer sourceFile.Close()

	tmpFile, err := os.CreateTemp("", "*.jpg")
	if err != nil {
		panic("Failed to create temporary file: " + err.Error())
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, sourceFile)
	if err != nil {
		panic("Failed to copy image data: " + err.Error())
	}

	return tmpFile.Name()
}

func TestCreateCustomer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.POST("/customers", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.CreateCustomer(c)
	})

	ktpPath := createTestJPEG()
	selfiePath := createTestJPEG()
	defer os.Remove(ktpPath)
	defer os.Remove(selfiePath)

	w := httptest.NewRecorder()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("customer_nik", "1234567890123456")
	_ = writer.WriteField("customer_full_name", "John Doe")
	_ = writer.WriteField("customer_legal_name", "John D")
	_ = writer.WriteField("customer_birth_place", "Jakarta")
	_ = writer.WriteField("customer_birth_date", "1990-01-01")
	_ = writer.WriteField("customer_salary", "1000000")

	addRealFile(writer, "customer_ktp_photo", ktpPath)
	addRealFile(writer, "customer_selfie_photo", selfiePath)

	writer.Close()

	req, _ := http.NewRequest("POST", "/customers", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	customerUsecase.On("CreateCustomer", mock.Anything).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Customer created successfully"`)
}

func TestCreateCustomer_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.POST("/customers", customerHandler.CreateCustomer)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/customers", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
}

func TestCreateCustomer_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.POST("/customers", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.CreateCustomer(c)
	})

	w := httptest.NewRecorder()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("customer_birth_date", "1990-01-01")

	writer.Close()

	req, _ := http.NewRequest("POST", "/customers", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error"`, "Expected validation error message")
}

func TestCreateCustomer_InvalidBirthDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.POST("/customers", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.CreateCustomer(c)
	})

	w := httptest.NewRecorder()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("customer_nik", "1234567890123456")
	_ = writer.WriteField("customer_full_name", "John Doe")
	_ = writer.WriteField("customer_legal_name", "John D")
	_ = writer.WriteField("customer_birth_place", "Jakarta")
	_ = writer.WriteField("customer_birth_date", "invalid-date")
	_ = writer.WriteField("customer_salary", "1000000")

	writer.Close()

	req, _ := http.NewRequest("POST", "/customers", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error":"Invalid birth date format"`)
}

func TestCreateCustomer_FailedFileUpload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.POST("/customers", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.CreateCustomer(c)
	})

	w := httptest.NewRecorder()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("customer_nik", "1234567890123456")
	_ = writer.WriteField("customer_full_name", "John Doe")
	_ = writer.WriteField("customer_legal_name", "John D")
	_ = writer.WriteField("customer_birth_place", "Jakarta")
	_ = writer.WriteField("customer_birth_date", "1990-01-01")
	_ = writer.WriteField("customer_salary", "1000000")

	writer.Close()

	req, _ := http.NewRequest("POST", "/customers", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error":"Failed to upload KTP photo"`)
}

func TestCreateCustomer_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.POST("/customers", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.CreateCustomer(c)
	})

	w := httptest.NewRecorder()

	ktpPath := createTestJPEG()
	selfiePath := createTestJPEG()
	defer os.Remove(ktpPath)
	defer os.Remove(selfiePath)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("customer_nik", "1234567890123456")
	_ = writer.WriteField("customer_full_name", "John Doe")
	_ = writer.WriteField("customer_legal_name", "John D")
	_ = writer.WriteField("customer_birth_place", "Jakarta")
	_ = writer.WriteField("customer_birth_date", "1990-01-01")
	_ = writer.WriteField("customer_salary", "1000000")

	addRealFile(writer, "customer_ktp_photo", ktpPath)
	addRealFile(writer, "customer_selfie_photo", selfiePath)

	writer.Close()

	req, _ := http.NewRequest("POST", "/customers", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	customerUsecase.On("CreateCustomer", mock.Anything).Return(errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to create customer"`)
}

func TestGetCustomer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.GET("/customers", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.GetCustomer(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customers?limit=10&page=1", nil)

	mockCustomers := []domain.Customer{
		{CustomerNIK: "123456789", CustomerFullName: "John Doe"},
		{CustomerNIK: "987654321", CustomerFullName: "Jane Doe"},
	}
	customerUsecase.On("GetAllCustomers", 10, 0).Return(mockCustomers, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"page":1`)
	assert.Contains(t, w.Body.String(), `"limit":10`)
	assert.Contains(t, w.Body.String(), `"customers"`)
	assert.Contains(t, w.Body.String(), `"John Doe"`)
	assert.Contains(t, w.Body.String(), `"Jane Doe"`)
}

func TestGetCustomer_InvalidLimitOrPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.GET("/customers", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.GetCustomer(c)
	})

	invalidCases := []struct {
		queryParam  string
		expectedMsg string
	}{
		{"limit=-1&page=1", `"error":"Invalid limit"`},
		{"limit=0&page=1", `"error":"Invalid limit"`},
		{"limit=10&page=0", `"error":"Invalid page value"`},
		{"limit=10&page=-5", `"error":"Invalid page value"`},
	}

	for _, tc := range invalidCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/customers?"+tc.queryParam, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
		assert.Contains(t, w.Body.String(), tc.expectedMsg)
	}
}

func TestGetCustomer_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.GET("/customers", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.GetCustomer(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customers?limit=10&page=1", nil)

	customerUsecase.On("GetAllCustomers", 10, 0).Return(nil, errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to retrieve customers"`)
}

func TestGetCustomerByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.GET("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.GetCustomerByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customers/1", nil)

	customerUsecase.On("GetCustomerByID", uint(1)).Return(&domain.Customer{
		CustomerNIK:        "1234567890123456",
		CustomerFullName:   "John Doe",
		CustomerLegalName:  "John D",
		CustomerBirthPlace: "Jakarta",
		CustomerBirthDate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		CustomerSalary:     1000000,
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"customer_nik":"1234567890123456"`)
}

func TestGetCustomerByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.GET("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.GetCustomerByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customers/999", nil)

	customerUsecase.On("GetCustomerByID", uint(999)).Return(nil, errors.New("customer not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Customer not found"`)
}

func TestGetCustomerByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.GET("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.GetCustomerByID(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customers/invalidID", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error":"Invalid customer ID"`)
}

func TestUpdateCustomer_Success(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.PUT("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1}) // ✅ Mock authentication
		customerHandler.UpdateCustomer(c)
	})

	// Create a multipart form body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("customer_nik", "1234567890123456")
	_ = writer.WriteField("customer_full_name", "Updated Name")
	_ = writer.WriteField("customer_legal_name", "Updated Legal Name")
	_ = writer.WriteField("customer_birth_place", "Jakarta")
	_ = writer.WriteField("customer_birth_date", "1990-01-01")
	_ = writer.WriteField("customer_salary", "2000000")

	// Close the multipart writer to finalize the body
	writer.Close()

	// Create the request
	req, _ := http.NewRequest("PUT", "/customers/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType()) // Set the content type to multipart/form-data

	// ✅ Mock the existing customer
	customerUsecase.On("GetCustomerByID", uint(1)).Return(&domain.Customer{
		CustomerNIK: "1234567890123456",
	}, nil)

	// ✅ Mock Successful Update
	customerUsecase.On("UpdateCustomer", mock.Anything).Return(nil)

	// Serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ✅ Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Customer updated successfully"`)
}

func TestUpdateCustomer_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.PUT("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.UpdateCustomer(c)
	})

	w := httptest.NewRecorder()
	body := `{
		"customer_full_name": "Updated Name"
	}`
	req, _ := http.NewRequest("PUT", "/customers/999", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	// ✅ Mock Customer Not Found
	customerUsecase.On("GetCustomerByID", uint(999)).Return(nil, errors.New("customer not found"))

	router.ServeHTTP(w, req)

	// ✅ Assertions
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Customer not found"`)
}

func TestUpdateCustomer_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.PUT("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.UpdateCustomer(c)
	})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("customer_nik", "1234567890123456")
	_ = writer.WriteField("customer_full_name", "")
	_ = writer.WriteField("customer_legal_name", "Legal Name")
	_ = writer.WriteField("customer_birth_place", "Jakarta")
	_ = writer.WriteField("customer_birth_date", "1990-01-01")
	_ = writer.WriteField("customer_salary", "2000000")

	writer.Close()

	req, _ := http.NewRequest("PUT", "/customers/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	customerUsecase.On("GetCustomerByID", uint(1)).Return(&domain.Customer{
		CustomerNIK: "1234567890123456",
	}, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error":"Key: 'CustomerInput.CustomerFullName' Error:Field validation for 'CustomerFullName' failed on the 'required' tag"`)
}

func TestUpdateCustomer_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.PUT("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.UpdateCustomer(c)
	})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("customer_nik", "1234567890123456")
	_ = writer.WriteField("customer_full_name", "Updated Name")
	_ = writer.WriteField("customer_legal_name", "Legal Name")
	_ = writer.WriteField("customer_birth_place", "Jakarta")
	_ = writer.WriteField("customer_birth_date", "1990-01-01")
	_ = writer.WriteField("customer_salary", "2000000")

	writer.Close()

	req, _ := http.NewRequest("PUT", "/customers/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	customerUsecase.On("GetCustomerByID", uint(1)).Return(&domain.Customer{
		CustomerNIK: "1234567890123456",
	}, nil)

	customerUsecase.On("UpdateCustomer", mock.Anything).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"database error"`)
}

func TestUpdateCustomer_FailedFileUpload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.PUT("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.UpdateCustomer(c)
	})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("customer_nik", "1234567890123456")
	_ = writer.WriteField("customer_full_name", "Updated Name")
	_ = writer.WriteField("customer_legal_name", "Legal Name")
	_ = writer.WriteField("customer_birth_place", "Jakarta")
	_ = writer.WriteField("customer_birth_date", "1990-01-01")
	_ = writer.WriteField("customer_salary", "2000000")

	fileWriter, _ := writer.CreateFormFile("customer_ktp_photo", "ktp.jpg")
	fileWriter.Write([]byte("fake file content"))

	customerUsecase.On("SaveUploadedFile", mock.Anything, "customer_ktp_photo", "uploads/ktp").Return("", errors.New("Failed to upload KTP photo"))

	writer.Close()

	req, _ := http.NewRequest("PUT", "/customers/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	customerUsecase.On("GetCustomerByID", uint(1)).Return(&domain.Customer{
		CustomerNIK: "1234567890123456",
	}, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP 400 Bad Request")
	assert.Contains(t, w.Body.String(), `"error":"Failed to upload KTP photo"`)
}

func TestDeleteCustomer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.DELETE("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.DeleteCustomer(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/customers/1", nil)

	customerUsecase.On("GetCustomerByID", uint(1)).Return(&domain.Customer{
		CustomerNIK: "1234567890123456",
	}, nil)
	customerUsecase.On("DeleteCustomer", uint(1)).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"message":"Customer deleted successfully"`)
}

func TestDeleteCustomer_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.DELETE("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.DeleteCustomer(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/customers/999", nil)

	customerUsecase.On("GetCustomerByID", uint(999)).Return(nil, errors.New("customer not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP 404 Not Found")
	assert.Contains(t, w.Body.String(), `"error":"Customer not found"`)
}

func TestDeleteCustomer_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	customerUsecase := new(mocks.CustomerUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router.DELETE("/customers/:id", func(c *gin.Context) {
		c.Set("user", domain.User{UserID: 1})
		customerHandler.DeleteCustomer(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/customers/1", nil)

	customerUsecase.On("GetCustomerByID", uint(1)).Return(&domain.Customer{
		CustomerNIK: "1234567890123456",
	}, nil)
	customerUsecase.On("DeleteCustomer", uint(1)).Return(errors.New("database error"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP 500 Internal Server Error")
	assert.Contains(t, w.Body.String(), `"error":"Failed to delete customer"`)
}
