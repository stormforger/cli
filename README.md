# forge! The StormForger Command Line Client

**HEAVY WORK IN PROGRESS**

```
This is the StormForger command line client!

Usage:
  forge [command]

Available Commands:
  har         Convert HAR to test case
  ping        Ping the StormForger API

Flags:
      --api-endpoint string   API Endpoint (default "https://api.stormforger.com")
      --jwt-token string      JWT Token

Use "forge [command] --help" for more information about a command.
```


## Authentication
To use the StormForger CLI and access the API you need to authenticate with your personal JWT token and give it to the StormForger CLI.


### Get your personal JWT token
Please write to support@stormforger.com to request your token.


## Configuration
There are three ways to give your JWT token to the CLI. The source of the JWT token is prioritized, means the last option overrides the first option.


### 1. Configuration File

Copy the `.stormforger.toml` configuration file either to the root folder of the cli binary or in your `$HOME` directory and fill it with your JWT token

`cp .stormforger.toml.example $HOME/.stormforger.toml`


### 2. Environment Variable

Set the environment variable `STORMFORGER_JWT` to your JWT token to use the StormForger CLI in [Twelve-Factor-App](https://12factor.net/) setups like your CI/CD and Build Pipeline.

`export STORMFORGER_JWT="your-jwt-token"`


### 3. Command Line Flag
Run the StormForger CLI with the

`--jwt-token "your-jwt-token"`

flag.


## Release

Dependencies:

```
go get github.com/fatih/color
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
