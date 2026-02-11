package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type SSOHandler struct {
	Store *store.Store
}

// POST /admin/users/sync_sso
func (h *SSOHandler) SyncSSO(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	externalID, _ := body["external_id"].(string)
	email, _ := body["email"].(string)
	username, _ := body["username"].(string)
	name, _ := body["name"].(string)

	if externalID == "" || email == "" {
		writeError(w, http.StatusUnprocessableEntity, "external_id and email are required")
		return
	}

	u, err := h.Store.SyncSSO(externalID, email, username, name)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.UserResponse{User: *u})
}
