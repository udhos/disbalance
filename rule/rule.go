package rule

// HealthCheck defines health check parameters.
type HealthCheck struct {
	Interval int    // seconds between checks
	Timeout  int    // seconds for every check
	Minimum  int    // min consecutive events for up/down transition
	Address  string // if empty defaults to target address
}

// Target stores a load balancer target backend.
type Target struct {
	Check HealthCheck
}

// Rule defines a full load balancer rule.
type Rule struct {
	Protocol string
	Listener string
	Targets  map[string]Target
}

// Clone clones this rule data.
func (r *Rule) Clone() *Rule {
	clone := *r
	clone.Targets = map[string]Target{}
	for tn, t := range r.Targets {
		clone.Targets[tn] = t
	}
	return &clone
}

// ForceValidChecks replaces invalid health check parameters with valid values.
func (r *Rule) ForceValidChecks() {
	for tn, t := range r.Targets {
		t.Check = NewCheck(t.Check.Interval, t.Check.Timeout, t.Check.Minimum, t.Check.Address)
		r.Targets[tn] = t // update: overwrite
	}
}

// NewCheck creates a new HealthCheck with valid parameters.
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
