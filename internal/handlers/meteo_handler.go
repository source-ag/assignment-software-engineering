package handlers

import (
	"net/http"

	"meteo-api/internal/middleware"
	"meteo-api/internal/models"
	"meteo-api/internal/service"

	"github.com/gin-gonic/gin"
)

type MeteoHandler struct {
	service    service.MeteoServiceInterface
	authMiddleware *middleware.AuthMiddleware
}

func NewMeteoHandler(service service.MeteoServiceInterface, authMiddleware *middleware.AuthMiddleware) *MeteoHandler {
	return &MeteoHandler{
		service:    service,
		authMiddleware: authMiddleware,
	}
}

func (h *MeteoHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	
	// Write endpoints (require write permission)
	api.POST("/meteo/measurements", h.authMiddleware.RequireAuth(middleware.WriteRole), h.CreateMeasurement)
	api.POST("/meteo/measurements/bulk", h.authMiddleware.RequireAuth(middleware.WriteRole), h.CreateMeasurementsBulk)
	
	// Read endpoints (require read permission)
	api.GET("/meteo/current", h.authMiddleware.RequireAuth(middleware.ReadRole), h.GetCurrentMeasurement)
	api.GET("/meteo/history", h.authMiddleware.RequireAuth(middleware.ReadRole), h.GetHistory)
	api.GET("/meteo/aggregated", h.authMiddleware.RequireAuth(middleware.ReadRole), h.GetAggregatedData)
	api.GET("/meteo/average", h.authMiddleware.RequireAuth(middleware.ReadRole), h.GetAverage)
	
	// Health check endpoint (no auth required)
	api.GET("/health", h.HealthCheck)
}

// CreateMeasurement handles single measurement creation
func (h *MeteoHandler) CreateMeasurement(c *gin.Context) {
	var rawData models.RawMeteoData
	if err := c.ShouldBindJSON(&rawData); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid JSON format: " + err.Error(),
		})
		return
	}

	measurement, err := h.service.ProcessRawData(&rawData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to process measurement: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    measurement,
		Message: "Measurement created successfully",
	})
}

// CreateMeasurementsBulk handles bulk measurement creation
func (h *MeteoHandler) CreateMeasurementsBulk(c *gin.Context) {
	var rawDataList []models.RawMeteoData
	if err := c.ShouldBindJSON(&rawDataList); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid JSON format: " + err.Error(),
		})
		return
	}

	if len(rawDataList) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "No measurements provided",
		})
		return
	}

	measurements, err := h.service.ProcessRawDataBulk(rawDataList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to process measurements: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    measurements,
		Message: "Measurements created successfully",
	})
}

// GetCurrentMeasurement returns the latest measurement
func (h *MeteoHandler) GetCurrentMeasurement(c *gin.Context) {
	measurement, err := h.service.GetLatestMeasurement()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get current measurement: " + err.Error(),
		})
		return
	}

	if measurement == nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "No measurements found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    measurement,
	})
}

// GetHistory returns historical measurements
func (h *MeteoHandler) GetHistory(c *gin.Context) {
	var req models.HistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid query parameters: " + err.Error(),
		})
		return
	}

	measurements, err := h.service.GetHistory(req.Period, req.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Failed to get history: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    measurements,
	})
}

// GetAggregatedData returns aggregated measurements
func (h *MeteoHandler) GetAggregatedData(c *gin.Context) {
	var req models.AggregationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid query parameters: " + err.Error(),
		})
		return
	}

	aggregations, err := h.service.GetAggregatedData(req.Period, req.Interval)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Failed to get aggregated data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    aggregations,
	})
}

// GetAverage returns overall average for a period
func (h *MeteoHandler) GetAverage(c *gin.Context) {
	period := c.Query("period")
	if period == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Period parameter is required",
		})
		return
	}

	average, err := h.service.GetAverage(period)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Failed to get average: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    average,
	})
}

// HealthCheck endpoint for service health monitoring
func (h *MeteoHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Service is healthy",
	})
}