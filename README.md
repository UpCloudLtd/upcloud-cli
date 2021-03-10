# UpCloud CLI - upctl

[![UpCloud upctl test](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml/badge.svg)](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml)

`upctl` provides a command-line interface to UpCloud services. It allows you
to control your resources from the command line or any compatible interface.

* upctl uses [upcloud-go-api](https://github.com/UpCloudLtd/upcloud-go-api)
* Built on [Cobra](https://cobra.dev)

## Quick Start

Download the appropriate binary from the
[Releases](https://github.com/UpCloudLtd/upcloud-cli/releases) page.

Create a `.upctl` (yaml) config file with user credentials into your home
directory or the current directory.

```yaml
username: upcloud_username
password: upcloud_password
```

Run the command!

```bash
upctl -h
```

## Examples

### Create server

```bash
upctl server create --hostname test-server.io --zone es-mad1 --ssh-keys id_rsa.pub
```

> NOTE: You will have to specify a method for authentication by
>
> * ssh-keys `--ssh-keys id_rsa.pub`
> * or password delivery `--password-delivery email`
>
> Note: If you have a custom default operating system template, these cannot be used. Use `--os your-custom-img` to specify your template; it's expected you have the correct authentication already set up in your custom method. 

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

## Usage with go 

Requires Golang version 1.11+.

### with go get

```bash
go get github.com/UpCloudLtd/upcloud-cli
```

### From source code

Clone the repo at https://github.com/UpCloudLtd/upcloud-cli.

You can use upctl with go:

```bash
go run cmd/upctl/main.go --help
```

Build the binary with:

```bash
make build
./bin/upctl --help
```

## Contributing

Contributions from the community are welcome!

* Check GitHub issues and pull requests if you want to contribute
  * create a new issue if you find something missing
* Any improvement ideas for Documentation are more than welcome
  * Please create a PR for any additions or corrections

## Development

You need a Golang version 1.11+ installed on you development machine.
Use `make` to build and test the CLI. Makefile help can be found:

```
$ make help
build             Build program binary for current os/arch
build-all         Build all targets
build-linux       Build program binary for linux x86_64
build-darwin      Build program binary for darwin x86_64
build-windows     Build program binary for windows x86_64
test              Run tests
fmt               Run gofmt on all source files
clean             Cleanup everything
```

## License

[Apache License 2.0](LICENSE)
