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

@test "dispatch with --file and --type flags" {
    LOG=loaded_from_sh_file.log
    create_commands
    dispatch run \
      --jobs 1 --log $LOG \
      --file commands.sh --type sh
    assert-diff $LOG
}