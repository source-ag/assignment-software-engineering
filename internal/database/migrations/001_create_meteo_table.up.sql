CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS meteo_measurements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    external_temperature_c DOUBLE PRECISION,
    wind_speed_unmuted_m_s DOUBLE PRECISION,
    wind_speed_m_s DOUBLE PRECISION,
    wind_direction_degrees INTEGER,
    wind_direction_compass VARCHAR(10),
    radiation_intensity_unmuted_w_m2 DOUBLE PRECISION,
    radiation_intensity_w_m2 DOUBLE PRECISION,
    standard_radiation_intensity_w_m2 DOUBLE PRECISION,
    radiation_sum_j_cm2 DOUBLE PRECISION,
    radiation_from_plant_w_m2 DOUBLE PRECISION,
    precipitation DOUBLE PRECISION,
    relative_humidity_perc DOUBLE PRECISION,
    moisture_deficit_g_kg DOUBLE PRECISION,
    moisture_deficit_g_m3 DOUBLE PRECISION,
    dew_point_temperature_c DOUBLE PRECISION,
    abs_humidity_g_kg DOUBLE PRECISION,
    enthalpy_kj_kg DOUBLE PRECISION,
    enthalpy_kj_m3 DOUBLE PRECISION,
    atmospheric_pressure_hpa DOUBLE PRECISION,
    status_meteo_station VARCHAR(50),
    status_meteo_station_communication VARCHAR(50),
    raw_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_meteo_timestamp ON meteo_measurements(timestamp);
CREATE INDEX IF NOT EXISTS idx_meteo_created_at ON meteo_measurements(created_at);
CREATE INDEX IF NOT EXISTS idx_meteo_raw_data ON meteo_measurements USING GIN (raw_data);