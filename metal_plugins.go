package main

import (
	"fmt"
	"log"

	"github.com/n7st/metal/internal/pkg/util"
	"github.com/n7st/metal/pkg/command"
)

func main() {
	plugins, errors := util.LoadPlugins([]string{
		"/home/mike/go/src/github.com/n7st/Metal/plugins/greeter.so",
		"/home/mike/go/src/github.com/n7st/Metal/plugins/ping.so",
	})

	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
	}

	command := &command.Command{Command: "ping", Username: "Someone"}

	for _, p := range plugins {
		output := p.Run(command)

		for _, o := range output {
			fmt.Printf("%v", o)
		}
	}

	for _, p := range plugins {
		output := p.Run(command)

		for _, o := range output {
			fmt.Printf("%v", o)
		}
	}
}
