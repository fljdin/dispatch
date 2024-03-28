function setup_service {
	export PGPASSFILE=pgpass.conf
	export PGSERVICEFILE=pgservice.conf

	cat <<-EOF > $PGPASSFILE
	localhost:5432:*:postgres:postgres
	EOF

	cat <<-EOF > $PGSERVICEFILE
	[testing]
	host=localhost
	port=5432
	dbname=postgres
	user=postgres
	EOF

	chmod 600 $PGPASSFILE $PGSERVICEFILE
}

function teardown_files {
	rm -f *.log
	rm -f $PGPASSFILE $PGSERVICEFILE
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
