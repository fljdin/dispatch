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

```sh
Usage:
  dispatch run [flags]

Flags:
  -f, --file string   file containing SQL statements
  -j, --jobs int      number of workers (default 2)
  -t, --type string   parser type (default "sh")

Global Flags:
  -c, --config string   configuration file
  -d, --dbname string   database name to connect to
      --help            show help
  -h, --host string     database server host or socket directory
  -l, --log string      log file
  -W, --password        force password prompt
  -p, --port int        database server port
  -U, --user string     database user name
  -v, --verbose         verbose mode
```

### Examples

```sh
cat <<EOF | psql -At > statements.sql
SELECT format('VACUUM ANALYZE %I.%I;', schemaname, relname)
  FROM pg_stat_user_tables WHERE last_analyze IS NULL
EOF

dispatch run -j 2 -f statements.sql
```

```text
2023/05/22 18:19:08 Worker 1 completed Task 0 (query #1) (success: true, elapsed: 12ms)
2023/05/22 18:19:08 Worker 2 completed Task 0 (query #0) (success: true, elapsed: 12ms)
```

## Command parsing

Internal parsers are used to load commands from `sh` or `psql` invocation. 

**`sh` rules**

* commands are newline-separated

**`psql` rules**

* queries are semicolon-separated or could be termined by a meta-command, like
  `\g` or `\gexec`
* transaction blocks and anonymous code blocks are detected as entire queries

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
- `uri`: connection string used by `psql`'s database option (`-d`)
- `connection`: connection name as described below, overrides `uri`
- `depends_on`: a list of identifiers of others tasks declared upstream

> All PostgreSQL environment variables can be used in place of `uri` as it used
> `psql` client.

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
    uri: postgresql://localhost
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

#### Generator task

A generator is a extended task that dispatch instructions from a result command
or a file.

- `generated`: in place of `command`
  - `from`: source execution context
  - `command`: instruction to be executed
  - `file`: instruction to be loaded from a file

```yaml
# run queries from a file simultaneously
tasks:
  - id: 1
    type: psql
    name: dispatch queries from a file
    generated:
      file: queries.sql
```

```yaml
# run queries generated by another query in parallel
tasks:
  - id: 1
    type: sh
    generated:
      from: psql
      command: |
        SELECT format('reindexdb -v -t %I;', tablename) FROM pg_tables 
        WHERE schemaname = 'public' AND tablename NOT IN ('log') LIMIT 10
```

### Traces

* `logfile`: summary of the tasks execution (default: disabled)
  - must be a valid path

```yaml
logfile: result.out
```

### Named connections

* `connections`: declares named connections used by tasks
  * `name`: connection name (`default` applied to any unattached tasks)
  * `uri`: a valid connection URI, takes precedence over following values
  * `host`: database server host or socket directory
  * `port`: database server port
  * `dbname`: database name to connect to
  * `user`: database user name
  * `password`: user password

```yaml
connections:
  - name: db
    uri: postgresql://remote
  - name: default
    host: localhost
    dbname: postgres
    user: postgres

tasks:
  - id: 1
    type: psql
    command: \conninfo
    connection: db
```

### Parallelism

* `workers`: declares number of workers
  - explicit argument passed to command takes precedence
  - limited by the number of logical CPUs usable by the main process

```yaml
workers: 1

# run the following tasks sequentially
tasks:
  - id: 1
    command: echo foo
  - id: 2
    command: echo bar
```

## Testing

[Bats](https://bats-core.readthedocs.io) testing framework is used. End-to-end
tests are located under `t/` directory A local PostgreSQL instance is required
with `postgres/postgres` authentication or `trust` method in `pg_hba.conf`

```sh
bats t
```

Unit tests are provided under `internal` packages.

```sh
go test ./...
```