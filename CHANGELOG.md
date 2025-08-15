# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [3.22.0] - 2025-08-15

### Added

- Object storage service creation with `object-storage create` command supporting labels, networks, and configured status
- Object storage label management with dedicated commands: `object-storage label add`, `object-storage label remove`, and `object-storage label list`
- Object storage network management with `object-storage network attach` and `object-storage network detach` commands
- Object storage bucket management with `object-storage bucket create`, `object-storage bucket delete`, and `object-storage bucket list` commands
- Object storage user management with `object-storage user create`, `object-storage user delete`, and `object-storage user list` commands
- Object storage user access key management with `object-storage create-access-key`, `object-storage delete-access-key`, and `object-storage list-access-keys` commands
- Expose GPU limits in `account show` command
- Expose GPU model and amount in `server plans` command
- Add `audit-log export` command.
- Add support for customising storage tier and size for Kubernetes node groups utilising GPU and Cloud Native plans.

## [3.21.0] - 2025-07-15

### Added

- Add `zone devices` and `zone devices show` commands for listing available devices.

## [3.20.2] - 2025-07-09

### Changed

- Group GPU plans under a separate section in human readable output of `server plans` command.

### Fixed

- Run as non-root user in container.

## [3.20.1] - 2025-06-09

### Fixed

- In database create, parse numeric string values into strings, E.g., Postgres version property will be now correctly parsed as string value.

## [3.20.0] - 2025-05-26

### Added

- Support load-balancer resources in `all list` and `all purge` commands.
- Add `--wait` flag to `load-balancer delete` command.

### Fixed

- Set exit code to 103 if authentication fails also in commands that do not take arguments (e.g. `server list`).

## [3.19.1] - 2025-04-29

### Fixed

- Redact token from debug logs.

## [3.19.0] - 2025-04-11

### Added

- Add `database create` command.
- Support kubernetes, server, server-group, storage, and tag resources in `all list` and `all purge` commands.
- Add `--stop` flag to `server delete` command.
- Add `--backups` flag to `storage delete` command.
- Add `--wait` flag to `kubernetes delete` command.

### Fixed

- When running kubernetes commands, show deprecation message when using `uks` alias. The deprecation message is currently displayed when using `k8s` alias, which will not be removed.

## [3.18.0] - 2025-04-04

### Added

- Experimental support for listing and deleting all resources with `all list` and `all purge` commands. This initial version supports networks, network-peerings, routers, databases, and object-storages.
- Add `--wait` flag to `object-storage delete` command.
- Add `--disable-termination-protection` and `--wait` flags to `database delete` command.

### Fixed

- Kubernetes node-group subcommand not working

### Deprecated

- Deprecated `upctl kubernetes nodegroup` command, use `upctl kubernetes node-group` instead.

## [3.17.0] - 2025-03-14

### Added

- Server relocation support

## [3.16.1] - 2025-03-07

### Fixed

- Remove client side default value for kubernetes cluster plan and use default from API instead if no plan is defined.

## [3.16.0] - 2025-03-05

This release introduces GitHub artifact attestations for our release binary assets.

### Added

- Experimental support for reading password or token from system keyring.
- Experimental support from saving token to system keyring with `upctl account login --with-token`.

### Changed

- In human readable output of `kubernetes plans`, remove `server_number` column and hide deprecated plans.

## [3.15.0] - 2025-02-26

### Added

- Added support for Valkey properties
- Add termination_protection to upctl database show output
- Experimental support for token authentication by defining token in `UPCLOUD_TOKEN` environment variable.
- Experimental support for managing tokens with `account token` commands.
- New command names and aliases added to improve consistency:
    - Commands:
        - load-balancer
        - network-peering
        - object-storage
        - server-group
    - Aliases:
        - account: acc
        - gateway: gw
        - network-peering: np
        - object-storage: obs
        - partner: pr
        - router: rt
        - server: srv
        - server-group: sg
        - storage: st

