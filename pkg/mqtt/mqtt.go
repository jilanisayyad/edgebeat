package mqtt

import (
	"context"
	"fmt"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Publisher struct {
	client pahomqtt.Client
	topic  string
	qos    byte
	logger *zap.Logger
}

type Config struct {
	Broker   string
	ClientID string
	Topic    string
	Username string
	Password string
	QoS      byte
}

func NewPublisher(ctx context.Context, cfg Config, logger *zap.Logger) (*Publisher, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	opts := pahomqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(60 * time.Second)

	opts.SetOnConnectHandler(func(client pahomqtt.Client) {
		logger.Info("mqtt connected", zap.String("broker", cfg.Broker))
	})

	opts.SetConnectionLostHandler(func(client pahomqtt.Client, err error) {
		logger.Warn("mqtt connection lost", zap.Error(err))
	})

	client := pahomqtt.NewClient(opts)
	token := client.Connect()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled during mqtt connection")
	case <-token.Done():
		if token.Error() != nil {
			return nil, fmt.Errorf("mqtt connect: %w", token.Error())
		}
	}

	return &Publisher{
		client: client,
		topic:  cfg.Topic,
		qos:    cfg.QoS,
		logger: logger,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, payload []byte) error {
	if p == nil || p.client == nil {
		return fmt.Errorf("publisher not initialized")
	}

	if !p.client.IsConnected() {
		return fmt.Errorf("mqtt client not connected")
	}

	token := p.client.Publish(p.topic, p.qos, false, payload)

	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled during publish")
	case <-token.Done():
		if token.Error() != nil {
			return fmt.Errorf("mqtt publish: %w", token.Error())
		}
	}

	p.logger.Debug("mqtt published", zap.String("topic", p.topic), zap.Int("bytes", len(payload)))
	return nil
}

func (p *Publisher) Close() error {
	if p == nil || p.client == nil {
		return nil
	}

	p.client.Disconnect(1000)
	p.logger.Info("mqtt disconnected")
	return nil
}
