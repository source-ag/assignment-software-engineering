package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"meteo-api/internal/middleware"
	"meteo-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMeteoService is a mock implementation of MeteoService
type MockMeteoService struct {
	mock.Mock
}

func (m *MockMeteoService) ProcessRawData(rawData *models.RawMeteoData) (*models.MeteoMeasurement, error) {
	args := m.Called(rawData)
	return args.Get(0).(*models.MeteoMeasurement), args.Error(1)
}

func (m *MockMeteoService) ProcessRawDataBulk(rawDataList []models.RawMeteoData) ([]models.MeteoMeasurement, error) {
	args := m.Called(rawDataList)
	return args.Get(0).([]models.MeteoMeasurement), args.Error(1)
}

func (m *MockMeteoService) GetLatestMeasurement() (*models.MeteoMeasurement, error) {
	args := m.Called()
	return args.Get(0).(*models.MeteoMeasurement), args.Error(1)
}

func (m *MockMeteoService) GetHistory(period string, limit int) ([]models.MeteoMeasurement, error) {
	args := m.Called(period, limit)
	return args.Get(0).([]models.MeteoMeasurement), args.Error(1)
}

func (m *MockMeteoService) GetAggregatedData(period, interval string) ([]models.MeteoAggregation, error) {
	args := m.Called(period, interval)
	return args.Get(0).([]models.MeteoAggregation), args.Error(1)
}

func (m *MockMeteoService) GetAverage(period string) (*models.MeteoAggregation, error) {
	args := m.Called(period)
	return args.Get(0).(*models.MeteoAggregation), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockMeteoService) {
	gin.SetMode(gin.TestMode)
	
	mockService := new(MockMeteoService)
	authMiddleware := middleware.NewAuthMiddleware()
	handler := NewMeteoHandler(mockService, authMiddleware)
	
	router := gin.New()
	handler.RegisterRoutes(router)
	
	return router, mockService
}

func TestMeteoHandler_CreateMeasurement(t *testing.T) {
	router, mockService := setupTestRouter()

	// Test data
	rawData := models.RawMeteoData{
		Name: "_ws_source_meteo",
		Rows: [][]interface{}{
			{"external_temperature_c", 9.1387},
		},
		TS: "2021-05-01T12:07:50+02:00",
		PT: 0,
	}

	expectedMeasurement := &models.MeteoMeasurement{
		ID:        "test-id",
		Timestamp: time.Now(),
	}

	// Mock expectations
	mockService.On("ProcessRawData", &rawData).Return(expectedMeasurement, nil)

	// Prepare request
	jsonData, _ := json.Marshal(rawData)
	req, _ := http.NewRequest("POST", "/api/v1/meteo/measurements", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer write-key-456")

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Measurement created successfully", response.Message)
	
	mockService.AssertExpectations(t)
}

func TestMeteoHandler_CreateMeasurement_Unauthorized(t *testing.T) {
	router, _ := setupTestRouter()

	rawData := models.RawMeteoData{
		Name: "_ws_source_meteo",
		TS:   "2021-05-01T12:07:50+02:00",
	}

	jsonData, _ := json.Marshal(rawData)
	req, _ := http.NewRequest("POST", "/api/v1/meteo/measurements", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Authorization header is required", response.Error)
}

func TestMeteoHandler_CreateMeasurement_Forbidden(t *testing.T) {
	router, _ := setupTestRouter()

	rawData := models.RawMeteoData{
		Name: "_ws_source_meteo",
		TS:   "2021-05-01T12:07:50+02:00",
	}

	jsonData, _ := json.Marshal(rawData)
	req, _ := http.NewRequest("POST", "/api/v1/meteo/measurements", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer read-key-123") // Read key for write endpoint

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Insufficient permissions", response.Error)
}

func TestMeteoHandler_GetCurrentMeasurement(t *testing.T) {
	router, mockService := setupTestRouter()

	expectedMeasurement := &models.MeteoMeasurement{
		ID:        "test-id",
		Timestamp: time.Now(),
	}

	mockService.On("GetLatestMeasurement").Return(expectedMeasurement, nil)

	req, _ := http.NewRequest("GET", "/api/v1/meteo/current", nil)
	req.Header.Set("Authorization", "Bearer read-key-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestMeteoHandler_GetCurrentMeasurement_NotFound(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("GetLatestMeasurement").Return((*models.MeteoMeasurement)(nil), nil)

	req, _ := http.NewRequest("GET", "/api/v1/meteo/current", nil)
	req.Header.Set("Authorization", "Bearer read-key-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "No measurements found", response.Error)
	
	mockService.AssertExpectations(t)
}

func TestMeteoHandler_GetHistory(t *testing.T) {
	router, mockService := setupTestRouter()

	expectedMeasurements := []models.MeteoMeasurement{
		{
			ID:        "test-id",
			Timestamp: time.Now(),
		},
	}

	mockService.On("GetHistory", "24h", 0).Return(expectedMeasurements, nil)

	req, _ := http.NewRequest("GET", "/api/v1/meteo/history?period=24h", nil)
	req.Header.Set("Authorization", "Bearer read-key-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestMeteoHandler_GetAggregatedData(t *testing.T) {
	router, mockService := setupTestRouter()

	expectedAggregations := []models.MeteoAggregation{
		{
			Timestamp: time.Now(),
			Count:     10,
		},
	}

	mockService.On("GetAggregatedData", "24h", "15m").Return(expectedAggregations, nil)

	req, _ := http.NewRequest("GET", "/api/v1/meteo/aggregated?period=24h&interval=15m", nil)
	req.Header.Set("Authorization", "Bearer read-key-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestMeteoHandler_GetAverage(t *testing.T) {
	router, mockService := setupTestRouter()

	expectedAverage := &models.MeteoAggregation{
		Timestamp: time.Now(),
		Count:     100,
	}

	mockService.On("GetAverage", "24h").Return(expectedAverage, nil)

	req, _ := http.NewRequest("GET", "/api/v1/meteo/average?period=24h", nil)
	req.Header.Set("Authorization", "Bearer read-key-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	mockService.AssertExpectations(t)
}

func TestMeteoHandler_HealthCheck(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/health", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Service is healthy", response.Message)
}