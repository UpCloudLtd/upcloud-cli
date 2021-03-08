# UpCloud command line client - upctl

[![UpCloud upctl test](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml/badge.svg)](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml)

upctl is a command line client for UpCloud. It allows you to control your resources from the command line or any compatible interface. Contributions from the community are welcomed!

* Check Github issues
  * or create new issue
* Any improvement ideas for Documentation are more than welcome
  * Please create a PR for any additions or corrections
* The Cli uses [upcloud-go-api](https://github.com/UpCloudLtd/upcloud-go-api) as it's base api
* Built on [Cobra](https://cobra.dev)

## Quick Start

Create .upctl (yaml) config file for user details.

``` yaml
username: user
password: pass
```

You can use upctl with go:

``` bash
go run cmd/upctl.go --help
```

Or build a binary of the tool

``` bash
make build
./upctl --help
```

## Examples

### Create server

``` bash
./upctl server create --hostname test-server.io --zone es-mad1 --ssh-keys id_rsa.pub
```

> NOTE: You will have to specify a method for authentication by
>
> * ssh-keys `--ssh-keys id_rsa.pub`
> * or password delivery `--password-delivery email`
>
> Or maybe you have a non-default OS, that you have created. `--os your-custom-img`. Then you won't need to define authentication, since it's expected to have an authentication method.

Server title defaults to hostname. To set a different title, add `--title "Test server"`

### Create storage

``` bash
./upctl storage create --size 25 --title test-storage --zone es-mad1
```

Note: Storage size is in GB.

### Attach storage to server

``` bash
./upctl server storage attach <SERVER-UUID> --storage <STORAGE-UUID> 
```

## Development

Besides Golang, you'll need pre-commit and some other tools. Please [install pre-commit](https://pre-commit.com/#install) on your own machine, and then run the following commands within the repository folder:

``` bash
go get -u golang.org/x/lint/golint
go get -u github.com/go-critic/go-critic/cmd/gocritic
pre-commit install
```
