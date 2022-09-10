package service

import (
	"bytes"
	"fmt"
	"os"
)

var systemdBytes = []byte("/lib/systemd")

func GetInitSystem() (InitSystem, error) {
	f, err := os.ReadFile("/sbin/init")
	if err != nil {
		return nil, err
	}
	if bytes.Contains(f, systemdBytes) {
		return &Systemd{}, nil
	}
	return nil, fmt.Errorf("unable to guess init system")
}
