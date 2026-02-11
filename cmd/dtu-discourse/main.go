// Command dtu-discourse starts a Digital Twin Universe server that
// replicates the Discourse REST API surface. It is designed to be a
// drop-in replacement for a real Discourse instance during integration
// testing, with 100% wire-compatibility targeted against the three most
// popular public SDK clients:
//
//   - discourse_api  (Ruby, official)
//   - pydiscourse    (Python)
//   - discourse-api  (JavaScript)
//
// All state is held in memory and pre-seeded with realistic data.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lightcap/dtu-discourse/internal/handler"
	"github.com/lightcap/dtu-discourse/internal/middleware"
	"github.com/lightcap/dtu-discourse/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4200"
	}

	s := store.New()
	mux := BuildRouter(s)

	wrapped := middleware.Auth(s)(mux)

	log.Printf("DTU Discourse listening on :%s", port)
	log.Printf("Default API key: test_api_key (user: system)")
	log.Printf("Admin API key:   admin_api_key (user: admin)")
	if err := http.ListenAndServe(":"+port, wrapped); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

// BuildRouter creates the HTTP mux with all Discourse-compatible routes.
// Wildcard path segments (e.g. {username}) will match values with or without
// a .json suffix; handlers strip the suffix when extracting the value.
func BuildRouter(s *store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	users := &handler.UsersHandler{Store: s}
	cats := &handler.CategoriesHandler{Store: s}
	topics := &handler.TopicsHandler{Store: s}
	posts := &handler.PostsHandler{Store: s}
	groups := &handler.GroupsHandler{Store: s}
	search := &handler.SearchHandler{Store: s}
	tags := &handler.TagsHandler{Store: s}
	badges := &handler.BadgesHandler{Store: s}
	notifs := &handler.NotificationsHandler{Store: s}
	invites := &handler.InvitesHandler{Store: s}
	uploads := &handler.UploadsHandler{Store: s}
	pm := &handler.PrivateMessagesHandler{Store: s}
	admin := &handler.AdminHandler{Store: s}
	sso := &handler.SSOHandler{Store: s}

	// ---- Users ----
	// {username} matches "admin" or "admin.json" — handler strips suffix
	mux.HandleFunc("GET /users/check_username.json", users.CheckUsername)
	mux.HandleFunc("GET /users/by-external/{external_id}", users.GetUserByExternalID)
	mux.HandleFunc("GET /users/{username}", users.GetUser)
	mux.HandleFunc("GET /u/{username}", users.GetUser)
	mux.HandleFunc("GET /admin/users/list/{type}", users.ListUsers)
	mux.HandleFunc("GET /admin/users/{id}", users.GetUserByID)
	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("POST /users.json", users.CreateUser)
	mux.HandleFunc("PUT /u/{username}/preferences/email", users.UpdateEmail)
	mux.HandleFunc("PUT /u/{username}/preferences/username", users.UpdateUsername)
	mux.HandleFunc("PUT /u/{username}", users.UpdateUser)
	mux.HandleFunc("PUT /admin/users/{id}/activate", users.Activate)
	mux.HandleFunc("PUT /admin/users/{id}/deactivate", users.Deactivate)
	mux.HandleFunc("PUT /admin/users/{id}/trust_level", users.UpdateTrustLevel)
	mux.HandleFunc("PUT /admin/users/{id}/grant_admin", users.GrantAdmin)
	mux.HandleFunc("PUT /admin/users/{id}/revoke_admin", users.RevokeAdmin)
	mux.HandleFunc("PUT /admin/users/{id}/grant_moderation", users.GrantModeration)
	mux.HandleFunc("PUT /admin/users/{id}/revoke_moderation", users.RevokeModeration)
	mux.HandleFunc("PUT /admin/users/{id}/suspend", users.Suspend)
	mux.HandleFunc("PUT /admin/users/{id}/unsuspend", users.Unsuspend)
	mux.HandleFunc("PUT /admin/users/{id}/anonymize", users.Anonymize)
	mux.HandleFunc("POST /admin/users/{id}/log_out", users.LogOut)
	mux.HandleFunc("DELETE /admin/users/{id}", users.DeleteUser)

	// ---- Categories ----
	mux.HandleFunc("GET /categories.json", cats.List)
	mux.HandleFunc("GET /categories", cats.List)
	mux.HandleFunc("POST /categories.json", cats.Create)
	mux.HandleFunc("POST /categories", cats.Create)
	mux.HandleFunc("POST /categories/reorder", cats.Reorder)
	mux.HandleFunc("PUT /categories/{id}", cats.Update)
	mux.HandleFunc("DELETE /categories/{id}", cats.Delete)
	mux.HandleFunc("GET /c/{category_slug}/l/latest.json", cats.LatestTopics)
	mux.HandleFunc("GET /c/{category_slug}/l/top.json", cats.TopTopics)
	mux.HandleFunc("GET /c/{category_slug}/l/new.json", cats.NewTopics)
	mux.HandleFunc("GET /c/{id}/show", cats.Show)
	mux.HandleFunc("GET /c/{slug}/{id}", cats.ListTopics)
	mux.HandleFunc("POST /category/{category_id}/notifications", cats.SetNotificationLevel)

	// ---- Topics ----
	// /t/ routes use a dispatcher to avoid Go ServeMux wildcard conflicts
	topicRouter := &handler.TopicSubRouter{Topics: topics}
	mux.HandleFunc("GET /latest.json", topics.Latest)
	mux.HandleFunc("GET /top.json", topics.Top)
	mux.HandleFunc("GET /new.json", topics.New)
	mux.HandleFunc("GET /t/{rest...}", topicRouter.ServeGET)
	mux.HandleFunc("PUT /t/{rest...}", topicRouter.ServePUT)
	mux.HandleFunc("POST /t/{rest...}", topicRouter.ServePOST)
	mux.HandleFunc("DELETE /t/{rest...}", topicRouter.ServeDELETE)
	mux.HandleFunc("GET /topics/created-by/{username}", topics.TopicsByUser)

	// ---- Posts ----
	mux.HandleFunc("POST /posts", posts.Create)
	mux.HandleFunc("POST /posts.json", posts.Create)
	mux.HandleFunc("GET /posts.json", posts.List)
	mux.HandleFunc("PUT /posts/{id}/wiki", posts.Wikify)
	mux.HandleFunc("GET /posts/{id}", posts.Get)
	mux.HandleFunc("PUT /posts/{id}", posts.Update)
	mux.HandleFunc("DELETE /posts/{id}", posts.Delete)
	mux.HandleFunc("POST /post_actions", posts.CreateAction)
	mux.HandleFunc("POST /post_actions.json", posts.CreateAction)
	mux.HandleFunc("DELETE /post_actions/{id}", posts.DeleteAction)
	mux.HandleFunc("GET /post_action_users.json", posts.ActionUsers)

	// ---- Groups ----
	mux.HandleFunc("GET /groups.json", groups.List)
	mux.HandleFunc("GET /groups/{group_name}/members.json", groups.Members)
	mux.HandleFunc("GET /groups/{group_name}", groups.Get)
	mux.HandleFunc("POST /admin/groups", groups.Create)
	mux.HandleFunc("POST /admin/groups.json", groups.Create)
	mux.HandleFunc("PUT /admin/groups/{group_id}/members.json", groups.AddMembers)
	mux.HandleFunc("DELETE /admin/groups/{group_id}/members.json", groups.RemoveMembers)
	mux.HandleFunc("PUT /admin/groups/{group_id}/owners.json", groups.AddOwners)
	mux.HandleFunc("DELETE /admin/groups/{group_id}/owners.json", groups.RemoveOwners)
	mux.HandleFunc("DELETE /admin/groups/{group_id}", groups.Delete)
	mux.HandleFunc("PUT /groups/{group_id}", groups.Update)
	mux.HandleFunc("POST /groups/{group}/notifications", groups.SetNotificationLevel)

	// ---- Search ----
	mux.HandleFunc("GET /search", search.Search)
	mux.HandleFunc("GET /search.json", search.Search)

	// ---- Tags ----
	mux.HandleFunc("GET /tags.json", tags.List)
	mux.HandleFunc("GET /tag/{tag}", tags.Show)

	// ---- Badges ----
	mux.HandleFunc("GET /admin/badges.json", badges.List)
	mux.HandleFunc("POST /admin/badges.json", badges.Create)
	mux.HandleFunc("PUT /admin/badges/{id}", badges.Update)
	mux.HandleFunc("DELETE /admin/badges/{id}", badges.Delete)
	mux.HandleFunc("GET /user-badges/{username}", badges.UserBadges)
	mux.HandleFunc("POST /user_badges", badges.Grant)
	mux.HandleFunc("POST /user_badges.json", badges.Grant)

	// ---- Notifications ----
	mux.HandleFunc("GET /notifications.json", notifs.List)

	// ---- Invites ----
	mux.HandleFunc("POST /invites", invites.Create)
	mux.HandleFunc("POST /invites.json", invites.Create)
	mux.HandleFunc("GET /invites/retrieve.json", invites.Retrieve)
	mux.HandleFunc("PUT /invites/{invite_id}", invites.Update)
	mux.HandleFunc("DELETE /invites", invites.Destroy)
	mux.HandleFunc("POST /invites/destroy-all-expired", invites.DestroyAllExpired)
	mux.HandleFunc("POST /invites/reinvite-all", invites.ResendAll)
	mux.HandleFunc("POST /invites/reinvite", invites.Resend)
	mux.HandleFunc("POST /invite-token/generate", invites.GenerateToken)

	// ---- Uploads ----
	mux.HandleFunc("POST /uploads", uploads.Create)
	mux.HandleFunc("POST /uploads.json", uploads.Create)

	// ---- Private Messages ----
	mux.HandleFunc("GET /topics/private-messages/{username}", pm.Inbox)
	mux.HandleFunc("GET /topics/private-messages-sent/{username}", pm.Sent)

	// ---- Admin / Site ----
	mux.HandleFunc("GET /admin/site_settings.json", admin.GetSiteSettings)
	mux.HandleFunc("PUT /admin/site_settings/{name}", admin.UpdateSiteSetting)
	mux.HandleFunc("GET /admin/backups.json", admin.ListBackups)
	mux.HandleFunc("POST /admin/backups.json", admin.CreateBackup)
	mux.HandleFunc("GET /admin/dashboard.json", admin.Dashboard)
	mux.HandleFunc("GET /site.json", admin.SiteInfo)
	mux.HandleFunc("GET /session/csrf.json", admin.CSRFToken)

	// ---- SSO ----
	mux.HandleFunc("POST /admin/users/sync_sso", sso.SyncSSO)
	mux.HandleFunc("POST /admin/users/sync_sso.json", sso.SyncSSO)

	return mux
}
