package handler

import (
	"net/http"
	"time"

	"github.com/lightcap/dtu-discourse/internal/model"
	"github.com/lightcap/dtu-discourse/internal/store"
)

// ExtendedAdminHandler covers undocumented admin endpoints: webhooks, themes,
// color schemes, watched words, site texts, permalinks, screened items,
// staff action logs, embeddable hosts, custom user fields, review queue,
// flags, impersonation, silence/unsilence, reports, version check, etc.
type ExtendedAdminHandler struct {
	Store *store.Store
}

// ---- Webhooks ----

func (h *ExtendedAdminHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"web_hooks":           []interface{}{},
		"extras":              map[string]interface{}{"default_event_types": []interface{}{}},
		"total_rows_web_hooks": 0,
	})
}

func (h *ExtendedAdminHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"web_hook": map[string]interface{}{
			"id": 1, "payload_url": "", "content_type": 1,
			"secret": "", "wildcard_web_hook": false, "active": true,
			"verify_certificate": true, "created_at": time.Now().UTC(),
		},
	})
}

func (h *ExtendedAdminHandler) ShowWebhook(w http.ResponseWriter, r *http.Request) {
	h.CreateWebhook(w, r)
}

func (h *ExtendedAdminHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	h.CreateWebhook(w, r)
}

func (h *ExtendedAdminHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) WebhookEvents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"web_hook_events": []interface{}{}})
}

func (h *ExtendedAdminHandler) PingWebhook(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Themes ----

func (h *ExtendedAdminHandler) ListThemes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"themes": []interface{}{}})
}

func (h *ExtendedAdminHandler) CreateTheme(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"theme": map[string]interface{}{
			"id": 1, "name": "Default", "created_at": time.Now().UTC(),
			"theme_fields": []interface{}{},
		},
	})
}

func (h *ExtendedAdminHandler) ShowTheme(w http.ResponseWriter, r *http.Request) {
	h.CreateTheme(w, r)
}

func (h *ExtendedAdminHandler) UpdateTheme(w http.ResponseWriter, r *http.Request) {
	h.CreateTheme(w, r)
}

func (h *ExtendedAdminHandler) DeleteTheme(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) ImportTheme(w http.ResponseWriter, r *http.Request) {
	h.CreateTheme(w, r)
}

func (h *ExtendedAdminHandler) ExportTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/gzip")
	w.WriteHeader(http.StatusOK)
}

// ---- Color Schemes ----

func (h *ExtendedAdminHandler) ListColorSchemes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

func (h *ExtendedAdminHandler) CreateColorScheme(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id": 1, "name": "Default", "is_base": true, "colors": []interface{}{},
	})
}

func (h *ExtendedAdminHandler) UpdateColorScheme(w http.ResponseWriter, r *http.Request) {
	h.CreateColorScheme(w, r)
}

func (h *ExtendedAdminHandler) DeleteColorScheme(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Watched Words ----

func (h *ExtendedAdminHandler) ListWatchedWords(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"words":   []interface{}{},
		"actions": []string{"block", "censor", "require_approval", "flag", "replace", "tag", "silence", "link"},
	})
}

func (h *ExtendedAdminHandler) CreateWatchedWord(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id": 1, "word": "", "action": "block",
	})
}

func (h *ExtendedAdminHandler) DeleteWatchedWord(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) UploadWatchedWords(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) ClearWatchedWordsAction(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Site Texts ----

func (h *ExtendedAdminHandler) ListSiteTexts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"site_texts": []interface{}{},
		"extras":     map[string]interface{}{"locale": "en", "has_more": false},
	})
}

func (h *ExtendedAdminHandler) ShowSiteText(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"site_text": map[string]interface{}{"id": "", "value": ""},
	})
}

func (h *ExtendedAdminHandler) UpdateSiteText(w http.ResponseWriter, r *http.Request) {
	h.ShowSiteText(w, r)
}

func (h *ExtendedAdminHandler) RevertSiteText(w http.ResponseWriter, r *http.Request) {
	h.ShowSiteText(w, r)
}

// ---- Permalinks ----

func (h *ExtendedAdminHandler) ListPermalinks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

func (h *ExtendedAdminHandler) CreatePermalink(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"permalink": map[string]interface{}{"id": 1, "url": "", "topic_id": nil, "post_id": nil, "category_id": nil},
	})
}

func (h *ExtendedAdminHandler) UpdatePermalink(w http.ResponseWriter, r *http.Request) {
	h.CreatePermalink(w, r)
}

