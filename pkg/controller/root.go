package controller

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, payload []byte) error
}

func Run(ctx context.Context, logger *zap.Logger, frequency time.Duration, store *Store, publisher Publisher) {
	if logger == nil {
		logger = zap.NewNop()
	}

	collectAndPublish(ctx, logger, store, publisher)

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("shutting down", zap.String("reason", ctx.Err().Error()))
			return
		case <-ticker.C:
			collectAndPublish(ctx, logger, store, publisher)
		}
	}
}

func collectAndPublish(ctx context.Context, logger *zap.Logger, store *Store, publisher Publisher) {
	info := collectSystemInfo()
	payload, err := json.Marshal(info)
	if err != nil {
		logger.Error("marshal system info", zap.Error(err))
		return
	}

	if store != nil {
		store.Set(payload)
	}

	if publisher != nil {
		if err := publisher.Publish(ctx, payload); err != nil {
			logger.Error("mqtt publish failed", zap.Error(err))
		}
	}

	if len(info.Errors) > 0 {
		logger.Warn("system info collected with errors", zap.Int("error_count", len(info.Errors)))
		return
	}

	logger.Info("system info collected")
}
