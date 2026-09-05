# Contributing to ffrelayctl

Hello! Thank you for your interest in contributing to ffrelayctl!

This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:
- [Go](https://go.dev/doc/install) (version 1.27 or later, as set in `go.mod`)
- [Git](https://git-scm.com/install/)
- [Docker](https://docs.docker.com/get-started/get-docker/) (required for commit hooks)

Additionally, you'll need an account and API key from [Firefox Relay](https://relay.firefox.com).

## Development Setup

After cloning the repository, set up your development environment:

```bash
$ make setup
```

This command will:
- Check for required dependencies (Docker)
- Configure git hooks for commit message linting
- Set up the development environment

### Building

Build the project locally:

```bash
$ go build -o ffrelayctl .
```

### Running the CLI

You can run the built binary directly:

```bash
$ ./ffrelayctl --help
```

Or install it to your Go bin directory:

```bash
$ go install .
```

### Testing

Run the unit tests and `go vet` the way CI does:

```bash
$ make check
```

The individual targets are also available:

```bash
$ make vet
$ make test
```

CI runs the tests with the race detector on Linux and macOS, so it is worth
doing the same locally before opening a pull request:

```bash
$ go test -race ./...
```

To exercise the CLI against the live API, set a Firefox Relay API key as an
environment variable:

```bash
$ export FFRELAYCTL_KEY=replace-me
```

Then run commands manually:

```bash
$ ./ffrelayctl profiles list
$ ./ffrelayctl masks list
```

## Making Changes

### Creating a Branch

Create a new branch for your changes:

```bash
$ git checkout -b feature/feature-name
# or
$ git checkout -b fix/fix-name
```

### Code Style

- Follow standard Go conventions
- Format your code with `$ gofmt -w .`; CI fails if `$ gofmt -l .` reports anything
- Ensure your code passes `$ go vet ./...`
- Keep `go.mod` and `go.sum` tidy with `$ go mod tidy`
- Keep things [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

## Commit Guidelines

[Conventional Commits](https://www.conventionalcommits.org/) are enforced by commitlint.

### Commit Scopes

Common scopes in this project:
- `cli`: CLI interface and commands
- `api`: API client code
- `masks`: Email masks
- `phones`: Phone masks
- `profiles`: Profile management
- `contacts`: Contact management
- `users`: Account users
- `export`: Data export
- `ci`: GitHub Actions workflows
- `docker`: Dockerfile and image builds
- `deps`: Dependency updates
- `goreleaser`: GoReleaser configuration
- `dev`: Development tools, setup, experience
- `release`: Release process
- `readme`: README documentation

### Commit Examples

```bash
feat(masks): add support for filtering masks by status
fix(api): handle rate limiting errors correctly
docs(readme): update installation instructions
chore(deps): update dependencies to latest versions
```

## Pull Requests

1. Update your branch with the latest changes from upstream:
   ```bash
   $ git fetch upstream
   $ git rebase upstream/main
   ```

2. Ensure all commits follow the commit guidelines described above

3. Open a Pull Request on GitHub with:
   - A clear title following conventional commit format
   - A brief summary of the changes and why they're needed
   - Any related issues or resources for reference

### PR Requirements

- All commits must follow conventional commit format
- Code must be properly formatted (`gofmt -l .` reports nothing)
- No new warnings from `go vet ./...`
- Tests must pass on Linux, macOS, and Windows
- `go mod tidy` must leave `go.mod` and `go.sum` unchanged
- The code must cross-compile and pass `govulncheck`
- The Docker image must still build
- The PR should focus on a single feature or fix
- Keep PRs reasonably sized for easier review

CI enforces all of the above on every push and pull request.

## Reporting Issues

Feel free to report issues (bugs, feature requests, etc) using GitHub.

## License

By contributing to ffrelayctl, you agree that your contributions will be licensed under the [MIT License](LICENSE).
