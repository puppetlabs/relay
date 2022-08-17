<p align="center">
  <img src="docs/relay-logo.svg" alt="Relay by Puppet" width="50%">
</p>

Relay is a service that lets you connect tools, APIs, and infrastructure to automate common tasks through simpler, smarter workflows. It links infrastructure events to workflow execution, so that for example, when a new JIRA ticket or GitHub issue comes in, your workflow can trigger deployments or send notifications.

This repo contains the source for the CLI tool.

## Installation

For Macs, install via homebrew:

```bash
brew install puppetlabs/puppet/relay
```

For other platforms, install directly via GitHub Releases:

[Get the latest version](https://github.com/puppetlabs/relay/releases)

The program is just a single binary, so you can simply download the one that matches your architecture and copy it to a location in your `$PATH`.

```bash
mv ./relay-v4*-linux-arm64 /usr/local/bin/relay
```

## Getting started

Once it's installed, you'll need to authenticate with the service, then you'll be able to work with the default set of workflows that are enabled on your account:

```bash
relay auth login
relay workflow list
```

### Config

Relay uses [viper](https://github.com/spf13/viper) for customizable config. The following config values may be set in a yaml file at `$HOME/.config/relay/config.yaml` or as environment variables with corresponding names in all caps, prefixed with `RELAY_`:

- `debug`: Run Relay in debug mode. Overridden by global `--debug` flag.
- `out=(text|json)`: Output mode. Overridden by global `--out` flag.
- `yes`: Skip confirmation prompts. Overridden by global `--yes` flag.
