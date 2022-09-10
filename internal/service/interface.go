package service

import "context"

type InitSystem interface {
	CreateService(context.Context) error
	TeardownService(context.Context) error
	StopUnit(context.Context, string) error
	StartUnit(context.Context, string) error
	StopDisplayManager(context.Context) error
	StartDisplayManager(context.Context) error
	IsDisplayManagerStopped(context.Context) (bool, error)
}
