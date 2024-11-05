# Contributing

## Packaging

The project is packaged using [GoReleaser] and especially the [nFPM] package
for creating RedHat and Debian packages. Before creating a new release, make
sure to update the version in the `nfpm.yaml` file.

[GoReleaser]: https://goreleaser.com/
[nFPM]: https://nfpm.goreleaser.com/

```yaml
# Version. (required)
# Hence, you should not prefix the version with 'v'.
version: 0.x.y
```

Use the `nfpm` command to create the packages.

```console
$ nfpm package --config nfpm.yaml --target dist/ --packager deb
```

A dedicated [Github Action] is available for convienience.

[Github Action]: https://github.com/fljdin/dispatch/actions/workflows/package.yml

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
