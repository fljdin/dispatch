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

* `tasks`: a list of tasks to run
  - `id` (required)
  - `command` (required): will be executed
  - `name`: attach a name or description to a task
  - `type`: defines the execution context 
    + `sh` (default)
    + `psql`: needs PostgreSQL `psql` client to be installed
  - `uri`: connection string used by `psql`'s database option (`-d`)
  - `connection`: connection name as described below

```yaml
tasks:
  - id: 1
    command: echo test
  - id: 2
    type: psql
    name: run this statement
    command: SELECT user;
    uri: postgresql://localhost
```

> All PostgreSQL environment variables can be used in place of `uri` as it used
> `psql` client.

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

# perform the following tasks sequentially
tasks:
  - id: 1
    command: echo foo
  - id: 2
    command: echo bar
```