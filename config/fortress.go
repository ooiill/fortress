package config

import (
    "fmt"
    tools "fortress/utils"
    "strings"
)

type Fortress struct {
    Mission  string
    Type     string
    Hosts    []string
    Commands []string
    Normal
    Ssh
    Sync
}

func (f *Fortress) Work(hosts map[string]Hosts, async bool, args map[string]interface{}, log chan string) {

    var cmdList map[string]string
    switch f.Type {
    case "ssh":
        cmdList = new(Ssh).CmdList(f, hosts, args)
    case "sync":
        cmdList = new(Sync).CmdList(f, hosts, args)
    case "normal":
    default:
        cmdList = new(Normal).CmdList(f, hosts, args)
    }

    for machine, command := range cmdList {
        if async {
            f.execute(machine, command, args, log)
        } else {
            f.execute(machine, command, args)
        }
    }
}

func (f *Fortress) execute(machine string, command string, args map[string]interface{}, log ...chan string) {

    line := strings.Repeat("-", 10)
    blue := tools.BlueBold.SprintFunc()
    cyan := tools.Cyan.SprintFunc()

    var logger chan string
    var async bool

    if len(log) > 0 {
        async = true
        logger = log[0]
    }

    handler := func(msg string) {
        if async {
            // logger <- msg
        } else {
            fmt.Print(msg)
        }
    }

    handler(tools.Cyan.Sprintf("\n%s%s ^^^ %s\n", tools.Indent, line, line))
    handler(fmt.Sprintf("\n%s%s : %s", tools.Indent, blue("Machine"), cyan(machine)))
    handler(fmt.Sprintf("\n%s%s : %s\n", tools.Indent, blue("Command"), cyan(command)))

    if _, ok := args["debug"]; !ok {
        handler(tools.Yellow.Sprintf("\n%s%s ··· %s\n\n", tools.Indent, line, line))
        if async {
            out, err := tools.ExecShell(command)
            tools.Handler(err)
            if logger != nil {
                logger <- out
            }
        } else {
            err := tools.ExecShellFine(command)
            tools.Handler(err)
        }
    }

    handler(tools.Cyan.Sprintf("\n%s%s $$$ %s\n\n", tools.Indent, line, line))
}
