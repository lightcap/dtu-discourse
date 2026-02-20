package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type TopicsHandler struct {
	Store *store.Store
}

// GET /latest.json
func (h *TopicsHandler) Latest(w http.ResponseWriter, r *http.Request) {
	topics := h.Store.ListTopics("latest")
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		Users: h.Store.UsersForTopics(topics),
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}

// GET /top.json
func (h *TopicsHandler) Top(w http.ResponseWriter, r *http.Request) {
	h.Latest(w, r)
}

// GET /new.json
func (h *TopicsHandler) New(w http.ResponseWriter, r *http.Request) {
	h.Latest(w, r)
}

// GET /t/{id}.json
func (h *TopicsHandler) GetTopic(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid topic id")
		return
	}
	t := h.Store.GetTopic(id)
	if t == nil {
		writeError(w, http.StatusNotFound, "topic not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// GET /t/external_id/{external_id}
func (h *TopicsHandler) GetTopicByExternalID(w http.ResponseWriter, r *http.Request) {
	extID := pathParam(r, "external_id")
	t := h.Store.GetTopicByExternalID(extID)
	if t == nil {
		writeError(w, http.StatusNotFound, "topic not found")
		return
	}
	// Discourse redirects to the topic URL
	w.Header().Set("Location", "/t/"+t.Slug+"/"+strconv.Itoa(t.ID))
	w.WriteHeader(http.StatusFound)
}

// GET /topics/created-by/{username}.json
func (h *TopicsHandler) TopicsByUser(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	username = strings.TrimSuffix(username, ".json")
	topics := h.Store.TopicsByUser(username)
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		Users: h.Store.UsersForTopics(topics),
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}

// PUT /t/{id}.json  — rename / recategorize
func (h *TopicsHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid topic id")
		return
	}
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// Handle nested "topic" key used by some SDKs
	if topicData, ok := body["topic"].(map[string]interface{}); ok {
		for k, v := range topicData {
			body[k] = v
		}
	}
	t, err := h.Store.UpdateTopic(id, body)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"basic_topic": t})
}

// PUT /t/{topic_id}/status
func (h *TopicsHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "topic_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid topic id")
		return
	}
	body, _ := decodeBody(r)
	status, _ := body["status"].(string)
	enabled := false
	switch v := body["enabled"].(type) {
	case bool:
		enabled = v
	case string:
		enabled = v == "true"
	case float64:
		enabled = v != 0
	}

	t, err := h.Store.UpdateTopicStatus(id, status, enabled)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.StatusResponse{Status: "ok", Topic: t})
}

// PUT /t/{topic_id}/change-timestamp
func (h *TopicsHandler) ChangeTimestamp(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /t/{id}.json
func (h *TopicsHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid topic id")
		return
	}
	if err := h.Store.DeleteTopic(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /t/{topic_id}/posts.json
func (h *TopicsHandler) TopicPosts(w http.ResponseWriter, r *http.Request) {
	topicID, ok := pathParamInt(r, "topic_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid topic id")
		return
	}

	var postIDs []int
	for _, s := range r.URL.Query()["post_ids[]"] {
		if id, err := strconv.Atoi(s); err == nil {
			postIDs = append(postIDs, id)
		}
	}

	posts := h.Store.GetTopicPosts(topicID, postIDs)
	stream := make([]int, len(posts))
	for i, p := range posts {
		stream[i] = p.ID
	}
	writeJSON(w, http.StatusOK, model.TopicPostsResponse{
		PostStream: model.PostStream{Posts: posts, Stream: stream},
		ID:         topicID,
	})
}

// POST /t/{topic_id}/change-owner.json
func (h *TopicsHandler) ChangeOwner(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /t/{topic_id}/notifications
func (h *TopicsHandler) SetNotificationLevel(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{topic_id}/bookmark.json
func (h *TopicsHandler) Bookmark(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /t/{topic_id}/remove_bookmarks.json
func (h *TopicsHandler) RemoveBookmark(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /t/{topic_id}/invite
func (h *TopicsHandler) InviteToTopic(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}
