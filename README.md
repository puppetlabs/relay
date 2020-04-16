<p align="center">
  <img src="docs/relay_logotype.png" alt="Relay by Puppet" width="50%">
</p>

Relay is a service that lets you connect tools, APIs, and infrastructure to automate common tasks through simpler, smarter workflows. It links infrastructure events to workflow execution, so that for example, when a new JIRA ticket or Github issue comes in, your workflow can trigger deployments or send notifications.

This repo contains the source for the CLI tool which interacts with the Relay service and also provides the issue tracker for the product as a whole. 

## Installation

Relay evolved from an incubation project at Puppet called Project Nebula, and some of the tooling and documentation still say "Nebula" while we get everything switched over.

You'll need an account on the service to use this tool. [Sign up here](https://puppet.com/products/project-nebula#nebula-form)!

Once you're up and running, you can install the CLI a couple of different ways:

For Macs, install via homebrew:

```bash
brew install puppetlabs/puppet/relay
```

For other platforms, install directly via Github Releases:

[Get the latest version](https://github.com/puppetlabs/relay/releases)

The program is just a single binary, so you can simply download the one that matches your architecture and copy it to a location in your `$PATH`. Note these binaries are named 'nebula', for the time being simply rename it for consistency:

```bash
mv ./nebula-v3.4.0-linux-arm64 /usr/local/bin/relay
```

## Getting started

Once it's installed, you'll need to authenticate with the service, then you'll be able to work with the default set of workflows that are enabled on your account:

```bash
relay login
relay workflow list
```

For more about workflows and further onboarding information, check out the [documentation website](https://puppet.com/docs/nebula/beta/overview.html)

## Build

To build run

```bash
./scripts/generate
./scripts/build
```

The resulting binaries will be in `./bin/relay`.

## Development

The CLI is built entirely using go. You can run locally with

```
go run ./cmd/relay
```

### Config

Relay uses [viper](https://github.com/spf13/viper) for customizable config. The following config values may be set in a yaml file at `$HOME/.config/relay/config.yaml` or as environment variables with corresponding names in all caps, prefixed with `RELAY_`:
* `debug`: Run relay in debug mode. Overrides global `--debug` flag.
* `out=(text|json)`: Output mode. Overrides global `--out` flag.
* `api_domain`: Relay api domain to connect to for all api operations.
* `ui_domain`: Relay ui domain, mainly used in generated links.
* `web_domain`: Relay web domain, mainly used in generated links.
* `cache_dir`: Cache directory.
* `token_path`: Path to token storage location.

## Getting help

If you have questions about Relay, you can [file a github issue](https://github.com/puppetlabs/relay/issues) or join us on the [Puppet community slack](https://slack.puppet.com) in #relay. 
