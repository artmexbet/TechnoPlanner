package service

import (
	"broker"
)

type TechManager struct {
	broker *broker.NATSBroker
}

func NewTechManager(broker *broker.NATSBroker) *TechManager {
	return &TechManager{
		broker: broker,
	}
}
