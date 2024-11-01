# UpCloud CLI - upctl

[![upcloud-cli test](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml/badge.svg)](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml)

`upctl` provides a command-line interface to [UpCloud](https://upcloud.com/) services. It allows you
to control your resources from the command line or any compatible interface.

## Getting started

For instructions on how to install `upctl`, configure credentials, and run commands, see [Getting started](https://upcloudltd.github.io/upcloud-cli/) instructions in the documentation.

## Install with go install

Install the latest version of `upctl` with `go install` by running:

```sh
go install github.com/UpCloudLtd/upcloud-cli/v3/...@latest
```

Run `upctl version` to verify that the tool was installed successfully and `upctl help` to print usage instructions.

```sh
upctl version
upctl help
```

## Exit codes

Exit code communicates success or the type and number of failures. Possible exit codes of `upctl` are:

Exit code | Description
--------- | -----------
0         | Command(s) executed successfully.
1 - 99    | Number of failed executions. For example, if stopping four servers and API returns error for one of the request, exit code will be 1.
100 -     | Other, non-execution related, errors. For example, required flag missing.

## Examples

Every command has a `--help` parameter that can be used to print detailed usage instructions and examples on how to use the command. For example, run `upctl network list --help`, to display usage instructions and examples for `upctl network list` command.

See [examples](./examples/) directory for examples on more complex use-cases.

## Documentation

The detailed documentation is available in [GitHub pages](https://upcloudltd.github.io/upcloud-cli/).

To generate markdown version of command reference, run `make md-docs`. Command reference will then be generated into `docs/commands_reference`.

```sh
make md-docs
```

To run the MkDocs documentation locally, run make docs and start static http server (e.g., `python3 -m http.server 8000`) in `site/` directory or run mkdocs serve in repository root.

```sh
make docs
mkdocs serve
```

## Contributing

Contributions from the community are much appreciated! Please note that all features using our
API should be implemented with [UpCloud Go API SDK](https://github.com/UpCloudLtd/upcloud-go-api).
If something is missing from there, add an issue or PR in that repository instead before implementing it here.

* Check GitHub issues and pull requests before creating new ones
  * If the issue isn't yet reported, you can [create a new issue](https://github.com/UpCloudLtd/upcloud-cli/issues/new).
* Besides bug reports, all improvement ideas and feature requests are more than welcome and can be submitted through GitHub issues.
  * New features and enhancements can be submitted by first forking the repository and then sending your changes back as a pull request.
* Following [semantic versioning](https://semver.org/), we won't accept breaking changes within the major version (1.x.x, 2.x.x etc).
  * Such PRs can be open for some time and are only accepted when the next major version is being created.

## Development

* `upctl` uses [UpCloud Go API SDK](https://github.com/UpCloudLtd/upcloud-go-api)
* `upctl` is built on [Cobra](https://cobra.dev)

You need a Go version 1.20+ installed on your development machine.

Use `make` to build and test the CLI. Makefile help can be found by running `make help`.

```sh
make help
```

### Debugging
Environment variables `UPCLOUD_DEBUG_API_BASE_URL` and `UPCLOUD_DEBUG_SKIP_CERTIFICATE_VERIFY` can be used for HTTP client debugging purposes. More information can be found in the client's [README](https://github.com/UpCloudLtd/upcloud-go-api/blob/986ca6da9ca85ff51ecacc588215641e2e384cfa/README.md#debugging) file.

### Requirements

This repository uses [pre-commit](https://pre-commit.com/#install) and [go-critic](https://github.com/go-critic/go-critic)
for maintaining code quality. Installing them is not mandatory, but it helps in avoiding the problems you'd
otherwise encounter after opening a pull request as they are run by automated tests for all PRs.

### Development quickstart

To begin development, first fork the repository to your own account, clone it and begin making changes.
```bash
git clone git@github.com/username/upcloud-cli.git
cd upcloud-cli
pre-commit install
```

Make the changes with your favorite editor. Once you're done, create a new branch and push it back to GitHub.
```bash
git checkout -b <branch-name>
<add your changes, "git status" helps>
git commit -m "New feature: create a new server in the nearest zone if not specified"
git push --set-upstream <branch-name>
```

After pushing the new branch, browse to your fork of the repository in GitHub and create a pull request from there.
Once the pull request is created, please make changes to your branch based on the comments & discussion in the PR.

## License

[MIT license](LICENSE)
