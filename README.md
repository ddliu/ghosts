# GHosts

[![Build Status](https://travis-ci.org/ddliu/ghosts.png)](https://travis-ci.org/ddliu/ghosts)

Efficient hosts switcher

## Features

- Switch hosts profiles
- Load remote hosts file
- Easy to config

## Supported systems

- Unix(link)
- Windows

Tested on Ubuntu linux.

## Usage

```
ghosts prod
ghosts base remote
```

You need to run this command as root(sudo) or administrator.

### Options

- --config=/custom/path/to/ghosts.yml
- --list Show groups
- --raw Keep comments, spaces and empty lines
- --target=/custom/path/to/hosts
- --print Do not write hosts file, just show result


## ghosts.yml

The `ghosts.yml` file contains groups of hosts definitions. 

By default it located at:

- Unix(Mac OS X, Linux): ~/.ghosts/ghosts.yml
- Windows: HOME\ghosts\ghosts.yml

Config example:

```yml
# GHost Sample Config File

# Line by line
group1:
 - 192.168.0.100 www.example.com
 - 192.168.0.101 static.example.com

# Line by line
group2:
 - 192.168.0.102 db.example.com
 - 192.168.0.103 admin.example.com

# Multi-line scope
group3: |
 192.168.0.104 help.example.com
 192.168.0.105 mail.example.com

# Combination
base:
 - 127.0.0.1 localhost
 - group1
 - group2

# Load a remote host file
remote: https://raw.githubusercontent.com/ddliu/hosts/master/ghosts.example.hosts

# Load a local host file
local: ./hosts

# Include groups
test:
 - base
 - local

prod:
 - base
 - remote
```

## Todo

- Flush DNS cache
- Prebuilt packages/binaries for different systems

## Changelog

### v0.1.0 (2014-03-30)

First release