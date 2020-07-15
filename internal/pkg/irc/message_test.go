// Package irc is a library for processing user input from IRC.
//
// This file contains tests for message processing.
package irc

import (
	"reflect"
	"testing"

	irc "github.com/thoj/go-ircevent"
)

func TestTryProcessCommand(t *testing.T) {
	tests := []struct {
		name   string
		config *CommandConfig
		event  *irc.Event
		want   *Command
	}{
		{
			name:   "empty IRC event should not return a command",
			config: &CommandConfig{},
			event:  &irc.Event{},
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.TryProcessCommand(tt.event)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TryProcessCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
