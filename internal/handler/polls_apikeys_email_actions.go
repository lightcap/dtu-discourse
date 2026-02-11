package handler

import (
	"net/http"
	"strconv"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type PollsHandler struct {
	Store *store.Store
}

// PUT /polls/vote
func (h *PollsHandler) Vote(w http.ResponseWriter, r *http.Request) {
	body, _ := decodeBody(r)
	postID := 0
	if v, ok := body["post_id"].(float64); ok {
		postID = int(v)
	}
	pollName, _ := body["poll_name"].(string)
	if pollName == "" {
		pollName = "poll"
	}

	var options []string
	if v, ok := body["options"].([]interface{}); ok {
		for _, opt := range v {
			if s, ok := opt.(string); ok {
				options = append(options, s)
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"poll": map[string]interface{}{
			"name":    pollName,
			"type":    "regular",
			"status":  "open",
			"post_id": postID,
			"options": options,
			"voters":  1,
		},
		"vote": options,
	})
}

// PUT /polls/toggle_status
func (h *PollsHandler) ToggleStatus(w http.ResponseWriter, r *http.Request) {
	body, _ := decodeBody(r)
	postID := 0
	if v, ok := body["post_id"].(float64); ok {
		postID = int(v)
	}
	pollName, _ := body["poll_name"].(string)
	status, _ := body["status"].(string)
	if status == "" {
		status = "closed"
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"poll": map[string]interface{}{
			"name":    pollName,
			"type":    "regular",
			"status":  status,
			"post_id": postID,
			"options": []interface{}{},
			"voters":  0,
		},
	})
}

// GET /polls/voters.json
func (h *PollsHandler) Voters(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"voters": map[string]interface{}{},
	})
}

// ---------- API Key Management ----------

type APIKeysHandler struct {
	Store *store.Store
}

// GET /admin/api/keys
func (h *APIKeysHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"id": 1, "key": "test_api_key", "truncated_key": "test_...",
				"description": "System API Key",
				"user": map[string]interface{}{"id": -1, "username": "system"},
				"revoked_at": nil,
			},
			{
				"id": 2, "key": "admin_api_key", "truncated_key": "admin...",
				"description": "Admin API Key",
				"user": map[string]interface{}{"id": 1, "username": "admin"},
				"revoked_at": nil,
			},
		},
	})
}

// POST /admin/api/keys
func (h *APIKeysHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := decodeBody(r)
	desc, _ := body["key"].(map[string]interface{})
	description := ""
	if desc != nil {
		description, _ = desc["description"].(string)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"key": map[string]interface{}{
			"id":          3,
			"key":         "newly_generated_api_key_" + strconv.Itoa(3),
			"description": description,
			"revoked_at":  nil,
		},
	})
}

// POST /admin/api/keys/{id}/revoke
func (h *APIKeysHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"key": map[string]interface{}{"revoked_at": "2024-01-01T00:00:00.000Z"},
	})
}

// POST /admin/api/keys/{id}/undo-revoke
func (h *APIKeysHandler) UndoRevoke(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"key": map[string]interface{}{"revoked_at": nil},
	})
}

// DELETE /admin/api/keys/{id}
func (h *APIKeysHandler) Delete(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /admin/api/keys/scopes
func (h *APIKeysHandler) Scopes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"scopes": map[string]interface{}{
			"topics":     []string{"read", "write", "update"},
			"posts":      []string{"read", "write", "update"},
			"users":      []string{"read", "write", "update"},
			"categories": []string{"read", "write", "update"},
			"uploads":    []string{"write"},
			"email":      []string{"read"},
			"badges":     []string{"read", "write"},
			"global":     []string{"read", "write"},
		},
	})
}

// ---------- Email Admin ----------

type EmailHandler struct {
	Store *store.Store
}

// GET /admin/email.json
func (h *EmailHandler) Settings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"delivery_method":       "smtp",
		"smtp_address":          "localhost",
		"smtp_port":             25,
		"smtp_domain":           "localhost",
		"smtp_authentication":   "plain",
		"smtp_enable_start_tls": true,
	})
}

// GET /admin/email/{filter}.json
func (h *EmailHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// POST /admin/email/test
func (h *EmailHandler) Test(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /admin/email/server-settings
func (h *EmailHandler) ServerSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"delivery_method": "smtp",
		"settings":        map[string]interface{}{},
	})
}

// GET /admin/email/preview-digest
func (h *EmailHandler) PreviewDigest(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"html": "<p>Digest preview</p>"})
}

// ---------- User Actions ----------

type UserActionsHandler struct {
	Store *store.Store
}

// GET /user_actions.json
func (h *UserActionsHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_actions": []interface{}{},
	})
}

// ---------- Topics Timings ----------

type TopicTimingsHandler struct {
	Store *store.Store
}

// POST /topics/timings
func (h *TopicTimingsHandler) Record(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /topics/similar_to
func (h *TopicTimingsHandler) SimilarTo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"similar_topics": []interface{}{},
	})
}

// PUT /topics/bulk
func (h *TopicTimingsHandler) Bulk(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"topic_ids": []interface{}{},
	})
}

// PUT /topics/reset-new
func (h *TopicTimingsHandler) ResetNew(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"message": "OK"})
}

// PUT /topics/pm-reset-new
func (h *TopicTimingsHandler) PMResetNew(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"message": "OK"})
}
