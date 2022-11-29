package main

import "strings"

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
	return nil
}

func (me *Script) SetDefaults() {
	if me.Type == "" {
		me.Type = DefaultScriptType
	}
	me.Type = strings.ToLower(me.Type)
}

func (me *Script) Run() error {
	log.Tracef("running script %s", me.Id)
	return nil
}
