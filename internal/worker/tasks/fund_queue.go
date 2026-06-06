// Package tasks contains the periodic background jobs run by the worker binary.
package tasks

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/internal/service"
)

// FundQueueTask is the core economic engine of the platform.
//
// On each run it:
//  1. Reads the current pool balance.
//  2. Walks the queue in position order (FIFO).
//  3. For each queued request whose estimated cost ≤ available balance,
//     marks it as funded and deducts the cost from the tracked available balance.
//
// Funding is FIFO: if the top request can't be covered, the run stops.
// This prevents lower-cost requests from jumping the queue.
//
// The pool balance is NOT debited here; that happens later in
// FulfillmentService.UpdateDetails when the admin sets the actual cost.
// The available balance is tracked locally within this run only to avoid
// over-committing funds before the admin acts.
type FundQueueTask struct {
	pool     *service.PoolService
	queue    *service.QueueService
	requests *service.RequestService
}

func NewFundQueueTask(
	pool *service.PoolService,
	queue *service.QueueService,
	requests *service.RequestService,
) *FundQueueTask {
	return &FundQueueTask{pool: pool, queue: queue, requests: requests}
}

func (t *FundQueueTask) Run(ctx context.Context) error {
	poolState, err := t.pool.GetStatus(ctx)
	if err != nil {
		return fmt.Errorf("fundQueue: get pool: %w", err)
	}

	entries, err := t.queue.ListAllOrdered(ctx)
	if err != nil {
		return fmt.Errorf("fundQueue: list queue: %w", err)
	}
	if len(entries) == 0 {
		return nil
	}

	available := poolState.Balance
	funded := 0

	for _, entry := range entries {
		req, err := t.requests.GetByID(ctx, entry.RequestID)
		if err != nil {
			slog.Warn("fundQueue: skipping entry, cannot fetch request",
				"request_id", entry.RequestID, "err", err)
			continue
		}

		// Guard against stale queue entries whose request moved out of queued state.
		if req.Status != domain.RequestQueued {
			continue
		}

		// Strict FIFO: stop at the first request we can't afford.
		if available < req.EstimatedCost {
			slog.Debug("fundQueue: insufficient balance, stopping",
				"request_id", req.ID,
				"needed", req.EstimatedCost,
				"available", available,
			)
			break
		}

		if err := t.requests.MarkFunded(ctx, req.ID); err != nil {
			slog.Error("fundQueue: mark funded failed",
				"request_id", req.ID, "err", err)
			continue
		}

		available -= req.EstimatedCost
		funded++

		slog.Info("fundQueue: request funded",
			"request_id", req.ID,
			"item", req.ItemName,
			"amount", req.EstimatedCost,
			"pool_remaining", available,
		)
	}

	if funded > 0 {
		slog.Info("fundQueue: cycle complete", "funded_count", funded)
	}

	return nil
}
