package main

import (
	"reflect"
	"testing"

	"github.com/n7st/metal/pkg/command"
	"github.com/n7st/metal/pkg/plugin"
)

func TestNewPlugin(t *testing.T) {
	plugin.CheckImplementsInterface(t, NewPlugin())
}

func TestCommands(t *testing.T) {
	plugin.CheckRunCommand(t, NewPlugin(), "ping", &command.Command{})
}

func TestParse(t *testing.T) {
	pinger := NewPlugin()

	if pinger.Parse(&command.Command{}) != nil {
		t.Fail()
	}
}

func TestPingCommand(t *testing.T) {
	expected := &command.Response{Message: "Pong!"}

	got := plugin.CheckRunCommand(t, NewPlugin(), "ping", &command.Command{})

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("pingCommand() got: %v, want: %v", got, expected)
	}
}
