[ ![Travis CI Status](https://travis-ci.org/stormforger/cli.svg?branch=master)](https://travis-ci.org/stormforger/cli)
[ ![Go Report Card](https://goreportcard.com/badge/github.com/stormforger/cli)](https://goreportcard.com/report/github.com/stormforger/cli)


# forge! The StormForger Command Line Client

Please note that this tool is still **HEAVY WORK IN PROGRESS**. If you have any questions, don't hesitate to [contact](https://stormforger.com/support).

Using `forge` you can:

* use StormForger's HAR converter
* work with test runs (only fetching logs currently)
* work with data sources (private beta currently)

You can use `--help` to get usage information on all commands.


## Installation

Download the latest release from [GitHub releases](https://github.com/stormforger/cli/releases) page.

In case you are on macOS and using [Homebrew](https://brew.sh/) you can:

```
brew tap stormforger/forge
brew install forge
```


## Getting Started

Most actions require authentication. So in case you don't have a StormForger account yet, you have to [sign up](https://app.stormforger.com) first - no worries, it's free!

When done, you can login via

```
forge login your-email@example.com
```

You will be asked for your credentials. On successful authentication, you will be presented with a JWT.

For following requests you have multiple options to provide your token:

1. TOML configuration: `.stormforger.toml` or `$HOME/.stormforger.toml`:

```
jwt = "your-jwt-token"
```

2. Environment: `export STORMFORGER_JWT="your-jwt-token"`

3. Command Line Flag: `--jwt "your-jwt-token"`



## Build

In case you want to use `forge`, you can stop reading now. This section describes how to build `forge` from source.

### Dependencies

Build dependencies:

```
go get -u github.com/tools/godep
go get -u golang.org/x/tools/cmd/goimports
go get -u github.com/golang/lint/golint
```

We use [`godep`](https://github.com/tools/godep) to vendor dependencies.

To add a new dependency, use `go get` to install it, use it (as in import it) and use
`godep save` to add the dependency to the `vendor` directory. **Make sure you
add dependencies in a dedicated commit!**.

### Release

Releases are done via Travis CI.

When ready, tag a new release and push the new tag

```
git tag vA.B.C
git push --tag
```

Travis will make a build and on success automatically publish a release to [GitHub releases](https://github.com/stormforger/cli/releases).

Now go to the [releases page](https://github.com/stormforger/cli/releases) and add release notes.


#### Local Release

In case there is an issue with the normal release process, a manual (or local) release can be done as well.

Releases are done with `goreleaser`:

```
go get -u github.com/goreleaser/goreleaser
```

In order to publish releases to GitHub, you need a personal access token, which you can acquire here: https://github.com/settings/tokens.

Now you can make a release with

```
GITHUB_TOKEN="geheim" make test local_release
```
