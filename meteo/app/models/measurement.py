"""
SQLAlchemy models for weather measurements.
"""
from sqlalchemy import Column, Integer, Float, String, DateTime, Enum

from meteo.app.models.enums.hortimax_synopta_enum import HortimaxSynoptaEnum
from meteo.app.schemas.sensor_data import HortimaxEnumValue

from meteo.app.core.database import Base


class Measurement(Base):
    """
    Weather measurement model.
    
    Stores processed weather data from sensor readings.
    """
    __tablename__: str = "measurements"

    # Primary key and timestamp
    id = Column(Integer, primary_key=True, index=True)
    timestamp = Column(DateTime, nullable=False, index=True)

    # Range
    range_row1 = Column(Integer, nullable=True)
    range_row2 = Column(Integer, nullable=True)
    range_col1 = Column(Integer, nullable=True)
    range_col2 = Column(Integer, nullable=True)

    # Temperature measurements
    external_temperature_c = Column(Float, nullable=True)
    dew_point_temperature_c = Column(Float, nullable=True)
    
    # Wind measurements
    wind_speed_m_s = Column(Float, nullable=True)
    wind_speed_unmuted_m_s = Column(Float, nullable=True)
    wind_direction_degrees = Column(Integer, nullable=True)
    wind_direction_compass = Column(Enum(HortimaxSynoptaEnum), nullable=True)
    
    # Humidity and pressure
    relative_humidity_perc = Column(Float, nullable=True)
    atmospheric_pressure_hpa = Column(Float, nullable=True)
    abs_humidity_g_kg = Column(Float, nullable=True)
    moisture_deficit_g_kg = Column(Float, nullable=True)
    moisture_deficit_g_m3 = Column(Float, nullable=True)
    
    # Radiation measurements
    radiation_intensity_w_m2 = Column(Float, nullable=True)
    radiation_intensity_unmuted_w_m2 = Column(Float, nullable=True)
    standard_radiation_intensity_w_m2 = Column(Float, nullable=True)
    radiation_sum_j_cm2 = Column(Float, nullable=True)
    radiation_from_plant_w_m2 = Column(Float, nullable=True)
    
    # Enthalpy measurements
    enthalpy_kj_kg = Column(Float, nullable=True)
    enthalpy_kj_m3 = Column(Float, nullable=True)
    
    # Weather conditions
    precipitation = Column(Float, nullable=True)
    
    # Station status
    status_meteo_station = Column(Enum(HortimaxSynoptaEnum), nullable=True)
    status_meteo_station_communication = Column(Enum(HortimaxSynoptaEnum), nullable=True)

    # Metadata
    name = Column(String(100), nullable=False)
    pt = Column(Integer, nullable=True)

    @classmethod
    def from_sensor_data(cls, sensor_data) -> "Measurement":
        data_dict = {
            "name": sensor_data.name,
            "timestamp": sensor_data.ts,
            "pt": sensor_data.pt,
            "range_row1": sensor_data.range.row1,
            "range_row2": sensor_data.range.row2,
            "range_col1": sensor_data.range.col1,
            "range_col2": sensor_data.range.col2
        }

        # Skip header row from sensor data
        for row in sensor_data.rows[1:]:
            variable_name = row[0]
            value = row[1]
                
            if isinstance(value, (str, int, float)):
                data_dict[variable_name] = value
            
            # Handle Hortimax Synopta enum values
            elif isinstance(value, HortimaxEnumValue) and value.type == "hortimax.synopta.enum":
                data_dict[variable_name] = HortimaxSynoptaEnum(value.key)

        return cls(**data_dict)