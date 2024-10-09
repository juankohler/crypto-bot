package server

import (
	"sync"
)

var mux sync.Mutex
var loggingEnabled bool

func EnableLogging() {
	mux.Lock()
	defer mux.Unlock()

	loggingEnabled = true
}

func DisableLogging() {
	mux.Lock()
	defer mux.Unlock()

	loggingEnabled = false
}
