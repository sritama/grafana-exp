# Go Application with Prometheus and Grafana Monitoring

This project demonstrates how to collect Go runtime metrics using Prometheus and visualize them with Grafana.

## Features

- **Go Runtime Metrics**: Collects memory usage, goroutine count, and garbage collection metrics
- **HTTP Metrics**: Tracks request counts, durations, and status codes
- **Prometheus Integration**: Exposes metrics on `/metrics` endpoint
- **Grafana Dashboard**: Pre-configured dashboard for monitoring Go applications
- **Docker Compose**: Easy setup with containerized services

## Metrics Collected

### Custom Metrics
- `http_requests_total`: Total HTTP requests with method, endpoint, and status labels
- `http_request_duration_seconds`: HTTP request duration histogram
- `goroutines_total`: Current number of goroutines
- `memory_alloc_bytes`: Current memory allocation
- `memory_heap_bytes`: Current heap memory usage
- `gc_duration_seconds`: Last garbage collection duration

### Go Runtime Metrics
The application also exposes standard Go runtime metrics provided by the Prometheus client library.

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.24+ (for local development)

### Running with Docker Compose

1. **Start all services**:
   ```bash
   docker-compose up -d
   ```

2. **Access the services**:
   - **Go Application**: http://localhost:8080
   - **Metrics Endpoint**: http://localhost:8080/metrics
   - **Health Check**: http://localhost:8080/health
   - **Prometheus**: http://localhost:9090
   - **Grafana**: http://localhost:3000 (admin/admin)

3. **View the Dashboard**:
   - Login to Grafana with `admin/admin`
   - The "Go Runtime Metrics" dashboard should be automatically loaded
   - If not, you can import the dashboard from `grafana/dashboards/go-runtime-metrics.json`

### Running Locally

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Run the application**:
   ```bash
   go run main.go
   ```

3. **Start Prometheus and Grafana**:
   ```bash
   docker-compose up -d prometheus grafana
   ```

## Project Structure

```
.
├── main.go                          # Go application with metrics
├── Dockerfile                       # Container configuration
├── docker-compose.yml              # Service orchestration
├── prometheus.yml                  # Prometheus configuration
├── grafana/
│   ├── provisioning/
│   │   ├── datasources/
│   │   │   └── prometheus.yml     # Grafana datasource config
│   │   └── dashboards/
│   │       └── dashboard.yml      # Dashboard provisioning
│   └── dashboards/
│       └── go-runtime-metrics.json # Grafana dashboard
└── README.md                       # This file
```

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics endpoint

## Monitoring Dashboard

The Grafana dashboard includes the following panels:

1. **HTTP Requests per Second** - Request rate over time
2. **HTTP Request Duration (95th percentile)** - Response time performance
3. **Number of Goroutines** - Active goroutine count
4. **Memory Usage** - Allocated and heap memory
5. **Garbage Collection Duration** - GC performance
6. **HTTP Requests by Status Code** - Request success/failure rates

## Customization

### Adding Custom Metrics

To add custom metrics, import the Prometheus client and create new metrics:

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    customCounter = promauto.NewCounter(prometheus.CounterOpts{
        Name: "custom_metric_total",
        Help: "Description of your metric",
    })
)
```

### Modifying the Dashboard

1. Edit `grafana/dashboards/go-runtime-metrics.json`
2. Restart the Grafana container: `docker-compose restart grafana`

## Troubleshooting

### Metrics Not Appearing
- Check that the application is running: `curl http://localhost:8080/health`
- Verify metrics endpoint: `curl http://localhost:8080/metrics`
- Check Prometheus targets: http://localhost:9090/targets

### Grafana Dashboard Issues
- Verify Prometheus datasource is configured correctly
- Check that the dashboard JSON is valid
- Restart Grafana container if needed

### Container Issues
- Check logs: `docker-compose logs [service-name]`
- Restart services: `docker-compose restart`

## Development

### Adding New Endpoints

When adding new endpoints, they will automatically be monitored by the metrics middleware. The middleware tracks:
- Request count
- Request duration
- Status codes

### Testing Metrics

Generate some load to see metrics in action:

```bash
# Generate HTTP requests
for i in {1..100}; do
  curl http://localhost:8080/health
  sleep 0.1
done
```

