package handler

import (
	"net/http"
	"strconv"
	"strings"
)

// TopicSubRouter handles all routes under /t/{...}/... where the Go ServeMux
// can't resolve ambiguity between patterns like /t/{id}/posts.json and
// /t/external_id/{external_id}. It dispatches to the correct handler.
type TopicSubRouter struct {
	Topics *TopicsHandler
}

// ServeGET handles GET /t/{rest...}
func (d *TopicSubRouter) ServeGET(w http.ResponseWriter, r *http.Request) {
	rest := r.PathValue("rest")
	parts := strings.Split(rest, "/")

	if len(parts) == 0 {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	first := parts[0]

	// /t/external_id/{external_id}
	if first == "external_id" && len(parts) >= 2 {
		r.SetPathValue("external_id", parts[1])
		d.Topics.GetTopicByExternalID(w, r)
		return
	}

	// Strip .json from first segment for topic ID
	idStr := strings.TrimSuffix(first, ".json")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	if len(parts) == 1 {
		// /t/{id} or /t/{id}.json
		r.SetPathValue("id", first)
		d.Topics.GetTopic(w, r)
		return
	}

	second := parts[1]

	// /t/{topic_id}/posts.json
	if second == "posts.json" || second == "posts" {
		r.SetPathValue("topic_id", strconv.Itoa(topicID))
		d.Topics.TopicPosts(w, r)
		return
	}

	writeError(w, http.StatusNotFound, "not found")
}

// ServePUT handles PUT /t/{rest...}
func (d *TopicSubRouter) ServePUT(w http.ResponseWriter, r *http.Request) {
	rest := r.PathValue("rest")
	parts := strings.Split(rest, "/")

	if len(parts) == 0 {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	first := parts[0]
	idStr := strings.TrimSuffix(first, ".json")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	tid := strconv.Itoa(topicID)

	if len(parts) == 1 {
		// /t/{id} or /t/{id}.json — update topic
		r.SetPathValue("id", first)
		d.Topics.UpdateTopic(w, r)
		return
	}

	second := parts[1]

	switch second {
	case "status":
		r.SetPathValue("topic_id", tid)
		d.Topics.UpdateStatus(w, r)
	case "change-timestamp":
		r.SetPathValue("topic_id", tid)
		d.Topics.ChangeTimestamp(w, r)
	case "bookmark.json", "bookmark":
		r.SetPathValue("topic_id", tid)
		d.Topics.Bookmark(w, r)
	case "remove_bookmarks.json", "remove_bookmarks":
		r.SetPathValue("topic_id", tid)
		d.Topics.RemoveBookmark(w, r)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

// ServePOST handles POST /t/{rest...}
func (d *TopicSubRouter) ServePOST(w http.ResponseWriter, r *http.Request) {
	rest := r.PathValue("rest")
	parts := strings.Split(rest, "/")

	if len(parts) < 2 {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	first := parts[0]
	idStr := strings.TrimSuffix(first, ".json")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	tid := strconv.Itoa(topicID)
	second := parts[1]

	switch {
	case second == "change-owner.json" || second == "change-owner":
		r.SetPathValue("topic_id", tid)
		d.Topics.ChangeOwner(w, r)
	case second == "notifications":
		r.SetPathValue("topic_id", tid)
		d.Topics.SetNotificationLevel(w, r)
	case second == "invite":
		r.SetPathValue("topic_id", tid)
		d.Topics.InviteToTopic(w, r)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

// ServeDELETE handles DELETE /t/{rest...}
func (d *TopicSubRouter) ServeDELETE(w http.ResponseWriter, r *http.Request) {
	rest := r.PathValue("rest")
	r.SetPathValue("id", rest)
	d.Topics.DeleteTopic(w, r)
}
