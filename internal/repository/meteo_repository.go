package repository

import (
	"database/sql"
	"fmt"
	"time"

	"meteo-api/internal/models"

	"github.com/google/uuid"
)

type MeteoRepository struct {
	db *sql.DB
}

func NewMeteoRepository(db *sql.DB) *MeteoRepository {
	return &MeteoRepository{db: db}
}

func (r *MeteoRepository) CreateMeasurement(measurement *models.MeteoMeasurement) error {
	query := `
		INSERT INTO meteo_measurements (
			id, timestamp, external_temperature_c, wind_speed_unmuted_m_s, wind_speed_m_s,
			wind_direction_degrees, wind_direction_compass, radiation_intensity_unmuted_w_m2,
			radiation_intensity_w_m2, standard_radiation_intensity_w_m2, radiation_sum_j_cm2,
			radiation_from_plant_w_m2, precipitation, relative_humidity_perc, moisture_deficit_g_kg,
			moisture_deficit_g_m3, dew_point_temperature_c, abs_humidity_g_kg, enthalpy_kj_kg,
			enthalpy_kj_m3, atmospheric_pressure_hpa, status_meteo_station, 
			status_meteo_station_communication, raw_data
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
		)`

	if measurement.ID == "" {
		measurement.ID = uuid.New().String()
	}

	_, err := r.db.Exec(query,
		measurement.ID, measurement.Timestamp, measurement.ExternalTemperatureC,
		measurement.WindSpeedUnmutedMS, measurement.WindSpeedMS, measurement.WindDirectionDegrees,
		measurement.WindDirectionCompass, measurement.RadiationIntensityUnmutedWM2,
		measurement.RadiationIntensityWM2, measurement.StandardRadiationIntensityWM2,
		measurement.RadiationSumJCM2, measurement.RadiationFromPlantWM2, measurement.Precipitation,
		measurement.RelativeHumidityPerc, measurement.MoistureDeficitGKG, measurement.MoistureDeficitGM3,
		measurement.DewPointTemperatureC, measurement.AbsHumidityGKG, measurement.EnthalpyKJKG,
		measurement.EnthalpyKJM3, measurement.AtmosphericPressureHPA, measurement.StatusMeteoStation,
		measurement.StatusMeteoStationCommunication, measurement.RawData,
	)

	return err
}

