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

func (r *Rule) Clone() *Rule {
	clone := *r
	clone.Targets = map[string]Target{}
	for tn, t := range r.Targets {
		clone.Targets[tn] = t
	}
	return &clone
}

func (r *Rule) ForceValidChecks() {
	for tn, t := range r.Targets {
		t.Check = NewCheck(t.Check.Interval, t.Check.Timeout, t.Check.Minimum, t.Check.Address)
		r.Targets[tn] = t // update: overwrite
	}
}

func NewCheck(vInt, vTmout, vMin int, addr string) HealthCheck {
	if vInt < 1 {
		vInt = 5
	}

	if vTmout < 1 {
		vTmout = 5
	}

	if vMin < 1 {
		vMin = 3
	}

	c := HealthCheck{
		Interval: vInt,
		Timeout:  vTmout,
		Minimum:  vMin,
		Address:  addr,
	}

	return c
}
