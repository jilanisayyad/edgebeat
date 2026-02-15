# EdgeBeat - IoT Health Monitoring Agent

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.24.0-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![MQTT Support](https://img.shields.io/badge/MQTT-5.0-orange)](#mqtt-publishing)
[![REST API](https://img.shields.io/badge/REST-Enabled-brightgreen)](#rest-api)

A lightweight, production-ready system health monitoring agent for IoT devices and edge computing environments. Collects comprehensive system metrics and publishes them via REST API and MQTT in real-time.

</div>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [IoT Requirements](#iot-requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Feature Test Commands](#feature-test-commands)
- [REST API](#rest-api)
- [MQTT Publishing](#mqtt-publishing)
- [Metrics Collected](#metrics-collected)
- [Development](#development)
- [Release](#release)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

**EdgeBeat** is a high-performance system health monitoring agent designed for IoT devices and edge computing environments. It continuously collects comprehensive system metrics at configurable intervals (1-180 seconds) and makes them available through:

- **REST API** - Real-time HTTP endpoint for querying latest metrics
- **MQTT Publishing** - Stream metrics to MQTT brokers for centralized monitoring

Perfect for:

- IoT device fleet management
- Edge computing health monitoring
- Smart home automation
- Remote system diagnostics
- Time-series data collection

---

## Features

### Core Monitoring

- **CPU Metrics** - Per-core and total CPU usage, load times, frequencies
- **Memory Metrics** - Virtual memory, swap usage, detailed breakdown
- **Disk Metrics** - Partition info, disk usage, I/O statistics
- **Network Metrics** - Interface details, bandwidth usage, error/drop counts
- **System Metrics** - Hostname, OS, kernel version, uptime, boot time, virtualization info
- **Sensor Metrics** - Temperature readings from sensors

### Publishing and Access

- **REST API** - Lightweight HTTP endpoint for polling metrics
- **MQTT Publishing** - Stream metrics to brokers with configurable QoS
- **JSON Formatting** - Structured JSON output for easy integration
- **Error Handling** - Comprehensive error tracking and reporting

### Configuration and Operations

- **Configurable Frequency** - 1-180 second collection intervals
- **YAML Configuration** - Easy-to-read configuration files
- **Graceful Shutdown** - Clean signal handling (SIGTERM, SIGINT)
- **Logging** - Production-grade logging with Zap
- **Auto-Reconnect** - MQTT connection resilience with exponential backoff

---

## Project Structure

```
edgebeat/
|-- cmd/
|   `-- edgebeat/
|       `-- edgebeat.go           # Application entry point
|-- pkg/
|   |-- config/
|   |   `-- config.go             # Configuration loading and validation
|   |-- controller/
|   |   |-- controller.go         # System metrics collection
|   |   |-- root.go               # Collection loop and publishing
|   |   `-- store.go              # In-memory metrics storage
|   |-- mqtt/
|   |   `-- mqtt.go               # MQTT publisher implementation
|   `-- utils/
|       `-- utils.go              # Data structures and types
|-- configs/
|   `-- config.yaml               # Application configuration
|-- go.mod                        # Go module definition
|-- go.sum                        # Dependency checksums
|-- LICENSE                       # MIT License
`-- README.md                     # This file
```

---

## Prerequisites

- **Go** 1.24.0 or higher
- **Docker** (optional, for running MQTT broker)
- **MQTT Broker** (if using MQTT) - Mosquitto recommended for testing

### Verify Prerequisites

```bash
# Check Go version
go version

# Check Docker (optional)
docker --version
```

---

## IoT Requirements

### Hardware

- CPU: 1 core at 1 GHz or higher for standard intervals
- Memory: 512 MB minimum, 1 GB or more for high-frequency sampling
- Storage: 100 MB for binary and config, plus space for logs and rotation

### OS and Packages

- Linux, macOS, or Windows with access to system telemetry
- On Linux, ensure access to `/proc` and `/sys` for full metrics
- Recommended packages: CA certificates, timezone data, `curl`, `jq`

### Networking and MQTT

- Stable network connectivity to the MQTT broker (if enabled)
- Open outbound ports 1883 (MQTT) or 8883 (MQTTS)
- NTP time sync for accurate timestamps

### Security and Certificates

- Use TLS with broker certificates for production MQTT
- Store credentials in environment variables or secret stores
- Run the agent as a non-root user when possible

### Observability

- Forward logs to a centralized system (Fluent Bit, Vector, or syslog)
- Monitor the `/ping` endpoint for liveness checks

---

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/edgebeat.git
cd edgebeat
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Build the Application

```bash
# Build binary
go build -o edgebeat ./cmd/edgebeat

# Or run directly
go run ./cmd/edgebeat
```

### 4. Docker Installation (Optional)

```dockerfile
FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o edgebeat ./cmd/edgebeat

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/edgebeat .
COPY --from=builder /app/configs/config.yaml ./configs/
CMD ["./edgebeat"]
```

Build and run:

```bash
docker build -t edgebeat .
docker run -p 8080:8080 edgebeat
```

---

## Configuration

### Configuration File: `configs/config.yaml`

```yaml
# Metrics collection frequency (1-180 seconds)
frequency_seconds: 60

# REST API Configuration
rest:
  address: ":8080" # Server address and port
  path: "/health" # Endpoint path

# MQTT Publishing Configuration (optional)
mqtt:
  enabled: true # Enable/disable MQTT
  broker: "tcp://localhost:1883"
  client_id: "edgebeat"
  topic: "edgebeat/health"
  username: "" # Leave empty if not needed
  password: "" # Leave empty if not needed
  qos: 1 # QoS level: 0, 1, or 2

# Industrial Protocol Integrations (optional)
integrations:
  modbus:
    enabled: false
    mode: "tcp" # tcp or rtu
    host: "localhost"
    port: 502
    unit_id: 1
    notes: ""
  opcua:
    enabled: false
    endpoint: "opc.tcp://localhost:4840"
    security_policy: "None"
    security_mode: "None"
    username: ""
    password: ""
    notes: ""
```

### Configuration Parameters

| Parameter           | Type    | Range | Default   | Description                               |
| ------------------- | ------- | ----- | --------- | ----------------------------------------- |
| `frequency_seconds` | integer | 1-180 | 60        | How often to collect metrics (in seconds) |
| `rest.address`      | string  | -     | `:8080`   | REST server bind address                  |
| `rest.path`         | string  | -     | `/health` | REST endpoint path                        |
| `mqtt.enabled`      | boolean | -     | false     | Enable MQTT publishing                    |
| `mqtt.broker`       | string  | -     | -         | MQTT broker URL (tcp://host:port)         |
| `mqtt.qos`          | integer | 0-2   | 1         | MQTT QoS level                            |

#### Integration Parameters

| Parameter                            | Type    | Default                    | Description                        |
| ------------------------------------ | ------- | -------------------------- | ---------------------------------- |
| `integrations.modbus.enabled`        | boolean | false                      | Enable Modbus integration metadata |
| `integrations.modbus.mode`           | string  | `tcp`                      | Modbus mode: tcp or rtu            |
| `integrations.modbus.host`           | string  | `localhost`                | Modbus TCP host                    |
| `integrations.modbus.port`           | integer | 502                        | Modbus TCP port                    |
| `integrations.modbus.unit_id`        | integer | 1                          | Modbus unit id                     |
| `integrations.opcua.enabled`         | boolean | false                      | Enable OPC UA integration metadata |
| `integrations.opcua.endpoint`        | string  | `opc.tcp://localhost:4840` | OPC UA server endpoint             |
| `integrations.opcua.security_policy` | string  | `None`                     | OPC UA security policy             |
| `integrations.opcua.security_mode`   | string  | `None`                     | OPC UA security mode               |

### Configuration Examples

#### Minimal Configuration (REST Only)

```yaml
frequency_seconds: 30
rest:
  address: ":8080"
  path: "/health"
mqtt:
  enabled: false
```

#### Full Configuration (REST + MQTT)

```yaml
frequency_seconds: 60
rest:
  address: ":8080"
  path: "/metrics"
mqtt:
  enabled: true
  broker: "tcp://broker.example.com:1883"
  client_id: "edge-device-001"
  topic: "devices/edge-001/metrics"
  username: "mqttuser"
  password: "strongpassword"
  qos: 1
```

#### High-Frequency Monitoring

```yaml
frequency_seconds: 5
rest:
  address: ":9090"
  path: "/health"
mqtt:
  enabled: true
  broker: "tcp://localhost:1883"
  qos: 2
```

#### Integrations Metadata Example

```yaml
integrations:
  modbus:
    enabled: true
    mode: "tcp"
    host: "192.168.1.50"
    port: 502
    unit_id: 1
  opcua:
    enabled: true
    endpoint: "opc.tcp://192.168.1.60:4840"
    security_policy: "Basic256Sha256"
    security_mode: "SignAndEncrypt"
    username: "opcua-user"
    password: "opcua-pass"
```

---

## Usage

### Basic Usage

1. **Configure the application:**

   ```bash
   # Edit configs/config.yaml as needed
   nano configs/config.yaml
   ```

2. **Run the application:**

   ```bash
   go run ./cmd/edgebeat
   ```

3. **Query metrics via REST:**
   ```bash
   curl http://localhost:8080/health | jq
   ```

### Using with MQTT Broker

#### Option 1: Docker Mosquitto (Recommended for Testing)

```bash
# Start Mosquitto broker
docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto

# Enable MQTT in config
# mqtt.enabled: true
# mqtt.broker: "tcp://localhost:1883"

# Run edgebeat
go run ./cmd/edgebeat

# In another terminal, subscribe to metrics
docker exec mosquitto mosquitto_sub -t "edgebeat/health"
```

#### Option 2: Homebrew (macOS)

```bash
# Install Mosquitto
brew install mosquitto

# Start broker
brew services start mosquitto

# Run edgebeat
go run ./cmd/edgebeat

# Subscribe to messages
mosquitto_sub -h localhost -t "edgebeat/health"
```

#### Option 3: Manual Installation

Download from [mosquitto.org](https://mosquitto.org/download/) and run:

```bash
mosquitto -c /path/to/mosquitto.conf
```

### Monitoring Logs

The application uses structured logging (Zap). Look for messages like:

```
{"level":"info","ts":1708055232.123,"msg":"rest server started","address":":8080","path":"/health"}
{"level":"info","ts":1708055232.456,"msg":"mqtt connected","broker":"tcp://localhost:1883"}
{"level":"info","ts":1708055240.789,"msg":"system info collected"}
```

---

## Feature Test Commands

### Run the agent

```bash
go run ./cmd/edgebeat
```

### REST endpoints

```bash
curl http://localhost:8080/health | jq
curl http://localhost:8080/metrics/cpu | jq
curl http://localhost:8080/metrics/memory | jq
curl http://localhost:8080/metrics/disk | jq
curl http://localhost:8080/metrics/network | jq
curl http://localhost:8080/metrics/system | jq
curl http://localhost:8080/metrics/sensors | jq
curl http://localhost:8080/integrations | jq
```

### Payload fabrication

```bash
curl -o /tmp/payload-1Ki.bin "http://localhost:8080/data/fabricate?size=1Ki"
curl -o /tmp/payload-1Mi.bin "http://localhost:8080/data/fabricate?size=1Mi"
curl -I "http://localhost:8080/data/fabricate?size=0"
```

### MQTT publish and subscribe

```bash
docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto
mosquitto_sub -h localhost -t "edgebeat/health"
```

### CLI download example (local or container)

```bash
VERSION=0.1.0
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')
curl -L -o edgebeat.tar.gz \
  "https://github.com/yourusername/edgebeat/releases/download/v${VERSION}/edgebeat_${VERSION}_${OS}_${ARCH}.tar.gz"
tar -xzf edgebeat.tar.gz
./edgebeat --help
```

```bash
docker run --rm -it alpine:3.20 sh -c '
  apk add --no-cache ca-certificates curl tar && \
  VERSION=0.1.0 && \
  curl -L -o /tmp/edgebeat.tar.gz \
    "https://github.com/yourusername/edgebeat/releases/download/v${VERSION}/edgebeat_${VERSION}_linux_amd64.tar.gz" && \
  tar -xzf /tmp/edgebeat.tar.gz -C /usr/local/bin && \
  edgebeat --help
'
```

---

## REST API

### Available Endpoints

EdgeBeat provides multiple endpoints for accessing system metrics:

| Endpoint           | Method | Description                               |
| ------------------ | ------ | ----------------------------------------- |
| `/health`          | GET    | Full system metrics (all data)            |
| `/metrics`         | GET    | Full system metrics (synonym for /health) |
| `/metrics/cpu`     | GET    | CPU metrics only                          |
| `/metrics/memory`  | GET    | Memory metrics only                       |
| `/metrics/disk`    | GET    | Disk metrics only                         |
| `/metrics/network` | GET    | Network metrics only                      |
| `/metrics/system`  | GET    | System info only                          |
| `/metrics/sensors` | GET    | Temperature sensors only                  |
| `/integrations`    | GET    | Modbus and OPC UA configuration info      |
| `/data/fabricate`  | GET    | Generate synthetic payload bytes          |
| `/ping`            | GET    | Health check (minimal response)           |

### Full Metrics Endpoint

```
GET /health
GET /metrics
```

Returns comprehensive system metrics in JSON format.

#### Example Request

```bash
curl http://localhost:8080/health | jq
```

#### Example Response (Abbreviated)

```json
{
  "timestamp": "2024-02-15T10:30:45.123456789Z",
  "cpu": {
    "info": [
      {
        "model_name": "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz",
        "cores": 6,
        "mhz": 2600.0,
        "cache_size": 9216
      }
    ],
    "total_percent": 35.5,
    "per_cpu_percent": [30.2, 35.5, 40.1, 32.3, 36.7, 38.5],
    "total_times": {
      "user": 1234.5,
      "system": 234.5,
      "idle": 5432.1
    }
  },
  "memory": {
    "virtual": {
      "total": 17179869184,
      "available": 8589934592,
      "used": 8589934592,
      "used_percent": 50.0
    },
    "swap": {
      "total": 0,
      "used": 0
    }
  },
  "network": {
    "interfaces": [
      {
        "name": "eth0",
        "mtu": 1500,
        "hardware_addr": "aa:bb:cc:dd:ee:ff",
        "addrs": ["192.168.1.100/24"]
      }
    ],
    "totals": {
      "bytes_sent": 1073741824,
      "bytes_recv": 2147483648,
      "packets_sent": 1000000,
      "packets_recv": 2000000
    }
  },
  "disk": {
    "partitions": [
      {
        "device": "/dev/sda1",
        "mountpoint": "/",
        "fs_type": "ext4"
      }
    ],
    "usage": [
      {
        "device": "/dev/sda1",
        "mountpoint": "/",
        "total": 107374182400,
        "used": 53687091200,
        "free": 53687091200,
        "used_percent": 50.0
      }
    ]
  },
  "host": {
    "hostname": "edge-device",
    "os": "linux",
    "platform": "linux",
    "uptime_seconds": 123456,
    "boot_time": 1707968000,
    "procs": 256
  },
  "sensors": {
    "temperatures": [
      {
        "sensor_key": "coretemp_core_0",
        "value": 45.5,
        "critical": 100.0
      }
    ]
  }
}
```

### Distributed Metric Endpoints

#### CPU Metrics

Get detailed CPU usage and performance information.

```bash
curl http://localhost:8080/metrics/cpu | jq
```

#### Memory Metrics

Get RAM and swap memory information.

```bash
curl http://localhost:8080/metrics/memory | jq
```

#### Disk Metrics

Get storage and disk I/O information.

```bash
curl http://localhost:8080/metrics/disk | jq
```

#### Network Metrics

Get network interface and I/O statistics.

```bash
curl http://localhost:8080/metrics/network | jq
```

#### System Metrics

Get system information (hostname, OS, uptime, etc).

```bash
curl http://localhost:8080/metrics/system | jq
```

#### Sensor Metrics

Get temperature and hardware sensor data.

```bash
curl http://localhost:8080/metrics/sensors | jq
```

#### Health Check

Quick health check endpoint with minimal response.

```bash
curl http://localhost:8080/ping
```

### Status Codes

| Code | Meaning                                           |
| ---- | ------------------------------------------------- |
| 200  | Metrics successfully retrieved                    |
| 405  | Method not allowed (only GET allowed)             |
| 503  | No metrics available (collection not started yet) |

### Integration Capabilities

Retrieve Modbus and OPC UA integration metadata and required fields.

```bash
curl http://localhost:8080/integrations | jq
```

Example response:

```json
{
  "modbus": {
    "enabled": false,
    "mode": "tcp",
    "host": "localhost",
    "port": 502,
    "unit_id": 1,
    "required_fields": ["mode", "host", "port", "unit_id"]
  },
  "opcua": {
    "enabled": false,
    "endpoint": "opc.tcp://localhost:4840",
    "security_policy": "None",
    "security_mode": "None",
    "username_configured": false,
    "password_configured": false,
    "required_fields": ["endpoint", "security_policy", "security_mode"]
  }
}
```

### Payload Fabrication

Generate a synthetic byte payload for throughput and pipeline testing.

```bash
curl -o /tmp/payload.bin "http://localhost:8080/data/fabricate?size=128Ki"
```

Size rules:

- Range: 0 to 1Gi
- Supported suffixes: `Ki`, `Mi`, `Gi`, `K`, `M`, `G`, or bytes without a suffix
- Default size when omitted: 1Ki

### Error Response Example

```json
{
  "error": "no data available"
}
```

---

## MQTT Publishing

### Overview

When MQTT is enabled, metrics are automatically published at the configured frequency to the specified topic.

### Configuration

```yaml
mqtt:
  enabled: true
  broker: "tcp://localhost:1883"
  client_id: "edgebeat"
  topic: "edgebeat/health"
  qos: 1
```

### Message Format

MQTT messages contain the same JSON payload as the REST API response.

### Subscribe to Metrics

Using `mosquitto_sub`:

```bash
mosquitto_sub -h localhost -p 1883 -t "edgebeat/health"
```

Using Node.js MQTT client:

```javascript
const mqtt = require("mqtt");
const client = mqtt.connect("mqtt://localhost:1883");

client.on("connect", () => {
  client.subscribe("edgebeat/health");
});

client.on("message", (topic, message) => {
  console.log(JSON.parse(message.toString()));
});
```

### QoS Levels

- **QoS 0** - Fire and forget (fastest, no guarantees)
- **QoS 1** - At least once (default, recommended)
- **QoS 2** - Exactly once (slowest, most reliable)

### Monitoring MQTT Traffic

```bash
# Terminal 1: Start broker
docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto

# Terminal 2: Subscribe to all topics
mosquitto_sub -h localhost -v -t '#'

# Terminal 3: Run edgebeat
go run ./cmd/edgebeat
```

---

## Metrics Collected

### CPU Metrics

- Model name, cores, frequency, cache size
- Per-core and total CPU usage percentage
- CPU times (user, system, idle, nice, iowait, irq, etc.)

### Memory Metrics

- Virtual memory (total, used, available, cached, buffered)
- Swap memory (total, used, available)
- Usage percentages

### Disk Metrics

- Partition information
- Per-partition usage (total, used, free, percentage)
- I/O statistics (read/write bytes and counts)

### Network Metrics

- Interface information (name, MAC address, MTU, IPs)
- Total network I/O (sent/received bytes and packets)
- Error and drop statistics

### System Metrics

- Hostname, OS, platform, version information
- Kernel version and architecture
- Uptime and boot time
- Number of running processes
- Virtualization system and role
- Active user sessions

### Sensor Metrics

- Temperature readings from system sensors
- High and critical temperature thresholds

---

## Development

### Project Structure

```
edgebeat/
|-- cmd/edgebeat/         # Entry point
|-- pkg/
|   |-- config/           # Configuration management
|   |-- controller/       # Business logic
|   |-- mqtt/             # MQTT implementation
|   `-- utils/            # Data structures
```

### Dependencies

View all dependencies:

```bash
go mod graph
```

Key dependencies:

- `github.com/shirou/gopsutil/v4` - System metrics collection
- `github.com/eclipse/paho.mqtt.golang` - MQTT client
- `go.uber.org/zap` - Structured logging
- `gopkg.in/yaml.v3` - YAML configuration

### Building

```bash
# Debug build
go build -o edgebeat ./cmd/edgebeat

# Release build
CGO_ENABLED=0 go build -ldflags="-s -w" -o edgebeat ./cmd/edgebeat

# Cross-compile for Linux ARM64
GOOS=linux GOARCH=arm64 go build -o edgebeat-arm64 ./cmd/edgebeat
```

### Testing Locally

```bash
# Terminal 1: Start MQTT broker
docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto

# Terminal 2: Run application
go run ./cmd/edgebeat

# Terminal 3: Monitor REST API
watch -n 5 'curl -s http://localhost:8080/health | jq .cpu.total_percent'

# Terminal 4: Monitor MQTT
docker exec mosquitto mosquitto_sub -t "edgebeat/health" | jq .timestamp
```

---

## Release

### Taskfile (go-task)

```bash
task test
task test:all
task build
task release:snapshot
task release
```

### GoReleaser

```bash
goreleaser check
goreleaser release --snapshot --clean
goreleaser release --clean
```

### Conventional Commits

```bash
task hooks:install
```

### Install From GitHub Releases

```bash
VERSION=0.1.0
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')
curl -L -o edgebeat.tar.gz \
  "https://github.com/jilanisayyad/edgebeat/releases/download/v${VERSION}/edgebeat_${VERSION}_${OS}_${ARCH}.tar.gz"
tar -xzf edgebeat.tar.gz
sudo mv edgebeat /usr/local/bin/edgebeat
edgebeat --help
```

Windows (PowerShell):

```powershell
$Version = "0.1.0"
$Os = "windows"
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
$Uri = "https://github.com/jilanisayyad/edgebeat/releases/download/v$Version/edgebeat_${Version}_${Os}_${Arch}.zip"
Invoke-WebRequest -Uri $Uri -OutFile "edgebeat.zip"
Expand-Archive -Path "edgebeat.zip" -DestinationPath ".\edgebeat"
Move-Item ".\edgebeat\edgebeat.exe" "$env:LOCALAPPDATA\Microsoft\WindowsApps\edgebeat.exe"
edgebeat --help
```

Notes:

- Create a git tag before releasing: `git tag v0.1.0 && git push --tags`
- Set `GITHUB_TOKEN` for GitHub releases

---

## Troubleshooting

### Application won't start

**Problem:** `config file not found`

```
Solution: Ensure configs/config.yaml exists in the current directory
```

### REST endpoint not responding

**Problem:** `curl: (7) Failed to connect to localhost port 8080`

```bash
# Check if application is running
ps aux | grep edgebeat

# Check if port is in use
lsof -i :8080

# Try different port in config.yaml
rest:
  address: ":9090"
```

### MQTT connection fails

**Problem:** `mqtt connect: connection refused`

```bash
# Verify broker is running
# For Docker:
docker ps | grep mosquitto

# For Homebrew:
brew services list

# Test connection manually
mosquitto_pub -h localhost -t test -m "hello"
```

### High CPU usage

**Problem:** Application consuming excessive CPU

```yaml
# Increase collection interval
frequency_seconds: 180 # Maximum is 180 seconds
```

### Memory metrics showing zeros

**Problem:** Swap memory not available on system

```
This is normal on systems without swap. Value will be 0.
```

### Network interface not showing

**Problem:** Some interfaces missing from output

```
Some interfaces may be filtered by gopsutil. Check /proc/net/dev on Linux.
```

### Logging issues

**Problem:** Logs not visible or too verbose

```go
// Edit cmd/edgebeat/edgebeat.go to change log level
// zap.NewProduction() for production
// zap.NewDevelopment() for development
```

---

## Contributing

Contributions are welcome! Here's how to contribute:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines

- Use meaningful variable names
- Add comments for complex logic
- Test with both REST and MQTT enabled
- Update README if adding new features
- Follow Go conventions (gofmt, golint)

### Reporting Issues

Please include:

- OS and hardware information
- Go version (`go version`)
- Configuration file (sanitized)
- Error messages and logs
- Steps to reproduce

---

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## Author

Jilani Sayyad

## Support

For issues, questions, or suggestions:

- **GitHub Issues:** [Report an issue](https://github.com/yourusername/edgebeat/issues)
- **Discussions:** [Start a discussion](https://github.com/yourusername/edgebeat/discussions)

---

## Acknowledgments

- [gopsutil](https://github.com/shirou/gopsutil) for system metrics
- [Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang) for MQTT support
- [Zap](https://github.com/uber-go/zap) for structured logging
- [YAML Go](https://github.com/go-yaml/yaml) for configuration parsing

---

Made for the IoT community
