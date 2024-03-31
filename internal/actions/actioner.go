package actions

var (
	Shell        = "sh"
	PgSQL        = "psql"
	CommandTypes = []string{Shell, PgSQL}
)

type Actioner interface {
	Command() string
	Run() (Result, []Actioner)
	String() string
	Validate() error
}
