// Package util contains common functionality for "utilities" required by the
// bot.
//
// This file contains code for loading and running Lua modules.
package util

import (
	"fmt"

	"github.com/Shopify/go-lua"
	luaUtil "github.com/Shopify/goluago/util"
)

// ModuleLoadError contains debugging information for a module which failed to
// load.
type ModuleLoadError struct {
	Module *Module
	Error  error
}

// Module contains information about a bot module.
type Module struct {
	Filename string
	State    *lua.State
	Command  bool
	Message  bool
	Commands map[string]string
	Counter  int
}

// ToString converts a ModuleLoadError into a string for output.
func (e *ModuleLoadError) ToString() string {
	return fmt.Sprintf(`Module "%s" failed to load: %s`, e.Module.Filename, e.Error.Error())
}

// Next gets the next method call number.
func (m *Module) Next() int {
	m.Counter = m.Counter + 1

	return m.Counter
}

func (m *Module) loadOptions() (bool, error) {
	m.State.Global("Options")
	m.State.Call(0, 1)

	if m.State.IsTable(m.Next()) {
		table, err := luaUtil.PullTable(m.State, 1)

		if err != nil {
			return false, err
		}

		options := table.(map[string]interface{})

		if v, ok := options["command"].(bool); ok {
			m.Command = v
		}

		if v, ok := options["message"].(bool); ok {
			m.Message = v
		}

		if v, ok := options["commands"].(map[string]interface{}); ok {
			for command, method := range v {
				m.Commands[command] = fmt.Sprintf("%s", method)
			}
		}
	}

	return true, nil
}

// LoadModules loads Lua scripts into the VM.
func LoadModules(filenames []string) ([]*Module, []*ModuleLoadError) {
	var (
		modules []*Module
		errors  []*ModuleLoadError
	)

	for _, filename := range filenames {
		state := lua.NewState()
		module := &Module{
			Filename: filename,
			State:    state,
			Commands: make(map[string]string),
		}

		lua.OpenLibraries(state)

		if err := lua.DoFile(state, filename); err != nil {
			errors = append(errors, &ModuleLoadError{
				Error:  err,
				Module: module,
			})

			continue
		}

		module.loadOptions()

		modules = append(modules, module)
	}

	return modules, errors
}
