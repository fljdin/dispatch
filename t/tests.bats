setup() {
    export GOTEST=1
    go build

    DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"
    PATH="$DIR/..:$PATH"
    cd $DIR
}

teardown() {
    rm -f *.log *.sql *.sh
}

run_with_config() {
    dispatch run --config config/$1.yaml
    diff expected/$1.log $1.log
}

create_commands() {
    cat <<EOF > commands.sh
echo 1
echo 2
EOF
}

create_queries() {
    cat <<EOF > queries.sql
SELECT 1;
SELECT 2;
EOF
}

@test "task with default connection" {
    run_with_config "default_connection"
}

@test "tasks loaded from a file" {
    create_queries
    run_with_config "loaded_from_sql_file"
}

@test "tasks loaded from psql output" {
    run_with_config "loaded_from_sql_output"
}

@test "dispatch with --file and --type flags" {
    LOG=loaded_from_sh_file.log
    create_commands
    dispatch run \
      --jobs 1 --log $LOG \
      --file commands.sh --type sh
    diff expected/$LOG $LOG
}