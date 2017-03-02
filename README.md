# forge! The StormForger Command Line Client

**HEAVY WORK IN PROGRESS**

```
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



## Release

Dependencies:

```
go get github.com/fatih/color
go get github.com/howeyc/gopass
go get github.com/inconshreveable/mousetrap # required to cross compile for windows
go get github.com/laher/goxc
go get golang.org/x/tools/cmd/goimports
```

In order to publish releases, you need a personal access token, which you can acquire here: https://github.com/settings/tokens.

Copy the token an add it to your local goxc configuration (will be written to .goxc.local.json):

```
goxc -wlc default publish-github -apikey=$API_KEY
```

To make a release, run:

```
make release
```
