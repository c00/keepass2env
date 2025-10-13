# Keepass2nv

This little program extracts secrets from keepass and dumps them into a .env file.

Useful for setting up new environments and updating secrets.

When writing to file it will attempt to update old values rather than just appending. It leaves existing other values untouched.

## Install

```sh
go install github.com/c00/keepass2env/cmd/keepass2env@latest
```

## Configure

Copy `keepass2env.example.yaml` to `keepass2env.yaml`, and edit it to match your environment. Home directories will be expanded.

The most important thing is the `entries`. These are the mappings from the keepass database onto environment variables. e.g.

```yaml
entries: 
	- envName: DOCKER_TOKEN
    keepassPath: Personal/Docker Token
```

This will try to find an entry called `Docker Token` in the folder `Personal` in the database. These paths are case sensitive. It will then store the password into the output file as `DOCKER_TOKEN=thepassworditfound`.

## Notes on the password and keyfile

If you don't use a keyfile, just remove that from the config yaml.

For the password you can either use an environment variable, or enter the password interactively. if `passwordEnv` is set in the config and the variable is set in the environment, then it will use that password. If either `passwordEnv` is not set (or empty) or the environment variable itself is unset or empty, then you will still be asked for your password interactively.

## Examples

```sh
# minimal example with config set fully
keepass2env

# Using a different config file
keepass2env -c some/other/config/file.yaml

# Set or override values
keepass2env -k path/to/keyfile.key -d path/to/db.kdbx -o path/to/output.env
```