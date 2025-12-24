package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// データ構造をここに集約
type Todo struct {
	ID        int       `json:"id"`         // ID (PK)
	Name      string    `json:"name"`       // name varchar(40)
	Status    int       `json:"status"`     // status int(1) 0:未完了, 1:完了
	Created   time.Time `json:"created"`    // created timestamp
	Updated   time.Time `json:"updated"`    // updated timestamp
	Deleted   bool      `json:"deleted"`    // deleted bool
}

type JSONRepository struct {
	FilePath string
	mu       sync.RWMutex
}

func (r *JSONRepository) Load() ([]Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, err := os.Stat(r.FilePath); os.IsNotExist(err) {
		return []Todo{}, nil
	}

	data, err := os.ReadFile(r.FilePath)
	if err != nil {
		return nil, err
	}

	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return []Todo{}, nil
	}

	return todos, nil
}