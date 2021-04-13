# UpCloud CLI - upctl

[![upcloud-cli test](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml/badge.svg)](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml)

`upctl` provides a command-line interface to UpCloud services. It allows you
to control your resources from the command line or any compatible interface.

```bash
upctl a CLI tool for managing your UpCloud services.

Usage:
upctl [flags]
upctl [command]

Available Commands:
account     Manage account
completion  Generates shell completion
help        Help about any command
ip-address  Manage ip address
network     Manage network
router      Manage router
server      Manage servers
storage     Manage storages
version     Display software infomation

Options:
  -t, --client-timeout duration   CLI timeout when using interactive mode on some commands
                                  Default: 1m0s

  --colours bool                  Use terminal colours (supported: auto, true, false)
                                  Default: true

  --config string                 Config file

  -o, --output string             Output format (supported: json, yaml and human)
                                  Default: human

Use "upctl [command] --help" for more information about a command.
```

## Installation

As upctl is currently in beta, we aren't yet offering it through different packaging managers. Expect this
to change once we hit version 1.0.0.

To use upctl, download the upctl binary from the
[Releases](https://github.com/UpCloudLtd/upcloud-cli/releases) page. After downloading, verify that the client works.

### macOS

```bash
brew tap UpCloudLtd/tap
brew install upcloud-cli
upctl -h
```

Setting up bash completion requires a few commands more.

```bash
brew install bash-completion
sudo upctl completion bash > /usr/local/etc/bash_completion.d/upctl
echo "[ -f /usr/local/etc/bash_completion ] && . /usr/local/etc/bash_completion" >> ~/.bash_profile
. /usr/local/etc/bash_completion
```

###  Linux

## AUR
```
yay -S upcloud-cli
```

## Other distros

Use the package corresponding to your distro (deb, rpm, apk), example Debian like:

```bash
sudo curl -o upcloud.deb https://github.com/kaminek/upcloud-cli/releases/download/v<VERSION>/upcloud-cli-<VERSION>_amd64.deb
sudo chmod +x /usr/local/bin/upctl
upctl -h
```

Bash completion can also be set up with some extra commands. You should adapt this for your package manager.
```bash
sudo apt install bash-completion
sudo upctl completion bash > /etc/bash_completion.d/upctl
echo "[ -f /etc/bash_completion ] && . /etc/bash_completion" >> ~/.bash_profile
. /etc/bash_completion
```

### Windows
```bash
Invoke-WebRequest -Uri "https://github.com/UpCloudLtd/upcloud-cli/releases/download/v<VERSION>/upcloud-cli-<VERSION>_windows_x86_64.zip" -OutFile "upcloud-cli.zip"
unzip upcloud-cli.zip
upctl.exe -h
```

## Quick Start

Create a `upctl.yaml` config file with user credentials into your home
directory's .config dir ($HOME/.config/upctl.yaml) or the current directory.

```yaml
username: your_upcloud_username
password: your_upcloud_password
```

Credentials can also be stored at environment variables `UPCLOUD_USERNAME` and `UPCLOUD_PASSWORD`. If variables
are set, matching config file items are ignored.

Run something to test that the credentials are working.

```bash
$ upctl server list
 UUID                                   Hostname             Plan        Zone      State   
────────────────────────────────────── ──────────────────── ─────────── ───────── ─────────
 00229ddf-0e46-45b5-a8f7-cad2c8d11f6a   server1              2xCPU-4GB   de-fra1   stopped 
 003c9d77-0237-4ee7-b3a1-306efba456dc   server2              1xCPU-2GB   sg-sin1   started
```

## Examples

Every command has a help included and you can find all its options by adding `-h` at the end of the command,
like `upctl network list -h`. Below, you'll find a few common commands that have many other available options as well.

### Create a new server

```bash
upctl server create --hostname test-server.io --zone es-mad1 --ssh-keys ~/.ssh/id_rsa.pub
```

> NOTE: You will have to specify a method for authentication by
>
> * ssh-keys `--ssh-keys id_rsa.pub`
> * or password delivery `--password-delivery sms`
>
> Note: If you have a custom operating system template, these don't work by default. Use `--os your-custom-img` with `--metadata` to make the keys available through the [metadata service](https://developers.upcloud.com/1.3/8-servers/#metadata-service).

Server title defaults to hostname. To set a different title, add `--title "Test server"`

### Create storage

```bash
upctl storage create --size 25 --title test-storage --zone es-mad1
```

Note: Storage size is in GB.

### Attach storage to server

```bash
upctl server storage attach <SERVER-UUID> --storage <STORAGE-UUID>
```

## Documentation

The detailed documentation can be found [here](docs/upctl.md)


## Contributing

Contributions from the community are much appreciated! Please note that all features using our
API should be implemented with [UpCloud Golang API SDK](https://github.com/UpCloudLtd/upcloud-go-api).
If something is missing from there, add an issue or PR in that repository instead before implementing it here.

* Check GitHub issues and pull requests before creating new ones
  * If the issue isn't yet reported, you can [create a new issue](https://github.com/UpCloudLtd/upcloud-cli/issues/new).
* Besides bug reports, all improvement ideas and feature requests are more than welcome and can be submitted through GitHub issues.
  * New features and enhancements can be submitted by first forking the repository and then sending your changes back as a pull request.
* Following [semantic versioning](https://semver.org/), we won't accept breaking changes within the major version (1.x.x, 2.x.x etc).
  * Such PRs can be open for some time and are only accepted when the next major version is being created.

## Development

* upctl uses [UpCloud Golang API SDK](https://github.com/UpCloudLtd/upcloud-go-api)
* upctl is built on [Cobra](https://cobra.dev)

You need a Golang version 1.11+ installed on your development machine.
Use `make` to build and test the CLI. Makefile help can be found:

```
$ make help
build                Build program binary for current os/arch
doc                  Build documentation (markdown)
build-all            Build all targets
build-linux          Build program binary for linux x86_64
build-darwin         Build program binary for darwin x86_64
build-windows        Build program binary for windows x86_64
test                 Run tests
fmt                  Run gofmt on all source files
clean                Cleanup everything
```

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
* Tag a commit with the version you want to release e.g. `v1.2.3`
* Push the tag & commit to GitHub
  * GitHub actions will automatically set the version based on the tag, create a GitHub release, build the project, and upload binaries & SHA sum to the GitHub release
* [Edit the new release in GitHub](https://github.com/UpCloudLtd/upcloud-cli/releases) and add the changelog for this release
* Done!

## License

[MIT license](LICENSE)
