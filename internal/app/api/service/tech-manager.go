package service

import (
	"technoBro/pkg/broker"
)

type TechManager struct {
	broker *broker.NATSBroker
}

func NewTechManager(broker *broker.NATSBroker) *TechManager {
	return &TechManager{
		broker: broker,
	}
}
