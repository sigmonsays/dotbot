package main

import (
	"bytes"
	"fmt"
	"os"
	"os/user"

	yaml "gopkg.in/yaml.v3"
)

// main configuration structure
type AppConfig struct {
	HomeDir string `yaml:"homedir"`
	Clean    []string          `yaml:"clean"`
	Mkdirs   []string          `yaml:"mkdirs"`
	Symlinks map[string]string `yaml:"symlinks"`
	Script   []*Script         `yaml:"script"`
	Include  []string          `yaml:"include"`
	WalkDir  map[string]string `yaml:"walkdir"`
}

func (c *AppConfig) LoadYaml(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	_, err = b.ReadFrom(f)
	if err != nil {
		return err
	}

	if err := c.LoadYamlBuffer(b.Bytes()); err != nil {
		return err
	}

	if err := c.FixupConfig(); err != nil {
		return err
	}

	return nil
}

func (c *AppConfig) LoadYamlBuffer(buf []byte) error {
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (me *AppConfig) PrintConfig() {
	d, err := yaml.Marshal(me)
	if err != nil {
		fmt.Println("Marshal error", err)
		return
	}
	fmt.Println("-- Configuration --")
	fmt.Println(string(d))
}

func GetDefaultConfig() *AppConfig {
	cfg := &AppConfig{}
	cfg.Symlinks = make(map[string]string, 0)

	// get home dir
	usr, _ := user.Current()
	homedir := usr.HomeDir
	log.Tracef("homedir is %s", homedir)
	cfg.HomeDir = homedir

	return cfg
}

func (c *AppConfig) LoadDefault() {
	*c = *GetDefaultConfig()
}

// after loading configuration this gives us a spot to "fix up" any configuration
// or abort the loading process
func (c *AppConfig) FixupConfig() error {
	// var emptyConfig AppConfig

	for i, s := range c.Script {
		s.SetDefaults()
		if s.Id == "" {
			s.Id = fmt.Sprintf("script%d", i)
		}
	}

	return nil
}

func PrintDefaultConfig() {
	conf := GetDefaultConfig()
	conf.PrintConfig()
}
