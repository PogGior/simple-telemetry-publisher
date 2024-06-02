package model

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	LogConfig        LogConfig     `mapstructure:"log"`
	MetricConfig     MetricConfig  `mapstructure:"metric"`
	TraceConfig      TraceConfig   `mapstructure:"trace"`
	GracefulShutdown time.Duration `mapstructure:"graceful-shutdown"`
}

type LogConfig struct {
	Disable     bool              `mapstructure:"disable"`
	AddSource   bool              `mapstructure:"add-source"`
	JsonFormat  bool              `mapstructure:"json-format"`
	Interval    time.Duration     `mapstructure:"interval"`
	ExtraFields map[string]string `mapstructure:"extra-fields"`
}

type MetricConfig struct {
	Disable          bool             `mapstructure:"disable"`
	PrometheusConfig PrometheusConfig `mapstructure:"prometheus"`
	OtelMetricConfig OtelMetricConfig `mapstructure:"otel"`
}

type TraceConfig struct {
	Disable     bool          `mapstructure:"disable"`
	Endpoint    string        `mapstructure:"endpoint"`
	ServiceName string        `mapstructure:"service-name"`
	TracerName  string        `mapstructure:"tracer-name"`
	Interval    time.Duration `mapstructure:"interval"`
}

type PrometheusConfig struct {
	Disable  bool          `mapstructure:"disable"`
	Port     uint16        `mapstructure:"port"`
	Interval time.Duration `mapstructure:"interval"`
}

type OtelMetricConfig struct {
	Disable           bool          `mapstructure:"disable"`
	Endpoint          string        `mapstructure:"endpoint"`
	ServiceName       string        `mapstructure:"service-name"`
	MeterProviderName string        `mapstructure:"meter-provider-name"`
	Interval          time.Duration `mapstructure:"interval"`
}

func LoadConfig(flagSet *pflag.FlagSet) (Config, error) {

	var config Config

	viper.BindPFlags(flagSet)
	viper.SetConfigFile(viper.GetString("config"))

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	viper.SetDefault("trace.service-name", "simple-telemetry-publisher")
	viper.SetDefault("trace.tracer-name", "simple-telemetry-publisher")
	viper.SetDefault("trace.endpoint", "0.0.0.0:4318")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, errors.WithStack(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, errors.WithStack(err)
	}

	return config, nil
}
