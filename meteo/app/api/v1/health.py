"""
Health and system status API endpoints.
"""
from typing import Dict

from fastapi import APIRouter

router = APIRouter(
    prefix="/api/v1",
    tags=["health"]
)


@router.get("/health")
async def health_check() -> Dict[str, str]:
    """Health check endpoint."""
    return {"status": "healthy", "service": "meteo-api"}


@router.get("/")
async def api_root() -> Dict[str, str]:
    """API v1 root endpoint."""
    return {"message": "Meteo API v1", "version": "0.1.0"}
