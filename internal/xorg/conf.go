package xorg

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/hertg/egpu-switcher/internal/logger"
)

func RemoveEgpuFile(path string, verbose bool) error {
	f, _ := os.Stat(path)
	if f != nil {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("unable to remove file %s", path)
		}
		if verbose {
			logger.Debug("the file %s has been removed", path)
		}
		return nil
	}
	if verbose {
		logger.Debug("the file %s is already absent", path)
	}
	return nil
}

func CreateEgpuFile(path string, contents string, verbose bool) error {
	_, err := os.Stat(path)
	if err != nil {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("unable to create file %s", path)
		}
		_, err = f.Write([]byte(contents))
		if err != nil {
			return fmt.Errorf("unable to write config to file %s", path)
		}
		if verbose {
			logger.Debug("the file %s has been created", path)
		}
		return nil
	}
	if verbose {
		logger.Debug("the file %s already exists", path)
	}
	return nil
}

func RenderConf(id string, driver string, busid string) string {

	const confTemplate = `
Section "Module"
    Load           "modesetting"
EndSection

Section "Device"
    Identifier     "{{.Id}}"
    Driver         "{{.Driver}}"
    BusID          "{{.Bus}}"
    Option         "AllowEmptyInitialConfiguration"
    Option         "AllowExternalGpus" "True"
EndSection
`

	type conf struct {
		Id     string
		Driver string
		Bus    string
	}

	c := conf{
		Id:     id,
		Driver: driver,
		Bus:    busid,
	}

	buf := bytes.NewBuffer(nil)
	t := template.Must(template.New("conf").Parse(confTemplate))
	err := t.Execute(buf, c)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
