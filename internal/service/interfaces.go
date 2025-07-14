package service

import (
	"meteo-api/internal/models"
)

// MeteoServiceInterface defines the interface for meteo service operations
type MeteoServiceInterface interface {
	ProcessRawData(rawData *models.RawMeteoData) (*models.MeteoMeasurement, error)
	ProcessRawDataBulk(rawDataList []models.RawMeteoData) ([]models.MeteoMeasurement, error)
	GetLatestMeasurement() (*models.MeteoMeasurement, error)
	GetHistory(period string, limit int) ([]models.MeteoMeasurement, error)
	GetAggregatedData(period, interval string) ([]models.MeteoAggregation, error)
	GetAverage(period string) (*models.MeteoAggregation, error)
}