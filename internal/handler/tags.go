package handler

import (
	"net/http"
	"strings"

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

// GET /tag/{tag}/l/latest.json (and hot/top/new variants)
func (h *TagsHandler) TopicsByTag(w http.ResponseWriter, r *http.Request) {
	tagName := strings.TrimSuffix(pathParam(r, "tag"), ".json")
	topics := h.Store.TopicsByTag(tagName)
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		Users: h.Store.UsersForTopics(topics),
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}

// GET /tags/c/{category_slug}/{category_id}/{tag}/l/latest.json (and variants)
func (h *TagsHandler) TopicsByCategoryAndTag(w http.ResponseWriter, r *http.Request) {
	tagName := strings.TrimSuffix(pathParam(r, "tag"), ".json")
	catID, ok := pathParamInt(r, "category_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	// Filter topics that match both category and tag
	allTopics := h.Store.TopicsByTag(tagName)
	var filtered []model.Topic
	for _, t := range allTopics {
		if t.CategoryID == catID {
			filtered = append(filtered, t)
		}
	}
	if filtered == nil {
		filtered = []model.Topic{}
	}
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		Users: h.Store.UsersForTopics(filtered),
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         filtered,
		},
	})
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
