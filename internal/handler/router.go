package handler

import (
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	v1 "github.com/eulerbutcooler/virtus/internal/handler/v1"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// Services bundles everything the HTTP layer depends on.
type Services struct {
	Auth    *service.AuthService
	Pool    *service.PoolService
	Request *service.RequestService
	Queue   *service.QueueService
}

func NewRouter(svc Services) http.Handler {
	authn := middleware.NewAuthenticator(svc.Auth)
	authH := v1.NewAuthHandler(svc.Auth)
	poolH := v1.NewPoolHandler(svc.Pool)
	queueH := v1.NewQueueHandler(svc.Queue)
	requestH := v1.NewRequestHandler(svc.Request)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(middleware.RequestLogger)
	r.Use(chimw.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes.
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.Post("/refresh", authH.Refresh)
			r.Post("/logout", authH.Logout)
		})

		// Authenticated routes.
		r.Group(func(r chi.Router) {
			r.Use(authn.RequireAuth)

			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				response.OK(w, http.StatusOK, map[string]any{
					"user_id": middleware.UserIDFrom(r.Context()),
					"role":    middleware.UserRoleFrom(r.Context()),
				})
			})

			r.Get("/pool", poolH.Status)
			r.Get("/queue", queueH.List)
			r.Get("/queue/{requestID}", queueH.GetByRequest)

			r.Route("/requests", func(r chi.Router) {
				r.Post("/", requestH.Submit)
				r.Get("/", requestH.List)
				r.Get("/{id}", requestH.Get)
				r.Patch("/{id}", requestH.Update)
				r.Delete("/{id}", requestH.Delete)
			})

			// Admin-only routes.
			r.Group(func(r chi.Router) {
				r.Use(authn.RequireRole(domain.RoleAdmin))
				r.Get("/admin/requests", requestH.AdminList)
				r.Post("/admin/requests/{id}/verify", requestH.Verify)
				r.Post("/admin/requests/{id}/reject", requestH.Reject)
			})
		})
	})

	return r
}
