package xorg

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"

	"github.com/hertg/egpu-switcher/internal/logger"
)

//go:embed conf.template
var confTemplate string

func RemoveEgpuFile(path string, verbose bool) error {
	f, _ := os.Stat(path)
	if f != nil {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("unable to remove file %s", path)
		}
	}
	if verbose {
		logger.Debug("deleted '%s'", path)
	}
	logger.Info("egpu has been removed from X.Org config")
	return nil
}

func CreateEgpuFile(path string, contents string, verbose bool) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create file %s", path)
	}
	_, err = f.Write([]byte(contents))
	if err != nil {
		return fmt.Errorf("unable to write config to file %s", path)
	}
	if verbose {
		logger.Debug("written '%s'", path)
	}
	logger.Info("egpu has been added to X.Org config")
	return nil
}

func RenderConf(id string, driver string, busid string, modesetting bool) string {

	type conf struct {
		Id          string
		Driver      string
		Bus         string
		Modesetting bool
	}

	c := conf{
		Id:          id,
		Driver:      driver,
		Bus:         busid,
		Modesetting: modesetting,
	}

	buf := bytes.NewBuffer(nil)
	t := template.Must(template.New("conf").Parse(confTemplate))
	err := t.Execute(buf, c)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
