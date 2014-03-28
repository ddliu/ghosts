# GHosts

Make hosts file configurable

## Features

- Switch hosts profiles
- Load remote hosts file
- Easy to config

## Usage

```
sudo ghosts --config=/path/to/ghosts.yml prod
```

## ghosts.yml

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

- Locate config file automatically
- Flush DNS cache