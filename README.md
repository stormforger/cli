[ ![Travis CI Status](https://travis-ci.org/stormforger/cli.svg?branch=master)](https://travis-ci.org/stormforger/cli)
[ ![Go Report Card](https://goreportcard.com/badge/github.com/stormforger/cli)](https://goreportcard.com/report/github.com/stormforger/cli)



# forge! The StormForger Command Line Client

**HEAVY WORK IN PROGRESS**

```
$ forge --help

The command line client "forge" to StormForger offers a interface
to the StormForger API and several convenience methods
to handle load and performance tests.

Happy Load Testing :)

Usage:
  forge [command]

Available Commands:
  datasource  Work with and manage data sources
  har         Convert HAR to test case
  login       Login to StormForger
  ping        Ping the StormForger API
  test-run    Work with and manage test runs
  version     Show forge version

Flags:
      --endpoint string   API Endpoint (default "https://api.stormforger.com")
      --jwt string        JWT access token

Use "forge [command] --help" for more information about a command.
```



## Installation

Download the latest release from [GitHub](https://github.com/stormforger/cli/releases).



## Getting Started

Most actions require authentication. So in case you don't have a StormForger account yet, you have to [sign up](https://app.stormforger.com) first.

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


### Dependencies

Build dependencies:

```
go get -u github.com/tools/godep
go get -u golang.org/x/tools/cmd/goimports
go get -u github.com/golang/lint/golint
```

We use [`godep`](https://github.com/tools/godep) to vendor dependencies.

To add a new dependency, use `go get` to install it, use it (import it) and use
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

In order to publish releases to GitHub, you need a personal access token, which you can acquire here: https://github.com/settings/tokens.

Copy the token an add it to your local goxc configuration (will be written to .goxc.local.json):

```
goxc -wlc default publish-github -apikey=$API_KEY
```

Now you can make a release with

```
make test local_release
```
