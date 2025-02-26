# Releasing

1. Update `CHANGELOG.md`
2. Test GoReleaser config with `goreleaser check`
3. Tag a commit with the version you want to release e.g. `v1.2.3`
4. Push the tag & commit to GitHub
    - GitHub action automatically
      - sets the version based on the tag
      - creates a draft release to GitHub
      - populates the release notes from `CHANGELOG.md` with `make release-notes`
      - builds, uploads, and generates provenance for given release
5. Verify that [release notes](https://github.com/UpCloudLtd/upcloud-cli/releases) are in line with `CHANGELOG.MD`
6. Publish the drafted release
