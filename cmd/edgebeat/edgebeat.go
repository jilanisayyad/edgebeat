package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jilanisayyad/edgebeat/pkg/config"
	"github.com/jilanisayyad/edgebeat/pkg/controller"
	"github.com/jilanisayyad/edgebeat/pkg/handler"
	"github.com/jilanisayyad/edgebeat/pkg/mqtt"
	"go.uber.org/zap"
)

const defaultConfigPath = "configs/config.yaml"

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	cfg, err := config.Load(defaultConfigPath)
	if err != nil {
		logger.Fatal("load config", zap.String("path", defaultConfigPath), zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := controller.NewStore()

	// Initialize MQTT publisher if enabled
	var publisher *mqtt.Publisher
	if cfg.MQTT.Enabled {
		mqttCfg := mqtt.Config{
			Broker:   cfg.MQTT.Broker,
			ClientID: cfg.MQTT.ClientID,
			Topic:    cfg.MQTT.Topic,
			Username: cfg.MQTT.Username,
			Password: cfg.MQTT.Password,
			QoS:      cfg.MQTT.QoS,
		}
		var err error
		publisher, err = mqtt.NewPublisher(ctx, mqttCfg, logger)
		if err != nil {
			logger.Fatal("mqtt initialization failed", zap.Error(err))
		}
		defer func() {
			if err := publisher.Close(); err != nil {
				logger.Error("mqtt close failed", zap.Error(err))
			}
		}()
	}

	go controller.Run(ctx, logger, time.Duration(cfg.FrequencySeconds)*time.Second, store, publisher)

	// Setup HTTP handlers
	mux := http.NewServeMux()
	h := handler.New(store, cfg.Integrations)
	h.RegisterRoutes(mux, "")

	server := &http.Server{
		Addr:         cfg.Rest.Address,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Info("server started",
			zap.String("address", cfg.Rest.Address),
			zap.Strings("endpoints", []string{
				"/health",
				"/metrics",
				"/metrics/cpu",
				"/metrics/memory",
				"/metrics/disk",
				"/metrics/network",
				"/metrics/system",
				"/metrics/sensors",
				"/integrations",
				"/data/fabricate",
				"/ping",
			}),
		)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown", zap.Error(err))
	}
}
