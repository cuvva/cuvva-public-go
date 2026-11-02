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
- `cdep u service avocado all -b fix-it` (every service in avocado whose build is ready)
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
- `--force`: deploy services even if their build isn't in ECR yet, see [Only services whose build is ready](#only-services-whose-build-is-ready)

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

### Only services whose build is ready

Each service runs one docker image (`go-services`, `go-services-playwright`, `js-services`,
`js-quote`, `web-underwriter`, ...), and each image has its own CI build that finishes in its own
time. Some images, like `go-services-playwright`, aren't built for every commit at all. Deploying a
service before its image exists leaves its pods unable to start.

So before changing any file, `update` and `update-default` check ECR for each service's image at
the commit being deployed (`master-<commit>`, or `branch-<commit>` for other branches). Each image is
looked up once, however many services use it.

| You ran... | If a service's build isn't in ECR yet |
| --- | --- |
| `cdep update-default services {env}` | It's skipped and listed. Run the same command again once CI has finished to pick it up. |
| `cdep update service {env} all` | Same: skipped and listed. |
| `cdep update service {env} sms quote` | cdep stops with `build_not_found` and changes nothing. |

For example, straight after a merge, while the JS builds are still running:

```
cdep update-default services avocado
WARN avocado: skipping 12 service(s) until js-services:master-abc123... is built, run cdep again once CI has finished: ...
```

The Go services go out now; running the command again later deploys the rest.

`--dry-run` shows what would be updated and skipped without committing anything. `--force` deploys
everything regardless, including when ECR can't be reached.

How cdep finds a service's image:

- `docker_image_name` in JSON config, or `image.repository` in YAML config.
- A service that doesn't set one inherits it from the env's `_base.json` / `_base.yaml`, as it does
  when deployed. A service with no image cdep can find is skipped under `all` and `update-default`.
- `update-default` treats `_base.json` like any other file: it's updated only when its own image is
  built. If a skipped service takes its commit from `_base.json`, it moves along with it; cdep warns
  and names it.
- `all` on `update` covers every service file in the env but never `_base` files; use
  `update-default` for those.

## Common errors

### `build_not_found`

A service you named hasn't got an image for the commit in ECR yet. Wait for its CI build to
finish, or pass `--force` to deploy anyway.

### `ecr_unavailable`

cdep couldn't look the builds up in ECR, usually because your AWS credentials have expired. Log in
again, or pass `--force` to deploy without checking.

### `all_requires_services`

`all` only works for services. For lambdas, terra and agentcore, list the items, or use
`update-default` to update every lambda.

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
