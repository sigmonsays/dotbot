

[![Build](https://github.com/sigmonsays/dotbot/actions/workflows/release.yml/badge.svg)](https://github.com/sigmonsays/dotbot/actions/workflows/release.yml)

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [dotbot](#dotbot)
- [TLDR](#tldr)
- [Configuration](#configuration)
- [Scripts](#scripts)
- [WalkDirs](#walkdirs)
- [Includes](#includes)
- [AutoMode](#automode)
- [Usage](#usage)

<!-- markdown-toc end -->

# dotbot

simple tool to manage symlinks in $HOME for aiding in keeping dot files in sync on multiple
machines

Inspired by the python version at https://github.com/anishathalye/dotbot

# TLDR

Create a config file 

     cat << EOF > dotbot.yaml
     symlinks:
       ~/.bashrc: .bashrc
     EOF

Run dotbot to make links

     dotbot link

If no configuration file is given, the default is assumed to be dotbot.yaml
     
     dotbot -c dotbot.yaml link

To see what needs to be done (pretend mode)

     dotbot -c dotbot.yaml link -p

To remove links
     
     dotbot unlink

To see status 

     dotbot status    # Table output
     dotbot status -j # As json

# Configuration

A configuration file is optional if automode is used. 

dotbot -c is used to provide what configuration to run. Multiple configuration
files can be specified by passing -c with a file multiple times.

If no configuration file is provided, dotbot.yaml is assumed if it exists in the
current directory.

Additional configuration files can be provided as positional arguments. This makes for a convenient invocation:

     dotbot link stuff.yaml

the clean block indicates directories to clean up broken symlinks (aka, dangling symlinks)
the tilde (`~`) is automatically expanded to the uses home directory. The path is evaluated
as a glob so wildcards may be used.

Sample configuration

     clean:
       - '~'
     mkdirs:
       - ~/asdf
     symlinks:
         ~/.bash_profile: .bash_profile
         ~/.bashrc: .bashrc
     walkdirs:
         ~/bin: bin

# Scripts

A script is a series of shell commands that allow you to run commands
before or after creating symlinks.

Example

     script:
      - id: example1
        type: pre
        command: |
          date
      - id: example2
        command: |
          chmod 0400 ~/.ssh/config
          chmod 0400 ~/.ssh/config.d/*
      - id: example3
        shell: /bin/bash
        quiet: true
        command: |
          set -x
          date > /tmp/test.txt

For the above examples, example1 is a 'pre' script which runs before
symlinks and example2 is a 'post' script which runs after symlinks.
By default, when type is not provied, 'post' will be used.

# WalkDirs

WalkDirs are used to list a directory and apply symlinks for every
item (path or file) found. This is useful if symlinks need to come 
from multiple directories and need to be merged into one.

# Includes

Includes allow additional dotfiles to be processed as if they were
being invoked directly. Include files are processed in order given, 
after expanding the glob pattern.

Includes are processed after evaluating the contents of the 
configuration file.

Example

     include:
       - stuff/foo.yaml
       - extra/bar*.yaml

# AutoMode

Auto mode is enabled by passing `--auto` or `-a` to `dotbot link`

With automode, no configuration file is required, instead a configuration is automatically generated at runtime. The current directory is used as a file list. 

The `.git` file name is ignored.

automode works with pretend mode

     dotbot link -a -p

# Usage

     NAME:
        dotbot - manage dot files
     
     USAGE:
        dotbot [global options] command [command options] [arguments...]
     
     COMMANDS:
        link, l    create symlinks
        unlink, u  remove symlinks
        status, s  print status table
        cleanup    show unreferenced files
        help, h    Shows a list of commands or help for one command
     
     GLOBAL OPTIONS:
        --loglevel value, -l value                             set log level (default: "info")
        --config value, -c value [ --config value, -c value ]  config file
        --help, -h                                             show help (default: false)
     
