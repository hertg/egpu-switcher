package service

import "fmt"

type Systemd struct {
}

func (s Systemd) CreateService() error {
	// todo: create service file (dest: /etc/systemd/system/egpu.service)
	// todo: systemctl daemon-reload
	// todo: enable egpu.service
	return fmt.Errorf("not implemented")
}

func (s Systemd) TeardownService() error {
	return fmt.Errorf("not implemented")
}

func createService() {
	// todo
}

func daemonReload() {
	// todo
}

func enable() {

}
