package main

import (
	"encoding/json"
	"testing"
)

// ============================================================
// Polls
// ============================================================

func TestExtended_Polls_Vote(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/polls/vote", map[string]interface{}{
		"post_id": float64(1), "poll_name": "poll", "options": []interface{}{"opt1"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Polls_ToggleStatus(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/polls/toggle_status", map[string]interface{}{
		"post_id": float64(1), "poll_name": "poll", "status": "closed",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Polls_Voters(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/polls/voters.json?post_id=1&poll_name=poll")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// API Keys
// ============================================================

func TestExtended_APIKeys_List(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/admin/api/keys")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["keys"]; !ok {
		t.Fatal("missing keys")
	}
}

func TestExtended_APIKeys_Create(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiRequest(ts, "POST", "/admin/api/keys", map[string]interface{}{
		"key": map[string]interface{}{"description": "test key"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	if _, ok := data["key"]; !ok {
		t.Fatal("missing key")
	}
}

func TestExtended_APIKeys_Scopes(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/admin/api/keys/scopes")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["scopes"]; !ok {
		t.Fatal("missing scopes")
	}
}

func TestExtended_APIKeys_Revoke(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/admin/api/keys/1/revoke", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_APIKeys_Delete(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/admin/api/keys/1", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Email Admin
// ============================================================

func TestExtended_Email_Settings(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/email.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Email_Test(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/admin/email/test", map[string]interface{}{
		"email_address": "test@example.com",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Email_ServerSettings(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/email/server-settings")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// User Actions
// ============================================================

func TestExtended_UserActions_List(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/user_actions.json?username=admin")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["user_actions"]; !ok {
		t.Fatal("missing user_actions")
	}
}

// ============================================================
// Topic Timings
// ============================================================

func TestExtended_TopicTimings_Record(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/topics/timings", map[string]interface{}{
		"topic_id": float64(1), "topic_time": float64(1000),
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_TopicTimings_SimilarTo(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/topics/similar_to?title=Welcome")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_TopicTimings_Bulk(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/topics/bulk", map[string]interface{}{
		"topic_ids": []interface{}{float64(1)}, "operation": map[string]interface{}{"type": "close"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_TopicTimings_ResetNew(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/topics/reset-new", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended Topics (via /t/ sub-router)
// ============================================================

func TestExtended_Topics_ArchiveMessage(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1/archive-message", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Topics_ClearPin(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1/clear-pin", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Topics_Mute(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1/mute", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Topics_Recover(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1/recover", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Topics_UpdateTags(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1/tags", map[string]interface{}{
		"tags": []interface{}{"welcome", "announcement"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Topics_MovePosts(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/t/1/move-posts", map[string]interface{}{
		"title": "Moved Posts", "post_ids": []interface{}{float64(1)},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Topics_PostIDs(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/t/1/post_ids")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["post_ids"]; !ok {
		t.Fatal("missing post_ids")
	}
}

func TestExtended_Topics_IDForSlug(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/t/id_for/welcome-to-discourse")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_id"]; !ok {
		t.Fatal("missing topic_id")
	}
}

func TestExtended_Topics_FeatureStats(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/topics/feature_stats")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended Posts (via /posts/ sub-router)
// ============================================================

func TestExtended_Posts_Raw(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/posts/1/raw")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_Cooked(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/posts/1/cooked")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_ReplyHistory(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/posts/1/reply-history")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_ReplyIDs(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/posts/1/reply-ids")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_Recover(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/posts/1/recover", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_Rebake(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/posts/1/rebake", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_Locked(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/posts/1/locked", map[string]interface{}{
		"locked": true,
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_RevisionsLatest(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/posts/1/revisions/latest")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_ByNumber(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/posts/by_number/1/1")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Posts_RawByNumber(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/raw/1/1")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Admin Extended
// ============================================================

func TestExtended_Admin_ListWebhooks(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/api/web_hooks.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_CreateWebhook(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/admin/api/web_hooks.json", map[string]interface{}{
		"web_hook": map[string]interface{}{"payload_url": "https://example.com/hook"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListThemes(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/admin/themes.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var themes []interface{}
	if err := json.Unmarshal(body, &themes); err != nil {
		// might be wrapped in an object
		data := parseJSON(t, body)
		if _, ok := data["themes"]; !ok {
			t.Fatal("unexpected response structure")
		}
	}
}

func TestExtended_Admin_CreateTheme(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/admin/themes.json", map[string]interface{}{
		"theme": map[string]interface{}{"name": "Test Theme"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListColorSchemes(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/color_schemes.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListWatchedWords(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/customize/watched_words.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListSiteTexts(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/customize/site_texts.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListPermalinks(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/permalinks.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListStaffActionLogs(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/logs/staff_action_logs.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListScreenedEmails(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/logs/screened_emails.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListScreenedIPs(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/logs/screened_ip_addresses.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListCustomUserFields(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/customize/user_fields.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListReviewables(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/review.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ReviewableCount(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/review/count.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListFlags(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/config/flags.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListReports(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/reports.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_VersionCheck(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/version_check.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_Impersonate(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/admin/impersonate/1", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_Search(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/search/all?term=admin")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListEmailTemplates(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/customize/email_templates.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_GetRobots(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/customize/robots.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_ListFormTemplates(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/customize/form-templates.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_BadgeTypes(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/badge_types.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Admin_BackupsStatus(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/backups/status")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Misc
// ============================================================

func TestExtended_Misc_HotTopics(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/hot.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_Filter(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/filter")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_DirectoryItems(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/directory_items?period=weekly")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_About(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/about")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["about"]; !ok {
		t.Fatal("missing about")
	}
}

func TestExtended_Misc_SiteBasicInfo(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/site/basic-info.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_ListDrafts(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/drafts")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["drafts"]; !ok {
		t.Fatal("missing drafts")
	}
}

func TestExtended_Misc_CreateDraft(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/drafts", map[string]interface{}{
		"draft_key": "new_topic", "data": "draft content",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_CreateBookmark(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/bookmarks", map[string]interface{}{
		"bookmarkable_id": float64(1), "bookmarkable_type": "Post",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_ListSidebarSections(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/sidebar_sections")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_TrackClick(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/clicks/track", map[string]interface{}{
		"url": "https://example.com", "post_id": float64(1),
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_Onebox(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/onebox?url=https://example.com")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_GenerateSlug(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/slugs", map[string]interface{}{
		"name": "Test Slug Name",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_EmbedTopics(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/embed/topics")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_PresenceUpdate(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/presence/update", map[string]interface{}{
		"present": true,
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_GetUserStatus(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/user-status")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_EnableDND(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/do-not-disturb", map[string]interface{}{
		"duration": float64(60),
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_ListEmojis(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/emojis")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_LookupHashtags(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/hashtags")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_ComposerMentions(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/composer/mentions?names=admin")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_ServerStatus(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/srv/status")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Misc_PermalinkCheck(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/permalink-check?path=/test")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Session / Auth
// ============================================================

func TestExtended_Session_Login(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/session", map[string]interface{}{
		"login": "admin", "password": "password",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Session_Current(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/session/current.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["current_user"]; !ok {
		t.Fatal("missing current_user")
	}
}

func TestExtended_Session_ForgotPassword(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/session/forgot_password", map[string]interface{}{
		"login": "admin",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Session_Honeypot(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/session/hp")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Session_TwoFactorAuth(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/session/2fa", map[string]interface{}{
		"second_factor_token": "123456",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Session_TwoFactorStatus(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/session/2fa.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended Users
// ============================================================

func TestExtended_Users_Summary(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/u/admin/summary")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["user_summary"]; !ok {
		t.Fatal("missing user_summary")
	}
}

func TestExtended_Users_Activity(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/u/admin/activity")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Users_Bookmarks(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/u/admin/bookmarks")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Users_Emails(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/u/admin/emails")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Users_SearchUsers(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/u/search/users?term=admin")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Users_PickAvatar(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/u/admin/preferences/avatar/pick", map[string]interface{}{
		"type": "system",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended PMs
// ============================================================

func TestExtended_PM_Unread(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/topics/private-messages-unread/admin")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_list"]; !ok {
		t.Fatal("missing topic_list")
	}
}

func TestExtended_PM_Archive(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/topics/private-messages-archive/admin")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_list"]; !ok {
		t.Fatal("missing topic_list")
	}
}

func TestExtended_PM_New(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/topics/private-messages-new/admin")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_list"]; !ok {
		t.Fatal("missing topic_list")
	}
}

// ============================================================
// Tag Groups
// ============================================================

func TestExtended_TagGroups_List(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/tag_groups")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["tag_groups"]; !ok {
		t.Fatal("missing tag_groups")
	}
}

func TestExtended_TagGroups_Create(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/tag_groups", map[string]interface{}{
		"name": "Test Tag Group", "tag_names": []interface{}{"welcome"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_TagGroups_Show(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/tag_groups/1")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_TagGroups_Update(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/tag_groups/1", map[string]interface{}{
		"name": "Updated Tag Group",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_TagGroups_Delete(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/tag_groups/1", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended Tags
// ============================================================

func TestExtended_Tags_Search(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/tags/filter/search?q=welcome")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Tags_Create(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/tags", map[string]interface{}{
		"tag": map[string]interface{}{"name": "newtag"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Tags_Info(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/tags/welcome/info")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Tags_Unused(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/tags/unused")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended Notifications
// ============================================================

func TestExtended_Notifications_MarkRead(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/notifications/mark-read", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Notifications_Totals(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/notifications/totals")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended Groups
// ============================================================

func TestExtended_Groups_Join(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/groups/staff/join", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Groups_Logs(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/groups/staff/logs")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Groups_Permissions(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/groups/staff/permissions")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended Categories
// ============================================================

func TestExtended_Categories_Search(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/categories/search?term=general")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Categories_CategoriesAndLatest(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/categories_and_latest")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["category_list"]; !ok {
		t.Fatal("missing category_list")
	}
}

// ============================================================
// Extended Uploads
// ============================================================

func TestExtended_Uploads_LookupMetadata(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/uploads/lookup-metadata?url=https://example.com/upload.png")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Uploads_GeneratePresignedPut(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/uploads/generate-presigned-put", map[string]interface{}{
		"file_name": "test.png", "type": "composer",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// Extended Backups
// ============================================================

func TestExtended_Backups_Status(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/backups/status")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExtended_Backups_Logs(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/backups/logs")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
