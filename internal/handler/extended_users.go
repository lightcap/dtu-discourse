package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// ExtendedUsersHandler covers user profile endpoints not in the core API:
// avatar, preferences, summary, card, activity, bookmarks, user search, etc.
type ExtendedUsersHandler struct {
	Store *store.Store
}

// ---- Avatar ----

// PUT /u/{username}/preferences/avatar/pick
func (h *ExtendedUsersHandler) PickAvatar(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- User Summary ----

// GET /u/{username}/summary
func (h *ExtendedUsersHandler) Summary(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSuffix(r.PathValue("username"), ".json")
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_summary": map[string]interface{}{
			"likes_given":          0,
			"likes_received":       0,
			"topics_entered":       0,
			"posts_read_count":     0,
			"days_visited":         1,
			"topic_count":          0,
			"post_count":           0,
			"time_read":            0,
			"recent_time_read":     0,
			"bookmark_count":       0,
			"can_see_summary_stats": true,
			"top_categories":       []interface{}{},
			"top_replies":          []interface{}{},
			"top_topics":           []interface{}{},
			"badges":               []interface{}{},
			"top_links":            []interface{}{},
			"most_liked_by_users":  []interface{}{},
			"most_liked_users":     []interface{}{},
			"most_replied_to_users": []interface{}{},
		},
		"users": []interface{}{},
	})
}

// ---- User Card ----

// GET /u/{username}/card.json
func (h *ExtendedUsersHandler) Card(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSuffix(r.PathValue("username"), ".json")
	username = strings.TrimSuffix(username, "/card")
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":              u.ID,
			"username":        u.Username,
			"name":            u.Name,
			"avatar_template": u.AvatarTemplate,
			"trust_level":     u.TrustLevel,
			"admin":           u.Admin,
			"moderator":       u.Moderator,
			"created_at":      u.CreatedAt,
			"last_seen_at":    u.LastSeenAt,
			"badge_count":     0,
			"time_read":       0,
			"recent_time_read": 0,
			"primary_group_id": nil,
			"flair_group_id":   nil,
			"featured_topic":   nil,
		},
	})
}

// ---- User Activity ----

// GET /u/{username}/activity
func (h *ExtendedUsersHandler) Activity(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_actions": []interface{}{},
	})
}

// GET /u/{username}/activity/topics
func (h *ExtendedUsersHandler) ActivityTopics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		TopicList: model.TopicList{Topics: []model.Topic{}},
	})
}

// GET /u/{username}/activity/replies
func (h *ExtendedUsersHandler) ActivityReplies(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_actions": []interface{}{},
	})
}

// ---- User Bookmarks ----

// GET /u/{username}/bookmarks
func (h *ExtendedUsersHandler) Bookmarks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_bookmark_list": map[string]interface{}{
			"bookmarks":     []interface{}{},
			"more_bookmarks_url": nil,
		},
	})
}

// ---- User Badges ----

// GET /u/{username}/badges
func (h *ExtendedUsersHandler) Badges(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"badges":      []interface{}{},
		"user_badges": []interface{}{},
	})
}

// ---- User Preferences ----

// PUT /u/{username}/preferences/badge_title
func (h *ExtendedUsersHandler) SetBadgeTitle(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /u/{username}/preferences/categories
func (h *ExtendedUsersHandler) SetCategoryPreferences(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /u/{username}/preferences/tags
func (h *ExtendedUsersHandler) SetTagPreferences(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /u/{username}/clear-featured-topic
func (h *ExtendedUsersHandler) ClearFeaturedTopic(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /u/{username}/feature-topic
func (h *ExtendedUsersHandler) FeatureTopic(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- User Emails ----

// GET /u/{username}/emails
func (h *ExtendedUsersHandler) Emails(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSuffix(r.PathValue("username"), ".json")
	u := h.Store.GetUserByUsername(username)
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"email":             u.Email,
		"secondary_emails":  []interface{}{},
		"unconfirmed_emails": []interface{}{},
		"associated_accounts": []interface{}{},
	})
}

// ---- User Search ----

// GET /u/search/users
func (h *ExtendedUsersHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	users := h.Store.ListAllUsers()
	var matched []map[string]interface{}
	for _, u := range users {
		if term == "" || strings.Contains(strings.ToLower(u.Username), strings.ToLower(term)) ||
			strings.Contains(strings.ToLower(u.Name), strings.ToLower(term)) {
			matched = append(matched, map[string]interface{}{
				"username":        u.Username,
				"name":            u.Name,
				"avatar_template": u.AvatarTemplate,
			})
		}
	}
	if matched == nil {
		matched = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users":  matched,
		"groups": []interface{}{},
	})
}

// ---- User Notifications ----

// GET /u/{username}/notifications
func (h *ExtendedUsersHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"notifications": []interface{}{},
	})
}

// GET /u/{username}/messages
func (h *ExtendedUsersHandler) Messages(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.TopicListResponse{
		TopicList: model.TopicList{Topics: []model.Topic{}},
	})
}

// GET /u/{username}/drafts
func (h *ExtendedUsersHandler) Drafts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"drafts": []interface{}{},
	})
}

// ---- Password / Token ----

// POST /u/password-reset/{token}
func (h *ExtendedUsersHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /u/confirm-email/{token}
func (h *ExtendedUsersHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /u/second_factors
func (h *ExtendedUsersHandler) SecondFactors(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"totps":        []interface{}{},
		"security_keys": []interface{}{},
	})
}

// PUT /u/second_factor
func (h *ExtendedUsersHandler) UpdateSecondFactor(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /u/create_second_factor_totp
func (h *ExtendedUsersHandler) CreateTOTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"key": "JBSWY3DPEHPK3PXP",
		"qr":  "data:image/png;base64,placeholder",
	})
}

// POST /u/enable_second_factor_totp
func (h *ExtendedUsersHandler) EnableTOTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// PUT /u/security_key
func (h *ExtendedUsersHandler) SecurityKey(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /u/create_second_factor_security_key
func (h *ExtendedUsersHandler) CreateSecurityKey(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"challenge": "placeholder_challenge",
		"rp_id":     "dtu-discourse",
		"rp_name":   "DTU Discourse",
		"supported_algorithms": []int{-7, -257},
	})
}

// POST /u/register_second_factor_security_key
func (h *ExtendedUsersHandler) RegisterSecurityKey(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// POST /u/second_factors_backup
func (h *ExtendedUsersHandler) BackupCodes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"backup_codes": []string{"abc123", "def456", "ghi789", "jkl012", "mno345"},
	})
}

// PUT /u/{username}/preferences/second-factor
func (h *ExtendedUsersHandler) SecondFactorPref(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// DELETE /u/{username}/preferences/user_image
func (h *ExtendedUsersHandler) DeleteUserImage(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- User Stat / Profile Fields ----

// GET /user-stat/{username}
func (h *ExtendedUsersHandler) UserStat(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"days_visited":     10,
		"time_read":        3600,
		"recent_time_read": 600,
		"topics_entered":   5,
		"posts_read_count": 20,
		"likes_given":      3,
		"likes_received":   7,
		"new_since":        time.Now().Add(-30 * 24 * time.Hour).UTC(),
		"read_faq":         false,
		"first_post_created_at": nil,
		"post_count":       0,
		"topic_count":      0,
	})
}
