package main

import (
	"fmt"
	"strings"
)

var (
	DefaultScriptType = "post"
)

type Script struct {
	Id       string
	Command  string
	Disabled bool

	// post or pre; default is 'post'
	Type string
}

func (me *Script) Validate() error {
	if me.Command == "" {
		return fmt.Errorf("command is required")
	}
	return nil
}

func (me *Script) SetDefaults() {
	if me.Type == "" {
		me.Type = DefaultScriptType
	}
	me.Type = strings.ToLower(me.Type)
}

func (me *Script) Run() (*ScriptResult, error) {
	ret := &ScriptResult{}
	log.Tracef("running script %s", me.Id)

	return ret, nil
}

type ScriptResult struct {
	ExitCode int
}

func (me *ScriptResult) String() string {
	return fmt.Sprintf("exitcode:%d", me.ExitCode)
}
