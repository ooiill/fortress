package main

import (
	"fmt"
	"fortress/config"
	tools "fortress/utils"
	"strings"
)

func main() {

	args := tools.ParseArgs(map[string]string{
		"--index": "Missions index, multiple split by comma.",
		"--mode": "Mode for run missions.",
		"-debug":  "Without execute shell.",
	})

	// read config
	data := new(config.Cnf).Bootstrap()

	// default index and menu
	numbersStr,  _ := args["index"]
	numbers, err := tools.Menu(data.ListMenu(), true, numbersStr)
	tools.Handler(err)

	if numbersStr == nil || numbersStr == "" {
		fmt.Printf("\n\t%s\n", strings.Repeat("⬇ ", 5))
	}

	// default mode and menu
	modesStr, _ := args["mode"]
	listMode := []string{"one by one (sync)", "scramble for (async)"}
	modes, err := tools.Menu(listMode, false, modesStr)
	tools.Handler(err)

	var async bool
	if modes[0] == 2 {
		async = true
	}

	// run
	log := make(chan string)
	for _, number := range numbers {
		current := data.Fortress[number-1]

		if async {
			go current.Work(data.Hosts, async, args, log)
		} else {
			current.Work(data.Hosts, async, args, log)
		}
	}

	// async
	if async {
		line := strings.Repeat("-", 10)
		fmt.Printf("\n%s%s 异步执行 ^ 无限等待 %s\n\n", tools.Indent, line, line)
		for item := range log {
			fmt.Print(item)
		}
	}
}
