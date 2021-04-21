## upctl server load

Load a CD-ROM into the server

```
upctl server load <UUID/Title/Hostname...> [flags]
```

### Examples

```
upctl server load my_server4 --storage 01000000-0000-4000-8000-000080030101
```

### Options

```
      --storage string   The UUID of the storage to be loaded in the CD-ROM device.
  -h, --help             help for load
```

### Options inherited from parent commands

```
  -t, --client-timeout duration   CLI timeout when using interactive mode on some commands (default 1m0s)
      --colours                   Use terminal colours (default true)
      --config string             Config file
  -o, --output string             Output format (supported: json, yaml and human) (default "human")
```

### SEE ALSO

* [upctl server](upctl_server.md)	 - Manage servers

###### Auto generated by spf13/cobra on 21-Apr-2021