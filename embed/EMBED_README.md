# Embedded configuration

This directory along with readme is mandatory if no embedded configuration is used, since embedding empty directory seems to be impossible. Configuration placed in embed directory has priority over all other configuration files.

## Embedding configuration

Peeker2 reads cfg.yaml file placed in embed directory. It has to be placed in embed directory, along with all other required data before compiliation (obviously). You can also place ssh private keys inside embed directory. In order to use embedded ssh keys, use `embed/<keynamehere>` path in cfg.yaml
