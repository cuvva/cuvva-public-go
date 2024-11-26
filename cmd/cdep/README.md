# cdep (cuvva-deploy)

## Global flags

- `--verbose` (`-v`): verbose logging
- `--dry-run`: dry run only, do the file changes but don't commit or run anything

## Commands

### update

A CLI tool to deploy things to a specific environment. Covers both orchestrated services and lambdas handled by AWS.

This tool is designed to speed up quick day-to-day updates, and is not meant to be all encompassing.

`cdep update {type} {env|all} {...services}`

For example:

- `cdep u service avocado -b fix-it sms ltm email web-underwriter`
- `cdep u service all sms`
- `cdep u lambda avocado -b fix-it marketing-consent stm-policy-sale`
- `cdep u service prod ltm --prod`
- `cdep u service avocado ltm -b fix-it`

Or with some flags

- `--prod`: work on prod, where `nonprod` is the default
- `--branch {name}` (`-b`): define a branch to use, where `master` is the default

### update-default

Specifically for fast forwarding the default commit on an environment, or all of them.

This command will find all services on the `master` branch and remove that from the individual config of that service, and instead update `_default.json` to the latest commit hash.

`cdep update-default {type} {env|all}`

For example

- `cdep update-default lambdas avocado`
- `cdep update-default services all`

## Common errors

### `no_mono_variable`

You need an environment variable (`CUVVA_CODE_REPO`) with the location of where you have the monorepo cloned.

`export CUVVA_CODE_REPO=/Users/duffleman/Source/cuvva`

### `Error: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain`

Due to the tool doing a git fetch/pull/push, you need to have your SSH keys in your SSH Agent otherwise it won't know which to use to contact GitHub.

You can see what identities are already linked by going into your shell and entering `ssh-add -l`. If none are there, you'll need to add them by entering `ssh-add`.

### `context deadline exceeded`

This error originates from the third party package we are using to list the references on the remote repository.
We are currently hard coding the context timeout to 30 seconds and at the time of writing this duration seems to be
working fine. However, if it takes longer in future (which shouldn't really happen) you can bump it in the
[monorepo.go](tools/cdep/git/monorepo.go) file where `remote.ListContext` function is called.
