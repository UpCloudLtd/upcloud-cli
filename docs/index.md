# Getting started

`upctl` provides a command-line interface to [UpCloud](https://upcloud.com/) services. It allows you
to control your resources from the command line or any compatible interface.

## Install upctl

`upctl` can be installed from a pre-built package in the repositories [GitHub releases](https://github.com/UpCloudLtd/upcloud-cli/releases) using a package manager or from sources with `go install`:

=== "Linux"

    Download pre-built package from GitHub releases and install it with your package manager.

    On Ubuntu or Debian, use the `.deb` package.

    ```sh
    curl -Lo upcloud-cli_{{ latest_release }}_amd64.deb https://github.com/UpCloudLtd/upcloud-cli/releases/download/v{{ latest_release }}/upcloud-cli_{{ latest_release }}_amd64.deb
    sudo apt install ./upcloud-cli_{{ latest_release }}_amd64.deb
    ```

    On RHEL based distributions, use the `.rpm` package.

    ```sh
    curl -Lo upcloud-cli-{{ latest_release }}-1.x86_64.rpm https://github.com/UpCloudLtd/upcloud-cli/releases/download/v{{ latest_release }}/upcloud-cli-{{ latest_release }}-1.x86_64.rpm
    sudo dnf install ./upcloud-cli-{{ latest_release }}-1.x86_64.rpm
    ```

=== "macOS"

    Use homebrew to install `upctl` from [UpCloudLtd tap](https://github.com/UpCloudLtd/homebrew-tap).

    ```sh
    brew tap UpCloudLtd/tap
    brew install upcloud-cli
    ```

=== "Windows"

    First, download the archived binary from GitHub releases to current folder and extract the binary from the archive.

    ```pwsh
    Invoke-WebRequest -Uri "https://github.com/UpCloudLtd/upcloud-cli/releases/download/v{{ latest_release }}/upcloud-cli_{{ latest_release }}_windows_x86_64.zip" -OutFile "upcloud-cli_{{ latest_release }}_windows_x86_64.zip"
    Expand-Archive -Path "upcloud-cli_{{ latest_release }}_windows_x86_64.zip"

    # Print current location
    Get-Location
    ```

    Then, close the current PowerShell session and open a new session as an administrator. Move the binary to `upcloud-cli` folder in _Program Files_, add the `upcloud-cli` folder in _Program Files_ to `Path`.

    ```pwsh
    # Open the PowerShell with Run as Administrator option.
    # Use Set-Location to change into folder that you used in previous step.

    New-Item -ItemType Directory $env:ProgramFiles\upcloud-cli\ -Force
    Move-Item -Path upcloud-cli_{{ latest_release }}_windows_x86_64\upctl.exe -Destination $env:ProgramFiles\upcloud-cli\ -Force

    # Setting the Path is required only on first install.
    # Thus, this step can be skipped when updating to a more recent version.
    [Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::Machine) + ";$env:ProgramFiles\upcloud-cli\", [EnvironmentVariableTarget]::Machine)
    ```

    After running the above commands, close the administrator PowerShell session and open a new PowerShell session to verify installation succeeded.

=== "`go install`"

    Install the latest version of `upctl` with `go install`, by running:

    ```sh
    go install github.com/UpCloudLtd/upcloud-cli/v3/...@latest
    ```

=== "Docker"

    Pull the latest build from GHCR by running:

    ```sh
    docker pull ghcr.io/upcloudltd/upctl:latest
    ```

---

After installing `upctl`, you can run `upctl version` command to verify that the tool was installed successfully.

```sh
upctl version
```

### Verify assets

Project release packages' SHA-256 checksums are available in the project releases,
files `checksums.txt`. They can be checked for example with:

```sh
# make sure at least one downloaded package and checksums.txt are in the current directory
sha256sum -c --ignore-missing checksums.txt
```

Project release checksum files and Docker images are signed using the
[Sigstore framework](https://www.sigstore.dev), and can be verified with
[cosign](https://docs.sigstore.dev/cosign/).

For example, to verify package checksums file for a release:

```sh
project_url=https://github.com/UpCloudLtd/upcloud-cli
release=vX.Y.Z
cosign verify-blob \
    --certificate-identity ${project_url}/.github/workflows/publish.yml@refs/tags/${release} \
    --certificate-oidc-issuer https://token.actions.githubusercontent.com \
    --certificate ${project_url}/releases/download/${release}/checksums.txt-keyless.pem \
    --signature ${project_url}/releases/download/${release}/checksums.txt-keyless.sig \
    ${project_url}/releases/download/${release}/checksums.txt
```

And to verify the Docker image for a release:

```sh
project=UpCloudLtd/upcloud-cli
release=vX.Y.Z
cosign verify \
   ghcr.io/${project}:${release} \
   --certificate-identity https://github.com/${project}/.github/workflows/publish.yml@refs/tags/${release} \
   --certificate-oidc-issuer https://token.actions.githubusercontent.com |
   jq .
```

### Configure shell completions

`upctl` provides shell completions for multiple shells. Run `upctl completion --help` to list the supported shells.

```sh
upctl completion --help
```

To configure the shell completions, follow the instructions provided in the help output of the command matching the shell you are using. For example, if you are using zsh, run `upctl completion zsh --help` to print the configuration instructions.

#### Bash completions

On bash, the completions depend on `bash-completion` package. Install and configure the package according to your OS:

=== "Linux"

    First, install `bash-completion` package, if it has not been installed already, and add command to source the completions to your `.bashrc`.

    On Ubuntu or Debian, use `apt` command to install the package:

    ```sh
    sudo apt install bash-completion
    echo "[ -f /etc/bash_completion ] && . /etc/bash_completion" >> ~/.bashrc
    ```

    On RHEL based distributions, use `dnf` command to install the package.

    ```sh
    sudo dnf install bash-completion
    echo "[ -f /etc/bash_completion ] && . /etc/bash_completion" >> ~/.bashrc
    ```

    Finally, configure the shell completions for `upctl` by either sourcing `upctl completion bash` output in your bash `.bashrc` or by saving the output of that command in `upctl` file under `/etc/bash_completion.d/`:

    ```sh
    # First alternative
    echo 'source <(upctl completion bash)' >>~/.bashrc

    # Second alternative
    upctl completion bash | sudo tee /etc/bash_completion.d/upctl > /dev/null

    # Source completions to current shell session
    . /etc/bash_completion
    ```

=== "macOS"

    First, install `bash-completion` package, if it has not been installed already, and add command to source the completions to your `.bash_profile`.

    ```sh
    brew install bash-completion
    echo '[ -f "$(brew --prefix)/etc/bash_completion" ] && . "$(brew --prefix)/etc/bash_completion"' >> ~/.bash_profile
    ```

    Then configure the shell completions for `upctl` by saving the output of `upctl completion bash` in `upctl` file under `/etc/bash_completion.d/`:

    ```sh
    upctl completion bash > $(brew --prefix)/etc/bash_completion.d/upctl
    . $(brew --prefix)/etc/bash_completion
    ```

## Configure credentials

To be able to manage your UpCloud resources, you need to configure credentials for `upctl` and enable API access for these credentials.

Define the credentials by setting `UPCLOUD_USERNAME` and `UPCLOUD_PASSWORD` environment variables.

API access can be configured in the UpCloud Hub on [Account page](https://hub.upcloud.com/account/overview) for the main-account and on the [Permissions tab](https://hub.upcloud.com/people/permissions) of the People page for sub-accounts. We recommend you to set up a sub-account specifically for the API usage with its own username and password, as it allows you to assign specific permissions for increased security.

## Execute commands

To verify you are able to access the UpCloud API, you can, for example, run `upctl account show` command to print your current balance and resource limits.

```sh
upctl account show
```

For usage examples, see the _Examples_ section of the documentation.

For reference on how to use each sub-command, see the _Commands reference_ section of the documentation. The same information is also available in `--help` output of each command.
