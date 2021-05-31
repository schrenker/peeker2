# peeker2

*sometimes, you just need to take a peek...*

## What

Configurable, self-contained, no dependency monitoring tool, that provides simple critical info about hosts of your choice.

## Why

I felt need for a very simple alternative to big blown-out monitoring tools, that require agents, their own servers dedicated servers, databases, and other technical overhead. This tool simply requires yaml configuration and ssh access to hosts.

## How

Based on yaml configuration provided, this tool makes parallel ssh calls to hosts, pulling necessary information and then displays it on terminal screen.

Program looks for configuration file `cfg.yaml` in following places, and applies first one found:

1. embed directory (for embedded configs)
2. command line argument
3. cfg.yaml in current directory

In order to use embedded configuration, the program has to be precompiled with configuration applied, along with required ssh keys. Every embedded piece of configuration has to be placed in embed directory.

### Example config

Below is example cfg.yaml config file:

```yaml
---
interval: 60
hosts:
  - hostname: server1.example.com
    port: 25025
    user: root
    key: ./user/.ssh/id_rsa
    services:
      - httpd
      - mysqld
      - php-fpm
    disks:
      - /
      - home
      - /var/log
```

Alternatively, you can embed config into programs binary, by putting cfg.yaml and all required keys to embed directory, and then compiling the program. In order to use embedded ssh keys, simply point them to embed path, like below:

```yaml
---
interval: 60
hosts:
  - hostname: embedded.example.com
    port: 22
    user: ssh-user
    key: embed/id_rsa
    services:
      - httpd
    disks:
      - /
```

## What now?

This tool is still under development. Some configuration methods may change.

### TBD

* [ ] Colored output based on host parameter status
* [ ] Sorting
* [ ] Specify config file as program argument
