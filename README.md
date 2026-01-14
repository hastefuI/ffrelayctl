# ffrelayctl [![Build](https://github.com/hastefuI/ffrelayctl/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/hastefuI/ffrelayctl/actions/workflows/ci.yml) [![Release](https://img.shields.io/github/v/release/hastefuI/ffrelayctl)](https://github.com/hastefuI/ffrelayctl/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/hastefuI/ffrelayctl)](https://goreportcard.com/report/github.com/hastefuI/ffrelayctl) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/hastefuI/ffrelayctl/blob/main/LICENSE)

A CLI for [Firefox Relay](https://relay.firefox.com) written in Go.

## Overview
Firefox Relay is a privacy service from Mozilla that provides email and phone number masks to help keep your identity private.

`ffrelayctl` is a command-line tool for managing Firefox Relay masks directly from the terminal.

<img src="./demo.gif" alt="Demo" style="width:100%; max-width:900px;" />

## Features
- **Contact Management**: Manage inbound contacts (premium only)
- **Email Mask Management**: Manage both random and custom domain email masks
- **Phone Management**: Manage phone masks and forwarding number (premium only)
- **Profile Management**: View your Relay profile and subscription status
- **Data Export**: Export your Relay data for backup purposes

## Installation

### Homebrew (macOS and Linux)

```bash
$ brew tap hastefui/tap
$ brew install ffrelayctl
```

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

### Docker

```bash
$ docker build -t ffrelayctl .
```

### Verify Installation

Verify that the installation for ffrelayctl was successful:
```bash
$ ffrelayctl --version
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
  contacts list                  # List phone contacts (premium only)
  contacts update                # Update a phone contact (premium only)
  masks list                     # List all masks
  masks get                      # Get a mask
  masks create                   # Create a new mask
  masks update                   # Update a mask
  masks delete                   # Delete a mask
  phones list                    # List phone masks (premium only)
  phones discover                # Discover phone masks available (premium only)
  phones search                  # Search phone masks by area code (premium only)
  phones update                  # Update a phone mask (premium only)
  phones forward list            # List forwarding numbers (premium only)
  phones forward get             # Get forwarding number (premium only)
  phones forward register        # Register forwarding number (premium only)
  phones forward verify          # Verify forwarding number (premium only)
  phones forward delete          # Delete forwarding number (premium only)
  profiles list                  # List available Relay profiles
  users list                     # List users for Relay account
  export                         # Export all Firefox Relay account data

Use "ffrelayctl [command] --help" for more information about a command.
```

## Examples

```bash
# Fetch the custom domain in use (premium only)
$ ffrelayctl profiles list | jq '.[].subdomain'

# Generate a random mask
$ ffrelayctl masks create --description "GitHub" --generated-for "github.com"

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

# List phone masks
$ ffrelayctl phones list | jq

# List all phone numbers that have texted your Relay number
$ ffrelayctl contacts list | jq '[.[] | select(.last_inbound_type == "text")]'

# List all masks using Docker
$ docker run --rm -e FFRELAYCTL_KEY=<replace-me> ffrelayctl profiles list
```

## Development

### Setup

After cloning this repository, run:
```bash
$ make setup
```

## Disclaimer

This is an unofficial CLI not affiliated with or endorsed by Mozilla or Firefox Relay.

## License

Licensed under [MIT License](https://opensource.org/licenses/MIT), see [LICENSE](./LICENSE) for details.

Copyright (c) 2026 hasteful.
