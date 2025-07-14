# Meteo API - Sensor Measurements Backend

A RESTful API service for ingesting and querying weather sensor measurements from greenhouse climate computers. Built with Go, PostgreSQL, and designed for scalability and maintainability.

## Features

- **Data Ingestion**: Accept raw sensor data from climate computers
- **Flexible Querying**: Get current conditions, historical data, and aggregated metrics
- **Authentication**: Role-based API key authentication (read/write permissions)
- **Bulk Operations**: Support for bulk data ingestion
- **Flexible Aggregation**: Configurable time periods and intervals
- **Persistent Storage**: PostgreSQL with JSONB for raw data and normalized columns
- **Comprehensive Testing**: Unit tests with mocks and integration tests
- **Docker Support**: Containerized deployment with Docker Compose

## Architecture

The application follows a clean architecture pattern:

```
cmd/
├── server/           # Application entry point
internal/
├── models/           # Data structures and DTOs
├── database/         # Database connection and migrations
├── repository/       # Data access layer
├── service/          # Business logic layer
├── handlers/         # HTTP handlers (controllers)
└── middleware/       # Authentication and other middleware
```

## API Endpoints

### Authentication
All endpoints (except health check) require authentication via API key in the Authorization header:
```
Authorization: Bearer <api-key>
```

**API Keys:**
- Read operations: `read-key-123` (default)
- Write operations: `write-key-456` (default)

### Write Endpoints (require write permission)

#### POST /api/v1/meteo/measurements
Ingest a single measurement from climate computer.

**Request Body:**
```json
{
  "name": "_ws_source_meteo",
  "range": {"col1": 1, "row1": 1, "col2": 2, "row2": 22},
  "rows": [
    ["Variable", "Value"],
    ["external_temperature_c", 9.1387],
    ["wind_speed_m_s", 3.06032],
    ["relative_humidity_perc", 73]
  ],
  "ts": "2021-05-01T12:07:50+02:00",
  "pt": 0
}
```

#### POST /api/v1/meteo/measurements/bulk
Ingest multiple measurements in a single request.

**Request Body:**
```json
[
  {
    "name": "_ws_source_meteo",
    "rows": [...],
    "ts": "2021-05-01T12:07:50+02:00",
    "pt": 0
  },
  ...
]
```

### Read Endpoints (require read permission)

#### GET /api/v1/meteo/current
Get the latest weather measurement.

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "timestamp": "2021-05-01T12:07:50+02:00",
    "external_temperature_c": 9.1387,
    "wind_speed_m_s": 3.06032,
    "relative_humidity_perc": 73,
    ...
  }
}
```

#### GET /api/v1/meteo/history
Get historical measurements for a specified period.

**Query Parameters:**
- `period` (required): Time period (e.g., "24h", "7d", "30d")
- `limit` (optional): Maximum number of results

**Example:** `GET /api/v1/meteo/history?period=24h&limit=100`

#### GET /api/v1/meteo/aggregated
Get aggregated measurements over specified intervals.

**Query Parameters:**
- `period` (required): Time period to aggregate over (e.g., "24h", "7d")
- `interval` (required): Aggregation interval (e.g., "15m", "1h", "1d")

**Example:** `GET /api/v1/meteo/aggregated?period=24h&interval=15m`

#### GET /api/v1/meteo/average
Get overall average for all parameters over a specified period.

**Query Parameters:**
- `period` (required): Time period (e.g., "24h", "7d")

**Example:** `GET /api/v1/meteo/average?period=24h`

#### GET /api/v1/health
Health check endpoint (no authentication required).

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Option 1: Using Docker Compose (Recommended)

1. **Clone and setup:**
   ```bash
   git clone <repository>
   cd meteo-api
   cp .env.example .env
   ```

2. **Start services:**
   ```bash
   make docker-up
   ```

3. **Run migrations:**
   ```bash
   make migrate-up
   ```

4. **Test the API:**
   ```bash
   curl -H "Authorization: Bearer read-key-123" http://localhost:8080/api/v1/health
   ```

### Option 2: Local Development

1. **Setup environment:**
   ```bash
   make dev-setup
   ```

2. **Run the application:**
   ```bash
   make run
   ```

## Development

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Build and run the application
make test              # Run tests
make test-coverage     # Run tests with coverage
make deps              # Download dependencies
make fmt               # Format code
make vet               # Vet code
make lint              # Lint code (requires golangci-lint)
make docker-up         # Start services with Docker
make docker-down       # Stop Docker services
make migrate-up        # Run database migrations
make migrate-down      # Rollback database migrations
make dev-setup         # Setup development environment
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./internal/service/...
```

