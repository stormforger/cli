[ ![Travis CI Status](https://travis-ci.org/stormforger/cli.svg?branch=master)](https://travis-ci.org/stormforger/cli)
[ ![Go Report Card](https://goreportcard.com/badge/github.com/stormforger/cli)](https://goreportcard.com/report/github.com/stormforger/cli)


# forge! The StormForger Command Line Client

Please note that this tool is still **HEAVY WORK IN PROGRESS**. If you have any questions, don't hesitate to get in [contact](https://stormforger.com/support).

Using `forge` you can:

* use StormForger's HAR converter
* work with `organisation`s (only list currently)
* work with `test-run`s (start, abort, reporting, listing and call logs)
* work with `test-case`s (list, create, update)
* work with `data-source`s (list, download, push, rename, show)

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

You will be asked for your credentials. On successful authentication your token will be written to `~/.stormforger.toml` (use `--no-save` to print the token instead).

Beside via `.stormforger.toml`, you can provide your JWT via

1. Environment: `export STORMFORGER_JWT="your-jwt-token"`
1. Flag: `--jwt "your-jwt-token"`

When you are done, you can check your token via `ping` which makes an authenticated ping request:

```
forge ping
```


## Usage

Help is available to all commands with `--help`.


### Test Cases

Test cases are scoped by organisation. You always have to provide at least the organisation (`acme-inc`) or a test case in form of `organisation-name/test-case-name`.

You can...
* list existing test cases in your `acme-inc` organisation: `forge test-case list acme-inc`
* validate a test definition without saving it: `forge test-case validate acme-inc cases/simple.js`
* create a new test case named `checkout` inside `acme-inc`: `forge test-case create acme-inc/checkout cases/simple.js`
* update an existing test case named `checkout` inside `acme-inc`: `forge test-case update acme-inc/checkout cases/simple.js`

Commands that take a file, also accept `-` to take input from stdin. Example:

```
sed 's/${target}/testapp.loadtest.party/g' tests/templated.js | forge tc update acme-inc/checkout -
```


### Test Runs

Test runs are executions of test cases. The subcommand is `test-run` or `tr`.

You can...

* launch `acme-inc/checkout`: `forge test-case launch acme-inc/checkout`
* watch a running test run: `forge test-run watch acme-inc/checkout/42`
* list all test runs of a test case: `forge test-run list acme-inc/checkout`
* show details: `forge test-run show acme-inc/checkout/42`
* view logs: `forge test-run logs acme-inc/checkout/42`
* view full traffic dump: `forge test-run dump acme-inc/checkout/42`
* check your requirements: `forge test-run nfr acme-inc/checkout/42 requirements/basic.yml`


### Data Sources

Data sources are scoped per organisation. Working with data sources can be done with the `datasource` (or short `ds`) sub command, e.g. `forge datasource ls acme-inc`.

You can...

* list available data sources: `forge datasource ls acme-inc`
* show details: `forge datasoure show acme-inc auth/users.csv`
* download originally uploaded file: `forge datasource get acme-inc auth/users.csv`
* rename data source: `forge datasource mv acme-inc users.csv auth/users.csv`
* create or update data sources: `forge datasource push acme-inc …`

`push` takes some more arguments: `forge datasource push acme-inc <list of files> [flags]`. It can be used to update or create new data sources (think update and create in one command: upsert). To get more information about available flags, use `--help`.

* `--delimiter`: Column delimiter for the structured file
* `--fields`: Name of the fields for columns, comma separated (can be edited later)
* `--name`: Name of the new data source. If not provided it will be inferred from the uploaded file name (optional)
* `--name-prefix-path`: Path prefix for new data sources (optional)
* `--raw`: Upload file as is (optional, default: false)


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
