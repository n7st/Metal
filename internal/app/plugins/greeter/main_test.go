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
	plugin.CheckRunCommand(t, NewPlugin(), "hello_world", &command.Command{})
}

func TestParse(t *testing.T) {
	t.Run("Message contains 'hi'", func(t *testing.T) {
		input := &command.Command{Message: "test hi test"}

		response := plugin.CheckRunParse(t, NewPlugin(), input)

		if reflect.DeepEqual(response, &command.Response{Message: "hi"}) {
			t.Fail()
		}
	})

	t.Run("Message doesn't contain 'hi'", func(t *testing.T) {
		input := &command.Command{Message: "test test"}

		response := plugin.CheckRunParse(t, NewPlugin(), input)

		if response != nil {
			t.Fail()
		}
	})
}

func TestHelloWorld(t *testing.T) {
	greeter := NewPlugin()

	input := &command.Command{Username: "Mike", Message: "hello_world"}

	t.Run("Not recently greeted", func(t *testing.T) {
		expected := &command.Response{Message: "Hello, world"}
		got := plugin.CheckRunCommand(t, greeter, "hello_world", input)

		if !reflect.DeepEqual(got, expected) {
			t.Fail()
		}
	})

	t.Run("Recently greeted", func(t *testing.T) {
		expected := &command.Response{Message: "You were greeted recently!"}
		got := plugin.CheckRunCommand(t, greeter, "hello_world", input)

		if !reflect.DeepEqual(got, expected) {
			t.Fail()
		}
	})

	t.Run("Previously greeted, but not recently", func(t *testing.T) {
		// This test requires a replacement for time ("clock") in order to pass
		// a specific time around.
		t.Skip("TODO")
	})
}
