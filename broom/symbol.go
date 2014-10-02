package broom

import (
	"fmt"
	"sync"
)

type symbolImpl struct {
	id int
}

type Pool struct {
	mutex sync.RWMutex
	lookup map[string] int
	xs []string
}

var pool *Pool

func initPool() {
	pool = &Pool{lookup:make(map[string]int), xs: make([]string, 0)}
}

func (p* Pool)Lock() { p.mutex.Lock() }
func (p* Pool)Unlock() { p.mutex.Unlock() }
func (p* Pool)RLock() { p.mutex.RLock() }
func (p* Pool)RUnlock() { p.mutex.RUnlock() }

func (p* Pool)MakeSymbol(s string) *symbolImpl{
	p.xs = append(p.xs, s)
	n := len(p.xs) - 1
	p.lookup[s] = n
	return &symbolImpl{id:n}
}

func (p* Pool)GenSymbol() *symbolImpl{
	return pool.MakeSymbol(fmt.Sprintf("_symol-%d", len(p.xs)))
}

func (p* Pool)LookUp(s string) (*symbolImpl, bool) {
	if n, ok := p.lookup[s]; ok {
		return &symbolImpl{id:n}, true
	}
	return nil, false
}

func (p* Pool)GetValue(n int) string {
	return p.xs[n]
}

func sym(t string) *symbolImpl {
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
	pool.RLock()
	defer pool.RUnlock()
	return pool.GetValue(s.id)
}

func (s *symbolImpl) Eq(other interface{}) bool {
	if t, ok := other.(*symbolImpl); ok {
		return s.id == t.id
	} else {
		return false
	}
}

func (s *symbolImpl) String() string {
	return s.GetValue()
}

func (s *symbolImpl) DetailedPrint() string {
	return fmt.Sprintf("#symbol-%d-:%s", s.id, s.GetValue())
}

func GenSym() *symbolImpl {
	pool.Lock()
	defer pool.Unlock()
	return pool.GenSymbol()
}