### Changed

- List cloud native plans in their own section in human readable `server plans` output.

### Fixed

- Prevent filename completion of flags that don't take filename args.

### Deprecated

- Deprecation of some commands and aliases ( new command names added to improve consistency )
    - Deprecated commands: 
        - loadbalancer
        - networkpeering
        - objectstorage
        - servergroup
    - Deprecated aliases:
        - object-storage: objsto


## [3.14.0] - 2025-01-08

### Added

- Allow using unix style glob pattern as an argument. For example, if there are two servers available with titles `server-1` and `server-2`, these servers can be stopped with `upctl server stop server-*` command.
- `--delete-buckets` option to `objectstorage delete` command.

### Fixed

- In `objectstorage delete` command, delete only user defined policies when `--delete-policies` flag is enabled as trying to delete system defined policy will cause an error.

## [3.13.0] - 2024-12-13

### Added

- Partner API support

### Changed

- Go version bump to 1.22

## [3.12.0] - 2024-11-18

### Changed

- Take server state into account in server completions. For example, do not offer started servers as completions for `server start` command.
- Allow using UUID prefix as an argument. For example, if there is only one network available that has a UUID starting with `0316`, details of that network can be listed with `upctl network show 0316` command.
- Match title and name arguments case-insensitively if the given parameter does not resolve with an exact match.

## [3.11.1] - 2024-08-12

### Fixed

- In `storage modify`, avoid segfault if the target storage does not have backup rule in the storage details. This would have happened, for example, when renaming private templates.

## [3.11.0] - 2024-07-23

### Added

- Add labels to `database show` output.
- In `router show` command, list static routes in the human readable output and add static route type field to all outputs.

## [3.10.0] - 2024-07-17

### Changed

- In `server create` command, use `Ubuntu Server 24.04 LTS (Noble Numbat)` as default value for `--os`. The new default template only supports SSH key based authentication. Use `--ssh-keys` option to provide the keys when creating a server with the default template.
- In `server create` command, enable metadata service by default when the selected (or default) template uses cloud-init (`template_type` is `cloud-init`) and thus requires it.

## [3.9.0] - 2024-07-04

### Added

- Add `gateway plans` command for listing gateway plans.
- Add `objectstorage regions` command for listing managed object storage regions.

### Changed

- In all outputs of `server plans`, sort plans by CPU count, memory amount, and storage size.
- In human readable output of `server plans`, group plans by type.

### Fixed

- Omit storage tier from server create payload to use plans default storage tier. This allows creating servers with developer plans that do not allow creating MaxIOPS storages with the server. Other plan types will continue to use MaxIOPS by default.

## [3.8.1] - 2024-05-24

### Changed

- Allow creating Kubernetes cluster without node-groups.

## [3.8.0] - 2024-04-30

### Added

- Add `host list` command for listing private cloud hosts.
- Add `--host` argument to `server restart` command.
- Add `--avoid-host` and `--host` arguments to `server start` command.
- Add `Parent zone` column to `zone list` output for users that have access to private zones.

### Changed

- Go version bump to 1.21

## [3.7.0] - 2024-04-04

### Added

- Add managed object storage name to `show` and `list` command outputs.
- Add completions for managed object storages and allow using managed object storage name (in addition to its UUID) as a positional argument.
- Add paging flags (`--limit` and `--page`) to `database`, `gateway`, `loadbalancer`, and `objectstorage` list commands

## [3.6.0] - 2024-03-07

### Added

- Support Kubernetes cluster labels: list labels with `show` commands and manage them with `create` and `modify` commands
- Add `networkpeering` commands (`delete`, `disable`, `list`) for network peering management.

## [3.5.0] - 2024-02-29

### Added

- Policies section in `objectstorage show` command
- Optional parameters `--delete-users` and `--delete-policies` to `objectstorage delete`

### Changed

