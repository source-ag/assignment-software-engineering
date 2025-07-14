package repository

import (
	"time"

	"meteo-api/internal/models"
)

// MeteoRepositoryInterface defines the interface for meteo repository operations
type MeteoRepositoryInterface interface {
	CreateMeasurement(measurement *models.MeteoMeasurement) error
	CreateMeasurementsBulk(measurements []models.MeteoMeasurement) error
	GetLatestMeasurement() (*models.MeteoMeasurement, error)
	GetMeasurementsByPeriod(period time.Duration, limit int) ([]models.MeteoMeasurement, error)
	GetAggregatedData(period time.Duration, interval time.Duration) ([]models.MeteoAggregation, error)
	GetOverallAverage(period time.Duration) (*models.MeteoAggregation, error)
}