"""
Measurement service for handling business logic.
"""
import logging
from typing import Dict, Any
from sqlalchemy.orm import Session

from meteo.app.models.measurement import Measurement
from meteo.app.schemas.sensor_data import SensorData

logger = logging.getLogger(__name__)


class MeasurementService:
    """Service for storing and loading measurements from and to the database."""
    
    def __init__(self, db: Session):
        """
        Initialize the measurement service.
        
        Args:
            db: Database session
        """
        self.db = db
    
    def create_measurement(self, sensor_data: SensorData) -> Measurement:
        """
        Create and save a measurement from sensor data.
        
        Args:
            sensor_data: Validated sensor data
            
        Returns:
            Measurement: The saved measurement object
            
        Raises:
            Exception: If saving fails
        """
        try:
            processed_data = Measurement.from_sensor_data(sensor_data)
            
            self.db.add(processed_data)
            self.db.commit()
            self.db.refresh(processed_data)
            
            logger.info(f"Successfully saved measurement with ID: {processed_data.id}")
            return processed_data
            
        except Exception as e:
            self.db.rollback()
            logger.error(f"Error saving measurement: {e}")
            raise e
    
    def get_latest_measurements(self) -> Measurement:
        """
        Get the latest measurements from the database.
            
        Returns:
            Measurement: The latest measurement
        """
        latest_measurements = self.get_measurements(limit=1)
        return latest_measurements[0] if latest_measurements else None

    def get_measurements(self, limit: int) -> list[Measurement]:
        """
        Get the latest measurements from the database.
        
        Args:
            limit: Maximum number of measurements to return
            
        Returns:
            list[Measurement]: List of latest measurements
        """
        return self.db.query(Measurement).order_by(Measurement.timestamp.desc()).limit(limit).all()


    def measurement_to_dict(self, measurement: Measurement) -> Dict[str, Any]:
        """
        Convert a measurement object to a dictionary for API response.
        
        Args:
            measurement: Measurement object
            
        Returns:
            Dict: Measurement data as dictionary
        """
        return {
            "id": measurement.id,
            "name": measurement.name,
            "timestamp": measurement.timestamp.isoformat(),
            "external_temperature_c": measurement.external_temperature_c,
            "wind_speed_m_s": measurement.wind_speed_m_s,
            "wind_direction_compass": measurement.wind_direction_compass.name if measurement.wind_direction_compass else None,
            "relative_humidity_perc": measurement.relative_humidity_perc,
            "atmospheric_pressure_hpa": measurement.atmospheric_pressure_hpa
        }

    def delete_all_measurements(self) -> int:
        """
        Delete all measurements from the database.
        
        Returns:
            int: Number of measurements deleted
            
        Raises:
            Exception: If deletion fails
        """
        try:
            count = self.db.query(Measurement).count()
            
            self.db.query(Measurement).delete()
            self.db.commit()

            logger.info(f"Successfully deleted all measurements")
            return count

        except Exception as e:
            self.db.rollback()
            logger.error(f"Failed to delete measurements: {e}")
            raise
