package storage

import (
	"context"
)

type repository interface {
	porterRepository
	equipmentRepository
	categoryRepository
	requestRepository
	requestHistoryRepository
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, payload interface{}) error
}

type Storage struct {
	Porters       *PorterStorage
	Equipment     *EquipmentStorage
	Categories    *CategoryStorage
	Requests      *RequestStorage
	StatusHistory *RequestHistoryStorage
}

func NewStorage(repo repository, publisher EventPublisher) *Storage {
	return &Storage{
		Porters:       NewPorterStorage(repo),
		Equipment:     NewEquipmentStorage(repo, publisher),
		Categories:    NewCategoryStorage(repo, publisher),
		Requests:      NewRequestStorage(repo, publisher),
		StatusHistory: NewRequestHistoryStorage(repo, publisher),
	}
}
