package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	DefaultScriptType  = "post"
	DefaultScriptShell = "/bin/bash"
)

type Script struct {
	Id       string
	Command  string
	Disabled bool

	// do not print commands stdout and stderr to logs
	Quiet bool

	// default shell
	Shell string

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
	if me.Shell == "" {
		me.Shell = DefaultScriptShell
	}
}

func (me *Script) Run() (*ScriptResult, error) {
	ret := &ScriptResult{}
	log.Tracef("running script %s", me.Id)
	ctx := context.Background()
	stdin := bytes.NewBufferString(me.Command)

	cmdline := []string{
		me.Shell,
		"-",
	}
	c := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	c.Stdin = stdin

	if me.Quiet == false {
		c.Stdout = os.Stdout
		c.Stderr = os.Stdout
	}
	c.Env = os.Environ()

	err := c.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			ret.Pid = exiterr.Pid()
			ret.ExitCode = exiterr.ExitCode()
		}
	}

	log.Tracef("%s script %s ran, exit code %d",
		me.Type, me.Id, ret.ExitCode)
	return ret, nil
}

type ScriptResult struct {
	Pid      int
	ExitCode int
}

func (me *ScriptResult) String() string {
	return fmt.Sprintf("pid:%d exitcode:%d",
		me.Pid, me.ExitCode)
}
