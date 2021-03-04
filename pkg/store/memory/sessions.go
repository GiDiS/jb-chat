package memory

import (
	"sync"
)

type sessionsMemoryStore struct {
	rwMx sync.RWMutex
}

func NewSessionsMemoryStore() *sessionsMemoryStore {
	return &sessionsMemoryStore{}
}
