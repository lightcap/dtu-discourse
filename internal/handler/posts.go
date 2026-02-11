package handler

import (
	"net/http"
	"strconv"

	"github.com/lightcap/dtu-discourse/internal/middleware"
	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type PostsHandler struct {
	Store *store.Store
}

// POST /posts  — create post or topic
func (h *PostsHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	raw, _ := body["raw"].(string)
	if raw == "" {
		writeError(w, http.StatusUnprocessableEntity, "raw is required")
		return
	}

	username := middleware.GetUsername(r)
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusForbidden, "user not found")
		return
	}

	// If title is present, this creates a new topic
	title, hasTitle := body["title"].(string)
	if hasTitle && title != "" {
		categoryID := 1
		if v, ok := body["category"].(float64); ok {
			categoryID = int(v)
		} else if v, ok := body["category"].(string); ok {
			if id, err := strconv.Atoi(v); err == nil {
				categoryID = id
			}
		}

		var tags []string
		if v, ok := body["tags"].([]interface{}); ok {
			for _, t := range v {
				if s, ok := t.(string); ok {
					tags = append(tags, s)
				}
			}
		}

		archetype, _ := body["archetype"].(string)
		targetUsernames, _ := body["target_usernames"].(string)
		if archetype == "" && targetUsernames != "" {
			archetype = "private_message"
		}

		_, post, err := h.Store.CreateTopic(title, raw, categoryID, u.ID, tags, archetype)
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, post)
		return
	}

	// Otherwise, reply to existing topic
	topicID := 0
	if v, ok := body["topic_id"].(float64); ok {
		topicID = int(v)
	} else if v, ok := body["topic_id"].(string); ok {
		topicID, _ = strconv.Atoi(v)
	}

	if topicID == 0 {
		writeError(w, http.StatusUnprocessableEntity, "topic_id is required for replies")
		return
	}

	var replyTo *int
	if v, ok := body["reply_to_post_number"].(float64); ok {
		i := int(v)
		replyTo = &i
	}

	post, err := h.Store.CreatePost(topicID, raw, u.ID, replyTo)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, post)
}

// GET /posts/{id}.json
func (h *PostsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}
	p := h.Store.GetPost(id)
	if p == nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// GET /posts.json
func (h *PostsHandler) List(w http.ResponseWriter, r *http.Request) {
	posts := h.Store.ListPosts()
	writeJSON(w, http.StatusOK, model.PostListResponse{LatestPosts: posts})
}

// PUT /posts/{id}
func (h *PostsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Handle nested "post" key
	if postData, ok := body["post"].(map[string]interface{}); ok {
		for k, v := range postData {
			body[k] = v
		}
	}

	raw, _ := body["raw"].(string)
	if raw == "" {
		writeError(w, http.StatusUnprocessableEntity, "raw is required")
		return
	}
	p, err := h.Store.UpdatePost(id, raw)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"post": p})
}

// DELETE /posts/{id}.json
func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}
	if err := h.Store.DeletePost(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /posts/{id}/wiki
func (h *PostsHandler) Wikify(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}
	body, _ := decodeBody(r)
	wiki := true
	if v, ok := body["wiki"].(bool); ok {
		wiki = v
	} else if v, ok := body["wiki"].(string); ok {
		wiki = v == "true"
	}
	p, err := h.Store.WikifyPost(id, wiki)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// POST /post_actions
func (h *PostsHandler) CreateAction(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	postID := 0
	if v, ok := body["id"].(float64); ok {
		postID = int(v)
	} else if v, ok := body["id"].(string); ok {
		postID, _ = strconv.Atoi(v)
	}
	actionType := 0
	if v, ok := body["post_action_type_id"].(float64); ok {
		actionType = int(v)
	} else if v, ok := body["post_action_type_id"].(string); ok {
		actionType, _ = strconv.Atoi(v)
	}

	pa, err := h.Store.CreatePostAction(postID, actionType)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pa)
}

// DELETE /post_actions/{id}.json
func (h *PostsHandler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid post action id")
		return
	}
	if err := h.Store.DeletePostAction(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /post_action_users.json
func (h *PostsHandler) ActionUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"post_action_users": []interface{}{}})
}
