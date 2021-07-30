package config

import (
	"fmt"
	tools "fortress/utils"
)

type Sync struct {
	Host  string
	Hosts []string
	From  string
	To    string
	Args  string
}

func (s *Sync) CmdList(current *Fortress, hosts map[string]Hosts, args map[string]interface{}) (list map[string]string) {
	s.Host = current.Host
	s.Hosts = current.Hosts
	s.From = current.From
	s.To = current.To
	s.Args = current.Args

	list = make(map[string]string)
	for _, machine := range s.Hosts {
		host := hosts[machine]

		command := fmt.Sprintf("rsync -avzh -e 'ssh -p %d' --delete --progress %s", host.Port, s.Args)

		if s.Host == "from" {
			command = fmt.Sprintf("%s %s@%s:%s %s", command, host.User, host.IP, s.From, s.To)
		} else {
			command = fmt.Sprintf("%s %s %s@%s:%s", command, s.From, host.User, host.IP, s.To)
		}

		command = tools.TplHandlerFromMap(command, args, "cli")
		command = tools.TplHandlerFromStructure(command, host, "host")

		if host.Password != "" {
			command = fmt.Sprintf("sshpass -p %s %s", host.Password, command)
		}

		list[machine] = command
	}
	return
}
