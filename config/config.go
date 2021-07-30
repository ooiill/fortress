package config

import (
	"encoding/json"
	tools "fortress/utils"
	"io/ioutil"
	"strings"
)

type Cnf struct {
	Hosts    map[string]Hosts
	Fortress []Fortress
}

func (c *Cnf) Bootstrap() (data Cnf) {

	cnfFile, err := tools.ExecShell("echo ~/.fortress.json")
	tools.Handler(err)

	cnfFile = strings.TrimRight(cnfFile, "\n")

	jsonStr, err := ioutil.ReadFile(cnfFile)
	tools.Handler(err)

	_ = json.Unmarshal([]byte(jsonStr), &data)
	return
}

func (c *Cnf) ListMenu() (fortressSlice []string) {
	for _, item := range c.Fortress {
		fortressSlice = append(fortressSlice, item.Mission)
	}
	return
}