func (h *ExtendedAdminHandler) DeletePermalink(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Staff Action Logs ----

func (h *ExtendedAdminHandler) ListStaffActionLogs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"staff_action_logs": []interface{}{},
		"extras":            map[string]interface{}{"user_history_actions": []interface{}{}},
	})
}

func (h *ExtendedAdminHandler) StaffActionLogDiff(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"side_by_side": "", "inline": ""})
}

// ---- Screened Items ----

func (h *ExtendedAdminHandler) ListScreenedEmails(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

func (h *ExtendedAdminHandler) DeleteScreenedEmail(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) ListScreenedIPs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

func (h *ExtendedAdminHandler) CreateScreenedIP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"screened_ip_address": map[string]interface{}{"id": 1, "ip_address": "", "action_type": 1},
	})
}

func (h *ExtendedAdminHandler) UpdateScreenedIP(w http.ResponseWriter, r *http.Request) {
	h.CreateScreenedIP(w, r)
}

func (h *ExtendedAdminHandler) DeleteScreenedIP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) ListScreenedURLs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// ---- Search Logs ----

func (h *ExtendedAdminHandler) ListSearchLogs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

func (h *ExtendedAdminHandler) SearchLogTerms(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// ---- Embeddable Hosts ----

func (h *ExtendedAdminHandler) ShowEmbedding(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"embeddable_hosts": []interface{}{},
	})
}

func (h *ExtendedAdminHandler) UpdateEmbedding(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) CreateEmbeddableHost(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"embeddable_host": map[string]interface{}{"id": 1, "host": ""},
	})
}

func (h *ExtendedAdminHandler) UpdateEmbeddableHost(w http.ResponseWriter, r *http.Request) {
	h.CreateEmbeddableHost(w, r)
}

func (h *ExtendedAdminHandler) DeleteEmbeddableHost(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Custom User Fields ----

func (h *ExtendedAdminHandler) ListCustomUserFields(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"user_fields": []interface{}{}})
}

func (h *ExtendedAdminHandler) CreateCustomUserField(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_field": map[string]interface{}{
			"id": 1, "name": "", "description": "", "field_type": "text",
			"editable": true, "required": false, "show_on_profile": false,
			"show_on_user_card": false, "position": 0,
		},
	})
}

func (h *ExtendedAdminHandler) UpdateCustomUserField(w http.ResponseWriter, r *http.Request) {
	h.CreateCustomUserField(w, r)
}

func (h *ExtendedAdminHandler) DeleteCustomUserField(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Review Queue / Reviewables ----

func (h *ExtendedAdminHandler) ListReviewables(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"reviewables":       []interface{}{},
		"meta":              map[string]interface{}{"total_rows_reviewables": 0},
		"__rest_serializer": "1",
	})
}

func (h *ExtendedAdminHandler) ShowReviewable(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"reviewable": map[string]interface{}{}})
}

func (h *ExtendedAdminHandler) ReviewableCount(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"reviewable_count": 0})
}

func (h *ExtendedAdminHandler) ReviewableTopics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"reviewable_topics": []interface{}{}})
}

func (h *ExtendedAdminHandler) ReviewSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"reviewable_score_types": []interface{}{},
		"reviewable_priorities":  map[string]interface{}{},
	})
}

func (h *ExtendedAdminHandler) UpdateReviewSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) PerformReviewAction(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"reviewable_perform_result": map[string]interface{}{"success": "OK"},
	})
}

func (h *ExtendedAdminHandler) UpdateReviewable(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) DeleteReviewable(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Admin Flags ----

func (h *ExtendedAdminHandler) ListFlags(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{
		map[string]interface{}{
			"id": 1, "name": "off_topic", "name_key": "off_topic",
			"description": "Off-Topic", "enabled": true, "position": 0, "is_flag": true,
		},
		map[string]interface{}{
			"id": 2, "name": "inappropriate", "name_key": "inappropriate",
			"description": "Inappropriate", "enabled": true, "position": 1, "is_flag": true,
		},
		map[string]interface{}{
			"id": 3, "name": "spam", "name_key": "spam",
			"description": "Spam", "enabled": true, "position": 2, "is_flag": true,
		},
	})
}

func (h *ExtendedAdminHandler) CreateFlag(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"flag": map[string]interface{}{"id": 4, "name": "", "enabled": true},
	})
}

func (h *ExtendedAdminHandler) UpdateFlag(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) DeleteFlag(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) ToggleFlag(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Reports ----

func (h *ExtendedAdminHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"reports": []interface{}{}})
}

func (h *ExtendedAdminHandler) ShowReport(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"report": map[string]interface{}{
			"type": pathParam(r, "type"), "title": "", "data": []interface{}{},
		},
	})
}

