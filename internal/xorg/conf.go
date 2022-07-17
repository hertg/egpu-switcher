package xorg

import (
	"bytes"
	"html/template"
)

func GenerateConf(id string, driver string, busid string) string {

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
