# cdep (cuvva-deploy)

## Global flags

- `--verbose` (`-v`): verbose logging
- `--dry-run`: dry run only, do the file changes but don't commit or run anything

## Commands

### update

A CLI tool to deploy things to a specific environment. Covers both orchestrated services and lambdas handled by AWS.

This tool is designed to speed up quick day-to-day updates, and is not meant to be all encompassing.

`cdep update {type} {env|all} {...services|all}`

Supported types are `service`, `lambda`, `terra` and `agentcore`.

For example:

- `cdep u service avocado -b fix-it sms ltm email web-underwriter`
- `cdep u service all sms`
- `cdep u service avocado all -b fix-it` (every service in avocado)
- `cdep u service avocado all --go-only -b fix-it` (every Go service in avocado)
- `cdep u lambda avocado -b fix-it marketing-consent stm-policy-sale`
- `cdep u service prod ltm --prod`
- `cdep u service avocado ltm -b fix-it`

#### agentcore

Agentcore config lives under the special `_system` environment (`{system}/_system/agentcore`),
so it is always deployed against `_system` rather than a per-service environment. It also tracks
the `main` branch by default instead of `master`, so `-b` is only needed to deploy another branch.

Unlike services and lambdas (which build from the monorepo), each agent lives in its own repo.
Every agentcore config file must therefore declare its source repo, and cdep resolves the latest
commit for the requested branch from that repo — no local clone required:

```json
{
	"repo": "git@github.com:example-org/example-agent.git",
	"branch": "main",
	"commit": "0000000000000000000000000000000000000000"
}
```

- `cdep u agentcore _system my-agent other-agent` (each resolves its own repo's `main` tip)
- `cdep u agentcore _system my-agent --prod`
- `cdep u agentcore _system my-agent -b my-feature`
- `cdep u agentcore _system my-agent -c <40-char-hash>` (pin one agent to a specific commit)

Because each agent has its own repo, a pinned `-c` can only target a single agent at a time.

Or with some flags

- `--prod`: work on prod, where `nonprod` is the default
- `--branch {name}` (`-b`): define a branch to use, where `master` is the default
- `--go-only` / `--js-only`: services only, see [Go only or JS only](#go-only-or-js-only)

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

For example

- `cdep update-default lambdas avocado`
- `cdep update-default services all`
- `cdep update-default services all --go-only`
- `cdep update-default services avocado --js-only`

### Go only or JS only

A change to a shared Go or JS library usually means redeploying every service in that language, and
the Go and JS builds finish at different times. `--go-only` and `--js-only` let you roll out one
language without waiting for, or touching, the other. They work with `update` and `update-default`,
for services only.

| You want to... | Run |
| --- | --- |
| Move every Go service in avocado to the latest master | `cdep update-default services avocado --go-only` |
| Move every JS service in every nonprod env to a commit | `cdep update-default services all --js-only -c <40-char-hash>` |
| Put every Go service in avocado on your branch | `cdep update service avocado all --go-only -b my-branch` |
| Only update the Go services out of a list | `cdep update service avocado sms quote web-underwriter --go-only` |

Add `--dry-run` first to see exactly what would change: cdep lists the services it will update and
the ones it skipped, per environment.

How cdep decides:

- A service's language comes from the docker image it runs: `docker_image_name` in JSON config, or
  `image.repository` in YAML config. A service that doesn't set one inherits it from the env's
  `_base.json` / `_base.yaml`, as it does when deployed.
- Images starting `go-services` (including `go-services-playwright`) are Go. Every other image
  (`js-services`, `js-quote`, `web-underwriter`, ...) is JS.
- `update-default` also updates `_base.json` when its image matches the filter. If a service of the
  other language takes its commit from that `_base.json`, it moves too; cdep warns and names it.
- `all` on `update` covers every service file in the env but never `_base` files; use
  `update-default` for those.
- cdep checks ECR once per image (e.g. once for `go-services`) and warns if the build for the
  commit isn't there yet.

## Common errors

### `conflicting_flags`

`--go-only` and `--js-only` can't be combined. Leave both off to update every service.

### `filter_requires_services` / `all_requires_services`

`--go-only`, `--js-only` and `all` only work for services. For lambdas, terra and agentcore, list
the items, or use `update-default` to update every lambda.

### `all_must_be_alone`

`all` already selects every service, so don't add service names after it.

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
