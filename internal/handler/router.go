package handler

import (
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	v1 "github.com/eulerbutcooler/virtus/internal/handler/v1"
	"github.com/eulerbutcooler/virtus/internal/service"
	stripepkg "github.com/eulerbutcooler/virtus/pkg/stripe"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Services bundles everything the HTTP layer depends on.
type Services struct {
	Auth         *service.AuthService
	User *service.UserService
	Pool         *service.PoolService
	Request      *service.RequestService
	Queue        *service.QueueService
	Contribution *service.ContributionService
	Fulfillment  *service.FulfillmentService
	Delivery     *service.DeliveryService
	Impact       *service.ImpactService
	Institution  *service.InstitutionService
	Stripe       *stripepkg.Provider
}

func NewRouter(svc Services) http.Handler {
	authn := middleware.NewAuthenticator(svc.Auth)
	userH :=v1.NewUserHandler(svc.User)
	authH := v1.NewAuthHandler(svc.Auth)
	poolH := v1.NewPoolHandler(svc.Pool)
	queueH := v1.NewQueueHandler(svc.Queue)
	requestH := v1.NewRequestHandler(svc.Request)
	contributionH := v1.NewContributionHandler(svc.Contribution, svc.Stripe)
	fulfillmentH := v1.NewFulfillmentHandler(svc.Fulfillment)
	deliveryH := v1.NewDeliveryHandler(svc.Delivery)
	impactH := v1.NewImpactHandler(svc.Impact)
	institutionH := v1.NewInstitutionHandler(svc.Institution)

	r := chi.NewRouter()

	// CORS — must be the first middleware so preflight OPTIONS requests are handled
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
		},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chimw.RequestID)
	r.Use(middleware.RequestLogger)
	r.Use(chimw.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Stripe webhook — public, no auth. Must read raw body before any middleware touches it.
	r.Post("/webhooks/stripe", contributionH.StripeWebhook)

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes.
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.Post("/refresh", authH.Refresh)
			r.Post("/logout", authH.Logout)
		})

		// All routes below require a valid access token.
		r.Group(func(r chi.Router) {
			r.Use(authn.RequireAuth)

			r.Get("/me", userH.GetMe)
			r.Patch("/me", userH.UpdateMe)
			r.Post("/me/password", userH.ChangePassword)

			// Pool
			r.Get("/pool", poolH.Status)

			// Queue
			r.Get("/queue", queueH.List)
			r.Get("/queue/{requestID}", queueH.GetByRequest)

			// Requests
			r.Route("/requests", func(r chi.Router) {
				r.Post("/", requestH.Submit)
				r.Get("/", requestH.List)
				r.Get("/{id}", requestH.Get)
				r.Patch("/{id}", requestH.Update)
				r.Delete("/{id}", requestH.Delete)
			})

			// Contributions (member)
			r.Route("/contributions", func(r chi.Router) {
				r.Post("/", contributionH.Initiate)
				r.Get("/", contributionH.List)
				r.Get("/total", contributionH.Total)
				r.Get("/{id}", contributionH.Get)
			})

			// Fulfillments (member read)
			r.Get("/fulfillments/{id}", fulfillmentH.Get)
			r.Get("/fulfillments/request/{requestID}", fulfillmentH.GetByRequest)

			// Deliveries (member read)
			r.Get("/deliveries/{id}", deliveryH.Get)
			r.Get("/deliveries/fulfillment/{fulfillmentID}", deliveryH.GetByFulfillment)

			// Impact records (member)
			r.Route("/impact", func(r chi.Router) {
				r.Post("/", impactH.Record)
				r.Get("/", impactH.ListMine)
				r.Get("/{id}", impactH.Get)
				r.Get("/delivery/{deliveryID}", impactH.ListByDelivery)
			})

			// Institutions (member/institution)
			r.Route("/institutions", func(r chi.Router) {
				r.Post("/", institutionH.Create)
				r.Get("/me", institutionH.GetMine)
				r.Get("/{id}", institutionH.Get)
				r.Patch("/{id}", institutionH.Update)
				r.Post("/{id}/contributions", institutionH.Contribute)
				r.Get("/{id}/contributions", institutionH.ListContributions)
			})

			// Admin-only routes.
			r.Group(func(r chi.Router) {
				r.Use(authn.RequireRole(domain.RoleAdmin))

				//Admin users
				r.Get("/admin/users", userH.ListUsers)
r.Get("/admin/users/{id}", userH.GetUser)
r.Post("/admin/users/{id}/verify", userH.VerifyUser)
r.Delete("/admin/users/{id}", userH.DeleteUser)

				// Requests admin
				r.Get("/admin/requests", requestH.AdminList)
				r.Post("/admin/requests/{id}/verify", requestH.Verify)
				r.Post("/admin/requests/{id}/reject", requestH.Reject)

				// Fulfillments admin
				r.Post("/admin/fulfillments", fulfillmentH.Begin)
				r.Get("/admin/fulfillments", fulfillmentH.List)
				r.Patch("/admin/fulfillments/{id}", fulfillmentH.Update)
				r.Post("/admin/fulfillments/{id}/cancel", fulfillmentH.Cancel)

				// Deliveries admin
				r.Post("/admin/deliveries", deliveryH.Ship)
				r.Post("/admin/deliveries/{id}/verify", deliveryH.Verify)
				r.Post("/admin/deliveries/{id}/fail", deliveryH.MarkFailed)

				// Institutions admin
				r.Post("/admin/institutions/{id}/verify", institutionH.Verify)
			})
		})
	})

	return r
}
