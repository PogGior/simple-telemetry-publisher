package model

import "context"

type TelemetryProvider interface {
	Start(ctx context.Context) error
}