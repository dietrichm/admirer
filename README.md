# Admirer

[![Go version](https://img.shields.io/github/go-mod/go-version/dietrichm/admirer)](go.mod)
[![Go CI](https://github.com/dietrichm/admirer/actions/workflows/go.yml/badge.svg)](https://github.com/dietrichm/admirer/actions/workflows/go.yml)
[![License](https://img.shields.io/github/license/dietrichm/admirer)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/dietrichm/admirer.svg)](https://pkg.go.dev/github.com/dietrichm/admirer)

A command line utility to sync song likes (loved tracks) between Spotify and Last.fm.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Installation](#installation)
  - [Building from source](#building-from-source)
- [Usage](#usage)
  - [Supported services](#supported-services)
  - [Authentication](#authentication)
- [Use cases](#use-cases)
  - [Listing recently loved or added tracks](#listing-recently-loved-or-added-tracks)
  - [Syncing recently loved tracks between services](#syncing-recently-loved-tracks-between-services)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation

Compiled binaries will be provided for version 1.0.0.

**Note**: Admirer is currently only tested on Linux.

### Building from source

Please [install Go 1.16+ first](https://golang.org/doc/install) and set the `GOPATH` environment variable.

```sh
go get -u github.com/dietrichm/admirer
```

## Usage

```
Usage:
  admirer [command]

Available Commands:
  help        Help about any command
  list        List loved tracks on specified service
  login       Log in on external service
  status      Retrieve status for services
  sync        Sync recently loved tracks from one service to another

Flags:
  -h, --help   help for admirer

Use "admirer [command] --help" for more information about a command.
```

### Supported services

Last.fm and Spotify have an initial implementation.

In the future, I would like to add support for [ListenBrainz](https://listenbrainz.org/) and CSV/JSON files as well.

### Authentication

Before using any of the provided services, you need to create **your own API application** on said service and export your new API client ID and secret as environment variables:

| Service | Internal name | Creating your app | Environment variables |
| ------- | ------------- | ----------------- | --------------------- |
| Last.fm | `lastfm` | [Create an account here](https://www.last.fm/api/account/create) | `LASTFM_CLIENT_ID` and `LASTFM_CLIENT_SECRET` |
| Spotify | `spotify` | [Manage and create an app here](https://developer.spotify.com/dashboard/applications) | `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET` |

When this is done, continue with the following steps.

1. Run `$GOPATH/bin/admirer login <service>` to retrieve an authentication URL.
1. By visiting this URL, the service will ask confirmation and redirect back to a non existing URL `https://admirer.test/...`.
1. Copy the code parameter from this URL's query parameters and pass it along as another parameter to `$GOPATH/bin/admirer login <service>`.
1. If all goes well, you will retrieve confirmation that you have been logged in.

**Warning**: please be aware that - for now - the authentication information will be saved in **plain text** in a file in `~/.config/admirer`. This file's permissions is set to `600`, however.

**Note**: in future versions of Admirer, I will add an internal HTTP server to retrieve the authentication callback automatically.

## Use cases

### Listing recently loved or added tracks

Using the `list` command, you can retrieve a list of your most recently loved or added tracks on said service.

### Syncing recently loved tracks between services

Using the `sync` command, you can synchronise recently loved tracks from one service to another.
For example to mark as loved on Last.fm the same tracks that were added to your library on Spotify, or vice versa.

## License

Copyright 2020, Dietrich Moerman.

Released under the terms of the [MIT License](LICENSE).
