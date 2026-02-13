package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/lightcap/dtu-discourse/internal/middleware"
	"github.com/lightcap/dtu-discourse/internal/store"
)

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := store.New()
	mux := BuildRouter(s, nil)
	wrapped := middleware.Auth(s)(mux)
	return httptest.NewServer(wrapped)
}

func apiGet(ts *httptest.Server, path string) (*http.Response, []byte) {
	req, _ := http.NewRequest("GET", ts.URL+path, nil)
	req.Header.Set("Api-Key", "test_api_key")
	req.Header.Set("Api-Username", "system")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp, body
}

func apiRequest(ts *httptest.Server, method, path string, body interface{}) (*http.Response, []byte) {
	var bodyReader io.Reader
	ct := "application/json"
	switch v := body.(type) {
	case map[string]interface{}:
		b, _ := json.Marshal(v)
		bodyReader = bytes.NewReader(b)
	case url.Values:
		bodyReader = strings.NewReader(v.Encode())
		ct = "application/x-www-form-urlencoded"
	case nil:
		bodyReader = nil
	}

	req, _ := http.NewRequest(method, ts.URL+path, bodyReader)
	req.Header.Set("Api-Key", "admin_api_key")
	req.Header.Set("Api-Username", "admin")
	req.Header.Set("Content-Type", ct)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp, b
}

func parseJSON(t *testing.T, data []byte) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to parse JSON: %v\nbody: %s", err, string(data))
	}
	return result
}

// ============================================================
// Authentication
// ============================================================

func TestAuth_RequiresAPIKey(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := http.Get(ts.URL + "/latest.json")
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", resp.StatusCode)
	}
}

func TestAuth_InvalidKey(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL+"/latest.json", nil)
	req.Header.Set("Api-Key", "invalid_key")
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", resp.StatusCode)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Users
// ============================================================

func TestRuby_GetUser(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/users/admin.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	user := data["user"].(map[string]interface{})
	if user["username"] != "admin" {
		t.Errorf("expected 'admin', got %v", user["username"])
	}
}

