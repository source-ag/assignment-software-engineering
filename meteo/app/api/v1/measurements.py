"""
Measurement data API endpoints.
"""
import logging
from typing import Dict, Any

from fastapi import APIRouter, HTTPException, Depends
from sqlalchemy.orm import Session

from meteo.app.schemas.sensor_data import SensorData
from meteo.app.core.database import get_db
from meteo.app.services import MeasurementService

logger = logging.getLogger(__name__)

router = APIRouter(
    prefix="/api/v1",
    tags=["measurements"]
)

# Let fastapi handle the validation of SensorData and raise 422 if invalid
@router.post("/measurement")
async def create_measurement(sensor_data: SensorData, db: Session = Depends(get_db)) -> Dict[str, Any]:
    """Creates a new measurement for data from the hortimax synopta controller."""
    # TODO: Validate access rights (stretch goal)
    
    try:
        # Create service instance with database session
        service = MeasurementService(db)
        measurement = service.create_measurement(sensor_data)
        
        return {
            "message": "Measurement created and saved successfully", 
            "id": measurement.id,
            "name": sensor_data.name,
            "timestamp": sensor_data.ts.isoformat(),
            "rows_count": len(sensor_data.rows)
        }
        
    except Exception as e:
        logger.error(f"Error saving measurement: {e}")
        # TODO: Beautify error for frontend display
        raise HTTPException(status_code=500, detail=f"Failed to save measurement: {str(e)}")

@router.get("/latest")
async def get_latest_measurements(db: Session = Depends(get_db)) -> Dict[str, Any]:
    """Get the latest weather measurements."""
    try:
        # Create service instance with database session
        service = MeasurementService(db)
        measurement = service.get_latest_measurements()

        # Convert to dictionaries using the service
        result = service.measurement_to_dict(measurement) if measurement else None

        return  result
            
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to retrieve measurements: {str(e)}")

@router.delete("/measurements")
async def delete_all_measurements(db: Session = Depends(get_db)) -> Dict[str, Any]:
    """Delete all measurements from the database. (Only for testing. In a proper application this would not be exposed.)"""
    try:
        # Create service instance with database session
        service = MeasurementService(db)
        deleted_count = service.delete_all_measurements()
        
        return {
            "message": "All measurements deleted successfully",
            "deleted_count": deleted_count
        }
        
    except Exception as e:
        logger.error(f"Error deleting measurements: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to delete measurements: {str(e)}")

# @router.get("/average")
# async def get_average_measurements() -> Dict[str, Any]:
#     """Get average measurements over a period."""
#     pass



# @router.post("/measurements/bulk")
# async def create_measurements_bulk() -> Dict[str, Any]:
#     """Create multiple measurements in bulk."""
#     pass
