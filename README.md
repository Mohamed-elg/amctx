# amctx: Context manager for amtool

![Written in go](https://img.shields.io/badge/written%20in-go-00ADD8.svg)

## Summary

**amctx** is a cli tool to switch between `Alertmanager` instances for `amtool`.
It works by editing the `~/.config/amtool/config.yml` file.
The aliases are stored in `~/.config/amctx/config.yml`

## Installation

```sh
go install github.com/Mohamed-elg/amctx@latest
```

## Examples

```sh
# switch to another Alertmanager that's in
$ amctx prod
switched to context 'prod'.

# add an alias
$ amctx prod=http://localhost:9093
alias 'prod' set with url: http://localhost:9093

# list aliases
$ amctx
prod
dev
```
