package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/libs/dto"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/models"
)

// RawRequestStorage — интерфейс хранилища для сырых запросов
type RawRequestStorage interface {
	ListRawRequests(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error)
	GetRawRequest(ctx context.Context, id uuid.UUID) (domain.RawRequest, error)
	ProcessRawRequest(ctx context.Context, req dto.RawRequestProcessRequest) (domain.Request, domain.RawRequest, error)
}

// RawRequestService — сервис для работы с сырыми запросами от бота
type RawRequestService struct {
	storage RawRequestStorage
}

// NewRawRequestService создаёт новый сервис для сырых запросов
func NewRawRequestService(storage RawRequestStorage) *RawRequestService {
	return &RawRequestService{storage: storage}
}

// List возвр��щает список сырых запросов (только для администратора)
func (s *RawRequestService) List(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	return s.storage.ListRawRequests(ctx, status, limit, offset)
}

// Get возвращает сырой запрос по ID (только для администратора)
func (s *RawRequestService) Get(ctx context.Context, id uuid.UUID) (domain.RawRequest, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.RawRequest{}, err
	}
	return s.storage.GetRawRequest(ctx, id)
}

// Process создаёт нормальную заявку на основе сырого запроса (только для администратора)
func (s *RawRequestService) Process(ctx context.Context, rawID uuid.UUID, body models.RawRequestProcessRequest) (domain.Request, domain.RawRequest, error) {
	if err := requireAdmin(ctx); err != nil {
		return domain.Request{}, domain.RawRequest{}, err
	}

	equipments := make([]dto.EquipmentInfo, len(body.Equipments))
	for i, eq := range body.Equipments {
		equipments[i] = dto.EquipmentInfo{
			ID:       eq.ID,
			Quantity: eq.Quantity,
		}
	}

	dtoReq := dto.RawRequestProcessRequest{
		RawRequestID:    rawID.String(),
		RequestText:     body.RequestText,
		ScheduleTime:    body.ScheduleTime,
		Address:         body.Address,
		EquipmentString: body.EquipmentString,
		Equipments:      equipments,
	}

	if body.EndTime != nil {
		t, err := time.Parse(time.RFC3339, *body.EndTime)
		if err != nil {
			return domain.Request{}, domain.RawRequest{}, fmt.Errorf("invalid end_time format: %w", err)
		}
		dtoReq.EndTime = &t
	}

	createdReq, updatedRaw, err := s.storage.ProcessRawRequest(ctx, dtoReq)
	if err != nil {
		return domain.Request{}, domain.RawRequest{}, fmt.Errorf("process raw request: %w", err)
	}

	return createdReq, updatedRaw, nil
}
