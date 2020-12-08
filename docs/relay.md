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

**`relay completion`** -- Generate shell completion scripts

**`relay dev cluster delete`** -- Delete the local cluster

**`relay dev cluster start [flags]`** -- Start the local cluster that can execute workflows
```
      --image-registry-name string   The name to use on the host and on the cluster nodes for the container image registry (default "docker-registry.docker-registry.svc.cluster.local")
      --image-registry-port int      The port to use on the host and on the cluster nodes for the container image registry (default 5000)
      --load-balancer-port int       The port to map from the host to the service load balancer (default 8080)
      --worker-count int             The number of worker nodes to create on the cluster
```

**`relay dev cluster stop`** -- Stop the local cluster

**`relay dev image import <image:tag>`** -- Imports a container image into the cluster

**`relay dev kubectl`** -- Run kubectl commands against the dev cluster

**`relay dev metadata [flags]`** -- Run a mock metadata service
```
  -i, --input string   Path to metadata mock file
  -r, --run string     Run ID of step to serve (default "1")
  -s, --step string    Step name to serve (default "default")
```

**`relay dev workflow run [flags]`** -- Run a workflow on the dev cluster
```
  -f, --file string             Path to Relay workflow file
  -p, --parameter stringArray   Parameters to invoke this workflow run with
```

**`relay dev workflow secret set [flags]`** -- Set a workflow secret
```
      --value-stdin   accept secret value from stdin
```

**`relay doc generate [flags]`** -- Generate markdown documentation to stdout
```
  -f, --file string   The path to a file to write the documentation to
```

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

**`relay workflow secret delete [workflow name] [secret name]`** -- Delete a Relay workflow secret

**`relay workflow secret list [workflow name]`** -- List Relay workflow secrets

**`relay workflow secret set [workflow name] [secret name] [flags]`** -- Set a Relay workflow secret
```
      --value-stdin   accept secret value from stdin
```

**`relay workflow validate [flags]`** -- Validate a local Relay workflow file
```
  -f, --file string   Path to Relay workflow file
```

### Global flags
```
  -d, --debug        print debugging information
  -o, --out string   output type: (text|json) (default "text")
  -y, --yes          skip confirmation prompts

```
