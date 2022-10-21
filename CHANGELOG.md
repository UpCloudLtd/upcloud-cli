# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Complete available types for `database plans`.
- Suppress positional argument filename completion for commands without specific completions.
- In `database list` output: if database has no title, database name is displayed in the title cell instead of leaving the cell empty, similarly than in the hub.

### Fixed
- Remove debug leftover print from IP address completions.

## [2.2.0] - 2022-10-17
### Added
- Include runtime operating system and architecture in `version` command output.
- Include instructions for defining credentials and API access in `upctl --help` output.

### Fixed
- Fix commands in `server delete` usage examples.
- Tune human output so that normal output is directed to `stdout`. Progess and error messages are still outputted to `stderr`.

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
- Print version info, instead of missing credentials error, when runnning `upctl version` without credentials
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
- fix(root): use default cobra behavious when called

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

[Unreleased]: https://github.com/UpCloudLtd/upcloud-cli/compare/v2.2.0...HEAD
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