- Users in `objectstorage show` now contain column `ARN` instead of `Updated`

## [3.4.0] - 2024-02-08

### Added

- Add `gateway` commands (`delete`, `list`) for network gateway management.
- In machine readable outputs of server list, add support for `--show-ip-addresses` parameter.
- Support for sub-account deletion via `account delete` command

## [3.3.0] - 2024-01-23

### Added

- Support for storage encryption to storage `create`, `clone`, `show`, and `list` commands as well as server `create` and `show` commands.
- _Managed object storages_ field to human readable output of `account show`.
- Commands for listing accounts and permissions.

### Removed

- From human output of `storage list`, _Created_ column. This field is still available in the machine readable outputs.

## [3.2.2] - 2024-01-02

### Added
- Support nested properties in `database properties *` and `database properties * show *` outputs. For example upctl `max_background_workers` sub-property of `timescaledb` PostgreSQL property is listed as `timescaledb.max_background_workers` in human output of `database properties pg` and its details can printed with `upctl database properties pg show timescaledb.max_background_workers` command.

### Fixed
- Do not return error on valid `dhcp-default-route` values in `network create` and `network modify` commands.

## [3.2.1] - 2023-11-29

### Added
- Add backend `TLS configs` field to `loadbalancer show` command

## [3.2.0] - 2023-11-15

### Added
- Add `objectstorage` commands (`delete`, `list`, `show`) for managed object storage management

## [3.1.0] - 2023-11-06

### Added

- Add `network_peerings`, `ntp_excess_gib`, `storage_maxiops` and `load_balancers` fields to `account show` outputs.
- Add `--version` parameter to `kubernetes create` and `version` field to `kubernetes show` output.

### Changed

- **Breaking**: Update return type of `kubernetes versions` from list of strings to list of objects. (No major version bump, because this end-point has not been included in the API docs)

### Fixed

- Use correct currency symbol, `€` instead of `$`, in human output of `account show`.
- Remove `v` prefix from version in `version` command output and User-Agent header

## [3.0.0] - 2023-10-18

This release updates output of `show` and `list` commands to return the API response as defined in the UpCloud Go SDK. See below for detailed list of changes.

In addition, `kubernetes create` will now, by default, block all access to the cluster. To be able to connect to the cluster, define list of allowed IP addresses and/or CIDR blocks or allow access from any IP.

### Added
- **Breaking**: Add `--kubernetes-api-allow-ip` argument to `kubernetes create` command. This changes default behavior from _allow access from any IP_ to _block access from all IPs_. To be able to connect to the cluster, define list of allowed IP addresses and/or CIDR blocks or allow access from any IP.
- Add `Kubernetes API allowed IPs` field to `kubernetes show` output.
- Add `kubernetes nodegroup show` for displaying node-group details. This also adds _Nodes_ table and _Anti-affinity_ field that were not available in previous `kubernetes show` output.
- Add `kubernetes modify` command for modifying IP addresses that are allowed to access cluster's Kubernetes API.
- Add `Kubernetes API allowed IPs` field to `kubernetes show` output.
- Add `database session list` for listing active database sessions.
- Add `database session cancel` for cancelling an active database session.