### Database Migrations

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Create new migration
make migrate-create name=add_new_field
```

## Configuration

Configuration is handled through environment variables. Copy `.env.example` to `.env` and modify as needed:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=meteo_db
DB_SSLMODE=disable

# Server Configuration
PORT=8080

# API Keys
READ_API_KEY=read-key-123
WRITE_API_KEY=write-key-456

# Gin Mode
GIN_MODE=debug
```

## Data Model

### Raw Data Format
The API accepts raw data in the format provided by climate computers:

```json
{
  "name": "_ws_source_meteo",
  "range": {"col1": 1, "row1": 1, "col2": 2, "row2": 22},
  "rows": [
    ["Variable", "Value"],
    ["external_temperature_c", 9.1387],
    ["wind_speed_m_s", 3.06032],
    ...
  ],
  "ts": "2021-05-01T12:07:50+02:00",
  "pt": 0
}
```

### Normalized Storage
Data is parsed and stored in normalized columns for efficient querying:

- `external_temperature_c`: External temperature in Celsius
- `wind_speed_m_s`: Wind speed in m/s
- `relative_humidity_perc`: Relative humidity percentage
- `atmospheric_pressure_hpa`: Atmospheric pressure in hPa
- And many more weather parameters...

Raw data is also preserved in JSONB format for future extensibility.

## Design Decisions

### Architecture
- **Clean Architecture**: Separation of concerns with distinct layers
- **Dependency Injection**: Loose coupling between components
- **Interface-based Design**: Easy testing and mocking

### Database Design
- **Hybrid Approach**: Normalized columns for performance + JSONB for flexibility
- **Indexes**: Optimized for time-based queries
- **Migrations**: Version-controlled schema changes

### API Design
- **RESTful**: Standard HTTP methods and status codes
- **Flexible Aggregation**: Configurable time periods and intervals
- **Consistent Response Format**: Standardized JSON responses
- **Authentication**: Role-based access control

### Error Handling
- **Comprehensive**: Detailed error messages and logging
- **HTTP Status Codes**: Proper status codes for different scenarios
- **Validation**: Input validation with clear error messages

## Performance Considerations

- **Database Indexes**: Optimized for time-based queries
- **Connection Pooling**: Efficient database connection management
- **Bulk Operations**: Support for bulk data ingestion
- **Aggregation**: Efficient SQL queries for time-based aggregations

## Security

- **API Key Authentication**: Simple but effective authentication
- **Role-based Access**: Separate read/write permissions
- **Input Validation**: Comprehensive input validation
- **SQL Injection Prevention**: Parameterized queries

## Monitoring and Observability

- **Health Check**: Endpoint for service health monitoring
- **Logging**: Structured logging with Gin middleware
- **Error Tracking**: Comprehensive error logging

## Future Improvements

### Short-term
- Add more comprehensive validation
- Implement rate limiting
- Add request/response logging
- Implement graceful shutdown

### Long-term
- Add caching layer (Redis)
- Implement real-time notifications
- Add metrics and monitoring (Prometheus)
- Implement data retention policies
- Add more sophisticated authentication (JWT, OAuth)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License.