package config

import (
	tools "fortress/utils"
	"strings"
)

type Normal struct {
	Commands []string
}

func (n *Normal) CmdList(current *Fortress, hosts map[string]Hosts, args map[string]interface{}) (list map[string]string) {
	n.Commands = current.Commands

	list = make(map[string]string)

	command := strings.Join(n.Commands, "; ")
	command = tools.TplHandlerFromMap(command, args, "cli")
	list["CURRENT"] = command

	return
}
