package main

import (
	"sync"
)

type pool struct {
	table   map[string]struct{}
	targets []string
	lock    sync.RWMutex
	next    int
}

func newPool() *pool {
	p := &pool{
		table: map[string]struct{}{},
	}
	return p
}

func (p *pool) update() {
	p.targets = make([]string, len(p.table), len(p.table))
	var i int
	for t := range p.table {
		p.targets[i] = t
		i++
	}
}

func (p *pool) add(target string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.table[target] = struct{}{}
	p.update()
}

func (p *pool) del(target string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.table, target)
	p.update()
}

func (p *pool) getNext() string {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.targets) < 1 {
		return ""
	}
	t := p.targets[p.next]
	p.next++
	if p.next >= len(p.targets) {
		p.next = 0
	}
	return t
}
