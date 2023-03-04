package service

import "context"

type InitSystem interface {
	CreateService(ctx context.Context, verbose bool) error
	TeardownService(ctx context.Context, verbose bool) error
	StopUnit(ctx context.Context, unit string, verbose bool) error
	StartUnit(ctx context.Context, unit string, verbose bool) error
	StopDisplayManager(ctx context.Context, verbose bool) error
	StartDisplayManager(ctx context.Context, verbose bool) error
	IsDisplayManagerStopped(ctx context.Context, verbose bool) (bool, error)
}
