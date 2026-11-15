package utils

import (
	"sync"
	"time"
)

var revokedTokens = struct {
	sync.RWMutex
	values map[string]time.Time
}{values: make(map[string]time.Time)}

func RevokeToken(token string, expiresAt time.Time) {
	revokedTokens.Lock()
	defer revokedTokens.Unlock()
	revokedTokens.values[token] = expiresAt
}

func IsTokenRevoked(token string) bool {
	now := time.Now()
	revokedTokens.Lock()
	defer revokedTokens.Unlock()
	for value, expiresAt := range revokedTokens.values {
		if !expiresAt.After(now) {
			delete(revokedTokens.values, value)
		}
	}
	_, revoked := revokedTokens.values[token]
	return revoked
}
