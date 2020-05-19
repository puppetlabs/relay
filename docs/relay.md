## relay

Relay by Puppet

### Synopsis

Relay connects your tools, APIs, and infrastructure
to automate common tasks through simple, event-driven workflows.

To get started, you'll need a relay.sh account - sign up for free
by following this link: 🔗 https://relay.sh/

Once you've signed up, run this to log in:
▶️   relay auth login

Use the 'workflow' subcommand to interact with workflows:
▶️   relay workflow


### Subcommand Usage

**`relay auth login [email] [flags]`** -- Log in to Relay
```
  -p, --password-stdin   accept password from stdin
```

**`relay auth logout`** -- Log out of Relay

**`relay doc generate`** -- Generate markdown documentation to stdout

**`relay workflow add [workflow name] [flags]`** -- Add a Relay workflow from a local file
```
  -f, --file string   Path to Relay workflow file
```

**`relay workflow delete [workflow name]`** -- Delete a Relay workflow

**`relay workflow download [workflow name] [flags]`** -- Download a workflow from the service
```
  -f, --file string   Filename to write workflow, relative to current working dir
```

**`relay workflow list`** -- Get a list of all your workflows

**`relay workflow replace [workflow name] [flags]`** -- Replace an existing Relay workflow
```
  -f, --file string   Path to Relay workflow file
```

**`relay workflow run [workflow name] [flags]`** -- Invoke a Relay workflow
```
  -p, --parameter stringArray   Parameters to invoke this workflow run with
```

### Global flags
```
  -d, --debug        print debugging information
  -o, --out string   output type: (text|json) (default "text")
  -y, --yes          skip confirmation prompts

```
