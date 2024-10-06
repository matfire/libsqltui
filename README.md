# LibSQLTui
> graphical tool for managing a self-hosted libsql server

![GitHub Release](https://img.shields.io/github/v/release/matfire/libsqltui?style=for-the-badge)

## Description

This tool can connect to a (for now) locally running sqld service and perform operations such as creating, deleting and forking databases

## Limitation (Temporary, hopefully)

LibSQLTUI makes a lot of assumptions for now:

- the user API http endpoint should listen on 127.0.0.1:8080 (or 0.0.0.0:8080 if you're feeling adventurous)
- the admin API http endpoint should listen on 127.0.0.1:8081 (or 0.0.0.0:8081 if you're feeling adventurous)

## Installation

You can find the latest release on the rekease section if you prefer downloading binaries.

You can also install it using homebrew by running:

```shell
    brew install matfire/matfire/libsqltui
```

Then you can run it using `libsqltui`
