package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// ExtendedTopicsHandler handles undocumented topic operations that are
// used by SDKs and integrations but not in the official OpenAPI spec.
type ExtendedTopicsHandler struct {
	Store *store.Store
}

// PUT /t/{id}/archive-message
func (h *ExtendedTopicsHandler) ArchiveMessage(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/move-to-inbox
func (h *ExtendedTopicsHandler) MoveToInbox(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/convert-topic/{type}
func (h *ExtendedTopicsHandler) ConvertTopic(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/publish
func (h *ExtendedTopicsHandler) Publish(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/reset-bump-date
func (h *ExtendedTopicsHandler) ResetBumpDate(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/clear-pin
func (h *ExtendedTopicsHandler) ClearPin(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/re-pin
func (h *ExtendedTopicsHandler) RePin(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/mute
func (h *ExtendedTopicsHandler) Mute(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/unmute
func (h *ExtendedTopicsHandler) Unmute(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/make-banner
func (h *ExtendedTopicsHandler) MakeBanner(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/remove-banner
func (h *ExtendedTopicsHandler) RemoveBanner(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/remove-allowed-user
func (h *ExtendedTopicsHandler) RemoveAllowedUser(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/remove-allowed-group
func (h *ExtendedTopicsHandler) RemoveAllowedGroup(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/recover
func (h *ExtendedTopicsHandler) Recover(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/tags
func (h *ExtendedTopicsHandler) UpdateTags(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{id}/slow_mode
func (h *ExtendedTopicsHandler) SlowMode(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /t/{id}/move-posts
func (h *ExtendedTopicsHandler) MovePosts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /t/{id}/merge-topic
func (h *ExtendedTopicsHandler) MergeTopic(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /t/{id}/invite-group
func (h *ExtendedTopicsHandler) InviteGroup(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /t/{id}/timings
func (h *ExtendedTopicsHandler) DestroyTimings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /t/{id}/post_ids
func (h *ExtendedTopicsHandler) PostIDs(w http.ResponseWriter, r *http.Request) {
	idStr := pathParam(r, "id")
	idStr = strings.TrimSuffix(idStr, ".json")
	topicID, _ := strconv.Atoi(idStr)
	posts := h.Store.GetTopicPosts(topicID, nil)
	ids := make([]int, len(posts))
	for i, p := range posts {
		ids[i] = p.ID
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"post_ids": ids})
}

// GET /t/{id}/excerpts
func (h *ExtendedTopicsHandler) Excerpts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// GET /t/id_for/{slug}
func (h *ExtendedTopicsHandler) IDForSlug(w http.ResponseWriter, r *http.Request) {
	slug := pathParam(r, "slug")
	topics := h.Store.ListTopics("")
	for _, t := range topics {
		if t.Slug == slug {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"topic_id":  t.ID,
				"slug":      t.Slug,
				"url":       "/t/" + t.Slug + "/" + strconv.Itoa(t.ID),
			})
			return
		}
	}
	writeError(w, http.StatusNotFound, "topic not found")
}

// GET /t/{id}/view-stats.json
func (h *ExtendedTopicsHandler) ViewStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"topic_view_stats": []interface{}{},
	})
}

// GET /topics/feature_stats
func (h *ExtendedTopicsHandler) FeatureStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
