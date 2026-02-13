package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// MiscHandler covers miscellaneous undocumented endpoints: hot topics,
// directory items, about, site basic info, drafts, bookmarks, published
// pages, sidebar sections, clicks, onebox, slugs, embed, presence, DND,
// emoji, hashtags, form templates, composer, etc.
type MiscHandler struct {
	Store *store.Store
}

// ---- Hot/Filter Topics ----

// GET /hot.json
func (h *MiscHandler) HotTopics(w http.ResponseWriter, r *http.Request) {
	topics := h.Store.ListTopics("hot")
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		TopicList: model.TopicList{
			CanCreateTopic: true,
			PerPage:        30,
			Topics:         topics,
		},
	})
}

// GET /filter
func (h *MiscHandler) FilterTopics(w http.ResponseWriter, r *http.Request) {
	h.HotTopics(w, r)
}

// ---- Directory Items ----

// GET /directory_items
func (h *MiscHandler) DirectoryItems(w http.ResponseWriter, r *http.Request) {
	users := h.Store.ListAllUsers()
	var items []map[string]interface{}
	for _, u := range users {
		if u.ID < 0 {
			continue // skip system user
		}
		items = append(items, map[string]interface{}{
			"id": u.ID,
			"user": map[string]interface{}{
				"id":              u.ID,
				"username":        u.Username,
				"name":            u.Name,
				"avatar_template": u.AvatarTemplate,
			},
			"likes_received": 0,
			"likes_given":    0,
			"topics_entered": 0,
			"topic_count":    0,
			"post_count":     0,
			"days_visited":   0,
			"posts_read":     0,
		})
	}
	if items == nil {
		items = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"directory_items":            items,
		"meta":                       map[string]interface{}{"total_rows_directory_items": len(items)},
		"total_rows_directory_items": len(items),
	})
}

// GET /directory-columns
func (h *MiscHandler) DirectoryColumns(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"directory_columns": []map[string]interface{}{
			{"id": 1, "name": "likes_received", "automatic": true, "position": 0},
			{"id": 2, "name": "likes_given", "automatic": true, "position": 1},
			{"id": 3, "name": "topic_count", "automatic": true, "position": 2},
			{"id": 4, "name": "post_count", "automatic": true, "position": 3},
		},
	})
}

// PUT /edit-directory-columns
func (h *MiscHandler) EditDirectoryColumns(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- About ----

// GET /about
func (h *MiscHandler) About(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"about": map[string]interface{}{
			"stats": map[string]interface{}{
				"topic_count": 3, "post_count": 4, "user_count": 3,
				"topics_7_days": 1, "topics_30_days": 3,
				"posts_7_days": 1, "posts_30_days": 4,
				"users_7_days": 0, "users_30_days": 3,
				"active_users_7_days": 3, "active_users_30_days": 3,
				"like_count": 7, "likes_7_days": 2, "likes_30_days": 7,
			},
			"description": "DTU Discourse — Digital Twin Universe",
			"title":       "DTU Discourse",
			"locale":      "en",
			"version":     "3.4.0",
			"admins":      []interface{}{},
			"moderators":  []interface{}{},
		},
	})
}

// GET /about/live_post_counts
func (h *MiscHandler) LivePostCounts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"post_count": 4, "topic_count": 3})
}

// GET /site/basic-info.json
func (h *MiscHandler) SiteBasicInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"title":                 "DTU Discourse",
		"description":           "Digital Twin Universe for Discourse",
		"logo_url":              "",
		"logo_small_url":        "",
		"apple_touch_icon_url":  "",
		"favicon_url":           "",
		"mobile_logo_url":       "",
		"external_auth_enabled": false,
	})
}

// GET /site/statistics
func (h *MiscHandler) SiteStatistics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"topic_count": 3, "post_count": 4, "user_count": 3,
		"topics_last_day": 0, "posts_last_day": 0, "users_last_day": 0,
		"active_users_last_day": 0, "likes_count": 7,
	})
}

// ---- Drafts ----

// GET /drafts
func (h *MiscHandler) ListDrafts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"drafts": []interface{}{}})
}

// POST /drafts
func (h *MiscHandler) CreateDraft(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"draft_sequence": 1})
}

// GET /drafts/{id}
func (h *MiscHandler) ShowDraft(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"draft": nil, "draft_sequence": 0})
}

// DELETE /drafts/{id}
func (h *MiscHandler) DeleteDraft(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Bookmarks (New API) ----

// POST /bookmarks
func (h *MiscHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id": 1, "created_at": time.Now().UTC(),
	})
}

