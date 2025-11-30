package service

import "github.com/google/uuid"

func uuidPtrFromString(id string) *uuid.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return &parsed
}
