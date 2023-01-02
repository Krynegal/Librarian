package telegram

import "sync"

type StateKeeper struct {
	mu     sync.RWMutex
	states map[int]string
}

func (sk *StateKeeper) state(userID int) string {
	sk.mu.Lock()
	defer sk.mu.Unlock()
	if _, ok := sk.states[userID]; !ok {
		sk.states[userID] = "init"
	}
	return sk.states[userID]
}

func (sk *StateKeeper) update(userID int, newState string) {
	sk.mu.Lock()
	defer sk.mu.Unlock()
	sk.states[userID] = newState
}
