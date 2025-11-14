# Versify

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Versify is a simple, yet powerful tool for Go projects that automates semantic versioning based on your commit messages. Spend less time tagging and more time coding.

## Features

- **Conventional Commits:** Leverages the [Conventional Commits](https://www.conventionalcommits.org/) specification to automatically determine version bumps.
- **Git-Powered:** Works seamlessly with your existing Git workflow. No extra configuration needed.
- **Simple & Fast:** A single binary that analyzes your commit history and suggests the next version in seconds.

## Installation

You can download the pre-compiled binary for your OS from the [**latest release page**](https://github.com/OpScaleHub/versify/releases/latest) or use the commands below.

### Linux & macOS

```bash
# Determine OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
fi

# Download the latest binary
curl -L "https://github.com/OpScaleHub/versify/releases/latest/download/versify-${OS}-${ARCH}" -o versify

# Make it executable
chmod +x versify

# Move it to a directory in your PATH (optional, requires sudo)
sudo mv versify /usr/local/bin/
```

### Windows (PowerShell)

```powershell
# Determine Architecture
$ARCH = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
  $ARCH = "arm64"
}

# Download the latest binary
$url = "https://github.com/OpScaleHub/versify/releases/latest/download/versify-windows-${ARCH}.exe"
$output = "versify.exe"
Invoke-WebRequest -Uri $url -OutFile $output

# Move it to a directory in your PATH (optional, requires admin rights)
# Example: Move-Item -Path ".\versify.exe" -Destination "C:\Windows\System32"
```

## Usage

Once installed, navigate to your Git repository and run the command with the desired flags:

```bash
versify [flags]
```

The tool analyzes your commit history since the last tag and outputs the suggested next semantic version to **stdout**, while logging its progress to **stderr**.

### Command-Line Flags

| Flag              | Description                                                                                     | Default      |
| ----------------- | ----------------------------------------------------------------------------------------------- | ------------ |
| `--prefix`        | The prefix for the version tag (e.g., `v`, `k8s`).                                              | `v`          |
| `--baseline`      | A specific version to use as the starting point, instead of the latest Git tag.                 | `(empty)`    |
| `--add-suffix`    | Always add a suffix to the new version (e.g., for pre-releases or build metadata).              | `false`      |
| `--suffix-format` | The format for the suffix when `--add-suffix` is used. Options: `short-hash`, `datetime`.       | `short-hash` |

### Example

```bash
# Get the next version with a 'v' prefix
versify
# v0.2.0

# Get the next version with a custom prefix and a datetime suffix
versify --prefix "k8s-" --add-suffix --suffix-format datetime
# k8s-0.2.0-20251114103000
```

**Example Output (stderr):**
```
--- SemVer Version Bumper (Conventional Commits) ---
Last released version: v0.1.2
Analyzing commits since v0.1.2...

Determined BUMP: MINOR (Found 'feat:' commit)
```

The final version is printed to **stdout**:
```
v0.2.0
```

### How it Works

- **`feat:`** commits result in a **minor** version bump (e.g., `v1.1.0` -> `v1.2.0`).
- **`fix:`** commits result in a **patch** version bump (e.g., `v1.1.0` -> `v1.1.1`).
- Commits with **`BREAKING CHANGE`** in the body or `!` after the type (e.g., `feat!:`) result in a **major** version bump (e.g., `v1.1.0` -> `v2.0.0`).
- Other commit types (`docs:`, `chore:`, `style:`, etc.) do not trigger a version bump.

## CI/CD Integration

Automate your release process by using Versify in your pipeline.

### GitHub Actions Example

This workflow runs on every push to the `main` branch. It calculates the next version and, if a version bump is needed, creates a new Git tag and a GitHub Release.

```yaml
name: Create Release

on:
  push:
    branches:
      - main

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3
        with:
          # Fetch all history for accurate version analysis
          fetch-depth: 0

      - name: Download and Install Versify
        run: |
          curl -L "https://github.com/OpScaleHub/versify/releases/latest/download/versify-linux-amd64" -o versify
          chmod +x versify
          sudo mv versify /usr/local/bin/

      - name: Get Next Version
        id: get_version
        run: |
          # Run versify and capture the new version from stdout
          VERSION=$(versify)
          
          # Check if a new version was determined
          if [[ -z "$VERSION" ]]; then
            echo "No version change detected."
            echo "version=" >> $GITHUB_OUTPUT
          else
            echo "New version found: $VERSION"
            echo "version=$VERSION" >> $GITHUB_OUTPUT
          fi

      - name: Create Git Tag and GitHub Release
        if: steps.get_version.outputs.version
        uses: softprops/action-gh-release@v1
        with:
          tag_name: ${{ steps.get_version.outputs.version }}
          name: Release ${{ steps.get_version.outputs.version }}
          body: "See CHANGELOG for details."
          draft: false
          prerelease: false
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
