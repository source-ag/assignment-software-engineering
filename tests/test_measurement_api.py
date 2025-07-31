"""
Integration tests for measurement API endpoints.
Tests the actual HTTP endpoints using the live database.
"""
import json
from pathlib import Path
from fastapi.testclient import TestClient

from meteo.app.main import app

# Create test client
client = TestClient(app)


class TestMeasurementAPIIntegration:
    """Integration tests for measurement API endpoints."""
    
    def setup_method(self):
        """Setup before each test method."""
        # Clear database before each test
        response = client.delete("/api/v1/measurements")
        print(f"Cleared database: {response.json()}")
    
    def teardown_method(self):
        """Cleanup after each test method."""
        # Clear database after each test
        response = client.delete("/api/v1/measurements")
        print(f"Cleaned up database: {response.json()}")
    
    def load_sample_data(self) -> dict:
        """Load sample data from the data folder."""
        sample_file = Path(__file__).parent / "data" / "test_measurement.json"
        
        with open(sample_file, 'r') as f:
            return json.load(f)
    
    def test_post_measurement_endpoint(self):
        # Arrange
        sample_data = self.load_sample_data()
        
        # Act
        response = client.post("/api/v1/measurement", json=sample_data)
        
        # Assert
        assert response.status_code == 200
        response_data = response.json()
        
        assert response_data["message"] == "Measurement created and saved successfully"
        assert "id" in response_data
        assert response_data["name"] == "_ws_source_meteo"
        assert "2021-05-01T02:02:50" in response_data["timestamp"]
        assert response_data["rows_count"] == len(sample_data["rows"])
    
    def test_get_latest_measurement_endpoint(self):
        # Arrange
        sample_data = self.load_sample_data()
        post_response = client.post("/api/v1/measurement", json=sample_data)
        assert post_response.status_code == 200
        created_id = post_response.json()["id"]
        
        # Act
        response = client.get("/api/v1/latest")
        
        # Assert
        assert response.status_code == 200
        response_data = response.json()
        
        assert response_data["id"] == created_id
        assert response_data["name"] == "_ws_source_meteo"
        assert response_data["external_temperature_c"] == 7.73957
        assert response_data["wind_speed_m_s"] == 1.81514
        assert response_data["relative_humidity_perc"] == 91
        assert response_data["atmospheric_pressure_hpa"] == 1013

    def test_post_measurement_with_invalid_data(self):
        # Arrange
        invalid_data = {
            "name": "test",
            "invalid_field": "value"
        }
        
        # Act
        response = client.post("/api/v1/measurement", json=invalid_data)
        
        # Assert
        assert response.status_code == 422
        response_data = response.json()
        assert "detail" in response_data
           
    def test_invalid_format(self):
        # Act
        response = client.post(
            "/api/v1/measurement", 
            data="invalid json",
            headers={"Content-Type": "application/json"}
        )

        # Assert
        assert response.status_code == 422

    def test_invalid_endpoint(self):
        # Act
        response = client.get("/api/v1/nonexistent")
  
        # Assert
        assert response.status_code == 404
