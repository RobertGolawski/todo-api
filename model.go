package main

import (
	"time"
)

type Priority int

const (
	Low Priority = iota
	Medium
	High
	Critical
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

type TodoItem struct {
	ID       int       `json:"id"`
	ToDo     string    `json:"todo"`
	DueDate  time.Time `json:"due_date"`
	Priority Priority  `json:"priority"`
	Done     bool      `json:"done"`
}

type TodoList struct {
	Todos  []*TodoItem `json:"items"`
	NextID int         `json:"next_id"`
}
