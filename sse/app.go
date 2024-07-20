package sse

import (
	"sync"
)

type SSEApp struct {
	sync.Map
	sync.RWMutex
}

func New() *SSEApp {
	return new(SSEApp)
}

// func (a *SSEApp) NewConn() *Conn {
// 	a.Lock()
// 	defer a.Unlock()
// 	a.Store()
// }