### Changed
- **Breaking**: In JSON and YAML output of `database list`: return the full API response. Value of `title` is not replaced with value from `name`, if `title` is empty.
- **Breaking**: In JSON and YAML output of `database types`: return the full API response. This changes the top level datatype from list to object, where keys are the available database type, e.g., `pg` and `mysql`.
- **Breaking**: In JSON and YAML output of `ip-address list`: return the full API response. This changes `partofplan` key to `part_of_plan` and `ptrrecord` key to `ptr_record`. The top level data-type changes from list to object.
- **Breaking**: In JSON and YAML output of `loadbalancer list`: return the full API response. This changes `state` field to `operational_state`.
- **Breaking**: In JSON and YAML output of `network list` and `network show`: return the full API response. Servers list will only contain server UUID and name. In `network list` output, the top level data-type changes from list to object.
- **Breaking**: In JSON and YAML output of `server list` and `server show`: return the full API response. This changes field `host_id` to `host`. `nics` is replaced with `networking` subfield `interfaces`. `storage` is replaced with `storage_devices`. `labels` contain subfield `label` which in turn contains the labels. In `server list` output, the top level data-type changes from list to object.
- **Breaking**: In JSON and YAML output of `server firewall show`: return the full API response. This removes fields `destination` and `source` fields in favor of `[destination|source]_address_start`, `[destination|source]_address_end`, `[destination|source]_port_start` and `[destination|source]_port_end`.
- **Breaking**: In JSON and YAML output of `server plans`: return the full API response. The top level data-type changes from list to object.
- **Breaking**: In JSON and YAML output of `storage list` and `storage show`: return the full API response. This changes `servers` field to contain `server` field, which in turn contains the servers. `labels` field will not be outputted if empty. In `storage list` output, the top level data-type changes from list to object.
- **Breaking**: In JSON and YAML output of `zone list`: return the full API response. The top level data-type changes from list to object.
- In JSON and YAML output of `kubernetes list`: return the full API response.
- In human readable output of `kubernetes show` command, show node-groups as table. Node-group details are available with `kubernetes nodegroup show` command.

### Fixed
- **Breaking**: In JSON and YAML output of `ip-address show`: use same JSON keys as in API documentation. This removes `credits` key that was used in place of `floating`.

### Removed
- **Breaking**: Remove `database connection list` and `database connection cancel` commands in favor of `database session` counterparts
- **Breaking**: In JSON and YAML output of `database properties * show`: pass-through the API response. This removes `key` field from the output.

## [2.10.0] - 2023-07-17
### Added
- Add `--disable-utility-network-access` for `kubernetes nodegroup create` command

### Fixed
- Use pending color (yellow) for kubernetes node group `scaling-down` and `scaling-up` states

