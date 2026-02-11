package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type InvitesHandler struct {
	Store *store.Store
}

// POST /invites
func (h *InvitesHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	email, _ := body["email"].(string)
	inv, err := h.Store.CreateInvite(email, nil, nil)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, inv)
}

// GET /invites/retrieve.json
func (h *InvitesHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.Invite{})
}

// PUT /invites/{invite_id}
func (h *InvitesHandler) Update(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /invites
func (h *InvitesHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /invites/destroy-all-expired
func (h *InvitesHandler) DestroyAllExpired(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /invites/reinvite-all
func (h *InvitesHandler) ResendAll(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /invites/reinvite
func (h *InvitesHandler) Resend(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /invite-token/generate
func (h *InvitesHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}
