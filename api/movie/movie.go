package movie

import (
	"sync"

	"gorm.io/gorm"
)

type handler struct {
	sync.RWMutex
	DB *gorm.DB
}

func NewController(db *gorm.DB) *handler {
	return &handler{
		DB: db,
	}
}
