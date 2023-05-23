# Simple task dispatcher

[![go-test](https://github.com/fljdin/dispatch/actions/workflows/go-test.yml/badge.svg)](https://github.com/fljdin/dispatch/actions/workflows/go-test.yml)

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

## Query parsing

An internal parser is used to load semicolon-separated queries as `psql`'s
tasks. It provides correct detection of transaction blocks and anonymous code
blocks.

## Configuration

Use a valid YAML file to describe tasks.

### Tasks declaration

* `tasks`: list of tasks to run
  - `id` (required)
  - `command` or `file`: instruction(s) to be executed or loaded from a file
  - `name`: as task description
  - `type`: execution context in following choices
    + `sh` (default)
    + `psql`: needs PostgreSQL `psql` client to be installed
  - `uri`: connection string used by `psql`'s database option (`-d`)
  - `connection`: connection name as described below
  - `output`: output file name as described below
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
# run queries from a file simultaneously
tasks:
  - id: 1
    type: psql
    name: dispatch queries from a file
    file: queries.sql
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

### Traces

* `logfile`: summary of the tasks execution (default: disabled)
  - must be a valid path

```yaml
logfile: result.out
```

* task `output`: writes command's output in a file
  - accept standard templating syntax on this task's context
  - does not interrupt others workers if file could not be created or written

```yaml
# write psql output in a dedicated file per query
tasks:
  - id: 1
    type: psql
    file: queries.sql
    output: result_{{.ID}}_{{.QueryID}}.out
```

### Named connections

* `connections`: declares named connections used by tasks
  * `name`: connection name
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
  - name: otherdb
    host: localhost
    dbname: otherdb

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