function setup() {
    load helpers/files.bash

    DIR="${BATS_TEST_FILENAME%/*}"
    PATH="$DIR/..:$PATH"
    cd $DIR

    setup_service
}

function teardown() {
    teardown_files
}

function assert-diff() {
    diff expected/$1 $1
}

@test "task with default connection" {
    dispatch --config config/default_connection.yaml
    assert-diff default_connection.log
}

@test "tasks loaded from a shell file" {
    create_commands

    dispatch --config config/loaded_from_sh_file.yaml
    assert-diff loaded_from_sh_file.log
}

@test "tasks loaded from a SQL file" {
    create_queries

    dispatch --config config/loaded_from_sql_file.yaml
    assert-diff loaded_from_sql_file.log
}

@test "tasks loaded from psql output" {
    dispatch --config config/loaded_from_sql_output.yaml
    assert-diff loaded_from_sql_output.log
}

@test "#35 task depends on a loader task" {
    dispatch --config config/depends_on_loader_task.yaml
    assert-diff depends_on_loader_task.log
}
