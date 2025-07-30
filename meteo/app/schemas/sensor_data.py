from pydantic import BaseModel, ConfigDict
from typing import List, Union, Optional
from datetime import datetime

class HortimaxEnumValue(BaseModel):
    type: str
    key: int
    value: str

class Range(BaseModel):
    col1: int
    row1: int
    col2: int
    row2: int

class SensorData(BaseModel):
    model_config = ConfigDict(
        json_encoders={
            datetime: lambda v: v.isoformat()
        }
    )
    
    name: str
    range: Range
    rows: List[List[Union[str, float, int, HortimaxEnumValue]]] 
    ts: datetime
    pt: int