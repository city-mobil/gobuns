package mysql

import "sync"

type failStats struct {
	mu      sync.RWMutex
	storage map[int]int
}

type BarberStats struct {
	CirculStats map[string]int
	FailStats   map[string]int
}
