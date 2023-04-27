# Simple task dispatcher

Provides an easy-to-use command to dispatch tasks described in a YAML file.

Common use cases:

* Launching multiple elementary tasks in parallel
* Add a condition with a task dependent on another
* Split SQL files to execute statements as elementary tasks

## Usage

```sh
Usage:
  dispatch -c config [-j 2] [flags]

Flags:
  -c, --config string   configuration file
  -h, --help            help for dispatch
  -j, --jobs int        number of workers (default 2)
```

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


### Named connections

* `connections`: declares named connections used by tasks
  * `name`: connection name
  * `uri`: connection string used by `psql`'s database option (`-d`)

```yaml
connections:
  - name: db
    uri: postgresql://remote

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