package queries

import (
	"testing"

	"github.com/google/uuid"
)

func TestRequestToDomain_PorterInfoNilWhenResponsibleNotAssigned(t *testing.T) {
	req := &Request{}

	got := req.ToDomain()
	if got.PorterInfo != nil {
		t.Fatalf("expected PorterInfo to be nil when PorterID is nil, got %#v", got.PorterInfo)
	}
}

func TestRequestToDomain_PorterInfoFilledWhenResponsibleAssigned(t *testing.T) {
	porterID := uuid.New()
	req := &Request{PorterID: &porterID}

	got := req.ToDomain()
	if got.PorterInfo == nil {
		t.Fatal("expected PorterInfo to be set")
	}
	if got.PorterInfo.UserID != porterID {
		t.Fatalf("expected porter id %s, got %s", porterID, got.PorterInfo.UserID)
	}
}
