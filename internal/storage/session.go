package storage

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SessionStore interface {
	AddMessage(ctx context.Context, sessionID string, role string, content string) error
	GetHistory(ctx context.Context, sessionID string) ([]Message, error)
}

// MemoryStore - Lưu trữ cục bộ trong RAM
type MemoryStore struct {
	mu       sync.Mutex
	sessions map[string][]Message
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{sessions: make(map[string][]Message)}
}

func (m *MemoryStore) AddMessage(ctx context.Context, sessionID string, role string, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[sessionID] = append(m.sessions[sessionID], Message{Role: role, Content: content})
	if len(m.sessions[sessionID]) > 12 {
		m.sessions[sessionID] = m.sessions[sessionID][len(m.sessions[sessionID])-12:]
	}
	return nil
}

func (m *MemoryStore) GetHistory(ctx context.Context, sessionID string) ([]Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msgs := make([]Message, len(m.sessions[sessionID]))
	copy(msgs, m.sessions[sessionID])
	return msgs, nil
}

// RedisStore - Lưu trữ phân tán trên Redis (Production-ready)
type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(redisURL string) (*RedisStore, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisStore{client: client}, nil
}

func (r *RedisStore) AddMessage(ctx context.Context, sessionID string, role string, content string) error {
	key := "session:" + sessionID
	data, err := json.Marshal(Message{Role: role, Content: content})
	if err != nil {
		return err
	}
	pipe := r.client.Pipeline()
	pipe.RPush(ctx, key, data)
	pipe.LTrim(ctx, key, -12, -1)
	pipe.Expire(ctx, key, 24*time.Hour)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisStore) GetHistory(ctx context.Context, sessionID string) ([]Message, error) {
	key := "session:" + sessionID
	vals, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var msgs []Message
	for _, v := range vals {
		var msg Message
		if err := json.Unmarshal([]byte(v), &msg); err == nil {
			msgs = append(msgs, msg)
		}
	}
	return msgs, nil
}
