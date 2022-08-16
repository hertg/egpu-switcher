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
	// todo: systemctl show -p LoadState egpu.service | sed 's/LoadState=//g') == "loaded" (only attempt stop if service is loaded)
	// todo: systemctl stop egpu.service
	// todo: systemctl disable egpu.service
	// todo: delete service file
	// todo: systemctl daemon-reload
	// todo: reset-failed
	return fmt.Errorf("not implemented")
}

func createService() {
	// todo
}

func daemonReload() {
	// todo
}

func enable() {
	// todo
}