## [2.9.1] - 2023-07-06
### Changed
- Release artifacts to follow package naming conventions provided by [nFPM](https://github.com/goreleaser/nfpm).
    ```
    upcloud-cli-2.9.0_x86_64.rpm  vs  upcloud-cli-2.9.1.x86_64.rpm  # no convention changes
    upcloud-cli-2.9.0_arm64.rpm   vs  upcloud-cli-2.9.1.aarch64.rpm # `arm64` -> `aarch64`

    upcloud-cli-2.9.0_amd64.apk   vs  upcloud-cli_2.9.1_x86_64.apk  # `cli-`  ->  `cli_`  & `amd64` -> `x86_64`
    upcloud-cli-2.9.0_arm64.apk   vs  upcloud-cli_2.9.1_aarch64.apk # `cli-`  ->  `cli_`  & `arm64` -> `aarch64`

    upcloud-cli-2.9.0_amd64.deb   vs  upcloud-cli_2.9.1_amd64.deb   # `cli-`  ->  `cli_`
    upcloud-cli-2.9.0_arm64.deb   vs  upcloud-cli_2.9.1_arm64.deb   # `cli-`  ->  `cli_`
    ```

## [2.9.0] - 2023-06-30
### Added
- Add `servergroup` commands (`create`, `delete`, `list`, `modify`, `show`) for server group management

## [2.8.0] - 2023-06-21
### Added
- Add support for OpenSearch database type
- Add `database index list` and `database index` commands for managing OpenSearch database indices
- Add completions for `--zone` arguments.
- Add `--private-node-groups` argument to `kubernetes create` command.
- Add _Private node groups_ field to `kubernetes show` output.
- Add `--label` flag to `server create` and `server modify` commands

## [2.7.1] - 2023-05-16
### Fixed
- Updated examples of `kubernetes create` command to use valid plans.

## [2.7.0] - 2023-05-02
### Added
- Add `ip` and `net` as aliases to `ip-address` and `network` commands, respectively.
- Add _Labels_ table to `loadbalancer show`, `network show`, `router show`, `server show`, and `storage show` outputs.
- Add `kubernetes plans` command for listing available plans.
- Add `--plan` argument to `kubernetes create` command for selecting cluster plan.
- Add `--wait` flag to `kubernetes create` command for waiting created cluster to reach running state.

## [2.6.0] - 2023-03-14
### Added
- The `upctl` container image now includes [jq](https://stedolan.github.io/jq/) tool for parsing values from JSON output.
- Add node-group states to `kubernetes show` output.
- Add completions for `--network` argument of `kubernetes create` and `server network-interface create`.
- Support also network name as input for `--network` argument of `kubernetes create` and `server network-interface create`

### Changed
- Completions will now only suggest private networks as arguments because names or UUIDs of public or utility networks are often not valid arguments.

## [2.5.0] - 2023-02-15
### Added
- Print warning about unknown resource state before exiting when execution is interrupted with SIGINT.
- Add `kubernetes nodegroup create`, `kubernetes nodegroup scale`, and `kubernetes nodegroup delete` commands (EXPERIMENTAL)
- Added support for all shell completions provided by `cobra`.
- Add `database properties <DB type>` command to list database properties for given database type and `database properties <DB type> show` command to show database property details.

### Changed
- Remove custom bash completion logic and replace it with `completion` command provided by `cobra`. To do this while supporting args with whitespace, whitespace in completions is replaced with non-breaking spaces.

### Fixed
- In `database show`: parse database version from metadata instead of properties. This enables displaying redis version instead of `<nil>`.

## [2.4.0] - 2022-12-19
### Added
- Add `kubernetes create`, `kubernetes config`, `kubernetes delete`, `kubernetes list`, `kubernetes show`, `kubernetes versions` commands (EXPERIMENTAL)
- Add `loadbalancer plans` command for listing available LB plans

## [2.3.0] - 2022-11-11
### Added
- Complete available types for `database plans`.
- Suppress positional argument filename completion for commands without specific completions.
- In `database list` output: if database has no title, database name is displayed in the title cell instead of leaving the cell empty, similarly than in the hub.
- Version information is parsed from `BuildInfo` when `upctl` binary was built without specifying `-ldflags` to define value for `.../config.Version`.
- Use alpine as base image for `upcloud/upctl` container image. This adds sh and other OS tools to the image and thus makes it more suitable for usage in CI systems.

### Fixed
- Remove debug leftover print from IP address completions.
- Added `/v2` postfix to module name in `go.mod`, this enables installing v2 versions of CLI with `go install`.

## [2.2.0] - 2022-10-17
### Added
- Include runtime operating system and architecture in `version` command output.
- Include instructions for defining credentials and API access in `upctl --help` output.

### Fixed
- Fix commands in `server delete` usage examples.
- Tune human output so that normal output is directed to `stdout`. Progress and error messages are still outputted to `stderr`.

## [2.1.0] - 2022-09-07
### Added
- Add `--wait` flag to `storage import` and `storage templatise` commands to wait until storage is in `online` state.
- In JSON and YAML output of `storage import`: information on target storage is now available under `storage` key.

### Fixed
- In human output of `storage list`: capitalize zone column header and color storage state similarly than in `storage show`.
- In human output of `storage import`: output UUID of created storage, instead of storage import operation. No UUID is outputted if existing storage was used.

## [2.0.0] - 2022-08-30
### Added
- Add `database delete` command.
- Add `loadbalancer delete` command.
- Add `Access` field to `storage show` output.
- Add fields `argument` and `resource` to JSON and YAML error outputs.

### Changed
- **Breaking**: Human output, including errors, is written to stderr instead of stdout.
- Refactor progress logging. This changes the appearance of progress logs. See [UpCloudLtd / progress](https://github.com/UpCloudLtd/progress) for the new implementation.

### Fixed
- **Breaking**: Set non-zero exit code if command execution fails.
- **Breaking**: Render servers IP addresses as array of objects, instead of previous pretty-printed string, in JSON and YAML outputs of `server show`.
- **Breaking**: Use key names from `json` field tag also in YAML output to have equal key names in JSON and YAML outputs. For example, `bootorder` key in server details will now be `boot_order` also in YAML output. As a side-effect data-types are limited to those supported by JSON. For example, timestamps will be presented as (double-quoted) strings. In addition, if command targets multiple resources, YAML output will now be a list, similarly than in JSON output, instead of previous multiple YAML documents.
- **Breaking**: In JSON and YAML output, `storage show` lists attached servers in `servers` list instead of `server` string.
- **Breaking**: In JSON and YAML output, `network show` lists DHCP DNS values in list instead of string.
- On `network show`, output server details as unknown instead of outputting an error, if fetching server details fails. This allows displaying network details for networks that contain a load balancer.
- Progress logging to non TTY output uses now 100 as text width instead of 0.

## [1.5.1] - 2022-07-15
### Fixed
- On `server create`, mount OS disk by default on `virtio` bus. Previously default OS storage address was not explicit and varyed depending on template type.
- Disable colors if user has set [NO_COLOR](https://no-color.org/) environment variable to non-empty value.

## [1.5.0] - 2022-07-05
### Added
- Add `--show-ip-addresses` flag to `server list` command to optionally include IP addresses in command output.
- Add `database connection list`, `database connection cancel`, `database start`, and `database stop` commands.

### Changed
- Make `--family` parameter of `server firewall create` command optional to allow editing the default rules.
- Update `cobra` to `v1.5.0` and refactor required flag validation code. This affects validation error messages.

### Fixed
- Complete shell input with uppercase letters (e.g., `Cap` to `CapitalizedName` will now work)
- Display UUID of created template in `storage templatise` output.

## [1.4.0] - 2022-06-15
### Added
- Add `database list`, `database show`,`database plans`, and `database types` commands.
- Add `loadbalancer list` and `loadbalancer show` commands.
- Add `db` and `lb` aliases to `database` and `loadbalancer`, respectively.

### Changed
- Color server state in `server list` output similarly than in `server show` output.
- Update Go version to 1.18
- Update `upcloud-go-api` to `v4.8.0`

## [1.3.0] - 2022-05-17
### Added
- Add `zone list` command that lists available zones.
- Add `--wait` flag to `server create` and `server stop` commands to wait until server is in `started` and `stopped` state, respectively.

### Changed
- Update `upcloud-go-api`to `v4.5.2`

### Fixed
- Do not display usage if execution fails because of missing credentials
- Mark error and warning livelogs finished when they will not be updated anymore: this stops the timer in the end of the row and stops livelog from refreshing these lines.

## [1.2.0] - 2022-04-29
### Added
- Include UUID (or address) of created resource in create command output
- `storage modify` command now accepts `enable-filesystem-autoresize` flag. When that flag is set upctl will attempt to resize partition and filesystem after storage size has been modified.

### Changed
- New go-api version v4.5.0

### Fixed
- Improved errors relating to argument resolver failures
- Print version info, instead of missing credentials error, when running `upctl version` without credentials
- Disable colors when outputting in JSON or YAML format
- Display both public and private addresses in `server create` output
- Render livelog messages of commands which execution takes less than render tick interval

## [1.1.3] - 2022-02-24
### Changed
- Update documentation

### Fixed
- Fix storage command attached-to-server key overrides zone

## [1.1.2] - 2022-01-21
### Fixed
- New release with no changes to fix [the Homebrew deprecation notice](https://github.com/goreleaser/goreleaser/pull/2591)

## [1.1.1] - 2021-09-30
### Changed
- Change password creation to be disabled by default in server creation
- Always create a password when password delivery method is chosen

## [1.1.0] - 2021-06-03
### Added
- Debug mode for finding root causes to problems
- Docker image and upload to Docker Hub
- Autogenerated documentation under /docs
- Terminal width handling

### Changed
- Global flags for colors refactored to `--force-colour` and `--no-colour`
- New go-api version
- Improved error messages

### Fixes
- No coloring of texts if stdout is not a terminal
- Detaching routers works now with the new `--detach-router` parameter in `network modify`

## [1.0.0] - 2021-04-16
First non-beta release! Includes all previous changes and fixes.

## [0.6.0] - 2021-04-16
### Changed
- Use goreleaser for releasing packages
- Move creation of service outside runcommand to facilitate testing

### Fixes
- fix(pre-commit): add missing golangci config file
- fix(root): use default cobra behaviour when called

## [0.5.0] - 2021-04-14
Initial public beta release :tada:

### Added
- Added commands for managing server firewall rules

### Fixed
- Shell auto-completion fixed

### Changed
- Input & output handling rewritten
- Several help text fixes & changes
- Default OS changed to Ubuntu 20.04
- Defaults for networks / routers in line with storages, so we'll show users' own resources by default


## [0.1.1] - 2021-03-12
### Fixed
- Load config files from the correct place on Windows
- Fix storage import failing on readerCounter not implementing io.Reader

## [0.1.0] - 2021-03-10

### Added
- Current feature set added! First internal release

[Unreleased]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.22.0...HEAD
[3.22.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.21.0...v3.22.0
[3.21.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.20.2...v3.21.0
[3.20.2]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.20.1...v3.20.2
[3.20.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.20.0...v3.20.1
[3.20.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.19.1...v3.20.0
[3.19.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.19.0...v3.19.1
[3.19.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.18.0...v3.19.0
[3.18.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.17.0...v3.18.0
[3.17.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.16.1...v3.17.0
[3.16.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.16.0...v3.16.1
[3.16.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.15.0...v3.16.0
[3.15.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.14.0...v3.15.0
[3.14.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.13.0...v3.14.0
[3.13.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.12.0...v3.13.0
[3.12.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.11.1...v3.12.0
[3.11.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.11.0...v3.11.1
[3.11.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.10.0...v3.11.0
[3.10.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.9.0...v3.10.0
[3.9.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.8.1...v3.9.0
[3.8.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.8.0...v3.8.1
[3.8.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.7.0...v3.8.0
[3.7.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.6.0...v3.7.0
[3.6.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.5.0...v3.6.0
[3.5.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.4.0...v3.5.0
[3.4.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.3.0...v3.4.0
[3.3.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.2.2...v3.3.0
[3.2.2]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.2.1...v3.2.2
[3.2.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.2.0...v3.2.1
[3.2.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.1.0...v3.2.0
[3.1.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v3.0.0...v3.1.0
[3.0.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.10.0...v3.0.0
[2.10.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.9.1...v2.10.0
[2.9.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.9.0...v2.9.1
[2.9.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.8.0...v2.9.0
[2.8.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.7.1...v2.8.0
[2.7.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.7.0...v2.7.1
[2.7.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.6.0...v2.7.0
[2.6.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.5.0...v2.6.0
[2.5.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.4.0...v2.5.0
[2.4.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.3.0...v2.4.0
[2.3.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.2.0...v2.3.0
[2.2.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.1.0...v2.2.0
[2.1.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.5.1...v2.0.0
[1.5.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.5.0...v1.5.1
[1.5.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.1.3...v1.2.0
[1.1.3]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.1.2...v1.1.3
[1.1.2]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v0.6.0...v1.0.0
[0.6.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/UpCloudLtd/upcloud-cli/compare/v0.1.1...v0.5.0
[0.1.1]: https://github.com/UpCloudLtd/upcloud-cli/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/UpCloudLtd/upcloud-cli/releases/tag/v0.1.0
