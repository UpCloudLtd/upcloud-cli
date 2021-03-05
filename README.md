# UpCloud command line client - upctl

[![UpCloud upctl test](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml/badge.svg)](https://github.com/UpCloudLtd/upctl/actions/workflows/test.yml)

UpCloud command line client is created to give our users more choices how to control their resources. Contributions from the community are welcomed!

* Check Github issues or create more issues
* Improve documentation
* Built on [Cobra](https://cobra.dev)

## Quick Start

```
make build

go run cmd/upctl.go
```

Create .upctl (yaml) config file for user details.

```
username: user
password: pass
```

## Examples

### Create server

```
upctl server create --hostname test-server.io --zone es-mad1
```

You will have to specify a method for authentication by password delivery `--password-delivery email` or ssh-keys `--ssh-keys id_rsa.pub`. Or maybe you have a non-default OS, that you have created. `--os your-custom-img`

Server title is by default the hostname. To set a different title, add `--title "Test server"`

### Create storage

```
upctl storage create --size 25 --title test-storage --zone es-mad1
```

Note: Storage size is in GB.

### Attach storage to server

```
upctl server storage attach <SERVER-UUID> --storage <STORAGE-UUID> 
```

## Development

Besides Golang, you'll need pre-commit and some other tools. Please [install pre-commit](https://pre-commit.com/#install) on your own machine, and then run the following commands within the repository folder:
```
go get -u golang.org/x/lint/golint
go get -u github.com/go-critic/go-critic/cmd/gocritic
pre-commit install
```
