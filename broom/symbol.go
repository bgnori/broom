package broom

import (
	"fmt"
	"sync"
)

type symbolImpl struct {
	id int
	value string
}

type Pool struct {
	mutex  sync.RWMutex
	lookup map[string]int
	xs     []Symbol
}

var pool *Pool

func initPool() {
	pool = &Pool{lookup: make(map[string]int), xs: make([]Symbol, 0)}
}

func (p *Pool) Lock()    { p.mutex.Lock() }
func (p *Pool) Unlock()  { p.mutex.Unlock() }
func (p *Pool) RLock()   { p.mutex.RLock() }
func (p *Pool) RUnlock() { p.mutex.RUnlock() }

func (p *Pool) MakeSymbol(s string) Symbol {
	n := len(p.xs)
	x := &symbolImpl{id: n-1, value: s}
	p.xs = append(p.xs, Symbol(x))
	p.lookup[s] = n
	return x
}

func (p *Pool) GenSymbol() Symbol {
	return pool.MakeSymbol(fmt.Sprintf("_symol-%d", len(p.xs)))
}

func (p *Pool) LookUp(s string) (Symbol, bool) {
	if n, ok := p.lookup[s]; ok {
		return p.xs[n], true
	}
	return nil, false
}

func sym(t string) Symbol {
	if pool == nil {
		initPool()
	}
	pool.RLock()
	if s, ok := pool.LookUp(t); ok {
		pool.RUnlock()
		return s
	}
	pool.RUnlock()

	pool.Lock()
	defer pool.Unlock()
	if s, ok := pool.LookUp(t); ok {
		return s
	}
	return pool.MakeSymbol(t)
}

func (s *symbolImpl) GetValue() string {
	return s.value
}

func (s *symbolImpl) String() string {
	return s.GetValue()
}

func (s *symbolImpl) DetailedPrint() string {
	return fmt.Sprintf("#symbol-%d-:%s", s.id, s.GetValue())
}

func GenSym() Symbol {
	pool.Lock()
	defer pool.Unlock()
	return pool.GenSymbol()
}
