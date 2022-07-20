package service

type InitSystem interface {
	CreateService() error
	TeardownService() error
}
