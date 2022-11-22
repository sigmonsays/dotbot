
![build status](https://github.com/sigmonsays/dotbot/actions/workflows/release.yml/badge.svg)

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [dotbot](#dotbot)
- [TLDR](#tldr)
- [Configuration](#configuration)
- [Usage](#usage)

<!-- markdown-toc end -->

# dotbot

simple tool to manage symlinks in $HOME for aiding in keeping dot files in sync on multiple
machines

# TLDR

Create a config file 

     cat << EOF > dotbot.yaml
     symlinks:
       ~/.bashrc: .bashrc
     EOF

Run dotbot to make links

     dotbot link

If no configuration file is given, the default is assumed to be
     
     dotbot -c dotbot.yaml link
     
# Configuration

dotbot -c is used to provide what configuration to run. Multiple configuration
files can be specified by passing -c with a file multiple times.

If no configuration file is provided, dotbot.yaml is assumed if it exists in the
current directory.

Sample configuration

     mkdirs:
       - ~/asdf
     symlinks:
         ~/.bash_profile: .bash_profile
         ~/.bashrc: .bashrc


# Usage

     NAME:
        dotbot - manage dot files
     
     USAGE:
        dotbot [global options] command [command options] [arguments...]
     
     COMMANDS:
        link, l    create symlinks
        status, s  print status as json
        cleanup    show unreferenced files
        help, h    Shows a list of commands or help for one command
     
     GLOBAL OPTIONS:
        --loglevel value, -l value                             set log level (default: "info")
        --config value, -c value [ --config value, -c value ]  config file
        --help, -h                                             show help (default: false)
     
