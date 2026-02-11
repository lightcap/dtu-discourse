package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type BadgesHandler struct {
	Store *store.Store
}

// GET /admin/badges.json
func (h *BadgesHandler) List(w http.ResponseWriter, r *http.Request) {
	badges := h.Store.ListBadges()
	writeJSON(w, http.StatusOK, model.BadgeListResponse{Badges: badges})
}

// POST /admin/badges.json
func (h *BadgesHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	name, _ := body["name"].(string)
	description, _ := body["description"].(string)
	badgeTypeID := 3
	if v, ok := body["badge_type_id"].(float64); ok {
		badgeTypeID = int(v)
	}

	badge, err := h.Store.CreateBadge(name, description, badgeTypeID)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.BadgeResponse{Badge: *badge})
}

// PUT /admin/badges/{id}.json
func (h *BadgesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid badge id")
		return
	}
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	badge, err := h.Store.UpdateBadge(id, body)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.BadgeResponse{Badge: *badge})
}

// DELETE /admin/badges/{id}.json
func (h *BadgesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid badge id")
		return
	}
	if err := h.Store.DeleteBadge(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /user-badges/{username}.json
func (h *BadgesHandler) UserBadges(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	username = strings.TrimSuffix(username, ".json")
	badges, userBadges := h.Store.GetUserBadges(username)
	writeJSON(w, http.StatusOK, model.UserBadgeResponse{
		Badges:     badges,
		UserBadges: userBadges,
	})
}

// POST /user_badges
func (h *BadgesHandler) Grant(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	username, _ := body["username"].(string)
	badgeID := 0
	if v, ok := body["badge_id"].(float64); ok {
		badgeID = int(v)
	} else if v, ok := body["badge_id"].(string); ok {
		badgeID, _ = strconv.Atoi(v)
	}

	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	ub, err := h.Store.GrantUserBadge(u.ID, badgeID)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ub)
}
