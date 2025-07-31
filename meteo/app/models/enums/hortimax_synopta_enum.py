from enum import Enum

class HortimaxSynoptaEnum(Enum):
    """Hortimax Synopta enum values"""

    # Wind direction compass values
    N   = 8772
    NNO = 8773
    NO  = 8774
    ONO = 8775
    O   = 8776
    OZO = 8777
    ZO  = 8778
    ZZO = 8779
    Z   = 8780
    ZZW = 8781
    ZW = 8782
    WZW = 8783
    W = 8784
    WNW = 8785
    NW  = 8786
    NNW = 8787
    
    # Meteorological station status values -> Left expandable for Inactive, offline etc.
    Actief = 8789
    
    # Meteorological station communication status values
    Online = 8796