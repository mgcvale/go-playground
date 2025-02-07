package config

import (
	"gorm.io/gorm"
	"sync"
)

type StateManager struct {
	mu sync.Mutex
	db *gorm.DB
}

var instance *StateManager
var once sync.Once

func GetApplicationState() *StateManager {
	once.Do(func() {
		instance = &StateManager{}
	})
	return instance
}

func (s *StateManager) SetDB(db *gorm.DB) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db = db
}

func (s *StateManager) GetDB() (db *gorm.DB) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db
}
