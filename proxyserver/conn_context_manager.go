package proxyserver

import "sync"

type ConnContextManager struct {
	ctxSet ConnContextSet
	ctxMu  sync.Mutex
	ctxWg  sync.WaitGroup
}

func NewConnContextManager() *ConnContextManager {
	return &ConnContextManager{
		ctxSet: make(ConnContextSet),
	}
}

type ConnContextSet map[*ConnContext]struct{}

func (m *ConnContextManager) Add(ctx *ConnContext) {
	m.ctxMu.Lock()
	m.ctxSet[ctx] = struct{}{}
	m.ctxMu.Unlock()
}

func (m *ConnContextManager) Remove(ctx *ConnContext) {
	m.ctxMu.Lock()
	delete(m.ctxSet, ctx)
	m.ctxMu.Unlock()
}

func (m *ConnContextManager) GetConnNum() int {
	var num int
	m.ctxMu.Lock()
	num = len(m.ctxSet)
	m.ctxMu.Unlock()
	return num
}

func (m *ConnContextManager) Destroy() {
	m.ctxMu.Lock()
	for ctx, _ := range m.ctxSet {
		ctx.Close()
	}
	m.ctxSet = nil
	m.ctxMu.Unlock()
}
