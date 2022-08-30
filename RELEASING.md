# Releasing

1. Update `CHANGELOG.md`
2. Test GoReleaser config with `goreleaser check`
3. Tag a commit with the version you want to release e.g. `v1.2.3`
4. Push the tag & commit to GitHub
    - GitHub actions will automatically set the version based on the tag, create a GitHub release, build the project, and upload binaries & SHA sum to the GitHub release
5. [Edit the new release in GitHub](https://github.com/UpCloudLtd/upcloud-cli/releases) and add the changelog for this release