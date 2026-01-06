# ffrelayctl

A CLI for [Firefox Relay](https://relay.firefox.com) written in Go.

## Overview
Firefox Relay is a privacy service from Mozilla that provides email masks to help keep your identity private.

`ffrelayctl` is a command-line tool for managing Firefox Relay masks directly from the terminal.

<img src="./demo.gif" alt="Demo" style="width:100%; max-width:900px;" />

## Features
- **Profile Management**: View your Relay profile and subscription status
- **Mask Management**: Manage both random and custom domain email masks

## Installation

### Pre-built Binaries

Download, extract, and install the latest release for your platform from the [releases page](https://github.com/hastefuI/ffrelayctl/releases).

### Build From Source

Clone this repository, and then:
```bash
# Build
$ go build -o ffrelayctl .

# Install
$ go install .
```

### Verify Installation

Verify that the installation for ffrelayctl was successful:
```bash
$ ffrelayctl --version
ffrelayctl version x.x.x
```

## Quick Start

### Prerequisites

A Firefox Relay account is required.

### Authenticating with Firefox Relay

To use `ffrelayctl`, you need to authenticate with Firefox Relay by providing a Relay API Key, which can be retrieved from the [Firefox Relay Settings Page](https://relay.firefox.com/accounts/settings/) after login.

Verify you're able to retrieve your Relay profile using the API Key for your account:
```bash
$ ffrelayctl profiles list --key <replace-me>
# or
$ export FFRELAYCTL_KEY=<replace-me> && ffrelayctl profiles list
```

## Usage

```bash
ffrelayctl is A CLI for Firefox Relay.

Usage:
  ffrelayctl [command]

Available Commands:
  help                           # Display help for any command
  profiles list                  # List your Relay profile(s)
  masks list                     # List all masks
  masks get                      # Get a mask
  masks create                   # Create a new mask
  masks update                   # Update a mask
  masks delete                   # Delete a mask

Use "ffrelayctl [command] --help" for more information about a command.
```

## Examples

```bash
# Generate a random mask
$ ffrelayctl masks create --description "GitHub" --generated-for "github.com"

# Fetch the custom domain in use (premium only)
$ ffrelayctl profiles list | jq '.[].subdomain'

# List all enabled masks
$ ffrelayctl masks list | jq '.[] | select(.mask.enabled == true)'

# List email addresses in use by all masks
$ ffrelayctl masks list | jq '.[].mask.full_address'

# List all masks containing "newsletter" in the description
$ ffrelayctl masks list | jq '.[] | select(.mask.description | test("newsletter"; "i"))'

# Count total forwarded emails from random masks
$ ffrelayctl masks list --random=true | jq '[.[].num_forwarded] | add'

# Count total masks
$ ffrelayctl masks list | jq '.[].mask.id' | wc -l
```

## Disclaimer

This is an unofficial CLI not affiliated with or endorsed by Mozilla or Firefox Relay.

## License

Licensed under [MIT License](https://opensource.org/licenses/MIT), see [LICENSE](./LICENSE) for details.

Copyright (c) 2026 hasteful.
