setup() {
    export GOTEST=1
    go build 
   
    DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"
    PATH="$DIR/..:$PATH"
    cd $DIR
}

teardown() {
    rm -f *.log
}

run_with_config() {
    dispatch run --config config/$1.yaml
    diff expected/$1.log $1.log
}

@test "task with default connection" {
    run_with_config "default_connection"
}

@test "worker forward uri to generated tasks" {
    true
}