package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// ============================================================================
// Tag Groups
// ============================================================================

type TagGroupsHandler struct {
	Store *store.Store
}

// GET /tag_groups
func (h *TagGroupsHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tag_groups": []map[string]interface{}{
			{
				"id": 1, "name": "Topic Types",
				"tag_names":    []string{"question", "discussion", "announcement"},
				"one_per_topic": true,
			},
		},
	})
}

// GET /tag_groups/{id}
func (h *TagGroupsHandler) Show(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tag_group": map[string]interface{}{
			"id": 1, "name": "Topic Types",
			"tag_names":    []string{"question", "discussion", "announcement"},
			"one_per_topic": true,
		},
	})
}

// POST /tag_groups
func (h *TagGroupsHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := decodeBody(r)
	tg, _ := body["tag_group"].(map[string]interface{})
	name := ""
	if tg != nil {
		name, _ = tg["name"].(string)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tag_group": map[string]interface{}{
			"id":           2,
			"name":         name,
			"tag_names":    []string{},
			"one_per_topic": false,
		},
	})
}

// PUT /tag_groups/{id}
func (h *TagGroupsHandler) Update(w http.ResponseWriter, r *http.Request) {
	h.Show(w, r)
}

// DELETE /tag_groups/{id}
func (h *TagGroupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ============================================================================
// Extended Notifications
// ============================================================================

type ExtendedNotificationsHandler struct {
	Store *store.Store
}

// PUT /notifications/mark-read
func (h *ExtendedNotificationsHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /notifications/totals
func (h *ExtendedNotificationsHandler) Totals(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"total_notifications":           0,
		"unread_notifications":          0,
		"unread_high_priority":          0,
		"read_first_notification":       true,
		"seen_notification_id":          0,
	})
}

// GET /notifications/{id}
func (h *ExtendedNotificationsHandler) Show(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"notification": map[string]interface{}{
			"id":                1,
			"notification_type": 1,
			"read":              false,
			"created_at":        "2024-01-01T00:00:00.000Z",
			"data":              map[string]interface{}{},
		},
	})
}

// PUT /notifications/{id}
func (h *ExtendedNotificationsHandler) Update(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /notifications/{id}
func (h *ExtendedNotificationsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ============================================================================
// Extended Groups
// ============================================================================

type ExtendedGroupsHandler struct {
	Store *store.Store
}

// PUT /groups/{group}/join
func (h *ExtendedGroupsHandler) Join(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /groups/{group}/leave
func (h *ExtendedGroupsHandler) Leave(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /groups/{group}/request_membership
func (h *ExtendedGroupsHandler) RequestMembership(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"relative_url": "/groups/staff",
	})
}

// GET /groups/{group}/requests
func (h *ExtendedGroupsHandler) ListRequests(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"requests":   []interface{}{},
		"meta":       map[string]interface{}{"total": 0, "limit": 50, "offset": 0},
	})
}

// PUT /groups/{group}/handle_membership_request
func (h *ExtendedGroupsHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /groups/{group}/logs
func (h *ExtendedGroupsHandler) Logs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"logs": []interface{}{},
		"meta": map[string]interface{}{"total": 0, "limit": 50, "offset": 0},
	})
}

// GET /groups/{group}/permissions
func (h *ExtendedGroupsHandler) Permissions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"permissions": []interface{}{},
	})
}

// GET /groups/{group}/mentionable
func (h *ExtendedGroupsHandler) Mentionable(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"mentionable": true,
	})
}

// GET /groups/{group}/messageable
func (h *ExtendedGroupsHandler) Messageable(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"messageable": true,
	})
}

// GET /groups/{group}/counts
func (h *ExtendedGroupsHandler) Counts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"counts": map[string]interface{}{
			"posts":   0,
			"topics":  0,
			"members": 0,
		},
	})
}

// GET /groups/{group}/topics
func (h *ExtendedGroupsHandler) Topics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		TopicList: model.TopicList{Topics: []model.Topic{}},
	})
}

// ============================================================================
// Extended Categories
// ============================================================================

type ExtendedCategoriesHandler struct {
	Store *store.Store
}

