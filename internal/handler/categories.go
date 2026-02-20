package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

type CategoriesHandler struct {
	Store *store.Store
}

// GET /categories.json
func (h *CategoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	cats := h.Store.ListCategories()
	writeJSON(w, http.StatusOK, model.CategoryListResponse{
		CategoryList: model.CategoryList{
			CanCreateCategory: true,
			CanCreateTopic:    true,
			Categories:        cats,
		},
	})
}

// POST /categories.json  or  POST /categories
func (h *CategoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	name, _ := body["name"].(string)
	slug, _ := body["slug"].(string)
	color, _ := body["color"].(string)
	textColor, _ := body["text_color"].(string)

	if name == "" {
		writeError(w, http.StatusUnprocessableEntity, "name is required")
		return
	}

	cat, err := h.Store.CreateCategory(name, slug, color, textColor)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.CategoryResponse{Category: *cat})
}

// PUT /categories/{id}.json  or  PUT /categories/{id}
func (h *CategoriesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	body, err := decodeBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cat, err := h.Store.UpdateCategory(id, body)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.CategoryResponse{Category: *cat})
}

// DELETE /categories/{id}
func (h *CategoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	if err := h.Store.DeleteCategory(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /c/{id}/show  or  /c/{id}/show.json
func (h *CategoriesHandler) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	cat := h.Store.GetCategory(id)
	if cat == nil {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	writeJSON(w, http.StatusOK, model.CategoryResponse{Category: *cat})
}

// GET /c/{slug}/{id}.json  — list topics in category
// Also handles /c/{slug}.json where {id} captures slug.json
func (h *CategoriesHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "id")
	if ok {
		topics := h.Store.TopicsByCategory(id)
		writeJSON(w, http.StatusOK, model.TopicListResponse{
			TopicList: model.TopicList{
				CanCreateTopic: true,
				PerPage:        30,
				Topics:         topics,
			},
		})
		return
	}
	// Try slug lookup (e.g. /c/general.json where {id} = "general.json")
	slug := strings.TrimSuffix(pathParam(r, "id"), ".json")
	cat := h.Store.GetCategoryBySlug(slug)
	if cat == nil {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	topics := h.Store.TopicsByCategory(cat.ID)
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}

// GET /c/{category_slug}/l/latest.json
// Also handles /c/{categoryId}/l/latest.json where category_slug is actually a numeric ID
func (h *CategoriesHandler) LatestTopics(w http.ResponseWriter, r *http.Request) {
	slug := pathParam(r, "category_slug")
	slug = strings.TrimSuffix(slug, ".json")
	// Try numeric ID first
	if id, err := strconv.Atoi(slug); err == nil {
		topics := h.Store.TopicsByCategory(id)
		writeJSON(w, http.StatusOK, model.TopicListResponse{
			Users: h.Store.UsersForTopics(topics),
			TopicList: model.TopicList{
				CanCreateTopic: true,
				PerPage:        30,
				Topics:         topics,
			},
		})
		return
	}
	cat := h.Store.GetCategoryBySlug(slug)
	if cat == nil {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	topics := h.Store.TopicsByCategory(cat.ID)
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		Users: h.Store.UsersForTopics(topics),
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}

// GET /c/{category_slug}/l/top.json
func (h *CategoriesHandler) TopTopics(w http.ResponseWriter, r *http.Request) {
	h.LatestTopics(w, r) // Same behavior for DTU
}

// GET /c/{category_slug}/l/new.json
func (h *CategoriesHandler) NewTopics(w http.ResponseWriter, r *http.Request) {
	h.LatestTopics(w, r) // Same behavior for DTU
}

// GET /c/{category_slug}/{category_id}/l/latest.json (and hot/top/new variants)
func (h *CategoriesHandler) LatestTopicsBySlugAndID(w http.ResponseWriter, r *http.Request) {
	id, ok := pathParamInt(r, "category_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	topics := h.Store.TopicsByCategory(id)
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		Users: h.Store.UsersForTopics(topics),
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}

// POST /categories/reorder
func (h *CategoriesHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /category/{category_id}/notifications
func (h *CategoriesHandler) SetNotificationLevel(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}
