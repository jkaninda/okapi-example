package session

import (
	"sync"
	"time"
)

// SessionManager handles SSE session tracking
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]time.Time
	count    int
}

func New() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]time.Time),
	}
}

// AddSession adds a new session and returns the current count
func (sm *SessionManager) AddSession(id string) int {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.sessions[id] = time.Now()
	sm.count++
	return sm.count
}

// RemoveSession removes a session and returns the remaining count
func (sm *SessionManager) RemoveSession(id string) int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, id)
	return len(sm.sessions)
}

// GetCount returns the current active session count
func (sm *SessionManager) GetCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

// GetTotalCount returns the total number of sessions since server start
func (sm *SessionManager) GetTotalCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.count
}
