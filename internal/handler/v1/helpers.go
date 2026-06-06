package v1

import (
	"context"
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func parseID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, "id"))
}

func canAccess(ctx context.Context, ownerID uuid.UUID) bool {
	return middleware.UserIDFrom(ctx) == ownerID || middleware.UserRoleFrom(ctx) == domain.RoleAdmin
}

func parseFulfillmentStatus(s string) domain.FulfillmentStatus {
	return domain.FulfillmentStatus(s)
}
