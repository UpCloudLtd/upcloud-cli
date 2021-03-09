# Upctl

`upctl` is a command line client for UpCloud. It allows you to control your
resources from the command line or any compatible interface. Contributions from
the community are welcome!
[![UpCloud upctl test](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml/badge.svg)](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml)

* Check GitHub issues
  * or create a new issue
* Any improvement ideas for Documentation are more than welcome
  * Please create a PR for any additions or corrections
* The Cli uses [upcloud-go-api](https://github.com/UpCloudLtd/upcloud-go-api)
* Built on [Cobra](https://cobra.dev)

## Quick Start

Create .upctl (yaml) config file with user credentials.

``` yaml
username: user
password: pass
```

You can use upctl with go:

``` bash
go run cmd/upctl/main.go --help
```

Or build the binary with:

``` bash
make build
./bin/upctl --help
```

## Examples

### Create server

``` bash
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

``` bash
upctl storage create --size 25 --title test-storage --zone es-mad1
```

Note: Storage size is in GB.

### Attach storage to server

``` bash
upctl server storage attach <SERVER-UUID> --storage <STORAGE-UUID> 
```

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
