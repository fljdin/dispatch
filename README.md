# Simple task dispatcher

[![go-test](https://github.com/fljdin/dispatch/actions/workflows/go-test.yml/badge.svg)](https://github.com/fljdin/dispatch/actions/workflows/go-test.yml)
[![go-e2e](https://github.com/fljdin/dispatch/actions/workflows/go-e2e.yml/badge.svg)](https://github.com/fljdin/dispatch/actions/workflows/go-e2e.yml)

Provides an easy-to-use command to dispatch tasks described in a YAML file.

Common use cases:

* Launching multiple elementary tasks in parallel
* Add a condition with a task dependent on another
* Split SQL files to execute statements as elementary tasks
* Behave as `\gexec` on multiple connections

## Usage

```text
Usage:
  dispatch [options]

Options:
  -c, --config=FILE      configuration file
  -h, --help             display this help and exit
  -o, --output=FILE      redirect output to file
  -P, --procs=(+)PROCS   number of processes (default 1)
  -v, --verbose          verbose mode
      --version          show version

The number of processes is limited to the number of CPU cores available
locally by default. In a remote execution context, where the number of
processes must not rely on the local machine, the sign "+" can be used to
by-pass this limitation. For example, "dispatch -P +16" will spawn 16
processes regardless of the number of CPU cores available locally.
```

## Configuration

Use a valid YAML file to describe tasks.

### Tasks declaration

* `tasks`: list of tasks to run
  * must be a valid array of tasks as described below

#### Elementary task

- `id` (required)
- `name`: as task description
- `type`: execution context in following choices
  + `sh` (default)
  + `psql`: needs PostgreSQL `psql` client to be installed
- `command`: instruction to be executed
- `env`: environment name as described below
- `variables`: a map of key-value used as environment variables, takes
  precedence over `env`
- `depends_on`: a list of identifiers of others tasks declared upstream

```yaml
# run the following shell commands simultaneously
tasks:
  - id: 1
    command: echo foo
  - id: 2
    command: echo bar
```

```yaml
# execute SQL statement with psql on localhost with default credentials
tasks:
  - id: 1
    type: psql
    name: run this statement
    command: SELECT user;
    variables:
      PGHOST: localhost
```

```yaml
# make a task dependent from another
tasks:
  - id: 1
    command: echo foo
  - id: 2
    command: echo bar
    depends_on: [1]
```

#### Loader tasks

A loader is an extended task that dispatch instructions from a result command or
a file. Delimiter detection is provided by [Fragment] package and only `PgSQL`
and `Shell` languages are supported.

[Fragment]: https://github.com/fljdin/fragment

To read and dispatch instructions from a file, use this:

- `file`: instructions to be loaded from a file

```yaml
# run queries from a file simultaneously
tasks:
  - id: 1
    type: psql
    name: dispatch queries from a file
    file: queries.sql
```

To dispatch commands from a specific result command, use the following
configuration:

- `loaded`: in place of `command`
  - `from`: source execution context
  - `command`: instruction to be executed
  - `env`: environment name as described below
  - `variables`: a map of key-value used as environment variables

```yaml
# run queries generated by another query in parallel
tasks:
  - id: 1
    type: sh
    name: execute reindexdb for all table except log
    loaded:
      from: psql
      command: |
        SELECT format('reindexdb -v -t %I;', tablename) FROM pg_tables
        WHERE schemaname = 'public' AND tablename NOT IN ('log')
```

### Named environments

* `environments`: declares named environment used by commands
  * `name`: environment name (`default` applied to all tasks)
  * `variables`: a map of key-value used as environment variables

```yaml
environments:
  - name: custom
    variables:
      PGHOST: remote.example.com
      PGUSER: alice
  - name: default
    variables:
      PGDATABASE: postgres

tasks:
  - id: 1
    name: Use variables, custom env and default env scopes
    env: custom
    variables:
      PGAPPNAME: my_app
```

### Parallelism

* `procs`: declares number of processes
  - option `--procs` takes precedence
* `remote`: defines the execution context
  - `false` (default): limit to the number of CPU cores available locally
  - `true`: no limit is applied to the number of processes

```yaml
procs: 1
remote: false

# run the following tasks sequentially
tasks:
  - id: 1
    command: echo foo
  - id: 2
    command: echo bar
```

### Traces

* `output`: redirect output to a file
* `verbose`: print debugging messages

```yaml
output: result.out
verbose: true
```

## Testing

[Bats](https://bats-core.readthedocs.io) testing framework is used. End-to-end
tests are located under `t/` directory A local PostgreSQL instance is required
with `postgres/postgres` authentication or `trust` method in `pg_hba.conf`

```sh
go build -tags testing
bats t

# or using
make test
```

Unit tests are provided under `internal` packages.

```sh
go test ./...
```
