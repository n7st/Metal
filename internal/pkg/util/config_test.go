package util

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

var noPluginIRCYaml = `---
irc:
  ident: 'Metal'
  modes: '+b'
  nickname: 'Metal'
  nickserv_account: 'myaccount'
  nickserv_password: 'mypassword'
  port: 6697
  server: 'irc.somewhere.org'
  use_tls: true
  max_reconnect: 5

  channels:
    - '##Metal'

log_level: 'debug'
`

var noPluginIRCConfig = &Config{
	IRC: &ircConfig{
		Channels:         []string{"##Metal"},
		Debug:            false,
		Ident:            "Metal",
		MaxReconnect:     5,
		ReconnectDelay:   10,
		Modes:            "+b",
		Nickname:         "Metal",
		NickservAccount:  "myaccount",
		NickservPassword: "mypassword",
		Port:             6697,
		RealName:         "Metal",
		Server:           "irc.somewhere.org",
		ServerPassword:   "",
		UseTLS:           true,
		Verbose:          false,

		Hostname:              "irc.somewhere.org 6697",
		ReconnectDelayMinutes: time.Duration(10),
	},
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *Config
	}{
		{
			name:  "no plugins",
			input: noPluginIRCYaml,
			want:  noPluginIRCConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "/tmp/metalconfig.yml"
			data := []byte(tt.input)
			ioutil.WriteFile(filename, data, 0777)

			if got := NewConfig(filename); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("%v", got)
				t.Errorf("NewConfig() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
