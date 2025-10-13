# Keepass2nv

This little program extracts secrets from keepass and dumps them into a .env file.

Useful for setting up new environments and updating secrets.

When writing to file it will attempt to update old values rather than just appending. It leaves existing other values untouched.

## Install

```sh
go install github.com/c00/keepass2env/cmd/keepass2env@latest
```

## Examples

```sh
keepass2env 
```