package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"requests/internal/domain"
	"requests/internal/postgres/queries"

	"github.com/google/uuid"
)

func (p *Postgres) CreateRequest(ctx context.Context, req domain.Request) (*domain.Request, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := p.q.WithTx(tx)

	// Create the request
	params := queries.CreateRequestParams{
		TelegramUserID: req.UserID,
		RequestText:    req.RequestText,
	}
	createdReq, err := q.CreateRequest(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	createdReqConverted := createdReq.ToDomain()

	// Assign technics to the request
	assignTechnicParams := make([]queries.AssignTechnicToRequestParams, len(req.Technics))
	for i, technic := range req.Technics {
		assignTechnicParams[i] = queries.AssignTechnicToRequestParams{
			RequestID: createdReq.ID,
			TechnicID: int32(technic.ID),
		}
	}
	br := q.AssignTechnicToRequest(ctx, assignTechnicParams)
	defer br.Close()

	// Collect assigned technics (need only ID and quantity)
	br.QueryRow(func(i int, row queries.TechnicsToRequest, _err error) {
		if _err != nil {
			slog.Error("error assigning technic to request", "error", _err)
			return
		}
		createdReqConverted.Technics = append(createdReqConverted.Technics, domain.Technic{
			ID:       int(row.TechnicID),
			Quantity: int(row.Quantity),
		})
	})

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}
	return createdReqConverted, nil
}

func (p *Postgres) GetRequestsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Request, error) {
	params := queries.GetRequestsByTelegramUserIDParams{
		TelegramUserID: userID,
		Limit:          limit,
		Offset:         offset,
	}
	requests, err := p.q.GetRequestsByTelegramUserID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error getting requests by user ID: %w", err)
	}

	var result []domain.Request
	for _, req := range requests {
		result = append(result, *req.ToDomain())
	}
	return result, nil
}

func (p *Postgres) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	q := p.q.WithTx(tx)
	params := queries.UpdateRequestStatusParams{
		ID:     requestID,
		Status: status,
	}
	_, err = q.UpdateRequestStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("error updating request status: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func (p *Postgres) GetRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.Request, error) {
	req, err := p.q.GetRequestByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("error getting request by ID: %w", err)
	}
	reqDomain := req.ToDomain()
	return reqDomain, nil
}

func (p *Postgres) AssignTechnicsToRequest(ctx context.Context, requestID uuid.UUID, technicIDs []int) []error {
	assignTechnicParams := make([]queries.AssignTechnicToRequestParams, len(technicIDs))
	for i, technicID := range technicIDs {
		assignTechnicParams[i] = queries.AssignTechnicToRequestParams{
			RequestID: requestID,
			TechnicID: int32(technicID),
		}
	}
	br := p.q.AssignTechnicToRequest(ctx, assignTechnicParams)
	defer br.Close()
	var resultErrors []error
	br.QueryRow(func(i int, t queries.TechnicsToRequest, err error) {
		if err != nil {
			slog.Error("error assigning technic to request", "err", err)
			resultErrors = append(resultErrors,
				fmt.Errorf("error assigning technic id:<%v> to request: %w", t.TechnicID, err))
		}
	})
	return resultErrors
}
