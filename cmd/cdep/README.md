# cdep (cuvva-deploy)

## Global flags

- `--verbose` (`-v`): verbose logging
- `--dry-run`: dry run only, do the file changes but don't commit or run anything

## Commands

### update

A CLI tool to deploy things to a specific environment. Covers both orchestrated services and lambdas handled by AWS.

This tool is designed to speed up quick day-to-day updates, and is not meant to be all encompassing.

`cdep update {type} {env|all} {...services|all}`

#### Updating All Services

You can now use `all` as the service name to update all services in an environment:

- `cdep update service avocado all` - Updates all services in avocado environment
- `cdep update service all all` - Updates all services in all environments

#### Service Filtering for Update Command

The `update` command now supports the same filtering options as `update-default`:

- `--go-only`: Only update Go services (services with `docker_image_name: go_services` or `go-services`)
- `--js-only`: Only update JS services (services with `docker_image_name` != `go_services`/`go-services`)

Examples:
- `cdep update service avocado all --go-only -c abc123` - Update only Go services in avocado
- `cdep update service all all --js-only -c abc123` - Update only JS services in all environments

For example:

- `cdep u service avocado -b fix-it sms ltm email web-underwriter`
- `cdep u service avocado all -c f1ec178befe6ed26ce9cec0aa419c763c203bc92`
- `cdep u service all sms`
- `cdep u lambda avocado -b fix-it marketing-consent stm-policy-sale`
- `cdep u service prod ltm --prod`
- `cdep u service avocado ltm -b fix-it`

Or with some flags

- `--prod`: work on prod, where `nonprod` is the default
- `--branch {name}` (`-b`): define a branch to use, where `master` is the default

### update-web

A CLI tool to deploy web applications on a specific system/environment.

`cdep update-web {type} {env|all} {...apps} <flags>`

For example:

- `cdep uw cf all website`
- `cdep uw cf all website -b new-ppc-landing`
- `cdep uw cf prod website --prod`

Flags are:

- `--prod`: work on the prod system, where `nonprod` is the default
- `--branch {name}` (`-b`): define a branch to use, where `master` is the default

### update-default

Specifically for fast forwarding the default commit on an environment, or all of them.

This command will find all services on the `master` branch and remove that from the individual config of that service, and instead update `_default.json` to the latest commit hash.

`cdep update-default {type} {env|all}`

#### Service Filtering

For services, you can filter by technology type:

- `--go-only`: Only update Go services (services with `docker_image_name: go_services` or `go-services`)
- `--js-only`: Only update JS services (services with `docker_image_name` != `go_services`/`go-services`)

**Note:** `_base.json` files are always updated regardless of filtering flags, as they are template/base configuration files that don't contain meaningful docker_image_name values for filtering.

For example

- `cdep update-default lambdas avocado`
- `cdep update-default services all`
- `cdep update-default services avocado --go-only`
- `cdep update-default services avocado --js-only`

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
