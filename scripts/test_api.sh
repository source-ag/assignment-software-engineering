#!/bin/bash

# Test script for Meteo API
# Usage: ./scripts/test_api.sh

set -e

API_URL="http://localhost:8080"
READ_KEY="read-key-123"
WRITE_KEY="write-key-456"

echo "Testing Meteo API..."
echo "===================="

# Test health check
echo "1. Testing health check..."
curl -s -X GET "${API_URL}/api/v1/health" | jq .
echo -e "\n"

# Test authentication failure
echo "2. Testing authentication failure..."
curl -s -X GET "${API_URL}/api/v1/meteo/current" | jq .
echo -e "\n"

# Test with read key
echo "3. Testing current measurement with read key..."
curl -s -H "Authorization: Bearer ${READ_KEY}" -X GET "${API_URL}/api/v1/meteo/current" | jq .
echo -e "\n"

# Test posting measurement
echo "4. Testing post measurement..."
curl -s -H "Authorization: Bearer ${WRITE_KEY}" \
     -H "Content-Type: application/json" \
     -X POST "${API_URL}/api/v1/meteo/measurements" \
     -d '{
       "name": "_ws_source_meteo",
       "range": {"col1": 1, "row1": 1, "col2": 2, "row2": 22},
       "rows": [
         ["Variable", "Value"],
         ["external_temperature_c", 9.1387],
         ["wind_speed_unmuted_m_s", 3.25534],
         ["wind_speed_m_s", 3.06032],
         ["wind_direction_degrees", 232],
         ["wind_direction_compass", {"type": "hortimax.synopta.enum", "key": 8782, "value": "ZW"}],
         ["radiation_intensity_unmuted_w_m2", 174.904],
         ["radiation_intensity_w_m2", 235.352],
         ["standard_radiation_intensity_w_m2", 174.904],
         ["radiation_sum_j_cm2", 437.738],
         ["radiation_from_plant_w_m2", 35],
         ["precipitation", 0],
         ["relative_humidity_perc", 73],
         ["moisture_deficit_g_kg", 1.95321],
         ["moisture_deficit_g_m3", 2.43753],
         ["dew_point_temperature_c", 4.5713],
         ["abs_humidity_g_kg", 5.23157],
         ["enthalpy_kj_kg", 22.3522],
         ["enthalpy_kj_m3", 27.8948],
         ["atmospheric_pressure_hpa", 1013],
         ["status_meteo_station", {"type": "hortimax.synopta.enum", "key": 8789, "value": "Actief"}],
         ["status_meteo_station_communication", {"type": "hortimax.synopta.enum", "key": 8796, "value": "Online"}]
       ],
       "ts": "2021-05-01T12:07:50+02:00",
       "pt": 0
     }' | jq .
echo -e "\n"

# Test getting current measurement
echo "5. Testing get current measurement..."
curl -s -H "Authorization: Bearer ${READ_KEY}" -X GET "${API_URL}/api/v1/meteo/current" | jq .
echo -e "\n"

# Test history
echo "6. Testing history (24h)..."
curl -s -H "Authorization: Bearer ${READ_KEY}" -X GET "${API_URL}/api/v1/meteo/history?period=24h&limit=5" | jq .
echo -e "\n"

# Test aggregated data
echo "7. Testing aggregated data..."
curl -s -H "Authorization: Bearer ${READ_KEY}" -X GET "${API_URL}/api/v1/meteo/aggregated?period=24h&interval=1h" | jq .
echo -e "\n"

# Test average
echo "8. Testing average..."
curl -s -H "Authorization: Bearer ${READ_KEY}" -X GET "${API_URL}/api/v1/meteo/average?period=24h" | jq .
echo -e "\n"

# Test forbidden access (read key for write endpoint)
echo "9. Testing forbidden access..."
curl -s -H "Authorization: Bearer ${READ_KEY}" \
     -H "Content-Type: application/json" \
     -X POST "${API_URL}/api/v1/meteo/measurements" \
     -d '{"name": "test", "ts": "2021-05-01T12:07:50+02:00", "pt": 0, "rows": []}' | jq .
echo -e "\n"

echo "API testing completed!"