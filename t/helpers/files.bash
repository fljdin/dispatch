function teardown_files {
    rm -f commands.sh queries.sql
}

function create_commands() {
    cat <<-EOF > commands.sh
	echo 1
	echo 2
	EOF
}

function create_queries() {
    cat <<-EOF > queries.sql
	SELECT 1;
	SELECT 2;
	EOF
}