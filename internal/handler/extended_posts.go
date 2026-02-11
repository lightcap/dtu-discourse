package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// ExtendedPostsHandler handles undocumented post operations.
type ExtendedPostsHandler struct {
	Store *store.Store
}

// PUT /posts/{id}/recover
func (h *ExtendedPostsHandler) Recover(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/{id}/rebake
func (h *ExtendedPostsHandler) Rebake(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/{id}/locked
func (h *ExtendedPostsHandler) Locked(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/{id}/post_type
func (h *ExtendedPostsHandler) PostType(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/{id}/unhide
func (h *ExtendedPostsHandler) Unhide(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/{id}/notice
func (h *ExtendedPostsHandler) Notice(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /posts/{id}/revisions/latest
func (h *ExtendedPostsHandler) LatestRevision(w http.ResponseWriter, r *http.Request) {
	id, _ := pathParamInt(r, "id")
	p := h.Store.GetPost(id)
	if p == nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"post_revision": map[string]interface{}{
			"post_id":        id,
			"version":        p.Version,
			"revision_number": p.Version,
			"created_at":     time.Now().UTC(),
			"body_changes": map[string]interface{}{
				"inline": p.Cooked,
			},
		},
	})
}

// GET /posts/{id}/revisions/{revision}
func (h *ExtendedPostsHandler) Revision(w http.ResponseWriter, r *http.Request) {
	h.LatestRevision(w, r)
}

// PUT /posts/{id}/revisions/{revision}/hide
func (h *ExtendedPostsHandler) HideRevision(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/{id}/revisions/{revision}/show
func (h *ExtendedPostsHandler) ShowRevision(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/{id}/revisions/{revision}/revert
func (h *ExtendedPostsHandler) RevertRevision(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /posts/{id}/revisions/permanently_delete
func (h *ExtendedPostsHandler) PermanentlyDeleteRevisions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /posts/destroy_many
func (h *ExtendedPostsHandler) DestroyMany(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/merge_posts
func (h *ExtendedPostsHandler) MergePosts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /posts/by_number/{topic_id}/{post_number}
func (h *ExtendedPostsHandler) ByNumber(w http.ResponseWriter, r *http.Request) {
	topicID, _ := pathParamInt(r, "topic_id")
	postNumStr := pathParam(r, "post_number")
	postNum, _ := strconv.Atoi(postNumStr)
	posts := h.Store.GetTopicPosts(topicID, nil)
	for _, p := range posts {
		if p.PostNumber == postNum {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}
	writeError(w, http.StatusNotFound, "post not found")
}

// GET /posts/{id}/reply-history
func (h *ExtendedPostsHandler) ReplyHistory(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// GET /posts/{id}/reply-ids
func (h *ExtendedPostsHandler) ReplyIDs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// GET /posts/{id}/cooked
func (h *ExtendedPostsHandler) Cooked(w http.ResponseWriter, r *http.Request) {
	id, _ := pathParamInt(r, "id")
	p := h.Store.GetPost(id)
	if p == nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"cooked": p.Cooked})
}

// GET /posts/{id}/raw
func (h *ExtendedPostsHandler) Raw(w http.ResponseWriter, r *http.Request) {
	id, _ := pathParamInt(r, "id")
	p := h.Store.GetPost(id)
	if p == nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(p.Raw))
}

// GET /raw/{topic_id}/{post_number}
func (h *ExtendedPostsHandler) RawByNumber(w http.ResponseWriter, r *http.Request) {
	topicID, _ := pathParamInt(r, "topic_id")
	postNumStr := pathParam(r, "post_number")
	postNum, _ := strconv.Atoi(strings.TrimSuffix(postNumStr, ".json"))
	if postNum == 0 {
		postNum = 1
	}
	posts := h.Store.GetTopicPosts(topicID, nil)
	for _, p := range posts {
		if p.PostNumber == postNum {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(p.Raw))
			return
		}
	}
	writeError(w, http.StatusNotFound, "post not found")
}

// GET /posts/{username}/deleted
func (h *ExtendedPostsHandler) Deleted(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// GET /posts/{username}/pending
func (h *ExtendedPostsHandler) Pending(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"pending_posts": []interface{}{}})
}

// GET /posts/{id}/replies
func (h *ExtendedPostsHandler) Replies(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// DELETE /posts/{id}/bookmark
func (h *ExtendedPostsHandler) RemoveBookmark(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}