func TestRuby_CreateUser(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiRequest(ts, "POST", "/users", map[string]interface{}{
		"name": "Test", "username": "testuser", "email": "test@example.com", "password": "pw123",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	if data["success"] != true {
		t.Errorf("expected success=true, got %v", data["success"])
	}
}

func TestRuby_ListUsers(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/admin/users/list/active.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var users []interface{}
	json.Unmarshal(body, &users)
	if len(users) == 0 {
		t.Error("expected users")
	}
}

func TestRuby_Activate(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/users/2/activate", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_Suspend(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/users/2/suspend", map[string]interface{}{"reason": "test"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_Unsuspend(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/users/2/unsuspend", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_GrantAdmin(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/users/2/grant_admin", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_RevokeAdmin(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/users/2/revoke_admin", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_GrantModeration(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/users/2/grant_moderation", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_RevokeModeration(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/users/2/revoke_moderation", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_DeleteUser(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/admin/users/3.json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_CheckUsername(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/users/check_username.json?username=admin")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if data["available"] != false {
		t.Errorf("expected taken")
	}
}

func TestRuby_LogOut(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/admin/users/1/log_out", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_GetUserByExternalID(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	// First sync an SSO user
	apiRequest(ts, "POST", "/admin/users/sync_sso", map[string]interface{}{
		"external_id": "ext-999", "email": "ext@example.com", "username": "extuser", "name": "Ext",
	})
	resp, body := apiGet(ts, "/users/by-external/ext-999")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Categories
// ============================================================

func TestRuby_ListCategories(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/categories.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	catList := data["category_list"].(map[string]interface{})
	cats := catList["categories"].([]interface{})
	if len(cats) == 0 {
		t.Fatal("expected categories")
	}
}

func TestRuby_CreateCategory(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiRequest(ts, "POST", "/categories", map[string]interface{}{
		"name": "Test Cat", "color": "FF0000",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	if _, ok := data["category"]; !ok {
		t.Fatal("missing category")
	}
}

func TestRuby_UpdateCategory(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/categories/1", map[string]interface{}{"name": "Updated"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_ShowCategory(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/c/1/show")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["category"]; !ok {
		t.Fatal("missing category")
	}
}

func TestRuby_CategoryLatestTopics(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/c/general/l/latest.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_list"]; !ok {
		t.Fatal("missing topic_list")
	}
}

func TestRuby_CategoryTopTopics(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/c/general/l/top.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_CategoryNewTopics(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/c/general/l/new.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_DeleteCategory(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/categories/3", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_ReorderCategories(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/categories/reorder", map[string]interface{}{})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_CategoryNotification(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/category/1/notifications", map[string]interface{}{"notification_level": float64(3)})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Topics
// ============================================================

func TestRuby_LatestTopics(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/latest.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	tl := data["topic_list"].(map[string]interface{})
	topics := tl["topics"].([]interface{})
	if len(topics) == 0 {
		t.Fatal("expected topics")
	}
}

func TestRuby_TopTopics(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/top.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_NewTopics(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/new.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_GetTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/t/1.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if data["title"] != "Welcome to Discourse" {
		t.Errorf("unexpected title: %v", data["title"])
	}
	if _, ok := data["post_stream"]; !ok {
		t.Error("expected post_stream")
	}
	if _, ok := data["details"]; !ok {
		t.Error("expected details")
	}
}

func TestRuby_CreateTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiRequest(ts, "POST", "/posts", map[string]interface{}{
		"title": "New Test Topic", "raw": "Body of the test topic.", "category": float64(1),
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	if data["topic_id"] == nil {
		t.Error("expected topic_id")
	}
}

func TestRuby_RenameTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1.json", map[string]interface{}{
		"topic": map[string]interface{}{"title": "Renamed"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_RecategorizeTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1.json", map[string]interface{}{
		"topic": map[string]interface{}{"category_id": float64(2)},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_UpdateTopicStatus(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	for _, status := range []string{"closed", "archived", "pinned", "visible"} {
		resp, _ := apiRequest(ts, "PUT", "/t/1/status", map[string]interface{}{
			"status": status, "enabled": true,
		})
		if resp.StatusCode != 200 {
			t.Errorf("status %s: expected 200, got %d", status, resp.StatusCode)
		}
	}
}

func TestRuby_TopicsByUser(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/topics/created-by/admin.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_list"]; !ok {
		t.Fatal("missing topic_list")
	}
}

func TestRuby_DeleteTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/t/3.json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_TopicPosts(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/t/1/posts.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["post_stream"]; !ok {
		t.Fatal("missing post_stream")
	}
}

func TestRuby_ChangeTimestamp(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1/change-timestamp", map[string]interface{}{"timestamp": "2024-01-01T00:00:00Z"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_ChangeOwner(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/t/1/change-owner.json", map[string]interface{}{"username": "alice", "post_ids": []interface{}{float64(1)}})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_TopicNotificationLevel(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/t/1/notifications", map[string]interface{}{"notification_level": float64(3)})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_BookmarkTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1/bookmark.json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_RemoveBookmark(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/t/1/remove_bookmarks.json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_InviteToTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/t/1/invite", map[string]interface{}{"user": "alice"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Posts
// ============================================================

func TestRuby_CreatePost(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiRequest(ts, "POST", "/posts", map[string]interface{}{
		"topic_id": float64(1), "raw": "A reply.",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
}

func TestRuby_GetPost(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/posts/1.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if data["username"] != "admin" {
		t.Errorf("expected 'admin', got %v", data["username"])
	}
}

func TestRuby_ListPosts(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/posts.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["latest_posts"]; !ok {
		t.Fatal("missing latest_posts")
	}
}

func TestRuby_UpdatePost(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/posts/1", map[string]interface{}{
		"post": map[string]interface{}{"raw": "Updated."},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_DeletePost(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/posts/4.json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_WikifyPost(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/posts/1/wiki", map[string]interface{}{"wiki": true})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_CreatePostAction(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/post_actions", map[string]interface{}{
		"id": float64(1), "post_action_type_id": float64(2),
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_PostActionUsers(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/post_action_users.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Groups
// ============================================================

func TestRuby_ListGroups(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/groups.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["groups"]; !ok {
		t.Fatal("missing groups")
	}
}

func TestRuby_GetGroup(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/groups/staff.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["group"]; !ok {
		t.Fatal("missing group")
	}
}

func TestRuby_CreateGroup(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiRequest(ts, "POST", "/admin/groups", map[string]interface{}{
		"group": map[string]interface{}{"name": "test_grp"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	if _, ok := data["basic_group"]; !ok {
		t.Fatal("missing basic_group")
	}
}

func TestRuby_UpdateGroup(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/groups/1", map[string]interface{}{
		"group": map[string]interface{}{"full_name": "Staff Group"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_GroupMembers(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/groups/staff/members.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["members"]; !ok {
		t.Fatal("missing members")
	}
}

func TestRuby_AddGroupMembers(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/groups/1/members.json", map[string]interface{}{
		"usernames": "alice",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_RemoveGroupMembers(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/admin/groups/1/members.json", map[string]interface{}{
		"usernames": "alice",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_AddGroupOwners(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/groups/1/owners.json", map[string]interface{}{
		"usernames": "alice",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_RemoveGroupOwners(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/admin/groups/1/owners.json", map[string]interface{}{
		"usernames": "alice",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_DeleteGroup(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/admin/groups/10.json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_GroupNotification(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/groups/staff/notifications", map[string]interface{}{"notification_level": float64(3)})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Search
// ============================================================

func TestRuby_Search(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/search?q=welcome")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topics"]; !ok {
		t.Fatal("missing topics")
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Tags
// ============================================================

func TestRuby_ListTags(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/tags.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["tags"]; !ok {
		t.Fatal("missing tags")
	}
}

func TestRuby_ShowTag(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/tag/welcome")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_list"]; !ok {
		t.Fatal("missing topic_list")
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Badges
// ============================================================

func TestRuby_ListBadges(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/admin/badges.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	badges := data["badges"].([]interface{})
	if len(badges) == 0 {
		t.Fatal("expected badges")
	}
}

func TestRuby_CreateBadge(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/admin/badges.json", map[string]interface{}{
		"name": "TestBadge", "description": "desc", "badge_type_id": float64(3),
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_UpdateBadge(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/badges/1.json", map[string]interface{}{"name": "Updated Basic"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_DeleteBadge(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "DELETE", "/admin/badges/2.json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_GrantUserBadge(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/user_badges", map[string]interface{}{
		"username": "alice", "badge_id": float64(1),
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_UserBadges(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	apiRequest(ts, "POST", "/user_badges", map[string]interface{}{
		"username": "alice", "badge_id": float64(1),
	})
	resp, body := apiGet(ts, "/user-badges/alice.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["user_badges"]; !ok {
		t.Fatal("missing user_badges")
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Notifications
// ============================================================

func TestRuby_Notifications(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/notifications.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["notifications"]; !ok {
		t.Fatal("missing notifications")
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Private Messages
// ============================================================

func TestRuby_PrivateMessages(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/topics/private-messages/admin.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_list"]; !ok {
		t.Fatal("missing topic_list")
	}
}

func TestRuby_SentPrivateMessages(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/topics/private-messages-sent/admin.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topic_list"]; !ok {
		t.Fatal("missing topic_list")
	}
}

func TestRuby_CreatePM(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiRequest(ts, "POST", "/posts", map[string]interface{}{
		"title": "DM Test", "raw": "pm body",
		"target_usernames": "alice", "archetype": "private_message",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Invites
// ============================================================

func TestRuby_CreateInvite(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/invites", map[string]interface{}{"email": "inv@example.com"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_RetrieveInvite(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/invites/retrieve.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_DestroyAllExpiredInvites(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/invites/destroy-all-expired", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_ResendAllInvites(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/invites/reinvite-all", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_ResendInvite(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/invites/reinvite", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_GenerateInviteToken(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/invite-token/generate", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Uploads
// ============================================================

func TestRuby_Upload(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "test.png")
	part.Write([]byte("fake"))
	writer.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/uploads", &buf)
	req.Header.Set("Api-Key", "admin_api_key")
	req.Header.Set("Api-Username", "admin")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — SSO
// ============================================================

func TestRuby_SyncSSO(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiRequest(ts, "POST", "/admin/users/sync_sso", map[string]interface{}{
		"external_id": "ext-1", "email": "sso@example.com", "username": "ssouser", "name": "SSO",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	if _, ok := data["user"]; !ok {
		t.Fatal("missing user")
	}
}

// ============================================================
// SDK: discourse_api (Ruby) — Admin
// ============================================================

func TestRuby_SiteSettings(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/admin/site_settings.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["site_settings"]; !ok {
		t.Fatal("missing site_settings")
	}
}

func TestRuby_UpdateSiteSetting(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "PUT", "/admin/site_settings/title.json", map[string]interface{}{"title": "New Title"})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_Backups(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/backups.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_CreateBackup(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiRequest(ts, "POST", "/admin/backups.json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_Dashboard(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, _ := apiGet(ts, "/admin/dashboard.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRuby_SiteInfo(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/site.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	for _, key := range []string{"notification_types", "post_types", "trust_levels", "groups", "categories"} {
		if _, ok := data[key]; !ok {
			t.Errorf("missing %s", key)
		}
	}
}

func TestRuby_CSRFToken(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/session/csrf.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["csrf"]; !ok {
		t.Fatal("missing csrf")
	}
}

// ============================================================
// SDK: pydiscourse (Python) — Form-encoded payloads
// ============================================================

func TestPython_FormCreateUser(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	form := url.Values{"name": {"Py"}, "username": {"pyuser"}, "email": {"py@example.com"}, "password": {"pw"}}
	resp, body := apiRequest(ts, "POST", "/users", form)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
}

func TestPython_FormCreateTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	form := url.Values{"title": {"Py Topic"}, "raw": {"Body from pydiscourse"}, "category": {"1"}}
	resp, body := apiRequest(ts, "POST", "/posts", form)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
}

func TestPython_FormCreateReply(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	form := url.Values{"topic_id": {"1"}, "raw": {"Reply from pydiscourse"}}
	resp, body := apiRequest(ts, "POST", "/posts", form)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}
}

// ============================================================
// SDK: discourse-api (JS)
// ============================================================

func TestJS_GetUser(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/users/alice")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	user := data["user"].(map[string]interface{})
	if user["username"] != "alice" {
		t.Errorf("expected alice, got %v", user["username"])
	}
}

func TestJS_Categories(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/categories.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	catList := data["category_list"].(map[string]interface{})
	if len(catList["categories"].([]interface{})) < 2 {
		t.Error("expected >= 2 categories")
	}
}

func TestJS_GetTopic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/t/2.json")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if data["title"] != "How to use the API" {
		t.Errorf("unexpected title: %v", data["title"])
	}
}

func TestJS_Search(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/search?q=API")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["topics"]; !ok {
		t.Fatal("missing topics")
	}
}

// ============================================================
// Integration: Full lifecycle
// ============================================================

func TestLifecycle_TopicCreateReplyEditDelete(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()

	// 1. Create topic
	resp, body := apiRequest(ts, "POST", "/posts", map[string]interface{}{
		"title": "Lifecycle", "raw": "Original.", "category": float64(1),
	})
	if resp.StatusCode != 200 {
		t.Fatalf("create topic: %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	topicID := data["topic_id"].(float64)
	postID := data["id"].(float64)

	// 2. Reply
	resp, body = apiRequest(ts, "POST", "/posts", map[string]interface{}{
		"topic_id": topicID, "raw": "Reply.",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("reply: %d: %s", resp.StatusCode, body)
	}
	reply := parseJSON(t, body)
	if reply["post_number"].(float64) != 2 {
		t.Errorf("expected post_number=2, got %v", reply["post_number"])
	}

	// 3. Edit post
	resp, _ = apiRequest(ts, "PUT", "/posts/"+strconv.Itoa(int(postID)), map[string]interface{}{
		"post": map[string]interface{}{"raw": "Edited."},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("edit: expected 200, got %d", resp.StatusCode)
	}

	// 4. Delete topic
	resp, _ = apiRequest(ts, "DELETE", "/t/"+strconv.Itoa(int(topicID)), nil)
	if resp.StatusCode != 200 {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}

	// 5. Verify gone
	resp, _ = apiGet(ts, "/t/"+strconv.Itoa(int(topicID))+".json")
	if resp.StatusCode != 404 {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestLifecycle_GroupCreateAddRemoveDelete(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()

	// Create
	resp, body := apiRequest(ts, "POST", "/admin/groups", map[string]interface{}{
		"group": map[string]interface{}{"name": "lc_group"},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("create: %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	group := data["basic_group"].(map[string]interface{})
	gid := strconv.Itoa(int(group["id"].(float64)))

	// Add member
	resp, _ = apiRequest(ts, "PUT", "/admin/groups/"+gid+"/members.json", map[string]interface{}{"usernames": "alice"})
	if resp.StatusCode != 200 {
		t.Fatalf("add member: expected 200, got %d", resp.StatusCode)
	}

	// Remove member
	resp, _ = apiRequest(ts, "DELETE", "/admin/groups/"+gid+"/members.json", map[string]interface{}{"usernames": "alice"})
	if resp.StatusCode != 200 {
		t.Fatalf("remove member: expected 200, got %d", resp.StatusCode)
	}

	// Delete group
	resp, _ = apiRequest(ts, "DELETE", "/admin/groups/"+gid+".json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}
}

func TestLifecycle_UserCreateSuspendDelete(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()

	// Create
	resp, body := apiRequest(ts, "POST", "/users", map[string]interface{}{
		"name": "LC", "username": "lcuser", "email": "lc@example.com", "password": "pw",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("create: %d: %s", resp.StatusCode, body)
	}
	data := parseJSON(t, body)
	uid := strconv.Itoa(int(data["user_id"].(float64)))

	// Suspend
	resp, _ = apiRequest(ts, "PUT", "/admin/users/"+uid+"/suspend", map[string]interface{}{"reason": "test"})
	if resp.StatusCode != 200 {
		t.Fatalf("suspend: expected 200, got %d", resp.StatusCode)
	}

	// Delete
	resp, _ = apiRequest(ts, "DELETE", "/admin/users/"+uid+".json", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}
}

// ============================================================
// SSO Round-Trip
// ============================================================

func TestSSO_RoundTrip(t *testing.T) {
	secret := "test_sso_secret_key"

	// Set up store with SSO config
	s := store.New()
	s.SSOSecret = secret
	// SSOCallbackURL will be set to the test server's own URL below (we use a placeholder for redirect)
	s.SSOCallbackURL = "http://eve.local/api/discourse/sso"

	mux := BuildRouter(s, nil)
	wrapped := middleware.Auth(s)(mux)
	ts := httptest.NewServer(wrapped)
	defer ts.Close()

	// Update callback URL to point to our test server's sso_login for a full round trip
	s.SSOCallbackURL = ts.URL + "/session/sso_login"

	// 1. GET /session/sso — should 302 redirect with sso+sig params
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // don't follow redirects
	}}
	resp, err := client.Get(ts.URL + "/session/sso")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 302 {
		t.Fatalf("expected 302, got %d", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if loc == "" {
		t.Fatal("missing Location header")
	}

	// Parse the redirect URL to extract sso payload and sig
	redirectURL, err := url.Parse(loc)
	if err != nil {
		t.Fatalf("parse redirect URL: %v", err)
	}
	ssoPayload := redirectURL.Query().Get("sso")
	sig := redirectURL.Query().Get("sig")
	if ssoPayload == "" || sig == "" {
		t.Fatal("missing sso or sig in redirect URL")
	}

	// Verify the HMAC signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ssoPayload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	if sig != expectedSig {
		t.Fatal("signature mismatch on redirect")
	}

	// Decode the payload and extract the nonce
	decoded, err := base64.StdEncoding.DecodeString(ssoPayload)
	if err != nil {
		t.Fatalf("decode sso payload: %v", err)
	}
	params, err := url.ParseQuery(string(decoded))
	if err != nil {
		t.Fatalf("parse sso params: %v", err)
	}
	nonce := params.Get("nonce")
	if nonce == "" {
		t.Fatal("missing nonce in payload")
	}

	// 2. Build the return payload (simulating what Eve would send back)
	returnPayload := url.Values{
		"nonce":       {nonce},
		"email":       {"ssotest@example.com"},
		"external_id": {"ext-sso-test"},
		"username":    {"ssotestuser"},
		"name":        {"SSO Test User"},
	}
	returnB64 := base64.StdEncoding.EncodeToString([]byte(returnPayload.Encode()))
	returnMAC := hmac.New(sha256.New, []byte(secret))
	returnMAC.Write([]byte(returnB64))
	returnSig := hex.EncodeToString(returnMAC.Sum(nil))

	// 3. GET /session/sso_login with the return payload
	loginURL := ts.URL + "/session/sso_login?sso=" + url.QueryEscape(returnB64) + "&sig=" + url.QueryEscape(returnSig)
	resp, err = client.Get(loginURL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 302 {
		t.Fatalf("sso_login: expected 302, got %d", resp.StatusCode)
	}

	// 4. Verify user was created in store
	u := s.GetUserByExternalID("ext-sso-test")
	if u == nil {
		t.Fatal("SSO user was not created in store")
	}
	if u.Username != "ssotestuser" {
		t.Errorf("expected username 'ssotestuser', got %q", u.Username)
	}
	if u.Email != "ssotest@example.com" {
		t.Errorf("expected email 'ssotest@example.com', got %q", u.Email)
	}
}

func TestSSO_InvalidSignatureRejected(t *testing.T) {
	s := store.New()
	s.SSOSecret = "real_secret"
	s.SSOCallbackURL = "http://localhost/callback"

	mux := BuildRouter(s, nil)
	wrapped := middleware.Auth(s)(mux)
	ts := httptest.NewServer(wrapped)
	defer ts.Close()

	// Build a payload with wrong secret
	payload := base64.StdEncoding.EncodeToString([]byte("nonce=fake&email=x@y.com&external_id=1&username=x"))
	mac := hmac.New(sha256.New, []byte("wrong_secret"))
	mac.Write([]byte(payload))
	badSig := hex.EncodeToString(mac.Sum(nil))

	resp, err := http.Get(ts.URL + "/session/sso_login?sso=" + url.QueryEscape(payload) + "&sig=" + url.QueryEscape(badSig))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403 for bad signature, got %d", resp.StatusCode)
	}
}

func TestSSO_DisabledReturnsFallback(t *testing.T) {
	// When SSO is not configured, /session/sso returns JSON stub
	ts := testServer(t)
	defer ts.Close()
	resp, body := apiGet(ts, "/session/sso")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	data := parseJSON(t, body)
	if _, ok := data["sso_url"]; !ok {
		t.Fatal("expected sso_url in fallback response")
	}
}
