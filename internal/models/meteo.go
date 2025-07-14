package models

import (
	"encoding/json"
	"time"
)

// RawMeteoData represents the raw data format from the climate computer
type RawMeteoData struct {
	Name  string `json:"name"`
	Range struct {
		Col1 int `json:"col1"`
		Row1 int `json:"row1"`
		Col2 int `json:"col2"`
		Row2 int `json:"row2"`
	} `json:"range"`
	Rows [][]interface{} `json:"rows"`
	TS   string          `json:"ts"`
	PT   int             `json:"pt"`
}

// EnumValue represents enum values in the raw data
type EnumValue struct {
	Type  string `json:"type"`
	Key   int    `json:"key"`
	Value string `json:"value"`
}

// MeteoMeasurement represents a normalized measurement
type MeteoMeasurement struct {
	ID                               string    `json:"id" db:"id"`
	Timestamp                        time.Time `json:"timestamp" db:"timestamp"`
	ExternalTemperatureC             *float64  `json:"external_temperature_c" db:"external_temperature_c"`
	WindSpeedUnmutedMS              *float64  `json:"wind_speed_unmuted_m_s" db:"wind_speed_unmuted_m_s"`
	WindSpeedMS                     *float64  `json:"wind_speed_m_s" db:"wind_speed_m_s"`
	WindDirectionDegrees            *int      `json:"wind_direction_degrees" db:"wind_direction_degrees"`
	WindDirectionCompass            *string   `json:"wind_direction_compass" db:"wind_direction_compass"`
	RadiationIntensityUnmutedWM2    *float64  `json:"radiation_intensity_unmuted_w_m2" db:"radiation_intensity_unmuted_w_m2"`
	RadiationIntensityWM2           *float64  `json:"radiation_intensity_w_m2" db:"radiation_intensity_w_m2"`
	StandardRadiationIntensityWM2   *float64  `json:"standard_radiation_intensity_w_m2" db:"standard_radiation_intensity_w_m2"`
	RadiationSumJCM2                *float64  `json:"radiation_sum_j_cm2" db:"radiation_sum_j_cm2"`
	RadiationFromPlantWM2           *float64  `json:"radiation_from_plant_w_m2" db:"radiation_from_plant_w_m2"`
	Precipitation                   *float64  `json:"precipitation" db:"precipitation"`
	RelativeHumidityPerc            *float64  `json:"relative_humidity_perc" db:"relative_humidity_perc"`
	MoistureDeficitGKG              *float64  `json:"moisture_deficit_g_kg" db:"moisture_deficit_g_kg"`
	MoistureDeficitGM3              *float64  `json:"moisture_deficit_g_m3" db:"moisture_deficit_g_m3"`
	DewPointTemperatureC            *float64  `json:"dew_point_temperature_c" db:"dew_point_temperature_c"`
	AbsHumidityGKG                  *float64  `json:"abs_humidity_g_kg" db:"abs_humidity_g_kg"`
	EnthalpyKJKG                    *float64  `json:"enthalpy_kj_kg" db:"enthalpy_kj_kg"`
	EnthalpyKJM3                    *float64  `json:"enthalpy_kj_m3" db:"enthalpy_kj_m3"`
	AtmosphericPressureHPA          *float64  `json:"atmospheric_pressure_hpa" db:"atmospheric_pressure_hpa"`
	StatusMeteoStation              *string   `json:"status_meteo_station" db:"status_meteo_station"`
	StatusMeteoStationCommunication *string   `json:"status_meteo_station_communication" db:"status_meteo_station_communication"`
	RawData                         json.RawMessage `json:"raw_data" db:"raw_data"`
}

// MeteoAggregation represents aggregated meteo data
type MeteoAggregation struct {
	Timestamp                        time.Time `json:"timestamp"`
	ExternalTemperatureC             *float64  `json:"external_temperature_c"`
	WindSpeedUnmutedMS              *float64  `json:"wind_speed_unmuted_m_s"`
	WindSpeedMS                     *float64  `json:"wind_speed_m_s"`
	WindDirectionDegrees            *float64  `json:"wind_direction_degrees"`
	RadiationIntensityUnmutedWM2    *float64  `json:"radiation_intensity_unmuted_w_m2"`
	RadiationIntensityWM2           *float64  `json:"radiation_intensity_w_m2"`
	StandardRadiationIntensityWM2   *float64  `json:"standard_radiation_intensity_w_m2"`
	RadiationSumJCM2                *float64  `json:"radiation_sum_j_cm2"`
	RadiationFromPlantWM2           *float64  `json:"radiation_from_plant_w_m2"`
	Precipitation                   *float64  `json:"precipitation"`
	RelativeHumidityPerc            *float64  `json:"relative_humidity_perc"`
	MoistureDeficitGKG              *float64  `json:"moisture_deficit_g_kg"`
	MoistureDeficitGM3              *float64  `json:"moisture_deficit_g_m3"`
	DewPointTemperatureC            *float64  `json:"dew_point_temperature_c"`
	AbsHumidityGKG                  *float64  `json:"abs_humidity_g_kg"`
	EnthalpyKJKG                    *float64  `json:"enthalpy_kj_kg"`
	EnthalpyKJM3                    *float64  `json:"enthalpy_kj_m3"`
	AtmosphericPressureHPA          *float64  `json:"atmospheric_pressure_hpa"`
	Count                           int       `json:"count"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// AggregationRequest represents query parameters for aggregation
type AggregationRequest struct {
	Period   string `form:"period" binding:"required"`   // e.g., "24h", "7d", "30d"
	Interval string `form:"interval" binding:"required"` // e.g., "15m", "1h", "1d"
}

// HistoryRequest represents query parameters for history
type HistoryRequest struct {
	Period string `form:"period" binding:"required"` // e.g., "24h", "7d", "30d"
	Limit  int    `form:"limit"`                     // optional limit
}