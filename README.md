# aps - AWS Profile Switcher

aps is a TUI-based tool for simply switching AWS profiles.

## Usage

```console
$ aps
AWS Profile Switcher
> default [current]
  production
  staging

(Use ↑/↓ or j/k to navigate, Enter to select, q to quit)
```

## Install

1. Install binary:

```shell
go install github.com/mi-wada/aps@latest
```

2. Add this function to your shell:

```shell
aps() {
  eval $(command aps)
}
```

## Development

For development, use:

```shell
eval $(go run .)
```
