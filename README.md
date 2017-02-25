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

Use "forge [command] --help" for more information about a command.
```


## Release

Dependencies:

```
go get github.com/laher/goxc
go get github.com/inconshreveable/mousetrap # required to cross compile for windows
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
