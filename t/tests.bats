function setup() {
    load helpers/files.bash

    DIR="${BATS_TEST_FILENAME%/*}"
    PATH="$DIR/..:$PATH"
    cd $DIR
}

function teardown() {
    teardown_files
}

function assert-diff() {
    diff expected/$1 $1
}

@test "task with default connection" {
    dispatch run --config config/default_connection.yaml
    assert-diff default_connection.log
}

@test "tasks loaded from a file" {
    create_queries

    dispatch run --config config/loaded_from_sql_file.yaml
    assert-diff loaded_from_sql_file.log
}

@test "tasks loaded from psql output" {
    dispatch run --config config/loaded_from_sql_output.yaml
    assert-diff loaded_from_sql_output.log
}

@test "exec with --file flags" {
    LOG=loaded_from_sh_file.log
    create_commands

    dispatch exec \
      --type sh --file commands.sh \
      --jobs 1 --log $LOG
    assert-diff $LOG
}

@test "exec with --to and --command flags" {
    LOG=loaded_from_sql_output.log

    dispatch exec \
      --type psql --to sh \
      --command "SELECT format('echo %s', i) FROM generate_series(1, 2) AS i" \
      --jobs 1 --log $LOG
    assert-diff $LOG
}