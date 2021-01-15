# ksuid

ksuid is a Go library that generated prefixed, k-sorted globally unique identifiers.

Each ksuid has a resource type and optionally an environment prefix (no environment prefix is for production use only). They are roughly sortable down to per-second resolution.

Properties of a ksuid:

  - resource type and environment prefixing
  - lexicographically, time sortable
  - no startup co-ordination
  - guaranteed unique relative to process/machine

## Usage

### API

ksuid is primarily a Go package to be consumed by Cuvva services, below are examples of its API usage.

To generate a ksuid with a custom resource type and for the production environment:

```go
id := ksuid.Generate("user")
/* => ID{
	Environment: "prod",
	Resource: "user",
	Timestamp: time.Time{"2018-03-05T15:03:52Z"},
	MachineID: net.HardwareAddr{"78:4f:43:84:fd:b8"},
	ProcessID: 39089,
	SequenceID: 1,
} */
```

To parse a single given ksuid:

```go
id, err := ksuid.Parse([]byte("user_0EoZhc2lK5BSLogCVb7UYL"))
/*
=> ID{
	Environment: "prod",
	Resource: "user",
	Timestamp: time.Time{"2018-03-05T15:03:52Z"},
	MachineID: net.HardwareAddr{"78:4f:43:84:fd:b8"},
	ProcessID: 39089,
	SequenceID: 1,
}, nil
*/
```

### Command Line Tool

ksuid provides a helper utility to generate and parse ksuid on the command line, it contains two subcommands: `parse` and `generate`.

To generate two ksuid with a custom resource type and for the production environment:

```sh
$ ksuid generate --resource=user --count=2
user_0EoZhc2lK5BSLogCVb7UYL
user_0EoZhc2lK5BSLogCVb7UYM
```

To parse a single given ksuid:

```sh
$ ksuid parse user_0EoZhc2lK5BSLogCVb7UYL
ID:          user_0EoZhc2lK5BSLogCVb7UYL
Environment: prod
Resource:    user
Timestamp:   2018-03-05T15:03:52Z
Machine ID:  78:4f:43:84:fd:b8
Process ID:  39089
Sequence ID: 1
```

## How They Work

ksuid are minimum 22 bytes long when Base62 encoded, consisting of 16 bytes decoded:

  - a 32-bit unix timestamp with a custom epoch of 2014-01-01T00:00:00Z
  - the 48-bit MAC address of the primary interface on the generating machine
  - the 16-bit process id of the generating service
  - a 32-bit incrementing counter, reset every second

Optionally a ksuid has two, underscore delimited prefixes. The first prefix is optional, and is the environment in which the ksuid was generated (test, dev, git commit etc), omitting the environment identifies production only. The second prefix is the resource type (user, profile, vehicle etc) and is required.
