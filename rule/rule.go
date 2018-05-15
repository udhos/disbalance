package rule

type HealthCheck struct {
	Interval int    // seconds between checks
	Timeout  int    // seconds for every check
	Minimum  int    // min consecutive events for up/down transition
	Address  string // if empty defaults to target address
}

type Target struct {
	Check HealthCheck
}

type Rule struct {
	Protocol string
	Listener string
	Targets  map[string]Target
}
