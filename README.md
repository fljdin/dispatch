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

```yaml
tasks:
  - id: 1
    command: echo test
```