package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/google/uuid"
)

type ImpactService struct {
	impact     domain.ImpactRepository
	deliveries domain.DeliveryRepository
}

func NewImpactService(impact domain.ImpactRepository, deliveries domain.DeliveryRepository) *ImpactService {
	return &ImpactService{
		impact:     impact,
		deliveries: deliveries,
	}
}

// This gets submitted when logging a follow-up impact survey
// Submitted either by the member themselves or by an admin on their behalf
type RecordInput struct {
	DeliveryID         uuid.UUID
	UserID             uuid.UUID
	IntervalLabel      string  // 30_day, 6_month
	OutcomeDescription *string // how it has helped the seeker
	SatisfactionScore  *int    // rating between 1-5
	Metrics            json.RawMessage
}

// Creates a new impact entry for a verified delivery
// The delivery must exist and be in 'delivered' status.
func (s *ImpactService) Record(ctx context.Context, in RecordInput) (*domain.ImpactRecord, error) {
	if in.IntervalLabel == "" {
		return nil, fmt.Errorf("%w: interval_label is required", domain.ErrInvalidInput)
	}
	if in.SatisfactionScore != nil && (*in.SatisfactionScore < 1 || *in.SatisfactionScore > 5) {
		return nil, fmt.Errorf("%w: satisfaction_score must be between 1 and 5", domain.ErrInvalidInput)
	}

	delivery, err := s.deliveries.GetByID(ctx, in.DeliveryID)
	if err != nil {
		return nil, fmt.Errorf("impactService.Record: fetch delivery: %w", err)
	}
	if delivery.Status != domain.DeliveryDelivered {
		return nil, fmt.Errorf("%w: impact can only be recorded for a delivered item (got %s)",
			domain.ErrInvalidState, delivery.Status)
	}
	record, err := s.impact.Create(ctx, domain.CreateImpactParams{
		DeliveryID:         in.DeliveryID,
		UserID:             in.UserID,
		IntervalLabel:      in.IntervalLabel,
		OutcomeDescription: in.OutcomeDescription,
		SatisfactionScore:  in.SatisfactionScore,
		Metrics:            in.Metrics,
	})
	if err != nil {
		return nil, fmt.Errorf("impactService.Record: %w", err)
	}
	return record, nil
}

// GetByID returns a single impact record by ID.
func (s *ImpactService) GetByID(ctx context.Context, id uuid.UUID) (*domain.ImpactRecord, error) {
	record, err := s.impact.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("impactService.GetByID: %w", err)
	}
	return record, nil
}

// Returns all impact records for a specific delivery.
// Used on the delivery detail page to show the full follow-up timeline.
func (s *ImpactService) ListByDelivery(ctx context.Context, deliveryID uuid.UUID) ([]*domain.ImpactRecord, error) {
	records, err := s.impact.ListByDelivery(ctx, deliveryID)
	if err != nil {
		return nil, fmt.Errorf("impactService.ListByDelivery: %w", err)
	}
	return records, nil
}

// Returns a paginated list of a member's impact records across all their deliveries.
// Used on the member profile / impact dashboard.
func (s *ImpactService) ListByUser(ctx context.Context, userID uuid.UUID, page pagination.Page) ([]*domain.ImpactRecord, pagination.Page, error) {
	records, total, err := s.impact.ListByUser(ctx, userID, page.Limit, page.Offset())
	if err != nil {
		return nil, pagination.Page{}, fmt.Errorf("impactService.ListByUser: %w", err)
	}
	return records, page.WithTotal(total), nil
}
