package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"meteo-api/internal/models"
	"meteo-api/internal/repository"
)

type MeteoService struct {
	repo repository.MeteoRepositoryInterface
}

func NewMeteoService(repo repository.MeteoRepositoryInterface) *MeteoService {
	return &MeteoService{repo: repo}
}

func (s *MeteoService) ProcessRawData(rawData *models.RawMeteoData) (*models.MeteoMeasurement, error) {
	measurement, err := s.parseRawData(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw data: %w", err)
	}

	if err := s.repo.CreateMeasurement(measurement); err != nil {
		return nil, fmt.Errorf("failed to store measurement: %w", err)
	}

	return measurement, nil
}

func (s *MeteoService) ProcessRawDataBulk(rawDataList []models.RawMeteoData) ([]models.MeteoMeasurement, error) {
	var measurements []models.MeteoMeasurement

	for _, rawData := range rawDataList {
		measurement, err := s.parseRawData(&rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse raw data: %w", err)
		}
		measurements = append(measurements, *measurement)
	}

	if err := s.repo.CreateMeasurementsBulk(measurements); err != nil {
		return nil, fmt.Errorf("failed to store measurements: %w", err)
	}

	return measurements, nil
}

func (s *MeteoService) GetLatestMeasurement() (*models.MeteoMeasurement, error) {
	return s.repo.GetLatestMeasurement()
}

func (s *MeteoService) GetHistory(period string, limit int) ([]models.MeteoMeasurement, error) {
	duration, err := s.parsePeriod(period)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	return s.repo.GetMeasurementsByPeriod(duration, limit)
}

func (s *MeteoService) GetAggregatedData(period, interval string) ([]models.MeteoAggregation, error) {
	periodDuration, err := s.parsePeriod(period)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	intervalDuration, err := s.parsePeriod(interval)
	if err != nil {
		return nil, fmt.Errorf("invalid interval: %w", err)
	}

	return s.repo.GetAggregatedData(periodDuration, intervalDuration)
}

func (s *MeteoService) GetAverage(period string) (*models.MeteoAggregation, error) {
	duration, err := s.parsePeriod(period)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	return s.repo.GetOverallAverage(duration)
}

func (s *MeteoService) parseRawData(rawData *models.RawMeteoData) (*models.MeteoMeasurement, error) {
	measurement := &models.MeteoMeasurement{}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, rawData.TS)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	measurement.Timestamp = timestamp

	// Store raw data as JSONB
	rawDataBytes, err := json.Marshal(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal raw data: %w", err)
	}
	measurement.RawData = rawDataBytes

	// Parse rows data
	for _, row := range rawData.Rows {
		if len(row) < 2 {
			continue
		}

		variable, ok := row[0].(string)
		if !ok {
			continue
		}

		value := row[1]

		switch variable {
		case "external_temperature_c":
			if v, ok := s.parseFloat(value); ok {
				measurement.ExternalTemperatureC = &v
			}
		case "wind_speed_unmuted_m_s":
			if v, ok := s.parseFloat(value); ok {
				measurement.WindSpeedUnmutedMS = &v
			}
		case "wind_speed_m_s":
			if v, ok := s.parseFloat(value); ok {
				measurement.WindSpeedMS = &v
			}
		case "wind_direction_degrees":
			if v, ok := s.parseInt(value); ok {
				measurement.WindDirectionDegrees = &v
			}
		case "wind_direction_compass":
			if v, ok := s.parseEnum(value); ok {
				measurement.WindDirectionCompass = &v
			}
		case "radiation_intensity_unmuted_w_m2":
			if v, ok := s.parseFloat(value); ok {
				measurement.RadiationIntensityUnmutedWM2 = &v
			}
		case "radiation_intensity_w_m2":
			if v, ok := s.parseFloat(value); ok {
				measurement.RadiationIntensityWM2 = &v
			}
		case "standard_radiation_intensity_w_m2":
			if v, ok := s.parseFloat(value); ok {
				measurement.StandardRadiationIntensityWM2 = &v
			}
		case "radiation_sum_j_cm2":
			if v, ok := s.parseFloat(value); ok {
				measurement.RadiationSumJCM2 = &v
			}
		case "radiation_from_plant_w_m2":
			if v, ok := s.parseFloat(value); ok {
				measurement.RadiationFromPlantWM2 = &v
			}
		case "precipitation":
			if v, ok := s.parseFloat(value); ok {
				measurement.Precipitation = &v
			}
		case "relative_humidity_perc":
			if v, ok := s.parseFloat(value); ok {
				measurement.RelativeHumidityPerc = &v
			}
		case "moisture_deficit_g_kg":
			if v, ok := s.parseFloat(value); ok {
				measurement.MoistureDeficitGKG = &v
			}
		case "moisture_deficit_g_m3":
			if v, ok := s.parseFloat(value); ok {
				measurement.MoistureDeficitGM3 = &v
			}
		case "dew_point_temperature_c":
			if v, ok := s.parseFloat(value); ok {
				measurement.DewPointTemperatureC = &v
			}
		case "abs_humidity_g_kg":
			if v, ok := s.parseFloat(value); ok {
				measurement.AbsHumidityGKG = &v
			}
		case "enthalpy_kj_kg":
			if v, ok := s.parseFloat(value); ok {
				measurement.EnthalpyKJKG = &v
			}
		case "enthalpy_kj_m3":
			if v, ok := s.parseFloat(value); ok {
				measurement.EnthalpyKJM3 = &v
			}
		case "atmospheric_pressure_hpa":
			if v, ok := s.parseFloat(value); ok {
				measurement.AtmosphericPressureHPA = &v
			}
		case "status_meteo_station":
			if v, ok := s.parseEnum(value); ok {
				measurement.StatusMeteoStation = &v
			}
		case "status_meteo_station_communication":
			if v, ok := s.parseEnum(value); ok {
				measurement.StatusMeteoStationCommunication = &v
			}
		}
	}

	return measurement, nil
}

func (s *MeteoService) parseFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func (s *MeteoService) parseInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case float32:
		return int(v), true
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i, true
		}
	}
	return 0, false
}

func (s *MeteoService) parseEnum(value interface{}) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case map[string]interface{}:
		if val, ok := v["value"].(string); ok {
			return val, true
		}
	}
	return "", false
}

func (s *MeteoService) parsePeriod(period string) (time.Duration, error) {
	if period == "" {
		return 0, fmt.Errorf("period cannot be empty")
	}

	// Handle common formats like "24h", "7d", "15m", etc.
	duration, err := time.ParseDuration(period)
	if err != nil {
		// Try parsing as number + unit (e.g., "24h", "7d")
		if len(period) < 2 {
			return 0, fmt.Errorf("invalid period format: %s", period)
		}

		unit := period[len(period)-1:]
		valueStr := period[:len(period)-1]
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return 0, fmt.Errorf("invalid period format: %s", period)
		}

		switch unit {
		case "m":
			duration = time.Duration(value) * time.Minute
		case "h":
			duration = time.Duration(value) * time.Hour
		case "d":
			duration = time.Duration(value) * 24 * time.Hour
		case "w":
			duration = time.Duration(value) * 7 * 24 * time.Hour
		default:
			return 0, fmt.Errorf("unsupported time unit: %s", unit)
		}
	}

	if duration <= 0 {
		return 0, fmt.Errorf("period must be positive")
	}

	return duration, nil
}