package config

import (
	"fmt"
	tools "fortress/utils"
	"strings"
)

type Ssh struct {
	Hosts    []string
	Commands []string
}

func (s *Ssh) CmdList(current *Fortress, hosts map[string]Hosts, args map[string]interface{}) (list map[string]string) {
	s.Hosts = current.Hosts
	s.Commands = current.Commands

	list = make(map[string]string)
	for _, machine := range s.Hosts {
		host := hosts[machine]

		command := strings.Join(s.Commands, "; ")
		command = fmt.Sprintf(`ssh -t -p %d %s@%s "%s"`, host.Port, host.User, host.IP, command)

		command = tools.TplHandlerFromMap(command, args, "cli")
		command = tools.TplHandlerFromStructure(command, host, "host")

		if host.Password != "" {
			command = fmt.Sprintf("sshpass -p %s %s", host.Password, command)
		}

		list[machine] = command
	}
	return
}
