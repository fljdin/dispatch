# Contributing

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
