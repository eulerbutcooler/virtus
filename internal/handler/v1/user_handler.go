package v1

import (
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

type updateProfileRequest struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=100"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=20"`
}

// Returns the authenticated user's profile
func (h *UserHandler) GetMe (w http.ResponseWriter, r *http.Request){
	userID:=middleware.UserIDFrom(r.Context())
	user,err:=h.users.GetByID(r.Context(),userID)
	if err!=nil{
		response.FromError(w,err)
		return
	}
	response.OK(w,http.StatusOK,user)
}

// UpdateMe patches the user's profile
func (h *UserHandler) UpdateMe (w http.ResponseWriter, r *http.Request){
	var req updateProfileRequest
	if err:=response.DecodeAndValidate(r,&req);err!=nil{
		response.FromError(w,err)
		return
	}
	userID:=middleware.UserIDFrom(r.Context())
	user,err:=h.users.UpdateProfile(r.Context(),userID,service.UpdateProfileInput{Name: req.Name})
	if err!=nil{
		response.FromError(w,err)
		return
	}
	response.OK(w, http.StatusOK,user)
}

// Updates the user's password
func (h *UserHandler) ChangePassword (w http.ResponseWriter, r *http.Request){
	var req changePasswordRequest
	if err:=response.DecodeAndValidate(r,&req); err!=nil{
		response.FromError(w,err)
		return
	}
	userID:=middleware.UserIDFrom(r.Context())
	if err:=h.users.ChangePassword(r.Context(),userID,service.ChangePasswordInput{
		CurrentPassword: req.CurrentPassword,
		NewPassword: req.NewPassword,
	});err!=nil{
		response.FromError(w,err)
		return
	}
	response.OK(w,http.StatusOK,map[string]string{"status": "password updated"})
}

// ADMIN ENDPOINTS
// These are admin only

// Returns a paginated list of all users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	users, pg, err := h.users.List(r.Context(), page)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.List(w, http.StatusOK, users, pg)
}

// Returns any user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid user id")
		return
	}
	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, user)
}


// Marks a user as verified.
func (h *UserHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid user id")
		return
	}
	user, err := h.users.VerifyUser(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, user)
}
// DeleteUser hard-deletes a user. Admin only.
// DELETE /admin/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid user id")
		return
	}
	if err := h.users.DeleteUser(r.Context(), id); err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, map[string]string{"status": "deleted"})
}