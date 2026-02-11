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
	Topics   *TopicsHandler
	Extended *ExtendedTopicsHandler
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

	// /t/id_for/{slug}
	if first == "id_for" && len(parts) >= 2 && d.Extended != nil {
		r.SetPathValue("slug", parts[1])
		d.Extended.IDForSlug(w, r)
		return
	}

	// Strip .json from first segment for topic ID
	idStr := strings.TrimSuffix(first, ".json")
	topicID, err := strconv.Atoi(idStr)

	// /t/{slug}/{id} — first segment is non-numeric slug, second is numeric ID
	if err != nil && len(parts) >= 2 {
		secondStr := strings.TrimSuffix(parts[1], ".json")
		if tid2, err2 := strconv.Atoi(secondStr); err2 == nil {
			r.SetPathValue("id", strconv.Itoa(tid2))
			d.Topics.GetTopic(w, r)
			return
		}
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	tid := strconv.Itoa(topicID)

	if len(parts) == 1 {
		// /t/{id} or /t/{id}.json
		r.SetPathValue("id", first)
		d.Topics.GetTopic(w, r)
		return
	}

	second := parts[1]

	// /t/{topic_id}/posts.json
	if second == "posts.json" || second == "posts" {
		r.SetPathValue("topic_id", tid)
		d.Topics.TopicPosts(w, r)
		return
	}

	// Numeric second segment means /t/{id}/{post_number} (topic with post)
	secondStr := strings.TrimSuffix(second, ".json")
	if _, err := strconv.Atoi(secondStr); err == nil {
		r.SetPathValue("id", tid)
		d.Topics.GetTopic(w, r)
		return
	}

	// Extended GET routes
	if d.Extended != nil {
		r.SetPathValue("id", tid)
		switch second {
		case "post_ids", "post_ids.json":
			d.Extended.PostIDs(w, r)
			return
		case "excerpts", "excerpts.json":
			d.Extended.Excerpts(w, r)
			return
		case "view-stats.json", "view-stats":
			d.Extended.ViewStats(w, r)
			return
		}
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

	// PUT /t/{slug}/{id} — first segment is non-numeric slug, second is numeric ID
	if err != nil && len(parts) >= 2 {
		secondStr := strings.TrimSuffix(parts[1], ".json")
		if tid2, err2 := strconv.Atoi(secondStr); err2 == nil {
			r.SetPathValue("id", strconv.Itoa(tid2))
			d.Topics.UpdateTopic(w, r)
			return
		}
		writeError(w, http.StatusNotFound, "not found")
		return
	}
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

	// Core routes
	switch second {
	case "status":
		r.SetPathValue("topic_id", tid)
		d.Topics.UpdateStatus(w, r)
		return
	case "change-timestamp":
		r.SetPathValue("topic_id", tid)
		d.Topics.ChangeTimestamp(w, r)
		return
	case "bookmark.json", "bookmark":
		r.SetPathValue("topic_id", tid)
		d.Topics.Bookmark(w, r)
		return
	case "remove_bookmarks.json", "remove_bookmarks":
		r.SetPathValue("topic_id", tid)
		d.Topics.RemoveBookmark(w, r)
		return
	}

	// Extended PUT routes
	if d.Extended != nil {
		r.SetPathValue("id", tid)
		switch second {
		case "archive-message":
			d.Extended.ArchiveMessage(w, r)
			return
		case "move-to-inbox":
			d.Extended.MoveToInbox(w, r)
			return
		case "convert-topic":
			if len(parts) >= 3 {
				r.SetPathValue("type", parts[2])
			}
			d.Extended.ConvertTopic(w, r)
			return
		case "publish":
			d.Extended.Publish(w, r)
			return
		case "reset-bump-date":
			d.Extended.ResetBumpDate(w, r)
			return
		case "clear-pin":
			d.Extended.ClearPin(w, r)
			return
		case "re-pin":
			d.Extended.RePin(w, r)
			return
		case "mute":
			d.Extended.Mute(w, r)
			return
		case "unmute":
			d.Extended.Unmute(w, r)
			return
		case "make-banner":
			d.Extended.MakeBanner(w, r)
			return
		case "remove-banner":
			d.Extended.RemoveBanner(w, r)
			return
		case "remove-allowed-user":
			d.Extended.RemoveAllowedUser(w, r)
			return
		case "remove-allowed-group":
			d.Extended.RemoveAllowedGroup(w, r)
			return
		case "recover":
			d.Extended.Recover(w, r)
			return
		case "tags":
			d.Extended.UpdateTags(w, r)
			return
		case "slow_mode":
			d.Extended.SlowMode(w, r)
			return
		}
	}

	writeError(w, http.StatusNotFound, "not found")
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

	// Core routes
	switch second {
	case "change-owner.json", "change-owner":
		r.SetPathValue("topic_id", tid)
		d.Topics.ChangeOwner(w, r)
		return
	case "notifications", "notifications.json":
		r.SetPathValue("topic_id", tid)
		d.Topics.SetNotificationLevel(w, r)
		return
	case "invite", "invite.json":
		r.SetPathValue("topic_id", tid)
		d.Topics.InviteToTopic(w, r)
		return
	}

	// Extended POST routes
	if d.Extended != nil {
		r.SetPathValue("id", tid)
		switch second {
		case "move-posts":
			d.Extended.MovePosts(w, r)
			return
		case "merge-topic":
			d.Extended.MergeTopic(w, r)
			return
		case "invite-group":
			d.Extended.InviteGroup(w, r)
			return
		}
	}

	writeError(w, http.StatusNotFound, "not found")
}

// ServeDELETE handles DELETE /t/{rest...}
func (d *TopicSubRouter) ServeDELETE(w http.ResponseWriter, r *http.Request) {
	rest := r.PathValue("rest")
	parts := strings.Split(rest, "/")

	if len(parts) >= 2 {
		first := parts[0]
		idStr := strings.TrimSuffix(first, ".json")
		second := parts[1]

		// Extended DELETE routes
		if d.Extended != nil && second == "timings" {
			r.SetPathValue("id", idStr)
			d.Extended.DestroyTimings(w, r)
			return
		}
	}

	r.SetPathValue("id", rest)
	d.Topics.DeleteTopic(w, r)
}
