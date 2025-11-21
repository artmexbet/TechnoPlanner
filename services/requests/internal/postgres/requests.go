package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"requests/internal/domain"
	"requests/internal/postgres/queries"
)

// CreateRequest creates a new request along with its associated technics in a transaction.
// It returns the created request with its ID and associated technics.
// Uses batching to assign technics to the request.
// In cause of need to do this in one transaction we have to write it this way.
func (p *Postgres) CreateRequest(ctx context.Context, req domain.Request) (*domain.Request, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	q := p.q.WithTx(tx)

	// Create the request
	params := queries.CreateRequestParams{
		TelegramUserID: req.Issuer.ID,
		RequestText:    req.RequestText,
		ScheduleTime:   req.ScheduleTime,
		Address:        req.Address,
	}
	createdReq, err := q.CreateRequest(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	createdReqConverted := createdReq.ToDomain()

	// Assign technics to the request
	assignTechnicParams := make([]queries.AssignEquipmentToRequestParams, len(req.Equipments))
	for i, technic := range req.Equipments {
		assignTechnicParams[i] = queries.AssignEquipmentToRequestParams{
			RequestID:   createdReq.ID,
			EquipmentID: int32(technic.ID),
			Quantity:    int32(technic.Quantity),
		}
	}
	br := q.AssignEquipmentToRequest(ctx, assignTechnicParams)
	defer br.Close() //nolint:errcheck

	// Collect assigned equipment (need only ID and quantity)
	br.QueryRow(func(_ int, row queries.EquipmentToRequest, _err error) {
		if _err != nil {
			slog.Error("error assigning technic to request", "error", _err)
			return
		}
		createdReqConverted.Equipments = append(createdReqConverted.Equipments, domain.Equipment{
			ID:       int(row.EquipmentID),
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

func (p *Postgres) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.StatusType) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	q := p.q.WithTx(tx)
	params := queries.UpdateRequestStatusParams{
		ID:     requestID,
		Status: queries.RequestStatus(status),
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

func (p *Postgres) AssignEquipmentToRequest(ctx context.Context, requestID uuid.UUID, technicIDs []domain.Equipment) []error {
	assignTechnicParams := make([]queries.AssignEquipmentToRequestParams, len(technicIDs))
	for i, technic := range technicIDs {
		assignTechnicParams[i] = queries.AssignEquipmentToRequestParams{
			RequestID:   requestID,
			EquipmentID: int32(technic.ID),
			Quantity:    int32(technic.Quantity),
		}
	}
	br := p.q.AssignEquipmentToRequest(ctx, assignTechnicParams)
	defer br.Close() //nolint:errcheck
	var resultErrors []error
	br.QueryRow(func(_ int, t queries.EquipmentToRequest, err error) {
		if err != nil {
			slog.Error("error assigning technic to request", "err", err)
			resultErrors = append(resultErrors,
				fmt.Errorf("error assigning technic id:<%v> to request: %w", t.EquipmentID, err))
		}
	})
	return resultErrors
}