func (r *MeteoRepository) CreateMeasurementsBulk(measurements []models.MeteoMeasurement) error {
	if len(measurements) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO meteo_measurements (
			id, timestamp, external_temperature_c, wind_speed_unmuted_m_s, wind_speed_m_s,
			wind_direction_degrees, wind_direction_compass, radiation_intensity_unmuted_w_m2,
			radiation_intensity_w_m2, standard_radiation_intensity_w_m2, radiation_sum_j_cm2,
			radiation_from_plant_w_m2, precipitation, relative_humidity_perc, moisture_deficit_g_kg,
			moisture_deficit_g_m3, dew_point_temperature_c, abs_humidity_g_kg, enthalpy_kj_kg,
			enthalpy_kj_m3, atmospheric_pressure_hpa, status_meteo_station, 
			status_meteo_station_communication, raw_data
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
		)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, measurement := range measurements {
		if measurement.ID == "" {
			measurement.ID = uuid.New().String()
		}

		_, err := stmt.Exec(
			measurement.ID, measurement.Timestamp, measurement.ExternalTemperatureC,
			measurement.WindSpeedUnmutedMS, measurement.WindSpeedMS, measurement.WindDirectionDegrees,
			measurement.WindDirectionCompass, measurement.RadiationIntensityUnmutedWM2,
			measurement.RadiationIntensityWM2, measurement.StandardRadiationIntensityWM2,
			measurement.RadiationSumJCM2, measurement.RadiationFromPlantWM2, measurement.Precipitation,
			measurement.RelativeHumidityPerc, measurement.MoistureDeficitGKG, measurement.MoistureDeficitGM3,
			measurement.DewPointTemperatureC, measurement.AbsHumidityGKG, measurement.EnthalpyKJKG,
			measurement.EnthalpyKJM3, measurement.AtmosphericPressureHPA, measurement.StatusMeteoStation,
			measurement.StatusMeteoStationCommunication, measurement.RawData,
		)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	return tx.Commit()
}

func (r *MeteoRepository) GetLatestMeasurement() (*models.MeteoMeasurement, error) {
	query := `
		SELECT id, timestamp, external_temperature_c, wind_speed_unmuted_m_s, wind_speed_m_s,
			wind_direction_degrees, wind_direction_compass, radiation_intensity_unmuted_w_m2,
			radiation_intensity_w_m2, standard_radiation_intensity_w_m2, radiation_sum_j_cm2,
			radiation_from_plant_w_m2, precipitation, relative_humidity_perc, moisture_deficit_g_kg,
			moisture_deficit_g_m3, dew_point_temperature_c, abs_humidity_g_kg, enthalpy_kj_kg,
			enthalpy_kj_m3, atmospheric_pressure_hpa, status_meteo_station, 
			status_meteo_station_communication, raw_data
		FROM meteo_measurements
		ORDER BY timestamp DESC
		LIMIT 1`

	var measurement models.MeteoMeasurement
	err := r.db.QueryRow(query).Scan(
		&measurement.ID, &measurement.Timestamp, &measurement.ExternalTemperatureC,
		&measurement.WindSpeedUnmutedMS, &measurement.WindSpeedMS, &measurement.WindDirectionDegrees,
		&measurement.WindDirectionCompass, &measurement.RadiationIntensityUnmutedWM2,
		&measurement.RadiationIntensityWM2, &measurement.StandardRadiationIntensityWM2,
		&measurement.RadiationSumJCM2, &measurement.RadiationFromPlantWM2, &measurement.Precipitation,
		&measurement.RelativeHumidityPerc, &measurement.MoistureDeficitGKG, &measurement.MoistureDeficitGM3,
		&measurement.DewPointTemperatureC, &measurement.AbsHumidityGKG, &measurement.EnthalpyKJKG,
		&measurement.EnthalpyKJM3, &measurement.AtmosphericPressureHPA, &measurement.StatusMeteoStation,
		&measurement.StatusMeteoStationCommunication, &measurement.RawData,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &measurement, nil
}

func (r *MeteoRepository) GetMeasurementsByPeriod(period time.Duration, limit int) ([]models.MeteoMeasurement, error) {
	query := `
		SELECT id, timestamp, external_temperature_c, wind_speed_unmuted_m_s, wind_speed_m_s,
			wind_direction_degrees, wind_direction_compass, radiation_intensity_unmuted_w_m2,
			radiation_intensity_w_m2, standard_radiation_intensity_w_m2, radiation_sum_j_cm2,
			radiation_from_plant_w_m2, precipitation, relative_humidity_perc, moisture_deficit_g_kg,
			moisture_deficit_g_m3, dew_point_temperature_c, abs_humidity_g_kg, enthalpy_kj_kg,
			enthalpy_kj_m3, atmospheric_pressure_hpa, status_meteo_station, 
			status_meteo_station_communication, raw_data
		FROM meteo_measurements
		WHERE timestamp >= NOW() - INTERVAL '%d seconds'
		ORDER BY timestamp DESC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.Query(fmt.Sprintf(query, int(period.Seconds())))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var measurements []models.MeteoMeasurement
	for rows.Next() {
		var measurement models.MeteoMeasurement
		err := rows.Scan(
			&measurement.ID, &measurement.Timestamp, &measurement.ExternalTemperatureC,
			&measurement.WindSpeedUnmutedMS, &measurement.WindSpeedMS, &measurement.WindDirectionDegrees,
			&measurement.WindDirectionCompass, &measurement.RadiationIntensityUnmutedWM2,
			&measurement.RadiationIntensityWM2, &measurement.StandardRadiationIntensityWM2,
			&measurement.RadiationSumJCM2, &measurement.RadiationFromPlantWM2, &measurement.Precipitation,
			&measurement.RelativeHumidityPerc, &measurement.MoistureDeficitGKG, &measurement.MoistureDeficitGM3,
			&measurement.DewPointTemperatureC, &measurement.AbsHumidityGKG, &measurement.EnthalpyKJKG,
			&measurement.EnthalpyKJM3, &measurement.AtmosphericPressureHPA, &measurement.StatusMeteoStation,
			&measurement.StatusMeteoStationCommunication, &measurement.RawData,
		)
		if err != nil {
			return nil, err
		}
		measurements = append(measurements, measurement)
	}

	return measurements, nil
}

func (r *MeteoRepository) GetAggregatedData(period time.Duration, interval time.Duration) ([]models.MeteoAggregation, error) {
	query := `
		SELECT 
			date_trunc('%s', timestamp) as time_bucket,
			AVG(external_temperature_c) as avg_external_temperature_c,
			AVG(wind_speed_unmuted_m_s) as avg_wind_speed_unmuted_m_s,
			AVG(wind_speed_m_s) as avg_wind_speed_m_s,
			AVG(wind_direction_degrees) as avg_wind_direction_degrees,
			AVG(radiation_intensity_unmuted_w_m2) as avg_radiation_intensity_unmuted_w_m2,
			AVG(radiation_intensity_w_m2) as avg_radiation_intensity_w_m2,
			AVG(standard_radiation_intensity_w_m2) as avg_standard_radiation_intensity_w_m2,
			AVG(radiation_sum_j_cm2) as avg_radiation_sum_j_cm2,
			AVG(radiation_from_plant_w_m2) as avg_radiation_from_plant_w_m2,
			AVG(precipitation) as avg_precipitation,
			AVG(relative_humidity_perc) as avg_relative_humidity_perc,
			AVG(moisture_deficit_g_kg) as avg_moisture_deficit_g_kg,
			AVG(moisture_deficit_g_m3) as avg_moisture_deficit_g_m3,
			AVG(dew_point_temperature_c) as avg_dew_point_temperature_c,
			AVG(abs_humidity_g_kg) as avg_abs_humidity_g_kg,
			AVG(enthalpy_kj_kg) as avg_enthalpy_kj_kg,
			AVG(enthalpy_kj_m3) as avg_enthalpy_kj_m3,
			AVG(atmospheric_pressure_hpa) as avg_atmospheric_pressure_hpa,
			COUNT(*) as count
		FROM meteo_measurements
		WHERE timestamp >= NOW() - INTERVAL '%d seconds'
		GROUP BY time_bucket
		ORDER BY time_bucket DESC`

	intervalStr := r.intervalToPostgresInterval(interval)
	rows, err := r.db.Query(fmt.Sprintf(query, intervalStr, int(period.Seconds())))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aggregations []models.MeteoAggregation
	for rows.Next() {
		var agg models.MeteoAggregation
		err := rows.Scan(
			&agg.Timestamp, &agg.ExternalTemperatureC, &agg.WindSpeedUnmutedMS,
			&agg.WindSpeedMS, &agg.WindDirectionDegrees, &agg.RadiationIntensityUnmutedWM2,
			&agg.RadiationIntensityWM2, &agg.StandardRadiationIntensityWM2, &agg.RadiationSumJCM2,
			&agg.RadiationFromPlantWM2, &agg.Precipitation, &agg.RelativeHumidityPerc,
			&agg.MoistureDeficitGKG, &agg.MoistureDeficitGM3, &agg.DewPointTemperatureC,
			&agg.AbsHumidityGKG, &agg.EnthalpyKJKG, &agg.EnthalpyKJM3, &agg.AtmosphericPressureHPA,
			&agg.Count,
		)
		if err != nil {
			return nil, err
		}
		aggregations = append(aggregations, agg)
	}

	return aggregations, nil
}

func (r *MeteoRepository) GetOverallAverage(period time.Duration) (*models.MeteoAggregation, error) {
	query := `
		SELECT 
			NOW() as timestamp,
			AVG(external_temperature_c) as avg_external_temperature_c,
			AVG(wind_speed_unmuted_m_s) as avg_wind_speed_unmuted_m_s,
			AVG(wind_speed_m_s) as avg_wind_speed_m_s,
			AVG(wind_direction_degrees) as avg_wind_direction_degrees,
			AVG(radiation_intensity_unmuted_w_m2) as avg_radiation_intensity_unmuted_w_m2,
			AVG(radiation_intensity_w_m2) as avg_radiation_intensity_w_m2,
			AVG(standard_radiation_intensity_w_m2) as avg_standard_radiation_intensity_w_m2,
			AVG(radiation_sum_j_cm2) as avg_radiation_sum_j_cm2,
			AVG(radiation_from_plant_w_m2) as avg_radiation_from_plant_w_m2,
			AVG(precipitation) as avg_precipitation,
			AVG(relative_humidity_perc) as avg_relative_humidity_perc,
			AVG(moisture_deficit_g_kg) as avg_moisture_deficit_g_kg,
			AVG(moisture_deficit_g_m3) as avg_moisture_deficit_g_m3,
			AVG(dew_point_temperature_c) as avg_dew_point_temperature_c,
			AVG(abs_humidity_g_kg) as avg_abs_humidity_g_kg,
			AVG(enthalpy_kj_kg) as avg_enthalpy_kj_kg,
			AVG(enthalpy_kj_m3) as avg_enthalpy_kj_m3,
			AVG(atmospheric_pressure_hpa) as avg_atmospheric_pressure_hpa,
			COUNT(*) as count
		FROM meteo_measurements
		WHERE timestamp >= NOW() - INTERVAL '%d seconds'`

	var agg models.MeteoAggregation
	err := r.db.QueryRow(fmt.Sprintf(query, int(period.Seconds()))).Scan(
		&agg.Timestamp, &agg.ExternalTemperatureC, &agg.WindSpeedUnmutedMS,
		&agg.WindSpeedMS, &agg.WindDirectionDegrees, &agg.RadiationIntensityUnmutedWM2,
		&agg.RadiationIntensityWM2, &agg.StandardRadiationIntensityWM2, &agg.RadiationSumJCM2,
		&agg.RadiationFromPlantWM2, &agg.Precipitation, &agg.RelativeHumidityPerc,
		&agg.MoistureDeficitGKG, &agg.MoistureDeficitGM3, &agg.DewPointTemperatureC,
		&agg.AbsHumidityGKG, &agg.EnthalpyKJKG, &agg.EnthalpyKJM3, &agg.AtmosphericPressureHPA,
		&agg.Count,
	)

	if err != nil {
		return nil, err
	}

	return &agg, nil
}

func (r *MeteoRepository) intervalToPostgresInterval(interval time.Duration) string {
	switch {
	case interval < time.Hour:
		return "minute"
	case interval < 24*time.Hour:
		return "hour"
	default:
		return "day"
	}
}