package status

var (
	Waiting     Status = 0
	Interrupted Status = 1
	Failed      Status = 2
	Ready       Status = 3
	Succeeded   Status = 4

	StatusTypes = []string{"waiting", "interrupted", "failed", "ready", "succeeded"}
)

type Status int

func (s Status) String() string {
	return StatusTypes[s]
}

func (s Status) IsSucceeded() bool {
	return s == Ready || s == Succeeded
}

func (s Status) IsFailed() bool {
	return s == Interrupted || s == Failed
}
