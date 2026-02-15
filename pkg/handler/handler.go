package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jilanisayyad/edgebeat/pkg/config"
	"github.com/jilanisayyad/edgebeat/pkg/controller"
)

// ResponseWithMetadata wraps the metric with timestamp
type ResponseWithMetadata struct {
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// Handler wraps the store and provides metric-specific endpoints
type Handler struct {
	store        *controller.Store
	integrations config.IntegrationConfig
}

// New creates a new handler with the given store
func New(store *controller.Store, integrations config.IntegrationConfig) *Handler {
	return &Handler{store: store, integrations: integrations}
}

// writeJSON handles common JSON response logic
func (h *Handler) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

// checkMethod validates HTTP method
func (h *Handler) checkMethod(w http.ResponseWriter, r *http.Request, allowed string) bool {
	if r.Method != allowed {
		w.Header().Set("Allow", allowed)
		h.writeJSON(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// getFullMetrics returns full system metrics
func (h *Handler) getFullMetrics(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	payload, ok := h.store.Get()
	if !ok {
		h.writeJSON(w, map[string]string{"error": "no data available"}, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}

// getCPUMetrics returns only CPU metrics
func (h *Handler) getCPUMetrics(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	data, ok := h.store.GetCPU()
	if !ok {
		h.writeJSON(w, map[string]string{"error": "no data available"}, http.StatusServiceUnavailable)
		return
	}

	h.writeJSON(w, ResponseWithMetadata{
		Timestamp: data.Timestamp,
		Data:      data.CPU,
	}, http.StatusOK)
}

// getMemoryMetrics returns only memory metrics
func (h *Handler) getMemoryMetrics(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	data, ok := h.store.GetMemory()
	if !ok {
		h.writeJSON(w, map[string]string{"error": "no data available"}, http.StatusServiceUnavailable)
		return
	}

	h.writeJSON(w, ResponseWithMetadata{
		Timestamp: data.Timestamp,
		Data:      data.Memory,
	}, http.StatusOK)
}

// getDiskMetrics returns only disk metrics
func (h *Handler) getDiskMetrics(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	data, ok := h.store.GetDisk()
	if !ok {
		h.writeJSON(w, map[string]string{"error": "no data available"}, http.StatusServiceUnavailable)
		return
	}

	h.writeJSON(w, ResponseWithMetadata{
		Timestamp: data.Timestamp,
		Data:      data.Disk,
	}, http.StatusOK)
}

// getNetworkMetrics returns only network metrics
func (h *Handler) getNetworkMetrics(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	data, ok := h.store.GetNetwork()
	if !ok {
		h.writeJSON(w, map[string]string{"error": "no data available"}, http.StatusServiceUnavailable)
		return
	}

	h.writeJSON(w, ResponseWithMetadata{
		Timestamp: data.Timestamp,
		Data:      data.Network,
	}, http.StatusOK)
}

// getSystemMetrics returns only system metrics
func (h *Handler) getSystemMetrics(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	data, ok := h.store.GetSystem()
	if !ok {
		h.writeJSON(w, map[string]string{"error": "no data available"}, http.StatusServiceUnavailable)
		return
	}

	h.writeJSON(w, ResponseWithMetadata{
		Timestamp: data.Timestamp,
		Data:      data.Host,
	}, http.StatusOK)
}

// getSensorMetrics returns only sensor metrics
func (h *Handler) getSensorMetrics(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	data, ok := h.store.GetSensors()
	if !ok {
		h.writeJSON(w, map[string]string{"error": "no data available"}, http.StatusServiceUnavailable)
		return
	}

	h.writeJSON(w, ResponseWithMetadata{
		Timestamp: data.Timestamp,
		Data:      data.Sensors,
	}, http.StatusOK)
}

type ModbusCapability struct {
	Enabled        bool     `json:"enabled"`
	Mode           string   `json:"mode"`
	Host           string   `json:"host"`
	Port           int      `json:"port"`
	UnitID         int      `json:"unit_id"`
	RequiredFields []string `json:"required_fields"`
	Notes          string   `json:"notes,omitempty"`
}

type OPCUACapability struct {
	Enabled            bool     `json:"enabled"`
	Endpoint           string   `json:"endpoint"`
	SecurityPolicy     string   `json:"security_policy"`
	SecurityMode       string   `json:"security_mode"`
	UsernameConfigured bool     `json:"username_configured"`
	PasswordConfigured bool     `json:"password_configured"`
	RequiredFields     []string `json:"required_fields"`
	Notes              string   `json:"notes,omitempty"`
}

type IntegrationCapabilities struct {
	Modbus ModbusCapability `json:"modbus"`
	OPCUA  OPCUACapability  `json:"opcua"`
}

func (h *Handler) getIntegrations(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	modbusRequired := modbusRequiredFields(h.integrations.Modbus.Mode)
	opcuaRequired := opcuaRequiredFields(
		h.integrations.OPCUA.SecurityMode,
		h.integrations.OPCUA.Username != "",
		h.integrations.OPCUA.Password != "",
	)

	resp := IntegrationCapabilities{
		Modbus: ModbusCapability{
			Enabled:        h.integrations.Modbus.Enabled,
			Mode:           h.integrations.Modbus.Mode,
			Host:           h.integrations.Modbus.Host,
			Port:           h.integrations.Modbus.Port,
			UnitID:         h.integrations.Modbus.UnitID,
			RequiredFields: modbusRequired,
			Notes:          h.integrations.Modbus.Notes,
		},
		OPCUA: OPCUACapability{
			Enabled:            h.integrations.OPCUA.Enabled,
			Endpoint:           h.integrations.OPCUA.Endpoint,
			SecurityPolicy:     h.integrations.OPCUA.SecurityPolicy,
			SecurityMode:       h.integrations.OPCUA.SecurityMode,
			UsernameConfigured: h.integrations.OPCUA.Username != "",
			PasswordConfigured: h.integrations.OPCUA.Password != "",
			RequiredFields:     opcuaRequired,
			Notes:              h.integrations.OPCUA.Notes,
		},
	}

	h.writeJSON(w, resp, http.StatusOK)
}

const (
	defaultFabricateBytes uint64 = 1024
	maxFabricateBytes     uint64 = 1 << 30
)

type repeatReader struct {
	pattern []byte
	offset  int
}

func (r *repeatReader) Read(p []byte) (int, error) {
	if len(r.pattern) == 0 {
		return 0, io.EOF
	}

	for i := range p {
		p[i] = r.pattern[r.offset]
		r.offset++
		if r.offset >= len(r.pattern) {
			r.offset = 0
		}
	}

	return len(p), nil
}

func (h *Handler) getFabricatedPayload(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	sizeBytes := defaultFabricateBytes
	sizeParam := strings.TrimSpace(r.URL.Query().Get("size"))
	if sizeParam != "" {
		parsedSize, err := parseByteSize(sizeParam)
		if err != nil {
			h.writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
		sizeBytes = parsedSize
	}

	if sizeBytes > maxFabricateBytes {
		h.writeJSON(w, map[string]string{"error": "size exceeds 1Gi limit"}, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatUint(sizeBytes, 10))
	w.Header().Set("X-Payload-Bytes", strconv.FormatUint(sizeBytes, 10))

	if sizeBytes == 0 {
		return
	}

	reader := &repeatReader{pattern: []byte("edgebeat")}
	_, _ = io.CopyN(w, reader, int64(sizeBytes))
}

func parseByteSize(input string) (uint64, error) {
	value := strings.TrimSpace(strings.ToLower(input))
	if value == "" {
		return 0, fmt.Errorf("size must be provided")
	}

	multiplier := uint64(1)
	suffixes := []struct {
		suffix     string
		multiplier uint64
	}{
		{"ki", 1024},
		{"mi", 1024 * 1024},
		{"gi", 1024 * 1024 * 1024},
		{"k", 1024},
		{"m", 1024 * 1024},
		{"g", 1024 * 1024 * 1024},
		{"b", 1},
	}

	for _, entry := range suffixes {
		if strings.HasSuffix(value, entry.suffix) && len(value) > len(entry.suffix) {
			multiplier = entry.multiplier
			value = strings.TrimSpace(strings.TrimSuffix(value, entry.suffix))
			break
		}
	}

	if value == "" {
		return 0, fmt.Errorf("invalid size format")
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size format")
	}

	if multiplier > 0 && parsed > (^(uint64(0))/multiplier) {
		return 0, fmt.Errorf("size is too large")
	}

	return parsed * multiplier, nil
}

func modbusRequiredFields(mode string) []string {
	if strings.EqualFold(mode, "rtu") {
		return []string{
			"mode",
			"serial_port",
			"baud_rate",
			"data_bits",
			"parity",
			"stop_bits",
			"unit_id",
		}
	}

	return []string{
		"mode",
		"host",
		"port",
		"unit_id",
	}
}

func opcuaRequiredFields(securityMode string, usernameSet bool, passwordSet bool) []string {
	required := []string{
		"endpoint",
		"security_policy",
		"security_mode",
	}

	if !strings.EqualFold(securityMode, "none") && (usernameSet || passwordSet) {
		required = append(required, "username", "password")
		return required
	}

	if usernameSet || passwordSet {
		required = append(required, "username", "password")
	}

	return required
}

// RegisterRoutes registers all metric endpoints to the mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux, basePrefix string) {
	if mux == nil {
		return
	}

	prefix := basePrefix

	// Full metrics endpoints
	mux.HandleFunc(prefix+"/health", h.getFullMetrics)
	mux.HandleFunc(prefix+"/metrics", h.getFullMetrics)

	// Individual metric endpoints
	mux.HandleFunc(prefix+"/metrics/cpu", h.getCPUMetrics)
	mux.HandleFunc(prefix+"/metrics/memory", h.getMemoryMetrics)
	mux.HandleFunc(prefix+"/metrics/disk", h.getDiskMetrics)
	mux.HandleFunc(prefix+"/metrics/network", h.getNetworkMetrics)
	mux.HandleFunc(prefix+"/metrics/system", h.getSystemMetrics)
	mux.HandleFunc(prefix+"/metrics/sensors", h.getSensorMetrics)
	mux.HandleFunc(prefix+"/integrations", h.getIntegrations)
	mux.HandleFunc(prefix+"/data/fabricate", h.getFabricatedPayload)

	// Health check endpoint (minimal)
	mux.HandleFunc(prefix+"/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}
	})
}