// GET /categories/search
func (h *ExtendedCategoriesHandler) Search(w http.ResponseWriter, r *http.Request) {
	cats := h.Store.ListCategories()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"categories": cats,
	})
}

// GET /categories/find
func (h *ExtendedCategoriesHandler) Find(w http.ResponseWriter, r *http.Request) {
	h.Search(w, r)
}

// GET /categories_and_latest
func (h *ExtendedCategoriesHandler) CategoriesAndLatest(w http.ResponseWriter, r *http.Request) {
	cats := h.Store.ListCategories()
	topics := h.Store.ListTopics("latest")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"category_list": map[string]interface{}{
			"can_create_category": true,
			"can_create_topic":    true,
			"categories":          cats,
		},
		"topic_list": map[string]interface{}{
			"can_create_topic": true,
			"per_page":         30,
			"topics":           topics,
		},
	})
}

// GET /categories_and_top
func (h *ExtendedCategoriesHandler) CategoriesAndTop(w http.ResponseWriter, r *http.Request) {
	h.CategoriesAndLatest(w, r)
}

// POST /categories/{id}/move
func (h *ExtendedCategoriesHandler) Move(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /c/{slug}/visible_groups
func (h *ExtendedCategoriesHandler) VisibleGroups(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"groups": []interface{}{},
	})
}

// ============================================================================
// Extended Tags
// ============================================================================

type ExtendedTagsHandler struct {
	Store *store.Store
}

// GET /tags/filter/search
func (h *ExtendedTagsHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	tags := h.Store.ListTags()
	var results []map[string]interface{}
	for _, t := range tags {
		if q == "" || containsCI(t.TagName, q) || containsCI(t.Name, q) {
			results = append(results, map[string]interface{}{
				"id":    t.TagName,
				"name":  t.TagName,
				"text":  t.TagName,
				"count": t.Count,
			})
		}
	}
	if results == nil {
		results = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"results": results,
	})
}

// POST /tags
func (h *ExtendedTagsHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := decodeBody(r)
	name, _ := body["name"].(string)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tag": map[string]interface{}{
			"id":    name,
			"text":  name,
			"count": 0,
		},
	})
}

// PUT /tags/{tag}
func (h *ExtendedTagsHandler) Update(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tag": map[string]interface{}{
			"id":    r.PathValue("tag"),
			"text":  r.PathValue("tag"),
			"count": 0,
		},
	})
}

// DELETE /tags/{tag}
func (h *ExtendedTagsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /tags/{tag}/info
func (h *ExtendedTagsHandler) Info(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tag_info": map[string]interface{}{
			"id":         tag,
			"name":       tag,
			"topic_count": 0,
			"staff":      false,
			"synonyms":   []interface{}{},
			"tag_group_names": []interface{}{},
		},
		"categories": []interface{}{},
	})
}

// GET /tag/{tag}/notifications
func (h *ExtendedTagsHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tag_notification": map[string]interface{}{
			"id":                 r.PathValue("tag"),
			"notification_level": 1,
		},
	})
}

// PUT /tag/{tag}/notifications
func (h *ExtendedTagsHandler) SetNotifications(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /tags/{tag}/synonyms
func (h *ExtendedTagsHandler) AddSynonym(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /tags/{tag}/synonyms/{synonym}
func (h *ExtendedTagsHandler) RemoveSynonym(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /tags/unused
func (h *ExtendedTagsHandler) Unused(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tags": []interface{}{},
	})
}

// helper
func containsCI(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) == 0 ||
		(func() bool {
			sl := make([]byte, len(s))
			for i := range s {
				if s[i] >= 'A' && s[i] <= 'Z' {
					sl[i] = s[i] + 32
				} else {
					sl[i] = s[i]
				}
			}
			tl := make([]byte, len(substr))
			for i := range substr {
				if substr[i] >= 'A' && substr[i] <= 'Z' {
					tl[i] = substr[i] + 32
				} else {
					tl[i] = substr[i]
				}
			}
			for i := 0; i <= len(sl)-len(tl); i++ {
				match := true
				for j := range tl {
					if sl[i+j] != tl[j] {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
			return false
		})())
}