// PUT /bookmarks/{id}
func (h *MiscHandler) UpdateBookmark(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /bookmarks/{id}
func (h *MiscHandler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /bookmarks/{id}/toggle_pin
func (h *MiscHandler) ToggleBookmarkPin(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /bookmarks/bulk
func (h *MiscHandler) BulkBookmarks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Published Pages ----

// GET /pub/check-slug
func (h *MiscHandler) CheckPublishedSlug(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"valid_slug": true})
}

// GET /pub/by-topic/{topic_id}
func (h *MiscHandler) GetPublishedPage(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"published_page": map[string]interface{}{"id": nil, "slug": "", "public": false},
	})
}

// PUT /pub/by-topic/{topic_id}
func (h *MiscHandler) UpdatePublishedPage(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"published_page": map[string]interface{}{"id": 1, "slug": "", "public": true},
	})
}

// DELETE /pub/by-topic/{topic_id}
func (h *MiscHandler) DeletePublishedPage(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Sidebar Sections ----

// GET /sidebar_sections
func (h *MiscHandler) ListSidebarSections(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"sidebar_sections": []interface{}{}})
}

// POST /sidebar_sections
func (h *MiscHandler) CreateSidebarSection(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"sidebar_section": map[string]interface{}{"id": 1, "title": "", "links": []interface{}{}},
	})
}

// PUT /sidebar_sections/{id}
func (h *MiscHandler) UpdateSidebarSection(w http.ResponseWriter, r *http.Request) {
	h.CreateSidebarSection(w, r)
}

// DELETE /sidebar_sections/{id}
func (h *MiscHandler) DeleteSidebarSection(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /sidebar_sections/reset/{id}
func (h *MiscHandler) ResetSidebarSection(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Clicks ----

// POST /clicks/track
func (h *MiscHandler) TrackClick(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Onebox ----

// GET /onebox
func (h *MiscHandler) Onebox(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// GET /inline-onebox
func (h *MiscHandler) InlineOnebox(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"inline-oneboxes": []interface{}{}})
}

// ---- Slugs ----

// POST /slugs
func (h *MiscHandler) GenerateSlug(w http.ResponseWriter, r *http.Request) {
	body, _ := decodeBody(r)
	name, _ := body["name"].(string)
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	writeJSON(w, http.StatusOK, map[string]interface{}{"slug": slug})
}

// ---- Embed ----

// GET /embed/topics
func (h *MiscHandler) EmbedTopics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"topic_list": map[string]interface{}{"topics": []interface{}{}}})
}

// GET /embed/comments
func (h *MiscHandler) EmbedComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div class='embed-comments'></div>"))
}

// GET /embed/count
func (h *MiscHandler) EmbedCount(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"count": 0})
}

// GET /embed/info
func (h *MiscHandler) EmbedInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// ---- Presence ----

// POST /presence/update
func (h *MiscHandler) PresenceUpdate(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /presence/get
func (h *MiscHandler) PresenceGet(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// ---- User Status / DND ----

// GET /user-status
func (h *MiscHandler) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// PUT /user-status
func (h *MiscHandler) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /user-status
func (h *MiscHandler) ClearUserStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /do-not-disturb
func (h *MiscHandler) EnableDND(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"ends_at": time.Now().Add(1 * time.Hour).UTC(),
	})
}

// DELETE /do-not-disturb
func (h *MiscHandler) DisableDND(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Emojis ----

// GET /admin/config/emoji
func (h *MiscHandler) ListCustomEmojis(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// POST /admin/config/emoji
func (h *MiscHandler) CreateCustomEmoji(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"name": "", "url": ""})
}

// DELETE /admin/config/emoji/{id}
func (h *MiscHandler) DeleteCustomEmoji(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// GET /emojis
func (h *MiscHandler) ListEmojis(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// ---- Hashtags ----

// GET /hashtags
func (h *MiscHandler) LookupHashtags(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// GET /hashtags/search
func (h *MiscHandler) SearchHashtags(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"results": []interface{}{}})
}

// ---- Composer ----

// GET /composer/mentions
func (h *MiscHandler) ComposerMentions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// GET /composer_messages
func (h *MiscHandler) ComposerMessages(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"composer_messages": []interface{}{}})
}

// ---- Export CSV ----

// POST /export_csv/export_entity
func (h *MiscHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Server Status ----

// GET /srv/status
func (h *MiscHandler) ServerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}

// ---- Permalink Check ----

// GET /permalink-check
func (h *MiscHandler) PermalinkCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"found":         false,
		"internal":      false,
		"target_url":    "",
	})
}

// ---- Pageview ----

// POST /pageview
func (h *MiscHandler) Pageview(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Search Extended ----

// GET /search/query
func (h *MiscHandler) SearchQuery(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("term")
	if q == "" {
		q = r.URL.Query().Get("q")
	}
	result := h.Store.Search(q)
	writeJSON(w, http.StatusOK, result)
}

// POST /search/click
func (h *MiscHandler) SearchClick(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}
