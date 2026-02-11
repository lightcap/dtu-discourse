package handler

import (
	"net/http"
	"strings"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type UsersHandler struct {
	Store *store.Store
}

// GET /users/{username}.json
func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	username = strings.TrimSuffix(username, ".json")
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, model.UserResponse{User: *u})
}

// GET /admin/users/{id}.json
func (h *UsersHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	u := h.Store.GetUser(id)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, *u)
}

// GET /users/by-external/{external_id}
func (h *UsersHandler) GetUserByExternalID(w http.ResponseWriter, r *http.Request) {
	extID := pathParam(r, "external_id")
	u := h.Store.GetUserByExternalID(extID)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, model.UserResponse{User: *u})
}

// POST /users
func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	name, _ := body["name"].(string)
	username, _ := body["username"].(string)
	email, _ := body["email"].(string)
	password, _ := body["password"].(string)

	if username == "" || email == "" {
		writeError(w, http.StatusUnprocessableEntity, "username and email are required")
		return
	}

	u, err := h.Store.CreateUser(name, username, email, password)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.CreateUserResponse{
		Success: true, Active: u.Active, Message: "user created", UserID: u.ID,
	})
}

// PUT /u/{username}/preferences/email
func (h *UsersHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	body, _ := decodeBody(r)
	if email, ok := body["email"].(string); ok {
		h.Store.UpdateUser(u.ID, map[string]interface{}{"email": email})
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /u/{username}/preferences/username
func (h *UsersHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	body, _ := decodeBody(r)
	if newName, ok := body["new_username"].(string); ok {
		h.Store.UpdateUser(u.ID, map[string]interface{}{"username": newName})
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /u/{username}
func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	body, _ := decodeBody(r)
	h.Store.UpdateUser(u.ID, body)
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/approve
func (h *UsersHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"approved": true})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/activate
func (h *UsersHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"active": true})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/deactivate
func (h *UsersHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"active": false})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/trust_level
func (h *UsersHandler) UpdateTrustLevel(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	body, _ := decodeBody(r)
	if v, ok := body["level"].(float64); ok {
		h.Store.UpdateUser(id, map[string]interface{}{"trust_level": v})
	} else {
		h.Store.UpdateUser(id, body)
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/grant_admin
func (h *UsersHandler) GrantAdmin(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"admin": true})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/revoke_admin
func (h *UsersHandler) RevokeAdmin(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"admin": false})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/grant_moderation
func (h *UsersHandler) GrantModeration(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"moderator": true})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/revoke_moderation
func (h *UsersHandler) RevokeModeration(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"moderator": false})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/suspend
func (h *UsersHandler) Suspend(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"suspended": true})
	writeJSON(w, http.StatusOK, map[string]interface{}{"suspension": map[string]interface{}{"suspended": true}})
}

// PUT /admin/users/{id}/unsuspend
func (h *UsersHandler) Unsuspend(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"suspended": false})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /admin/users/{id}/anonymize
func (h *UsersHandler) Anonymize(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	h.Store.UpdateUser(id, map[string]interface{}{"name": "anon"})
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /admin/users/{id}/log_out
func (h *UsersHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /admin/users/{id}.json
func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	if err := h.Store.DeleteUser(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"deleted": true})
}

// GET /admin/users/list/{type}.json
func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	listType := pathParam(r, "type")
	listType = strings.TrimSuffix(listType, ".json")
	users := h.Store.ListUsers(listType)
	writeJSON(w, http.StatusOK, users)
}

// GET /users/check_username.json
func (h *UsersHandler) CheckUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	u := h.Store.GetUserByUsername(username)
	if u != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"available": false, "suggestion": username + "1"})
	} else {
		writeJSON(w, http.StatusOK, map[string]interface{}{"available": true})
	}
}
