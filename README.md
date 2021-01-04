# cli
UpCloud command line client

## Usage

```
make build

go run cmd/upctl.go
```

Create .upctl (yaml) config file for user details.

```
username: user
password: pass
```

## Development

Besides Golang, you'll need pre-commit and some other tools. Please [install pre-commit](https://pre-commit.com/#install) on your own machine, and then run the following commands within the repository folder:
```
go get -u golang.org/x/lint/golint
go get -u github.com/go-critic/go-critic/cmd/gocritic
pre-commit install
```
