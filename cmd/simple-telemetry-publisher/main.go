package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"

	//"simple-telemetry-publisher/internal/otel"
	//"simple-telemetry-publisher/internal/prometheus"
	"simple-telemetry-publisher/internal/model"
	"simple-telemetry-publisher/internal/startup"
)

func main() {

	pflag.String("config", "./simple-publisher-config.yaml", "config file (default is ./simple-publisher-config.yaml)")
	pflag.Parse()
	config, err := model.LoadConfig(pflag.CommandLine)
    if err != nil {
        panic(fmt.Errorf("error loading config: %s", err))
    }

    fmt.Printf("config: %+v\n", config)

    ctx, cancel := context.WithCancel(context.Background())
    
    telemetryProviders, err := startup.Init(config)
    if err != nil {
        panic(fmt.Errorf("error starting: %s", err))
    }


	for _,telemetryProvider := range telemetryProviders {
		err = telemetryProvider.Start(ctx)
		if err != nil {
			panic(fmt.Errorf("error starting: %s", err))
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancel()
}
