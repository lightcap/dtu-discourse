package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type TagsHandler struct {
	Store *store.Store
}

// GET /tags.json
func (h *TagsHandler) List(w http.ResponseWriter, r *http.Request) {
	tags := h.Store.ListTags()
	writeJSON(w, http.StatusOK, model.TagListResponse{Tags: tags})
}

// GET /tag/{tag}
func (h *TagsHandler) Show(w http.ResponseWriter, r *http.Request) {
	tagName := pathParam(r, "tag")
	tag := h.Store.GetTag(tagName)
	if tag == nil {
		writeError(w, http.StatusNotFound, "tag not found")
		return
	}
	topics := h.Store.TopicsByTag(tagName)
	writeJSON(w, http.StatusOK, model.TagResponse{
		Tag: *tag,
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}
