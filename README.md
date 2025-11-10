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

Once installed, navigate to your Git repository and run the command:

```bash
versify
```

The tool will analyze your commit history since the last tag and output the suggested next semantic version.

**Example Output:**
```
--- SemVer Version Bumper (Conventional Commits) ---
Last released version: v0.1.2
Analyzing commits since v0.1.2...

Determined BUMP: MINOR (Found 'feat:' commit)
Next suggested version: v0.2.0
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
          # Run versify and extract the version number from the output
          # The output format is "Next suggested version: vX.Y.Z"
          VERSION=$(versify | grep 'Next suggested version:' | awk '{print $4}')
          
          # Check if a new version was determined
          if [[ -z "$VERSION" || "$VERSION" == *"No-Change"* ]]; then
            echo "No version change detected."
            echo "::set-output name=version::"
          else
            echo "New version found: $VERSION"
            echo "::set-output name=version::$VERSION"
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
