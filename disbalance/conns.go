package main

import (
	"sync"
)

type connTable struct {
	table map[string]int
	lock  sync.RWMutex
}

func newConnTable() *connTable {
	c := &connTable{
		table: map[string]int{},
	}
	return c
}

func (c *connTable) cloneTable() map[string]int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	tab := map[string]int{}
	for k, v := range c.table {
		tab[k] = v
	}

	return tab
}

func (c *connTable) add(target string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	count, ok := c.table[target]
	if ok {
		count++
	} else {
		count = 1
	}
	c.table[target] = count
}

func (c *connTable) del(target string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	count, ok := c.table[target]
	if !ok {
		return
	}

	if count == 1 {
		delete(c.table, target)
		return
	}

	count--
	c.table[target] = count
}
