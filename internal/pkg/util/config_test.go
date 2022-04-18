package util

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
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

log_level: 'info'
`

var noPluginIRCConfig = &Config{
	IRC: &ircConfig{
		Channels:         []string{"##Metal"},
		CommandTrigger:   "!",
		Debug:            false,
		Ident:            "Metal",
		MaxReconnect:     5,
		ReconnectDelay:   time.Duration(600 * time.Second),
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

		Hostname: "irc.somewhere.org:6697",
	},
	Plugins:          make(map[string]*pluginConfig),
	UnparsedLogLevel: "info",
	LogLevel:         logrus.InfoLevel,
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
				t.Errorf("NewConfig() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
