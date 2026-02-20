package main

import "time"

type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    int       `json:"status"` // 0: incomplete, 1: complete
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Deleted   bool      `json:"deleted"`
}