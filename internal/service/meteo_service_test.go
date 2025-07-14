package service

import (
	"testing"
	"time"

	"meteo-api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMeteoRepository is a mock implementation of MeteoRepository
type MockMeteoRepository struct {
	mock.Mock
}

func (m *MockMeteoRepository) CreateMeasurement(measurement *models.MeteoMeasurement) error {
	args := m.Called(measurement)
	return args.Error(0)
}

func (m *MockMeteoRepository) CreateMeasurementsBulk(measurements []models.MeteoMeasurement) error {
	args := m.Called(measurements)
	return args.Error(0)
}

func (m *MockMeteoRepository) GetLatestMeasurement() (*models.MeteoMeasurement, error) {
	args := m.Called()
	return args.Get(0).(*models.MeteoMeasurement), args.Error(1)
}

func (m *MockMeteoRepository) GetMeasurementsByPeriod(period time.Duration, limit int) ([]models.MeteoMeasurement, error) {
	args := m.Called(period, limit)
	return args.Get(0).([]models.MeteoMeasurement), args.Error(1)
}

func (m *MockMeteoRepository) GetAggregatedData(period time.Duration, interval time.Duration) ([]models.MeteoAggregation, error) {
	args := m.Called(period, interval)
	return args.Get(0).([]models.MeteoAggregation), args.Error(1)
}

func (m *MockMeteoRepository) GetOverallAverage(period time.Duration) (*models.MeteoAggregation, error) {
	args := m.Called(period)
	return args.Get(0).(*models.MeteoAggregation), args.Error(1)
}

func TestMeteoService_ProcessRawData(t *testing.T) {
	mockRepo := new(MockMeteoRepository)
	service := &MeteoService{repo: mockRepo}

	// Test data
	rawData := &models.RawMeteoData{
		Name: "_ws_source_meteo",
		Rows: [][]interface{}{
			{"Variable", "Value"},
			{"external_temperature_c", 9.1387},
			{"wind_speed_m_s", 3.06032},
			{"relative_humidity_perc", 73.0},
		},
		TS: "2021-05-01T12:07:50+02:00",
		PT: 0,
	}

	// Mock expectations
	mockRepo.On("CreateMeasurement", mock.AnythingOfType("*models.MeteoMeasurement")).Return(nil)

	// Execute
	result, err := service.ProcessRawData(rawData)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 9.1387, *result.ExternalTemperatureC)
	assert.Equal(t, 3.06032, *result.WindSpeedMS)
	assert.Equal(t, 73.0, *result.RelativeHumidityPerc)
	
	mockRepo.AssertExpectations(t)
}

func TestMeteoService_ParsePeriod(t *testing.T) {
	service := &MeteoService{}

	tests := []struct {
		name     string
		period   string
		expected time.Duration
		hasError bool
	}{
		{
			name:     "Valid hours",
			period:   "24h",
			expected: 24 * time.Hour,
			hasError: false,
		},
		{
			name:     "Valid days",
			period:   "7d",
			expected: 7 * 24 * time.Hour,
			hasError: false,
		},
		{
			name:     "Valid minutes",
			period:   "15m",
			expected: 15 * time.Minute,
			hasError: false,
		},
		{
			name:     "Valid weeks",
			period:   "2w",
			expected: 2 * 7 * 24 * time.Hour,
			hasError: false,
		},
		{
			name:     "Invalid format",
			period:   "invalid",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Empty period",
			period:   "",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.parsePeriod(tt.period)
			
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMeteoService_ParseFloat(t *testing.T) {
	service := &MeteoService{}

	tests := []struct {
		name     string
		value    interface{}
		expected float64
		valid    bool
	}{
		{"Float64", 123.45, 123.45, true},
		{"Float32", float32(123.45), float64(float32(123.45)), true},
		{"Int", 123, 123.0, true},
		{"Int64", int64(123), 123.0, true},
		{"String valid", "123.45", 123.45, true},
		{"String invalid", "invalid", 0, false},
		{"Nil", nil, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, valid := service.parseFloat(tt.value)
			assert.Equal(t, tt.valid, valid)
			if tt.valid {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMeteoService_ParseEnum(t *testing.T) {
	service := &MeteoService{}

	tests := []struct {
		name     string
		value    interface{}
		expected string
		valid    bool
	}{
		{
			name:     "String value",
			value:    "ZW",
			expected: "ZW",
			valid:    true,
		},
		{
			name: "Map with value",
			value: map[string]interface{}{
				"type":  "hortimax.synopta.enum",
				"key":   8782,
				"value": "ZW",
			},
			expected: "ZW",
			valid:    true,
		},
		{
			name:     "Invalid type",
			value:    123,
			expected: "",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, valid := service.parseEnum(tt.value)
			assert.Equal(t, tt.valid, valid)
			if tt.valid {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMeteoService_GetHistory(t *testing.T) {
	mockRepo := new(MockMeteoRepository)
	service := &MeteoService{repo: mockRepo}

	expectedMeasurements := []models.MeteoMeasurement{
		{
			ID:        "test-id",
			Timestamp: time.Now(),
		},
	}

	mockRepo.On("GetMeasurementsByPeriod", 24*time.Hour, 100).Return(expectedMeasurements, nil)

	result, err := service.GetHistory("24h", 100)

	assert.NoError(t, err)
	assert.Equal(t, expectedMeasurements, result)
	mockRepo.AssertExpectations(t)
}

func TestMeteoService_GetAggregatedData(t *testing.T) {
	mockRepo := new(MockMeteoRepository)
	service := &MeteoService{repo: mockRepo}

	expectedAggregations := []models.MeteoAggregation{
		{
			Timestamp: time.Now(),
			Count:     10,
		},
	}

	mockRepo.On("GetAggregatedData", 24*time.Hour, 15*time.Minute).Return(expectedAggregations, nil)

	result, err := service.GetAggregatedData("24h", "15m")

	assert.NoError(t, err)
	assert.Equal(t, expectedAggregations, result)
	mockRepo.AssertExpectations(t)
}