<p align="center">
  <img src="docs/relay-logo.svg" alt="Relay by Puppet" width="50%">
</p>

Relay is a service that lets you connect tools, APIs, and infrastructure to automate common tasks through simpler, smarter workflows. It links infrastructure events to workflow execution, so that for example, when a new JIRA ticket or GitHub issue comes in, your workflow can trigger deployments or send notifications.

This repo contains the source for the CLI tool which interacts with the Relay service and also provides the issue tracker for the product as a whole.

## Installation

You'll need an account on the service to use this tool. [Sign up here](https://app.relay.sh/signup). There's a [Getting Started guide](https://relay.sh/docs/getting-started/) to familiarize yourself with Relay concepts.

Once you're up and running, you can install the CLI a couple of different ways:

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

For more about workflows, check out the [Using Workflows](https://relay.sh/docs/using-workflows/) documentation.

## Build

To build run

```bash
./scripts/generate
./scripts/build
```

The resulting binaries will be in `./bin/relay-[version]-[architecture]`.

## Development

The CLI is built entirely using go. You can run locally with

```
go run ./cmd/relay
```

### Config

Relay uses [viper](https://github.com/spf13/viper) for customizable config. The following config values may be set in a yaml file at `$HOME/.config/relay/config.yaml` or as environment variables with corresponding names in all caps, prefixed with `RELAY_`:

- `debug`: Run Relay in debug mode. Overridden by global `--debug` flag.
- `out=(text|json)`: Output mode. Overridden by global `--out` flag.
- `yes`: Skip confirmation prompts. Overridden by global `--yes` flag.

## Getting help

If you have questions about Relay, you can [file a GitHub issue](https://github.com/puppetlabs/relay/issues) or join us on Slack in the [Puppet Community #relay channel](https://puppetcommunity.slack.com/archives/CMKBMAW2K).
