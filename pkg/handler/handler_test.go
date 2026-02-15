package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jilanisayyad/edgebeat/pkg/config"
	"github.com/jilanisayyad/edgebeat/pkg/controller"
	"github.com/jilanisayyad/edgebeat/pkg/utils"
)

type metaResponse struct {
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

func seedStore(t *testing.T) (*controller.Store, utils.SystemInfo, []byte) {
	t.Helper()
	info := utils.SystemInfo{
		Timestamp: "2026-02-15T00:00:00Z",
		CPU:       utils.CPUStats{TotalPercent: 12.3},
		Memory:    utils.MemoryStats{Virtual: utils.VirtualMemory{Total: 1}},
		Disk:      utils.DiskStats{Partitions: []utils.DiskPartition{{Device: "disk0"}}},
		Network:   utils.NetworkStats{Totals: utils.NetIO{BytesSent: 1}},
		Host:      utils.HostStats{Hostname: "test-host"},
		Sensors:   utils.SensorsStats{Temperatures: []utils.Temperature{{SensorKey: "cpu"}}},
	}
	payload, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	store := controller.NewStore()
	store.Set(payload)
	return store, info, payload
}

func TestParseByteSize(t *testing.T) {
	cases := []struct {
		input string
		want  uint64
		ok    bool
	}{
		{"1", 1, true},
		{"1k", 1024, true},
		{"2ki", 2048, true},
		{"1m", 1024 * 1024, true},
		{"1g", 1024 * 1024 * 1024, true},
		{"1b", 1, true},
		{"", 0, false},
		{"bad", 0, false},
	}
	for _, tc := range cases {
		got, err := parseByteSize(tc.input)
		if tc.ok && err != nil {
			t.Fatalf("parseByteSize(%q) error: %v", tc.input, err)
		}
		if !tc.ok && err == nil {
			t.Fatalf("parseByteSize(%q) expected error", tc.input)
		}
		if tc.ok && got != tc.want {
			t.Fatalf("parseByteSize(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestModbusRequiredFields(t *testing.T) {
	rtu := modbusRequiredFields("rtu")
	if len(rtu) == 0 || rtu[0] != "mode" {
		t.Fatalf("modbusRequiredFields(rtu) = %v", rtu)
	}
	tcp := modbusRequiredFields("tcp")
	if len(tcp) == 0 || tcp[1] != "host" {
		t.Fatalf("modbusRequiredFields(tcp) = %v", tcp)
	}
}

func TestOpcuaRequiredFields(t *testing.T) {
	base := opcuaRequiredFields("none", false, false)
	if len(base) != 3 {
		t.Fatalf("opcuaRequiredFields base = %v", base)
	}
	withCreds := opcuaRequiredFields("Sign", true, false)
	if len(withCreds) != 5 {
		t.Fatalf("opcuaRequiredFields creds = %v", withCreds)
	}
}

func TestRepeatReader(t *testing.T) {
	reader := &repeatReader{pattern: []byte("ab")}
	buf := make([]byte, 5)
	n, err := reader.Read(buf)
	if err != nil || n != 5 {
		t.Fatalf("Read = %d, err=%v", n, err)
	}
	if string(buf) != "ababa" {
		t.Fatalf("Read data = %q", string(buf))
	}

	empty := &repeatReader{}
	buf = make([]byte, 1)
	n, err = empty.Read(buf)
	if err != io.EOF || n != 0 {
		t.Fatalf("empty Read = %d, err=%v", n, err)
	}
}

func TestGetFullMetrics(t *testing.T) {
	store, _, payload := seedStore(t)
	h := New(store, config.IntegrationConfig{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.getFullMetrics(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != strings.TrimSpace(string(payload)) {
		t.Fatalf("payload mismatch")
	}
}

func TestGetFullMetricsNoData(t *testing.T) {
	h := New(controller.NewStore(), config.IntegrationConfig{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.getFullMetrics(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestGetCPUMetrics(t *testing.T) {
	store, info, _ := seedStore(t)
	h := New(store, config.IntegrationConfig{})
	req := httptest.NewRequest(http.MethodGet, "/metrics/cpu", nil)
	rec := httptest.NewRecorder()
	h.getCPUMetrics(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	var resp metaResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if resp.Timestamp != info.Timestamp {
		t.Fatalf("Timestamp = %q", resp.Timestamp)
	}
}

func TestGetIntegrations(t *testing.T) {
	integrations := config.IntegrationConfig{
		Modbus: config.ModbusConfig{Enabled: true, Mode: "tcp", Host: "localhost", Port: 502, UnitID: 1, Notes: "note"},
		OPCUA:  config.OPCUAConfig{Enabled: true, Endpoint: "opc.tcp://localhost:4840", SecurityPolicy: "None", SecurityMode: "None"},
	}
	h := New(controller.NewStore(), integrations)
	req := httptest.NewRequest(http.MethodGet, "/integrations", nil)
	rec := httptest.NewRecorder()
	h.getIntegrations(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestGetFabricatedPayload(t *testing.T) {
	h := New(controller.NewStore(), config.IntegrationConfig{})

	req := httptest.NewRequest(http.MethodGet, "/data/fabricate?size=5", nil)
	rec := httptest.NewRecorder()
	h.getFabricatedPayload(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Header().Get("Content-Length") != "5" {
		t.Fatalf("Content-Length = %q", rec.Header().Get("Content-Length"))
	}
	if rec.Body.String() != "edgeb" {
		t.Fatalf("payload = %q", rec.Body.String())
	}
}

func TestGetFabricatedPayloadInvalidSize(t *testing.T) {
	h := New(controller.NewStore(), config.IntegrationConfig{})

	req := httptest.NewRequest(http.MethodGet, "/data/fabricate?size=bad", nil)
	rec := httptest.NewRecorder()
	h.getFabricatedPayload(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestGetFabricatedPayloadTooLarge(t *testing.T) {
	h := New(controller.NewStore(), config.IntegrationConfig{})

	req := httptest.NewRequest(http.MethodGet, "/data/fabricate?size=2g", nil)
	rec := httptest.NewRecorder()
	h.getFabricatedPayload(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestRegisterRoutesPing(t *testing.T) {
	store := controller.NewStore()
	h := New(store, config.IntegrationConfig{})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux, "")

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != "{\"status\":\"ok\"}" {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	store, _, _ := seedStore(t)
	h := New(store, config.IntegrationConfig{})

	req := httptest.NewRequest(http.MethodPost, "/metrics/cpu", nil)
	rec := httptest.NewRecorder()
	h.getCPUMetrics(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d", rec.Code)
	}
}
