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

@test "task must be interrupter" {
    dispatch --config config/interrupted_task.yaml
    assert-diff interrupted_task.log
}

@test "#35 task depends on a loader task" {
    dispatch --verbose --config config/depends_on_loader_task.yaml
    assert-diff depends_on_loader_task.log
}

@test "#23 use environment variables" {
    dispatch --config config/env_variables.yaml
    assert-diff env_variables.log

    dispatch --config config/loader_with_env.yaml
    assert-diff loader_with_env.log
}
