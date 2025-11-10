# Versify

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Versify is a simple, yet powerful tool for Go projects that automates semantic versioning based on your commit messages. Spend less time tagging and more time coding.

## Features

- **Conventional Commits:** Leverages the [Conventional Commits](https://www.conventionalcommits.org/) specification to automatically determine version bumps.
- **Git-Powered:** Works seamlessly with your existing Git workflow. No extra configuration needed.
- **Simple & Fast:** A single binary that analyzes your commit history and suggests the next version in seconds.

## Usage

Using Versify is as simple as running a single command in your Go project's repository.

```bash
go run main.go
```

The tool will then analyze your commits and suggest the next version number.

### How it Works

The tool reads the Git history since the last tag and analyzes the commit messages. Based on the Conventional Commits specification, it determines the next version bump:

- `feat:` commits will result in a `minor` version bump.
- `fix:` commits will result in a `patch` version bump.
- Commits with `BREAKING CHANGE` in the body or `!` after the type will result in a `major` version bump.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


###
