package storage

import (
    "context"
    "sync"
    "time"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type SessionStore interface {
    AddMessage(ctx context.Context, sessionID string, msg Message) error
    GetMessages(ctx context.Context, sessionID string) ([]Message, error)
}

// MemoryStore dùng RAM trực tiếp, không cần Redis và không cần file go.sum
type MemoryStore struct {
    mu       sync.Mutex
    sessions map[string][]Message
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        sessions: make(map[string][]Message),
    }
}

func (s *MemoryStore) AddMessage(ctx context.Context, sessionID string, msg Message) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.sessions[sessionID] = append(s.sessions[sessionID], msg)
    return nil
}

func (s *MemoryStore) GetMessages(ctx context.Context, sessionID string) ([]Message, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.sessions[sessionID], nil
}
