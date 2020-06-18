# Contributing to Relay

Relay welcomes contributions! Read on if you're interested in getting involved with the project.

## Where to Contribute

There are several routes to contribution in and around Relay.

* The [Relay Workflows](https://github.com/puppetlabs/relay-workflows) repository has a collection of useful workflows written by the Relay team. We welcome any improvements or bug fixes to them, and if you write something that could be of use to others, send in a pull request! 
* Add steps and triggers to existing integrations XXX
* If you want to create a Relay integration with a new service or tool, there's an [Integration developer guide](https://relay.sh/docs/integrating-with-relay/) to walk you through it. (Make sure there's not one [already in the works](https://github.com/relay-integrations/) first though!)
* The [Relay CLI](https://github.com/puppetlabs/relay) and related workflow development tools are great targets for contributors comfortable with the Go Language.
* Our [documentation site](https://github.com/puppetlabs/relay-docs) is open source and can always use improvement. Fixing anything in the docs from typos to better examples can be a great way to get involved with the project.
* For the brave of heart, the [core of the workflow execution system](https://github.com/puppetlabs/relay-core) is also open source.

## Guidelines for contributions

All interactions between Puppet employees, contributors, and community members on Relay-related projects are subject to [Puppet Community Code of Conduct](https://puppet.com/community/community-guidelines/).

Make sure there's not some existing code or a discussion that covers the change you want to make by searching existing Github issues. However, since much of Relay's development history happened prior to its public launch,  it's safest to [file a Github issue in the Relay project](https://github.com/puppetlabs/relay/issues) to chat with the team before starting any complex work.

To make it easier to contribute while still staying in the good graces of our (super wonderful!) Legal department, we require a [Developer Certificate of Origin](https://developercertificate.org/) sign-off on contributions. See [this explanation](https://helm.sh/blog/helm-dco/) from the Helm project to understand the rationale behind the DCO.As a practical matter, this means adding the `-s | --signoff` flag to your commits.


## Making Changes

* Clone the repository into your own namespace
* Create a topic branch from where you want to base your work.
  * To quickly create a topic branch based on `main`, run `git checkout -b fix/my_fix origin/main`.
* Make commits of logical and atomic units.
* Check for unnecessary whitespace with `git diff --check` before committing. 
* Make sure your commit messages are in the proper format. We (try to!) follow the [codelikeagirl guidelines](https://code.likeagirl.io/useful-tips-for-writing-better-git-commit-messages-808770609503) for writing good commit messages: format for short lines, use the imperative mood ("Add X to Y"), describe before and after state in the commit message body. Remember to add the `-s` flag to commits to DCO-sign them!
* Make sure you have added the necessary tests for your changes.
* Submit a pull request per the usual github PR process.


## Additional Resources

* [Puppet community guidelines](https://puppet.com/community/community-guidelines)
* [Puppet community slack](https://slack.puppet.com)
* [Relay issue tracker](https://github.com/puppetlabs/relay/issues)
* [General GitHub documentation](https://help.github.com/)
* [GitHub pull request documentation](https://help.github.com/articles/creating-a-pull-request/)
