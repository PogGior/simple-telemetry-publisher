# Simple Telemetry Publisher

Simple Telemetry Publisher is a simple tool that allows you to publish basic logs, metrics and traces on a distributed system.

## Use case

The main use case for simple telemetry publisher is to test a central observability system in a distributed environment. 

The tool can be used to publish logs, metrics and traces to one or more telemetry collectors and verify that the data is correctly received and processed.

## How it works

The tool is a simple command line application that can be used to publish logs, metrics and traces. 

The tool can be configured to publish data to one or more telemetry collectors using different protocols and formats.

To see all the available options, see the [configuration section](#configuration).

## Installation

The plugin can be builded using the following command:

```bash
make compile
```

And runned locally with:

```bash
make run
```

To run the plugin in a container, you can use the following command:

```bash
make docker-image
```

And saved the image in a tar file with:

```bash
make docker-save
```

Examples of how to run the plugin on a container in a complex environment can be found [here](deploy).

## Configuration

The tool can be configured by providing a configuration file.:

```bash
./simple-telemetry-publisher --config ./fixtures/simple-publisher-config.yaml
```

The configuration file is a YAML file that contains the following fields:

```yaml
log:
  disable: false
  interval: 15s
  json-format: false
  extra-fields:
    foo: bar
metric:
  disable: false
  prometheus:
    port: 9004
    interval: 10s
trace:
  disable: false
  interval: 10s
  endpoint: "0.0.0.0:4318"
  service-name: simple-publisher
  tracer-name: jaeger
graceful-shutdown: 5s
```

- `log`: Configuration for log publishing.
  - `disable`: Disables log publishing.
  - `interval`: Interval between log messages.
  - `json-format`: Enables JSON format for log messages.
  - `extra-fields`: Extra fields to be added to log messages.
- `metric`: Configuration for metric publishing.
    - `disable`: Disables metric publishing.
    - `prometheus`: Configuration for Prometheus metrics.
        - `port`: Port where the Prometheus metrics will be exposed.
        - `interval`: Interval between metric updates.
- `trace`: Configuration for trace publishing
    - `disable`: Disables trace publishing
    - `interval`: Interval between trace spans.
    - `endpoint`: Endpoint where the trace spans will be sent.
    - `service-name`: Service name for the trace spans.
    - `tracer-name`: Tracer name for the trace spans.
- `graceful-shutdown`: Graceful shutdown timeout.



