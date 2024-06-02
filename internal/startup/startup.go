package startup

import (
	"simple-telemetry-publisher/internal/log"
	"simple-telemetry-publisher/internal/metric"
	"simple-telemetry-publisher/internal/model"
	"simple-telemetry-publisher/internal/trace"
)



func Init(config model.Config) ([]model.TelemetryProvider,error) {

	var telemetryProviders []model.TelemetryProvider
	
	if !config.LogConfig.Disable {
		telemetryProviders = append(telemetryProviders, log.LogProvider{Config: config.LogConfig})
	}
	if !config.MetricConfig.Disable {
		if !config.MetricConfig.PrometheusConfig.Disable {
			telemetryProviders = append(telemetryProviders, metric.PrometheusMetrics{GracefulShutdown: config.GracefulShutdown, Config: config.MetricConfig.PrometheusConfig})
		}
		if !config.MetricConfig.OtelMetricConfig.Disable {
			telemetryProviders = append(telemetryProviders, metric.OtelMetrics{GracefulShutdown: config.GracefulShutdown, Config: config.MetricConfig.OtelMetricConfig})
		}
	}
	if !config.TraceConfig.Disable {
		telemetryProviders = append(telemetryProviders, trace.OtelTrace{GracefulShutdown: config.GracefulShutdown, Config: config.TraceConfig})
	}

	return telemetryProviders, nil
}