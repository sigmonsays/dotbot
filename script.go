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
	DefaultScriptShell = "bash"
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
	ctx := context.Background()
	stdin := bytes.NewBufferString(me.Command)

	key := fmt.Sprintf("%s/%s", me.Type, me.Id)
	log.Tracef("running script %s (id:%s type:%s quiet:%v shell:%s)",
		key, me.Id, me.Type, me.Quiet, me.Shell)

	if log.IsTrace() {
		log.Tracef("command --- begin ---\n%s\n--- end ---\n", me.Command)
	}

	cmdline := []string{
		"/usr/bin/env",
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

	if c.ProcessState != nil {
		ret.Pid = c.ProcessState.Pid()
	}
	if ret.Pid == 0 {
		log.Warnf("%s script returned but no pid", key)
	}

	log.Tracef("%s script ran, exit code %d", key, ret.ExitCode)
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
