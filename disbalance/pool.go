package main

import (
	"sort"
	"sync"
	"time"
)

type checkStatus struct {
	Up    bool
	Since time.Time
}

type pool struct {
	table     map[string]checkStatus
	targetsUp []string
	lock      sync.RWMutex
	next      int
}

func newPool() *pool {
	p := &pool{
		table: map[string]checkStatus{},
	}
	return p
}

func (p *pool) cloneTable() map[string]checkStatus {
	p.lock.RLock()
	defer p.lock.RUnlock()

	tab := map[string]checkStatus{}
	for n, c := range p.table {
		tab[n] = c
	}

	return tab
}

func (p *pool) update() {
	p.targetsUp = make([]string, 0, len(p.table))
	for t, c := range p.table {
		if c.Up {
			p.targetsUp = append(p.targetsUp, t)
		}
	}
	sort.Strings(p.targetsUp)
}

func (p *pool) add(target string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.table[target] = checkStatus{
		Up:    true,
		Since: time.Now(),
	}
	p.update()
}

func (p *pool) del(target string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.table[target] = checkStatus{
		Up:    false,
		Since: time.Now(),
	}
	p.update()
}

func (p *pool) getNext() string {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.targetsUp) < 1 {
		return ""
	}
	p.next++
	if p.next >= len(p.targetsUp) {
		p.next = 0
	}
	t := p.targetsUp[p.next]
	return t
}
