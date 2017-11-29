[ ![Travis CI Status](https://travis-ci.org/stormforger/cli.svg?branch=master)](https://travis-ci.org/stormforger/cli)
[ ![Go Report Card](https://goreportcard.com/badge/github.com/stormforger/cli)](https://goreportcard.com/report/github.com/stormforger/cli)


# forge! The StormForger Command Line Client

Please note that this tool is still **HEAVY WORK IN PROGRESS**. If you have any questions, don't hesitate to get in [contact](https://stormforger.com/support).

Using `forge` you can:

* use StormForger's HAR converter
* work with `organisation`s (only list currently)
* work with `test-run`s (start, abort, reporting, listing and call logs)
* work with `test-case`s (list, create, update)

You can use `--help` to get usage information on all commands.

For more information on how to usage StormForger in general, visit our [documentation](https://docs.stormforger.com).

## Installation

Download the latest release from [GitHub releases](https://github.com/stormforger/cli/releases) page.

In case you are on macOS and using [Homebrew](https://brew.sh/) you can:

```
brew tap stormforger/forge
brew install forge
```

You can also use our published Docker image [`stormforger/cli`](https://hub.docker.com/r/stormforger/cli). We will publish the `latest` tag, so you can do:

```
docker pull stormforger/cli
docker run stormforger/cli
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

1. Environment: `export STORMFORGER_JWT="your-jwt-token"`

1. Flag: `--jwt "your-jwt-token"`



## Build

In case you want to use `forge`, you can stop reading now. This section describes how to build `forge` from source.

### Dependencies

Use `make setup` to install required go build dependencies.

We use [`dep`](https://github.com/golang/dep) to vendor dependencies.

To add a dependency, simply import it and then run `make dep` which will add the dependency to the manifest. **Make sure you add dependencies in a dedicated commit!**.

### Release

Releases are done via [Travis CI](https://travis-ci.org/stormforger/cli).

When ready for a release and pull requests are merged into master, just create and push a new tag:

```
git tag vA.B.C
git push --tag
```

Travis will make a build and on success automatically publish a release to [GitHub releases](https://github.com/stormforger/cli/releases), to [Docker hub](https://hub.docker.com/r/stormforger/cli) and also update our [homebrew tab](https://github.com/stormforger/homebrew-forge).

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
