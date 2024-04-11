package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

type Status struct {
	ctx *Context
}

func (me *Status) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "output JSON",
		},
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "show all file status",
		},
	}
}

func (me *Status) Run(c *cli.Context) error {
	configfiles := me.ctx.getConfigFiles(c)
	asJson := c.Bool("json")
	showAll := c.Bool("all")

	if len(configfiles) == 0 {
		log.Warnf("Nothing to do, try passing -c dotbot.yaml ")
		return nil
	}

	for _, filename := range configfiles {
		err := me.RunFile(filename, asJson, showAll)
		if err != nil {
			log.Warnf("RunFile %s: %s", filename, err)
		}
	}

	return nil
}

func (me *Status) RunFile(path string, asJson, showAll bool) error {
	log.Tracef("runfile %s", path)
	cfg := GetDefaultConfig()
	err := cfg.LoadYaml(path)
	if err != nil {
		return err
	}
	if log.IsTrace() {
		cfg.PrintConfig()
	}
	p := NewRunParamsConfig(cfg)
	run, err := CompileRun(path, p)
	if err != nil {
		return err
	}

	if asJson {
		buf, err := json.Marshal(run)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", buf)
		return nil
	}

	// table output

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{
		"Target",
		"Type",
		"Status",
	})
	var row []string

	for _, li := range run.Links {
		row = []string{}
		row = append(row, li.Target)

		if li.NeedsCreate == false && showAll == false {
			continue
		}

		row = append(row, li.FileType)

		status := "unknown"
		if li.DestExists {
			status = "created"
		} else if li.NeedsCreate {
			status = "missing"
		}
		row = append(row, status)

		table.Append(row)
	}
	table.Render()
	return nil

}
