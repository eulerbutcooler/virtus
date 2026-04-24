package service

import (
	"context"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	dbgen "github.com/eulerbutcooler/virtus/internal/repository/postgres/db"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This is returned when a payment intent is created
type PaymentIntent struct {
	// provider's id (Stripe: "jashdabs")
	// stored on the contribution row so webhooks can look it up
	IntentID string
	// goes to the frontend so it can complete the payment via stripe sdk
	// not stored or logged
	ClientSecret string
}

// Abstracts the payment gateway so stripe can be switched with dodopayments etc
type PaymentProvider interface {
	CreateIntent(
		ctx context.Context,
		amount float64,
		currency string,
		metadate map[string]string,
	) (*PaymentIntent, error)
}

type ContributionService struct {
	contributions domain.ContributionRepository
	pool          *PoolService
	db            *pgxpool.Pool
	queries       *dbgen.Queries
	payments      PaymentProvider
}

func NewContributionService(
	contributions domain.ContributionRepository,
	pool *PoolService,
	db *pgxpool.Pool,
	queries *dbgen.Queries,
	payments PaymentProvider,
) *ContributionService {
	return &ContributionService{
		contributions: contributions,
		pool:          pool,
		db:            db,
		queries:       queries,
		payments:      payments,
	}
}

// Sent by handler when a member wants to contribute
type InitialInput struct {
	UserID   uuid.UUID
	Amount   float64
	Currency string
}

// Returned to handler and the client secret goes to the frontend
type InitiateResult struct {
	Contribution *domain.Contribution
	ClientSecret string
}

// Creates a Stripe PaymentIntent and a pending contribution row
// Flow -> Call stripe to reserve the payment -> get back (intentID,clientSecret)
// -> Create a pending contribution row with stripe_pi_id set
// -> Return clientSecret to the frontend
// Frontend uses clientSecret+stripe.js to collect the card and confirm the payment.
// Stripe then fires a webhook which triggers func Complete()
func (s *ContributionService) Initiate(ctx context.Context, in InitialInput) (*InitiateResult, error) {
	if in.Amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be positive", domain.ErrInvalidInput)
	}
	intent, err := s.payments.CreateIntent(ctx, in.Amount, in.Currency, map[string]string{
		"user_id": in.UserID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("contributionService.Initiate: create intent: %w", err)
	}
	contribution, err := s.contributions.Create(ctx, domain.CreateContributionParams{
		UserID:     in.UserID,
		PoolID:     domain.GlobalPoolID,
		Amount:     in.Amount,
		Currency:   in.Currency,
		StripePIID: &intent.IntentID,
	})
	if err != nil {
		return nil, fmt.Errorf("contributionService.Initiate: persist: %w", err)
	}
	return &InitiateResult{
		Contribution: contribution,
		ClientSecret: intent.ClientSecret,
	}, nil
}

// Called by Stripe webhook worker when a payment succeeds
// Marks the contribution status - Completed
// Credit the pool balance
// If either step fails both roll back.
// This method is idempotent. Stripe may deliver the same webhook more than once
// If the contribution is already completee we return nil without doing anything
func (s *ContributionService) Complete(ctx context.Context, stripePaymentID, paymentRef string) error {
	contribution, err := s.contributions.GetByStripePI(ctx, stripePaymentID)
	if err != nil {
		return fmt.Errorf("contributionService.Complete: lookup %w", err)
	}
	if contribution.Status == domain.ContributionCompleted {
		return nil
	}
	if contribution.Status != domain.ContributionPending {
		return fmt.Errorf("%w: cannot complete a %s contribution", domain.ErrInvalidState, contribution.Status)
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("contributionService.Complete: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	txq := dbgen.New(tx)
	if err := txq.UpdateContributionStatus(ctx, dbgen.UpdateContributionStatusParams{
		ID:         contribution.ID,
		Status:     string(domain.ContributionCompleted),
		PaymentRef: &paymentRef,
	}); err != nil {
		return fmt.Errorf("contributionService.Complete: update status: %w", err)
	}
	if err := txq.CreditPool(ctx, dbgen.CreditPoolParams{
		Amount: contribution.Amount,
		ID:     domain.GlobalPoolID,
	}); err != nil {
		return fmt.Errorf("contributionService.Complete: credit pool: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("contributionService.Complete: commit: %w", err)
	}
	return nil
}

// Called by the Stripe webhook worker when a payment fails or is cancelled
// The pool is not touched, only complete payments reach the pool
func (s *ContributionService) Fail(ctx context.Context, stripePaymentIntentID string) error {
	contribution, err := s.contributions.GetByStripePI(ctx, stripePaymentIntentID)
	if err != nil {
		return fmt.Errorf("contributionService.Fail: lookup: %w", err)
	}
	if contribution.Status == domain.ContributionFailed {
		return nil
	}
	if contribution.Status != domain.ContributionPending {
		return fmt.Errorf("%w: cannot fail a %s contribution", domain.ErrInvalidInput, contribution.Status)
	}
	return s.contributions.UpdateStatus(ctx, contribution.ID, domain.ContributionFailed, nil)
}

func (s *ContributionService) Refund(ctx context.Context, contributionID uuid.UUID) error {
	contribution, err := s.contributions.GetByID(ctx, contributionID)
	if err != nil {
		return fmt.Errorf("contributionService.Refund: lookup: %w", err)
	}
	if contribution.Status != domain.ContributionCompleted {
		return fmt.Errorf("%w: only completed contributions can be refunded", domain.ErrInvalidState)
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("contributionService.Refund: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	txq := dbgen.New(tx)

	if err := txq.UpdateContributionStatus(ctx, dbgen.UpdateContributionStatusParams{
		ID:         contribution.ID,
		Status:     string(domain.ContributionRefunded),
		PaymentRef: nil,
	}); err != nil {
		return fmt.Errorf("contributionService.Refund: update status: %w", err)
	}
	n, err := txq.DebitPool(ctx, dbgen.DebitPoolParams{
		Amount: contribution.Amount,
		ID:     domain.GlobalPoolID,
	})
	if err != nil {
		return fmt.Errorf("contributionService.Refund: debit pool: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("contributionService.Refund: %w", domain.ErrInsufficientFunds)
	}
	return tx.Commit(ctx)
}

func (s *ContributionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Contribution, error) {
	c, err := s.contributions.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("contributionService.GetByID: %w", err)
	}
	return c, nil
}

func (s *ContributionService) ListByUser(ctx context.Context,
	userID uuid.UUID, page pagination.Page) ([]*domain.Contribution, pagination.Page, error) {
	contributions, total, err := s.contributions.ListByUser(ctx, userID, page.Limit, page.Offset())
	if err != nil {
		return nil, pagination.Page{}, fmt.Errorf("contributionService.ListByUser: %w", err)
	}
	return contributions, page.WithTotal(total), nil
}

func (s *ContributionService) GetUserTotal(ctx context.Context, userID uuid.UUID) (float64, error) {
	total, err := s.contributions.SumByUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("contributionService.GetUserTotal: %w", err)
	}
	return total, nil
}
