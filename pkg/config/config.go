package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFrequencySeconds = 60
	MinFrequencySeconds     = 1
	MaxFrequencySeconds     = 180
	DefaultRestAddress      = ":8080"
	DefaultRestPath         = "/health"
	DefaultMQTTQoS          = 1
	DefaultModbusMode       = "tcp"
	DefaultModbusPort       = 502
	DefaultModbusUnitID     = 1
	DefaultOpcuaEndpoint    = "opc.tcp://localhost:4840"
	DefaultOpcuaPolicy      = "None"
	DefaultOpcuaMode        = "None"
)

type Config struct {
	FrequencySeconds int               `yaml:"frequency_seconds"`
	Rest             RestConfig        `yaml:"rest"`
	MQTT             MQTTConfig        `yaml:"mqtt"`
	Integrations     IntegrationConfig `yaml:"integrations"`
}

type RestConfig struct {
	Address string `yaml:"address"`
	Path    string `yaml:"path"`
}

type MQTTConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Broker   string `yaml:"broker"`
	ClientID string `yaml:"client_id"`
	Topic    string `yaml:"topic"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	QoS      byte   `yaml:"qos"`
}

type IntegrationConfig struct {
	Modbus ModbusConfig `yaml:"modbus"`
	OPCUA  OPCUAConfig  `yaml:"opcua"`
}

type ModbusConfig struct {
	Enabled bool   `yaml:"enabled"`
	Mode    string `yaml:"mode"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	UnitID  int    `yaml:"unit_id"`
	Notes   string `yaml:"notes"`
}

type OPCUAConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Endpoint       string `yaml:"endpoint"`
	SecurityPolicy string `yaml:"security_policy"`
	SecurityMode   string `yaml:"security_mode"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	Notes          string `yaml:"notes"`
}

func Default() Config {
	return Config{
		FrequencySeconds: DefaultFrequencySeconds,
		Rest: RestConfig{
			Address: DefaultRestAddress,
			Path:    DefaultRestPath,
		},
		MQTT: MQTTConfig{
			Enabled: false,
			QoS:     DefaultMQTTQoS,
		},
		Integrations: IntegrationConfig{
			Modbus: ModbusConfig{
				Enabled: false,
				Mode:    DefaultModbusMode,
				Host:    "localhost",
				Port:    DefaultModbusPort,
				UnitID:  DefaultModbusUnitID,
			},
			OPCUA: OPCUAConfig{
				Enabled:        false,
				Endpoint:       DefaultOpcuaEndpoint,
				SecurityPolicy: DefaultOpcuaPolicy,
				SecurityMode:   DefaultOpcuaMode,
			},
		},
	}
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	cfg := Default()

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if cfg.FrequencySeconds < MinFrequencySeconds || cfg.FrequencySeconds > MaxFrequencySeconds {
		return Config{}, fmt.Errorf("frequency_seconds out of range: %d", cfg.FrequencySeconds)
	}

	if cfg.Rest.Address == "" {
		cfg.Rest.Address = DefaultRestAddress
	}
	if cfg.Rest.Path == "" {
		cfg.Rest.Path = DefaultRestPath
	}

	if cfg.Integrations.Modbus.Mode == "" {
		cfg.Integrations.Modbus.Mode = DefaultModbusMode
	}
	if cfg.Integrations.Modbus.Host == "" {
		cfg.Integrations.Modbus.Host = "localhost"
	}
	if cfg.Integrations.Modbus.Port == 0 {
		cfg.Integrations.Modbus.Port = DefaultModbusPort
	}
	if cfg.Integrations.Modbus.UnitID == 0 {
		cfg.Integrations.Modbus.UnitID = DefaultModbusUnitID
	}
	if cfg.Integrations.OPCUA.Endpoint == "" {
		cfg.Integrations.OPCUA.Endpoint = DefaultOpcuaEndpoint
	}
	if cfg.Integrations.OPCUA.SecurityPolicy == "" {
		cfg.Integrations.OPCUA.SecurityPolicy = DefaultOpcuaPolicy
	}
	if cfg.Integrations.OPCUA.SecurityMode == "" {
		cfg.Integrations.OPCUA.SecurityMode = DefaultOpcuaMode
	}

	return cfg, nil
}
