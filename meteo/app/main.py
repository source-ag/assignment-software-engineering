"""
Main FastAPI application entry point.
"""
from typing import Dict, Any

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .core.database import create_tables
from .models import Measurement  # Import models to register them
from .api.v1 import health_router, measurements_router

# Create database tables on startup
create_tables()

app: FastAPI = FastAPI(
    title="Meteo API",
    description="Sensor Measurements API for Weather Data",
    version="0.1.0",
    docs_url="/docs",
    redoc_url="/redoc",
)

# Configure CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include API routers
app.include_router(health_router)
app.include_router(measurements_router)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
