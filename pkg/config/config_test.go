package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestLoadDefaults(t *testing.T) {
	path := writeTempConfig(t, "frequency_seconds: 60\nrest:\n  address: ''\n  path: ''\nmqtt:\n  enabled: false\n")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.FrequencySeconds != DefaultFrequencySeconds {
		t.Fatalf("FrequencySeconds = %d, want %d", cfg.FrequencySeconds, DefaultFrequencySeconds)
	}
	if cfg.Rest.Address != DefaultRestAddress {
		t.Fatalf("Rest.Address = %q, want %q", cfg.Rest.Address, DefaultRestAddress)
	}
	if cfg.Rest.Path != DefaultRestPath {
		t.Fatalf("Rest.Path = %q, want %q", cfg.Rest.Path, DefaultRestPath)
	}
	if cfg.Integrations.Modbus.Mode != DefaultModbusMode {
		t.Fatalf("Modbus.Mode = %q, want %q", cfg.Integrations.Modbus.Mode, DefaultModbusMode)
	}
	if cfg.Integrations.Modbus.Port != DefaultModbusPort {
		t.Fatalf("Modbus.Port = %d, want %d", cfg.Integrations.Modbus.Port, DefaultModbusPort)
	}
	if cfg.Integrations.Modbus.UnitID != DefaultModbusUnitID {
		t.Fatalf("Modbus.UnitID = %d, want %d", cfg.Integrations.Modbus.UnitID, DefaultModbusUnitID)
	}
	if cfg.Integrations.OPCUA.Endpoint != DefaultOpcuaEndpoint {
		t.Fatalf("OPCUA.Endpoint = %q, want %q", cfg.Integrations.OPCUA.Endpoint, DefaultOpcuaEndpoint)
	}
	if cfg.Integrations.OPCUA.SecurityPolicy != DefaultOpcuaPolicy {
		t.Fatalf("OPCUA.SecurityPolicy = %q, want %q", cfg.Integrations.OPCUA.SecurityPolicy, DefaultOpcuaPolicy)
	}
	if cfg.Integrations.OPCUA.SecurityMode != DefaultOpcuaMode {
		t.Fatalf("OPCUA.SecurityMode = %q, want %q", cfg.Integrations.OPCUA.SecurityMode, DefaultOpcuaMode)
	}
}

func TestLoadOverrides(t *testing.T) {
	path := writeTempConfig(t, "frequency_seconds: 10\nrest:\n  address: ':9090'\n  path: '/metrics'\nmqtt:\n  enabled: true\n  broker: 'tcp://localhost:1883'\n  client_id: 'edgebeat-test'\n  topic: 'edgebeat/metrics'\n  qos: 2\nintegrations:\n  modbus:\n    enabled: true\n    mode: 'rtu'\n    host: '127.0.0.1'\n    port: 1502\n    unit_id: 2\n  opcua:\n    enabled: true\n    endpoint: 'opc.tcp://example:4840'\n    security_policy: 'Basic256'\n    security_mode: 'Sign'\n    username: 'user'\n    password: 'pass'\n")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.FrequencySeconds != 10 {
		t.Fatalf("FrequencySeconds = %d, want 10", cfg.FrequencySeconds)
	}
	if cfg.Rest.Address != ":9090" || cfg.Rest.Path != "/metrics" {
		t.Fatalf("Rest = %+v, want address :9090 and path /metrics", cfg.Rest)
	}
	if !cfg.MQTT.Enabled || cfg.MQTT.QoS != 2 {
		t.Fatalf("MQTT = %+v, want enabled and qos 2", cfg.MQTT)
	}
	if cfg.Integrations.Modbus.Mode != "rtu" || cfg.Integrations.Modbus.Port != 1502 {
		t.Fatalf("Modbus = %+v", cfg.Integrations.Modbus)
	}
	if cfg.Integrations.OPCUA.Endpoint != "opc.tcp://example:4840" {
		t.Fatalf("OPCUA.Endpoint = %q", cfg.Integrations.OPCUA.Endpoint)
	}
	if cfg.Integrations.OPCUA.SecurityMode != "Sign" {
		t.Fatalf("OPCUA.SecurityMode = %q", cfg.Integrations.OPCUA.SecurityMode)
	}
}

func TestLoadInvalidFrequency(t *testing.T) {
	path := writeTempConfig(t, "frequency_seconds: 500\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid frequency")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/path/does/not/exist.yaml")
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}
