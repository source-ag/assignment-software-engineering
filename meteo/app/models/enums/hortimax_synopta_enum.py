from enum import Enum
from typing import Optional

class HortimaxSynoptaEmun(Enum):
    """Hortimax Synopta enum values"""

    # Wind direction compass values
    W = 8784  

    # Meteorological station status values
    ACTIEF = 8789

    # Meteorological station communication status values
    ONLINE = 8796
        
    @classmethod
    def from_key(cls, key: int) -> Optional["HortimaxSynoptaEmun"]:
        try:
            return cls(key)
        except ValueError:
            print(f"Warning: Invalid HortimaxSynoptaEmun key: {key}")
            return None