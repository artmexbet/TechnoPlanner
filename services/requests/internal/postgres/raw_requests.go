package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
	"github.com/artmexbet/TechnoPlanner/services/requests/internal/postgres/queries"
)

// CreateRawRequest сохраняет сырой запрос от Telegram-бота
func (p *Postgres) CreateRawRequest(ctx context.Context, req domain.RawRequest) (*domain.RawRequest, error) {
	params := queries.CreateRawRequestParams{
		TelegramID: req.TelegramID,
		Username:   req.Username,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		RawText:    req.RawText,
	}
	created, err := p.q.CreateRawRequest(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create raw request: %w", err)
	}
	return created.ToDomain(), nil
}

// GetRawRequests возвращает список сырых запросов с фильтрацией по статусу и пагинацией
func (p *Postgres) GetRawRequests(ctx context.Context, status string, limit, offset int32) ([]domain.RawRequest, error) {
	params := queries.GetRawRequestsParams{
		Column1: status,
		Limit:   limit,
		Offset:  offset,
	}
	rows, err := p.q.GetRawRequests(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get raw requests: %w", err)
	}
	result := make([]domain.RawRequest, 0, len(rows))
	for _, row := range rows {
		result = append(result, *row.ToDomain())
	}
	return result, nil
}

// GetRawRequestByID возвращает сырой запрос по ID
func (p *Postgres) GetRawRequestByID(ctx context.Context, id uuid.UUID) (*domain.RawRequest, error) {
	row, err := p.q.GetRawRequestByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get raw request by id: %w", err)
	}
	return row.ToDomain(), nil
}

// MarkRawRequestProcessed помечает сырой запрос как обработанный и прикрепляет ID созданной заявки
func (p *Postgres) MarkRawRequestProcessed(ctx context.Context, id uuid.UUID, requestID uuid.UUID) (*domain.RawRequest, error) {
	params := queries.MarkRawRequestProcessedParams{
		ID:                 id,
		ProcessedRequestID: &requestID,
	}
	updated, err := p.q.MarkRawRequestProcessed(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mark raw request processed: %w", err)
	}
	return updated.ToDomain(), nil
}
