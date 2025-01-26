package xorg

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/hertg/egpu-switcher/internal/logger"
)

//go:embed conf.template
var embeddedTemplate string
var templatePath = "/usr/share/egpu-switcher/x11-template.conf"

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

func RenderConf(id string, driver string, busid string, modesetting bool, verbose bool) (string, bool, error) {

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

	var confTemplate string

	templateFile, err := os.OpenFile(templatePath, os.O_RDONLY, 0644)
	isCustom := false
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			// if we get an error, other than "not exists", print an error
			// the "not exists" error is an expected outcome when the config file is not customized
			logger.Error("unable to open custom x11 config template, using default template instead")
		}
		if verbose {
			logger.Debug("using default template for x11 conf")
		}
		confTemplate = embeddedTemplate
	} else {
		if verbose {
			logger.Debug("using custom template at '%s' for x11 conf", templatePath)
		}
		var buf bytes.Buffer
		io.Copy(&buf, templateFile)
		confTemplate = buf.String()
		isCustom = true
	}

	t := template.Must(template.New("conf").Parse(confTemplate))
	err = t.Execute(buf, c)
	if err != nil {
		return "", isCustom, err
	}

	return buf.String(), isCustom, nil
}
