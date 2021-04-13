# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
Initial public release :tada:

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

[Unreleased]: https://github.com/UpCloudLtd/upcloud-cli/compare/v0.1.1...HEAD
[0.1.0]: https://github.com/UpCloudLtd/upcloud-cli/releases/tag/v0.1.0
[0.1.1]: https://github.com/UpCloudLtd/upcloud-cli/releases/tag/v0.1.1
