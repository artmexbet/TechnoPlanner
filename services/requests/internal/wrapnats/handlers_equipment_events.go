package wrapnats

import (
	"encoding/json"
	"log/slog"

	"github.com/artmexbet/TechnoPlanner/libs/broker"
	"github.com/artmexbet/TechnoPlanner/libs/dto"
	"github.com/artmexbet/TechnoPlanner/services/requests/internal/domain"
)

// handleEquipmentCreatedEvent синхронизирует локальную копию equipment
// при создании нового оборудования в equipment сервисе.
func (w *NatsWrapper) handleEquipmentCreatedEvent(msg *broker.Msg) error {
	ctx := msg.Context()

	var ev dto.EquipmentSyncEvent
	if err := json.Unmarshal(msg.Data, &ev); err != nil {
		slog.ErrorContext(ctx, "handleEquipmentCreatedEvent: unmarshal error", "error", err)
		return err
	}

	eq := domain.Equipment{
		ID:          ev.ID,
		Name:        ev.Name,
		Description: strPtr(ev.Description),
		Quantity:    ev.Quantity,
	}

	if err := w.eqService.SyncCreate(ctx, eq); err != nil {
		slog.ErrorContext(ctx, "handleEquipmentCreatedEvent: sync error", "error", err, "id", ev.ID)
		return err
	}

	slog.InfoContext(ctx, "equipment synced (created)", "id", ev.ID)
	return nil
}

// handleEquipmentUpdatedEvent синхронизирует локальную копию equipment
// при обновлении оборудования в equipment сервисе.
func (w *NatsWrapper) handleEquipmentUpdatedEvent(msg *broker.Msg) error {
	ctx := msg.Context()

	var ev dto.EquipmentSyncEvent
	if err := json.Unmarshal(msg.Data, &ev); err != nil {
		slog.ErrorContext(ctx, "handleEquipmentUpdatedEvent: unmarshal error", "error", err)
		return err
	}

	eq := domain.Equipment{
		ID:          ev.ID,
		Name:        ev.Name,
		Description: strPtr(ev.Description),
		Quantity:    ev.Quantity,
	}

	if err := w.eqService.SyncUpdate(ctx, eq); err != nil {
		slog.ErrorContext(ctx, "handleEquipmentUpdatedEvent: sync error", "error", err, "id", ev.ID)
		return err
	}

	slog.InfoContext(ctx, "equipment synced (updated)", "id", ev.ID)
	return nil
}

// handleEquipmentDeletedEvent удаляет оборудование из локальной копии
// при удалении его в equipment сервисе.
func (w *NatsWrapper) handleEquipmentDeletedEvent(msg *broker.Msg) error {
	ctx := msg.Context()

	var ev dto.EquipmentDeletedEvent
	if err := json.Unmarshal(msg.Data, &ev); err != nil {
		slog.ErrorContext(ctx, "handleEquipmentDeletedEvent: unmarshal error", "error", err)
		return err
	}

	if err := w.eqService.SyncDelete(ctx, ev.ID); err != nil {
		slog.ErrorContext(ctx, "handleEquipmentDeletedEvent: sync error", "error", err, "id", ev.ID)
		return err
	}

	slog.InfoContext(ctx, "equipment synced (deleted)", "id", ev.ID)
	return nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
