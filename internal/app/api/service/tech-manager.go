package service

import (
	"technoBro/internal/broker"
)

type TechManager struct {
	broker *broker.NATSBroker
}

func NewTechManager(broker *broker.NATSBroker) *TechManager {
	return &TechManager{
		broker: broker,
	}
}