func (h *ExtendedAdminHandler) BulkReports(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"reports": []interface{}{}})
}

// ---- Dashboard sub-routes ----

func (h *ExtendedAdminHandler) DashboardGeneral(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"reports": []interface{}{}, "updated_at": nil,
	})
}

func (h *ExtendedAdminHandler) DashboardModeration(w http.ResponseWriter, r *http.Request) {
	h.DashboardGeneral(w, r)
}

func (h *ExtendedAdminHandler) DashboardSecurity(w http.ResponseWriter, r *http.Request) {
	h.DashboardGeneral(w, r)
}

func (h *ExtendedAdminHandler) DashboardProblems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"problems": []interface{}{}})
}

func (h *ExtendedAdminHandler) VersionCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"latest_version":    "3.4.0",
		"critical_updates":  false,
		"installed_version": "3.4.0",
		"installed_sha":     "abc123",
		"git_branch":        "main",
		"updated_at":        time.Now().UTC(),
	})
}

// ---- Silence / Unsilence ----

func (h *ExtendedAdminHandler) SilenceUser(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"silenced": true})
}

func (h *ExtendedAdminHandler) UnsilenceUser(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"silenced": false})
}

// ---- User admin ops ----

func (h *ExtendedAdminHandler) SetPrimaryGroup(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) DeletePostsBatch(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) MergeUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) ResetBounceScore(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) BulkApproveUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) BulkDestroyUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) IPInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"ip": "", "hostname": "", "location": "",
	})
}

func (h *ExtendedAdminHandler) GenerateAPIKeyForUser(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"api_key": map[string]interface{}{"id": 1, "key": "generated_key"},
	})
}

// ---- Impersonation ----

func (h *ExtendedAdminHandler) StartImpersonation(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) StopImpersonation(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Admin Search ----

func (h *ExtendedAdminHandler) AdminSearch(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users":      []interface{}{},
		"categories": []interface{}{},
		"tags":       []interface{}{},
	})
}

// ---- Email Templates ----

func (h *ExtendedAdminHandler) ListEmailTemplates(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

func (h *ExtendedAdminHandler) ShowEmailTemplate(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"email_template": map[string]interface{}{
			"id": "", "title": "", "subject": "", "body": "",
		},
	})
}

func (h *ExtendedAdminHandler) UpdateEmailTemplate(w http.ResponseWriter, r *http.Request) {
	h.ShowEmailTemplate(w, r)
}

func (h *ExtendedAdminHandler) RevertEmailTemplate(w http.ResponseWriter, r *http.Request) {
	h.ShowEmailTemplate(w, r)
}

// ---- Email Style ----

func (h *ExtendedAdminHandler) GetEmailStyle(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"email_style": map[string]interface{}{"compiled_css": "", "html_template": "", "css": ""},
	})
}

func (h *ExtendedAdminHandler) UpdateEmailStyle(w http.ResponseWriter, r *http.Request) {
	h.GetEmailStyle(w, r)
}

// ---- Robots.txt ----

func (h *ExtendedAdminHandler) GetRobots(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"robots_txt": "User-agent: *\nAllow: /"})
}

func (h *ExtendedAdminHandler) UpdateRobots(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) ResetRobots(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Form Templates ----

func (h *ExtendedAdminHandler) ListFormTemplates(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"form_templates": []interface{}{}})
}

func (h *ExtendedAdminHandler) CreateFormTemplate(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"form_template": map[string]interface{}{"id": 1, "name": "", "template": ""},
	})
}

func (h *ExtendedAdminHandler) UpdateFormTemplate(w http.ResponseWriter, r *http.Request) {
	h.CreateFormTemplate(w, r)
}

func (h *ExtendedAdminHandler) DeleteFormTemplate(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Reseed ----

func (h *ExtendedAdminHandler) GetReseed(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"categories": []interface{}{}, "topics": []interface{}{}})
}

func (h *ExtendedAdminHandler) PostReseed(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

// ---- Badge Admin Extended ----

func (h *ExtendedAdminHandler) BadgeTypes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"badge_types": []map[string]interface{}{
			{"id": 1, "name": "Gold"},
			{"id": 2, "name": "Silver"},
			{"id": 3, "name": "Bronze"},
		},
	})
}

func (h *ExtendedAdminHandler) BadgeGroupings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.SuccessResponse{Success: "OK"})
}

func (h *ExtendedAdminHandler) PreviewBadge(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"grant_count": 0, "sample": []interface{}{}})
}
