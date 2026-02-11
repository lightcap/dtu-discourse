package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type AdminHandler struct {
	Store *store.Store
}

// GET /admin/site_settings.json
func (h *AdminHandler) GetSiteSettings(w http.ResponseWriter, r *http.Request) {
	settings := h.Store.GetSiteSettings()
	writeJSON(w, http.StatusOK, model.SiteSettingsResponse{SiteSettings: settings})
}

// PUT /admin/site_settings/{name}.json
func (h *AdminHandler) UpdateSiteSetting(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	value, ok := body[name]
	if !ok {
		value = body["value"]
	}
	h.Store.UpdateSiteSetting(name, value)
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /admin/backups.json
func (h *AdminHandler) ListBackups(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.BackupListResponse{})
}

// POST /admin/backups.json
func (h *AdminHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /admin/dashboard.json
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"global_reports":    []interface{}{},
		"updated_at":        nil,
		"top_referred_topics": []interface{}{},
	})
}

// GET /site.json
func (h *AdminHandler) SiteInfo(w http.ResponseWriter, r *http.Request) {
	cats := h.Store.ListCategories()
	groups := h.Store.ListGroups()
	writeJSON(w, http.StatusOK, model.SiteInfo{
		DefaultArchetype: "regular",
		NotificationTypes: map[string]int{
			"mentioned":          1,
			"replied":            2,
			"quoted":             3,
			"edited":             4,
			"liked":              5,
			"private_message":    6,
			"invited_to_private_message": 7,
			"invitee_accepted":   8,
			"posted":             9,
			"moved_post":         10,
			"linked":             11,
			"granted_badge":      12,
			"invited_to_topic":   13,
			"custom":             14,
			"group_mentioned":    15,
			"group_message_summary": 16,
			"watching_first_post": 17,
			"topic_reminder":     18,
			"liked_consolidated": 19,
			"post_approved":      20,
			"code_review_commit_approved": 21,
		},
		PostTypes: map[string]int{
			"regular":        1,
			"moderator_action": 2,
			"small_action":   3,
			"whisper":        4,
		},
		TrustLevels: map[string]int{
			"newuser":  0,
			"basic":    1,
			"member":   2,
			"regular":  3,
			"leader":   4,
		},
		Groups:     groups,
		Categories: cats,
	})
}

// GET /session/csrf.json — CSRF token endpoint
func (h *AdminHandler) CSRFToken(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"csrf": "dtu-csrf-token"})
}
