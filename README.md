# UpCloud CLI - upctl

[![upcloud-cli test](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml/badge.svg)](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml)

`upctl` provides a command-line interface to UpCloud services. It allows you
to control your resources from the command line or any compatible interface.

```txt
upctl a CLI tool for managing your UpCloud services.

Usage:
upctl [command]

Available Commands:
account     Manage account
completion  Generates shell completion
database    Manage databases
help        Help about any command
ip-address  Manage ip address
network     Manage network
router      Manage router
server      Manage servers
storage     Manage storages
version     Display software information
zone        Display zone information

Options:
  -t, --client-timeout duration   CLI timeout when using interactive mode on some commands
                                  Default: 0s

  --config string                 Config file

  --debug bool                    Print out more verbose debug logs
                                  Default: false

  --force-colours                 force coloured output despite detected terminal support

  --no-colours                    disable coloured output despite detected terminal support

  -o, --output string             Output format (supported: json, yaml and human)
                                  Default: human

Use "upctl [command] --help" for more information about a command.
```

## Installation

To use upctl as a binary, download it from the
[Releases](https://github.com/UpCloudLtd/upcloud-cli/releases) page. After downloading, verify that the client works.

### macOS

```bash
brew tap UpCloudLtd/tap
brew install upcloud-cli
upctl -h
```

Setting up bash completion requires a few commands more.

First, install `bash-completion`, if it is not installed already.

```bash
brew install bash-completion
echo '[ -f "$(brew --prefix)/etc/bash_completion" ] && . "$(brew --prefix)/etc/bash_completion"' >> ~/.bash_profile
```

Then configure the shell completions for `upctl` by saving the output of `upctl completion bash` in `upctl` file under `/etc/bash_completion.d/`:

```bash
upctl completion bash > $(brew --prefix)/etc/bash_completion.d/upctl
. $(brew --prefix)/etc/bash_completion
```

### Linux

#### AUR
```
yay -S upcloud-cli
```

#### Other Linux distros

Use the package corresponding to your distro (deb, rpm, apk) from the [releases page](https://github.com/UpCloudLtd/upcloud-cli/releases), example for Debian like:

```bash
# Replace <VERSION> with the version you want to install
curl -Lo upcloud-cli.deb https://github.com/UpCloudLtd/upcloud-cli/releases/download/v<VERSION>/upcloud-cli-<VERSION>_amd64.deb
sudo dpkg -i upcloud-cli.deb
upctl -h
```

Bash completion can also be set up with some extra commands. You should adapt this for your package manager.

First, install `bash-completion`, if it is not installed already.

```bash
sudo apt install bash-completion
echo "[ -f /etc/bash_completion ] && . /etc/bash_completion" >> ~/.bashrc
```

Then configure the shell completions for `upctl` by either sourcing `upctl completion bash` output in your bash `.bashrc` or by saving the output of that command in `upctl` file under `/etc/bash_completion.d/`:

```bash
# First alternative
echo 'source <(upctl completion bash)' >>~/.bashrc

# Second alternative
upctl completion bash | sudo tee /etc/bash_completion.d/upctl > /dev/null
. /etc/bash_completion
```

### Windows
```bash
Invoke-WebRequest -Uri "https://github.com/UpCloudLtd/upcloud-cli/releases/download/v<VERSION>/upcloud-cli-<VERSION>_windows_x86_64.zip" -OutFile "upcloud-cli.zip"
Expand-Archive -Path upcloud-cli.zip -Destination 'C:\Program Files\Upcloud CLI'
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\Upcloud CLI", [System.EnvironmentVariableTarget]::Machine)
upctl.exe -h
```

## Quick Start

Create a `upctl.yaml` config file with user credentials into your home
directory's .config dir ($HOME/.config/upctl.yaml).

```yaml
username: your_upcloud_username
password: your_upcloud_password
```

Credentials can also be stored at environment variables `UPCLOUD_USERNAME` and `UPCLOUD_PASSWORD`. If variables
are set, matching config file items are ignored.

> NOTE: Make sure your account allows API connections. To do so, log into
> [UpCloud control panel](https://hub.upcloud.com/login) and go to **Account**
> -> **Permissions** -> **Allow API connections** checkbox.

Run something to test that the credentials are working.

```bash
$ upctl server list
 UUID                                   Hostname             Plan        Zone      State
────────────────────────────────────── ──────────────────── ─────────── ───────── ─────────
 00229ddf-0e46-45b5-a8f7-cad2c8d11f6a   server1              2xCPU-4GB   de-fra1   stopped
 003c9d77-0237-4ee7-b3a1-306efba456dc   server2              1xCPU-2GB   sg-sin1   started
```

## Examples

Every command has a help and examples included and you can find all its options
by adding `-h` at the end of the command, like `upctl network list -h`. Below,
you'll find a few common commands that have many other available options as
well.

* Create a new server

```bash
upctl server create --hostname test-server.io --zone de-fra1 --ssh-keys ~/.ssh/id_rsa.pub
```

* Create storage

```bash
upctl storage create --size 25 --title test-storage --zone es-mad1
```

> Note: Storage size is in GB.

* Attach storage to server

```bash
upctl server storage attach <SERVER-UUID> --storage <STORAGE-UUID>
```

## Documentation

The detailed documentation can be found [here](docs/upctl.md)


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

* upctl uses [UpCloud Go API SDK](https://github.com/UpCloudLtd/upcloud-go-api)
* upctl is built on [Cobra](https://cobra.dev)

You need a Go version 1.11+ installed on your development machine.
Use `make` to build and test the CLI. Makefile help can be found:

```
$ make help
build                Build program binary for current os/arch
doc                  Generate documentation (markdown)
build-all            Build all targets
build-linux          Build program binary for linux x86_64
build-darwin         Build program binary for darwin x86_64
build-windows        Build program binary for windows x86_64
build-freebsd        Build program binary for freebsd x86_64
test                 Run tests
fmt                  Run gofmt on all source files
clean                Cleanup everything
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

Make the changes with your favorite editor. Once you're done, create a new branch and push it back to Github.
```bash
git checkout -b <branch-name>
<add your changes, "git status" helps>
git commit -m "New feature: create a new server in the nearest zone if not specified"
git push --set-upstream <branch-name>
```

After pushing the new branch, browse to your fork of the repository in GitHub and create a pull request from there.
Once the pull request is created, please make changes to your branch based on the comments & discussion in the PR.

## Releasing for beta versions

Current release process:

* Update CHANGELOG.md
* Test GoReleaser config with `goreleaser check`
* Tag a commit with the version you want to release e.g. `v1.2.3`
* Push the tag & commit to GitHub
  * GitHub actions will automatically set the version based on the tag, create a GitHub release, build the project, and upload binaries & SHA sum to the GitHub release
* [Edit the new release in GitHub](https://github.com/UpCloudLtd/upcloud-cli/releases) and add the changelog for this release
* Done!

## License

[MIT license](LICENSE)
