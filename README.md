# LibSQLTui
> graphical tool for managing a self-hosted libsql server

![GitHub Release](https://img.shields.io/github/v/release/matfire/libsqltui?style=for-the-badge)

## Description

This tool can connect to a (for now) locally running sqld service and perform operations such as creating, deleting and forking databases

>[!IMPORTANT]
>this application requires your sqld server to run with the following flags enabled:
>--enable-namespaces
>--admin-listen-addr


## Configuration

add a configuration file either in same directory as the binary, or in `$HOME/.config/libsqltui`. The config file can be in either JSON, TOML or YAML and should be called `config.{format}` where `{format}` is one of the formats mentionded above.

The currently supported configuration values are:
- CLIENT_ENDPOINT: the http client endpoint exposed by the sqld binary (defaults to http://127.0.0.1:8080)
- ADMIN_ENDPOINT: the http admin endpoint exposed by the sqld binary (requires more flags passed to the sqld binary) (defaults to http://127.0.0.1:8081)

## Installation

You can find the latest release on the rekease section if you prefer downloading binaries.

You can also install it using homebrew by running:

```shell
    brew install matfire/matfire/libsqltui
```

Then you can run it using `libsqltui`

## Developing

In the repo I provided a flake.nix file that can create a devShell for x86 linux; to use it simply run `nix develop`
