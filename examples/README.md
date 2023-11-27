# Examples

This directory contains examples on more complex `upctl` use-cases. As `upctl` is often used in scripts the examples also aim to parse values from machine readable outputs. This allows using the examples also as end-to-end test cases and makes them more copy-pasteable.

## Testing

The examples in this directory are validated with [mdtest](https://github.com/UpCloudLtd/mdtest). It parses `env` and `sh` code-blocks from the markdown files and executes those as scripts.

The tool can be installed with go install.

```sh
go install github.com/UpCloudLtd/mdtest@latest
```

To test the examples, run `mdtest .`.

```sh
mdtest .
```
