## upctl database connection list

List current connections to specified databases

```
upctl database connection list <UUID/Title...> [flags]
```

### Examples

```
upctl database connection list
```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
  -t, --client-timeout duration   Client timeout to use in API calls.
      --config string             Configuration file path.
      --debug                     Print out more verbose debug logs.
      --force-colours[=true]      Force coloured output despite detected terminal support.
      --no-colours[=true]         Disable coloured output despite detected terminal support. Colours can also be disabled by setting NO_COLOR environment variable.
  -o, --output string             Output format (supported: json, yaml and human) (default "human")
```

### SEE ALSO

* [upctl database connection](upctl_database_connection.md)	 - Manage database connections

###### Auto generated by spf13/cobra on 5-Aug-2022