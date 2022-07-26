package service

type InitSystem interface {
	CreateService() error
	TeardownService() error
	StopDisplayManager() error
	StartDisplayManager() error
	IsDisplayManagerStopped() (bool, error)
}
