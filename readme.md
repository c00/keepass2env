# Keepass2env / Keepass2keyring

This little program extracts secrets from keepass and dumps them into a .env file or the system keyring.

Useful for setting up new environments and updating secrets.

## Install

If you have cloned the repo, just run `make install`. Otherwise, do go install:

```sh
# Keepass2env
go install github.com/c00/keepass2env/cmd/keepass2env@latest

# Keepass2keyring
go install github.com/c00/keepass2env/cmd/Keepass2keyring@latest
```

## Configure

Copy the contents of `keepass2env.example.yaml` to `~/.config/keepass2env.yaml`, and edit it to match your environment. 

The most important thing is the `entries`. These are the mappings from the keepass database onto environment variables. e.g.

```yaml
entries: 
  - envName: DOCKER_TOKEN
    keepassPath: Personal/Docker Token
```

This will try to find an entry called `Docker Token` in the folder `Personal` in the database. These paths are case sensitive.

### Attributes

By default it will extract the password from an entry, but in some cases you may want to get a different bit of data. For this you can use the `attribute` property. This can point to custom attributes in a Keepass entry:

```yaml
entries: 
  - envName: MY_PRIVATE_KEY
    keepassPath: Personal/Cert for important things
    attribute: private.key
  - envName: MY_PUBLIC_KEY
    keepassPath: Personal/Cert for important things
    attribute: public.key
```

## Notes on the password and keyfile

If you don't use a keyfile, just remove that from the config yaml.

For the password you can either use an environment variable, or enter the password interactively. if `passwordEnv` is set in the config and the variable is set in the environment, then it will use that password. If either `passwordEnv` is not set (or empty) or the environment variable itself is unset or empty, then you will still be asked for your password interactively.

## keepass2env - Writing to a file

When writing to file it will attempt to update old values rather than just appending. It leaves existing other values untouched.

### Examples for keepass2env

```sh
# minimal example with config set fully
keepass2env

# Using a different config file
keepass2env -c some/other/config/file.yaml

# Set or override values
keepass2env -k path/to/keyfile.key -d path/to/db.kdbx -o path/to/output.env
```

### Sourcing the .env file

Add something like this to your `.bashrc` or `.profile` to automatically load the secrest into your environment.

```sh
if [ -f ~/.secrets.env ]; then
  export $(grep -v '^#' ~/.secrets.env | xargs)
fi
```

## keepass2keyring - Writing to the system keyring

This should be compatible with any keyring solution. It will write to the default keyring. It uses the very excellent [go-keyring](https://github.com/zalando/go-keyring) library.

It will create entries in the keyring with the following attributes:

- `service`: `keepass2keyring` (this value can be set in the config)
- `application`: `[entry-env-name]` (e.g. `DOCKER_TOKEN`)

You can get the values written using the `secret-tool`:

```sh
secret-tool search --all service keepass2keyring
```
### Examples for keepass2keyring

```sh
# minimal example with config set fully
keepass2keyring

# Using a different config file
keepass2keyring -c some/other/config/file.yaml

# Set or override values
keepass2keyring -k path/to/keyfile.key -d path/to/db.kdbx -s my-service
```