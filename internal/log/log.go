package log

import (
	"context"
	"log/slog"
	"os"
	"simple-telemetry-publisher/internal/model"
	"time"
)

type LogProvider struct {
	Config model.LogConfig
}

func (l LogProvider)Start(ctx context.Context) error {

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler

	if l.Config.JsonFormat {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)

	for k, v := range l.Config.ExtraFields {
		logger = logger.With(k, v)
	}

	ticker := time.NewTicker(l.Config.Interval)

	go func() {
		loop:
			for {
				select {
				case <-ticker.C:
					logger.Debug("debug message")
					logger.Info("info message")
					logger.Warn("warn message")
					logger.Error("error message")
				case <-ctx.Done():
					logger.Info("Stopping log")
					ticker.Stop()
					break loop
				}
			}
	}()

	return nil
}
