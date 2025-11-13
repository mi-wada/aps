# aps - AWS Profile Switcher

aps is a TUI-based tool for simply switching AWS profiles.

## Install

1. Install binary:

```shell
go install github.com/mi-wada/aps@latest
```

2. Add alias to your shell:

```shell
alias aps='eval $(command aps)'
```

## Usage

```shell
aps
```

## Development

For development, use:

```shell
eval $(go run .)
```

Or add this alias to your shell:

```shell
alias aps-dev='eval $(go run /Users/wada.mitsuaki/github.com/mi-wada/aps)'
```
