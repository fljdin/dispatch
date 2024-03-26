package status

var (
	Waiting     Status = 0 // int=0
	Interrupted Status = 1 // int=1
	Failed      Status = 2 // int=2
	Ready       Status = 3 // int=3
	Running     Status = 4 // int=4
	Succeeded   Status = 5 // int=5

	StatusTypes = []string{"waiting", "interrupted", "failed", "ready", "running", "succeeded"}
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
