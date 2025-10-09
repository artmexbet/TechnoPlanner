package service

import (
	"technoBro/internal/broker"
)

type TaskManager struct {
	broker *broker.NATSBroker
}

func NewTaskManager(broker *broker.NATSBroker) *TaskManager {
	return &TaskManager{
		broker: broker,
	}
}
