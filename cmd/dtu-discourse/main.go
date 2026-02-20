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
	"github.com/lightcap/dtu-discourse/internal/webhook"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4200"
	}

	s := store.New()

	// SSO configuration (optional)
	s.SSOSecret = os.Getenv("DISCOURSE_CONNECT_SECRET")
	s.SSOCallbackURL = os.Getenv("SSO_CALLBACK_URL")

	// Webhook dispatcher (optional)
	webhookURL := os.Getenv("WEBHOOK_URL")
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	dispatcher := webhook.New(webhookURL, webhookSecret)

	mux := BuildRouter(s, dispatcher)

	wrapped := middleware.Auth(s)(mux)

	log.Printf("DTU Discourse listening on :%s", port)
	log.Printf("Default API key: test_api_key (user: system)")
	log.Printf("Admin API key:   admin_api_key (user: admin)")
	if s.SSOSecret != "" && s.SSOCallbackURL != "" {
		log.Printf("SSO enabled → callback: %s", s.SSOCallbackURL)
	}
	if webhookURL != "" {
		log.Printf("Webhooks enabled → %s", webhookURL)
	}
	if err := http.ListenAndServe(":"+port, wrapped); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

// BuildRouter creates the HTTP mux with all Discourse-compatible routes.
// Wildcard path segments (e.g. {username}) will match values with or without
// a .json suffix; handlers strip the suffix when extracting the value.
func BuildRouter(s *store.Store, dispatcher *webhook.Dispatcher) *http.ServeMux {
	mux := http.NewServeMux()

	// ---- Core handlers ----
	users := &handler.UsersHandler{Store: s}
	cats := &handler.CategoriesHandler{Store: s}
	topics := &handler.TopicsHandler{Store: s}
	posts := &handler.PostsHandler{Store: s, Webhook: dispatcher}
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

	// ---- Extended handlers ----
	extTopics := &handler.ExtendedTopicsHandler{Store: s}
	extPosts := &handler.ExtendedPostsHandler{Store: s}
	extAdmin := &handler.ExtendedAdminHandler{Store: s}
	misc := &handler.MiscHandler{Store: s}
	extUsers := &handler.ExtendedUsersHandler{Store: s}
	session := &handler.SessionHandler{Store: s}
	extPM := &handler.ExtendedPMHandler{Store: s}
	tagGroups := &handler.TagGroupsHandler{Store: s}
	extNotifs := &handler.ExtendedNotificationsHandler{Store: s}
	extGroups := &handler.ExtendedGroupsHandler{Store: s}
	extCats := &handler.ExtendedCategoriesHandler{Store: s}
	extTags := &handler.ExtendedTagsHandler{Store: s}
	extUploads := &handler.ExtendedUploadsHandler{Store: s}
	extBackups := &handler.ExtendedBackupsHandler{Store: s}
	polls := &handler.PollsHandler{Store: s}
	apiKeys := &handler.APIKeysHandler{Store: s}
	email := &handler.EmailHandler{Store: s}
	userActions := &handler.UserActionsHandler{Store: s}
	topicTimings := &handler.TopicTimingsHandler{Store: s}

	// ==================================================================
	// Users (core)
	// ==================================================================
	mux.HandleFunc("GET /users/check_username.json", users.CheckUsername)
	mux.HandleFunc("GET /users/hp.json", session.Honeypot)
	mux.HandleFunc("GET /users/search/users.json", extUsers.SearchUsers)
	mux.HandleFunc("GET /users/by-external/{external_id}", users.GetUserByExternalID)
	mux.HandleFunc("GET /users/{username}", users.GetUser)
	// GET /admin/users/ uses sub-router to avoid conflicts between
	// /admin/users/{id}/ip_info and /admin/users/list/{type}
	adminUsersRouter := &handler.AdminUsersSubRouter{Users: users, Extended: extAdmin}
	mux.HandleFunc("GET /admin/users/{rest...}", adminUsersRouter.ServeGET)
	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("POST /users.json", users.CreateUser)
	mux.HandleFunc("PUT /u/{username}/preferences/email", users.UpdateEmail)
	mux.HandleFunc("PUT /u/{username}/preferences/username", users.UpdateUsername)
	mux.HandleFunc("PUT /admin/users/{id}/approve", users.Approve)
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
	mux.HandleFunc("PUT /admin/users/{id}/anonymize.json", users.Anonymize)
	mux.HandleFunc("POST /admin/users/{id}/log_out", users.LogOut)
	mux.HandleFunc("POST /admin/users/{id}/log_out.json", users.LogOut)
	mux.HandleFunc("PUT /admin/users/{id}/suspend.json", users.Suspend)
	mux.HandleFunc("DELETE /admin/users/{id}", users.DeleteUser)

	// ==================================================================
	// Users (extended profile, preferences, 2FA)
	// ==================================================================
	mux.HandleFunc("GET /u/search/users", extUsers.SearchUsers)
	mux.HandleFunc("PUT /u/{username}/preferences/avatar/pick", extUsers.PickAvatar)
	mux.HandleFunc("PUT /u/{username}/preferences/avatar/pick.json", extUsers.PickAvatar)
	mux.HandleFunc("PUT /u/{username}/preferences/badge_title", extUsers.SetBadgeTitle)
	mux.HandleFunc("PUT /u/{username}/preferences/categories", extUsers.SetCategoryPreferences)
	mux.HandleFunc("PUT /u/{username}/preferences/tags", extUsers.SetTagPreferences)
	mux.HandleFunc("PUT /u/{username}/preferences/second-factor", extUsers.SecondFactorPref)
	mux.HandleFunc("DELETE /u/{username}/preferences/user_image", extUsers.DeleteUserImage)
	mux.HandleFunc("PUT /u/{username}/clear-featured-topic", extUsers.ClearFeaturedTopic)
	mux.HandleFunc("PUT /u/{username}/feature-topic", extUsers.FeatureTopic)
	// NOTE: /u/{username} GET must be registered after more-specific /u/ routes
	mux.HandleFunc("PUT /u/{username}", users.UpdateUser)
	mux.HandleFunc("GET /u/{username}/summary", extUsers.Summary)
	mux.HandleFunc("GET /u/{username}/summary.json", extUsers.Summary)
	mux.HandleFunc("GET /u/{username}/activity", extUsers.Activity)
	mux.HandleFunc("GET /u/{username}/activity.json", extUsers.Activity)
	mux.HandleFunc("GET /u/{username}/activity/topics", extUsers.ActivityTopics)
	mux.HandleFunc("GET /u/{username}/activity/replies", extUsers.ActivityReplies)
	mux.HandleFunc("GET /u/{username}/bookmarks", extUsers.Bookmarks)
	mux.HandleFunc("GET /u/{username}/bookmarks.json", extUsers.Bookmarks)
	mux.HandleFunc("GET /u/{username}/badges", extUsers.Badges)
	mux.HandleFunc("GET /u/{username}/badges.json", extUsers.Badges)
	mux.HandleFunc("GET /u/{username}/emails", extUsers.Emails)
	mux.HandleFunc("GET /u/{username}/emails.json", extUsers.Emails)
	mux.HandleFunc("GET /u/{username}/notifications", extUsers.Notifications)
	mux.HandleFunc("GET /u/{username}/notifications.json", extUsers.Notifications)
	mux.HandleFunc("GET /u/{username}/messages", extUsers.Messages)
	mux.HandleFunc("GET /u/{username}/messages.json", extUsers.Messages)
	mux.HandleFunc("GET /u/{username}/drafts", extUsers.Drafts)
	mux.HandleFunc("GET /u/{username}/drafts.json", extUsers.Drafts)
	mux.HandleFunc("GET /u/{username}/card.json", extUsers.Card)
	mux.HandleFunc("GET /u/{username}", users.GetUser)
	mux.HandleFunc("POST /u/password-reset/{token}", extUsers.PasswordReset)
	mux.HandleFunc("POST /u/confirm-email/{token}", extUsers.ConfirmEmail)
	mux.HandleFunc("POST /u/second_factors", extUsers.SecondFactors)
	mux.HandleFunc("PUT /u/second_factor", extUsers.UpdateSecondFactor)
	mux.HandleFunc("POST /u/create_second_factor_totp", extUsers.CreateTOTP)
	mux.HandleFunc("POST /u/enable_second_factor_totp", extUsers.EnableTOTP)
	mux.HandleFunc("PUT /u/security_key", extUsers.SecurityKey)
	mux.HandleFunc("POST /u/create_second_factor_security_key", extUsers.CreateSecurityKey)
	mux.HandleFunc("POST /u/register_second_factor_security_key", extUsers.RegisterSecurityKey)
	mux.HandleFunc("POST /u/second_factors_backup", extUsers.BackupCodes)
	mux.HandleFunc("GET /user-stat/{username}", extUsers.UserStat)

	// ==================================================================
	// Categories (core + extended)
	// ==================================================================
	mux.HandleFunc("GET /categories.json", cats.List)
	mux.HandleFunc("GET /categories", cats.List)
	mux.HandleFunc("GET /categories/search", extCats.Search)
	mux.HandleFunc("GET /categories/find", extCats.Find)
	mux.HandleFunc("GET /categories_and_latest", extCats.CategoriesAndLatest)
	mux.HandleFunc("GET /categories_and_top", extCats.CategoriesAndTop)
	mux.HandleFunc("POST /categories.json", cats.Create)
	mux.HandleFunc("POST /categories", cats.Create)
	mux.HandleFunc("POST /categories/reorder", cats.Reorder)
	mux.HandleFunc("POST /categories/{id}/move", extCats.Move)
	mux.HandleFunc("PUT /categories/{id}", cats.Update)
	mux.HandleFunc("DELETE /categories/{id}", cats.Delete)
	mux.HandleFunc("GET /c/{category_slug}/l/latest.json", cats.LatestTopics)
	mux.HandleFunc("GET /c/{category_slug}/l/top.json", cats.TopTopics)
	mux.HandleFunc("GET /c/{category_slug}/l/new.json", cats.NewTopics)
	mux.HandleFunc("GET /c/{category_slug}/l/hot.json", cats.LatestTopics)
	// Eve app uses /c/{slug}/{id}/l/latest.json pattern
	mux.HandleFunc("GET /c/{category_slug}/{category_id}/l/latest.json", cats.LatestTopicsBySlugAndID)
	mux.HandleFunc("GET /c/{category_slug}/{category_id}/l/hot.json", cats.LatestTopicsBySlugAndID)
	mux.HandleFunc("GET /c/{category_slug}/{category_id}/l/top.json", cats.LatestTopicsBySlugAndID)
	mux.HandleFunc("GET /c/{category_slug}/{category_id}/l/new.json", cats.LatestTopicsBySlugAndID)
	mux.HandleFunc("GET /c/{slug}/visible_groups", extCats.VisibleGroups)
	mux.HandleFunc("GET /c/{id}/show", cats.Show)
	mux.HandleFunc("GET /c/{id}/show.json", cats.Show)
	mux.HandleFunc("GET /c/{id}", cats.ListTopics)
	mux.HandleFunc("GET /c/{slug}/{id}", cats.ListTopics)
	mux.HandleFunc("POST /category/{category_id}/notifications", cats.SetNotificationLevel)

	// ==================================================================
	// Topics (core + extended via sub-router)
	// ==================================================================
	topicRouter := &handler.TopicSubRouter{Topics: topics, Extended: extTopics}
	mux.HandleFunc("GET /latest.json", topics.Latest)
	mux.HandleFunc("GET /top.json", topics.Top)
	mux.HandleFunc("GET /top/all.json", topics.Top)
	mux.HandleFunc("GET /new.json", topics.New)
	mux.HandleFunc("GET /t/{rest...}", topicRouter.ServeGET)
	mux.HandleFunc("PUT /t/{rest...}", topicRouter.ServePUT)
	mux.HandleFunc("POST /t/{rest...}", topicRouter.ServePOST)
	mux.HandleFunc("DELETE /t/{rest...}", topicRouter.ServeDELETE)
	mux.HandleFunc("GET /topics/created-by/{username}", topics.TopicsByUser)
	mux.HandleFunc("GET /topics/feature_stats", extTopics.FeatureStats)

	// Topics timings/bulk/similar/reset
	mux.HandleFunc("POST /topics/timings", topicTimings.Record)
	mux.HandleFunc("GET /topics/similar_to", topicTimings.SimilarTo)
	mux.HandleFunc("PUT /topics/bulk", topicTimings.Bulk)
	mux.HandleFunc("PUT /topics/reset-new", topicTimings.ResetNew)
	mux.HandleFunc("PUT /topics/pm-reset-new", topicTimings.PMResetNew)

	// ==================================================================
	// Posts (core + extended)
	// ==================================================================
	postRouter := &handler.PostSubRouter{Posts: posts, Extended: extPosts}
	mux.HandleFunc("POST /posts", posts.Create)
	mux.HandleFunc("POST /posts.json", posts.Create)
	mux.HandleFunc("GET /posts.json", posts.List)
	// GET /posts/... uses sub-router to avoid ServeMux conflicts between
	// /posts/{id}/revisions/... and /posts/by_number/{topic_id}/{post_number}
	mux.HandleFunc("GET /posts/{rest...}", postRouter.ServeGET)
	mux.HandleFunc("DELETE /posts/destroy_many", extPosts.DestroyMany)
	mux.HandleFunc("PUT /posts/merge_posts", extPosts.MergePosts)
	mux.HandleFunc("PUT /posts/{id}/wiki", posts.Wikify)
	mux.HandleFunc("PUT /posts/{id}/recover", extPosts.Recover)
	mux.HandleFunc("PUT /posts/{id}/rebake", extPosts.Rebake)
	mux.HandleFunc("PUT /posts/{id}/locked", extPosts.Locked)
	mux.HandleFunc("PUT /posts/{id}/post_type", extPosts.PostType)
	mux.HandleFunc("PUT /posts/{id}/unhide", extPosts.Unhide)
	mux.HandleFunc("PUT /posts/{id}/notice", extPosts.Notice)
	mux.HandleFunc("DELETE /posts/{id}/revisions/permanently_delete", extPosts.PermanentlyDeleteRevisions)
	mux.HandleFunc("PUT /posts/{id}/revisions/{revision}/hide", extPosts.HideRevision)
	mux.HandleFunc("PUT /posts/{id}/revisions/{revision}/show", extPosts.ShowRevision)
	mux.HandleFunc("PUT /posts/{id}/revisions/{revision}/revert", extPosts.RevertRevision)
	mux.HandleFunc("DELETE /posts/{id}/bookmark", extPosts.RemoveBookmark)
	mux.HandleFunc("PUT /posts/{id}", posts.Update)
	mux.HandleFunc("DELETE /posts/{id}", posts.Delete)
	mux.HandleFunc("POST /post_actions", posts.CreateAction)
	mux.HandleFunc("POST /post_actions.json", posts.CreateAction)
	mux.HandleFunc("DELETE /post_actions/{id}", posts.DeleteAction)
	mux.HandleFunc("GET /post_action_users.json", posts.ActionUsers)
	mux.HandleFunc("GET /raw/{topic_id}/{post_number}", extPosts.RawByNumber)

	// ==================================================================
	// Groups (core + extended)
	// ==================================================================
	mux.HandleFunc("GET /groups.json", groups.List)
	mux.HandleFunc("GET /groups/search.json", groups.List)
	mux.HandleFunc("GET /groups/{group_name}/members.json", groups.Members)
	mux.HandleFunc("PUT /groups/{group}/join", extGroups.Join)
	mux.HandleFunc("DELETE /groups/{group}/leave", extGroups.Leave)
	mux.HandleFunc("POST /groups/{group}/request_membership", extGroups.RequestMembership)
	mux.HandleFunc("GET /groups/{group}/requests", extGroups.ListRequests)
	mux.HandleFunc("PUT /groups/{group}/handle_membership_request", extGroups.HandleRequest)
	mux.HandleFunc("GET /groups/{group}/logs", extGroups.Logs)
	mux.HandleFunc("GET /groups/{group}/permissions", extGroups.Permissions)
	mux.HandleFunc("GET /groups/{group}/mentionable", extGroups.Mentionable)
	mux.HandleFunc("GET /groups/{group}/messageable", extGroups.Messageable)
	mux.HandleFunc("GET /groups/{group}/counts", extGroups.Counts)
	mux.HandleFunc("GET /groups/{group}/topics", extGroups.Topics)
	mux.HandleFunc("POST /groups/{group}/notifications", groups.SetNotificationLevel)
	mux.HandleFunc("GET /groups/{group_name}", groups.Get)
	mux.HandleFunc("POST /admin/groups", groups.Create)
	mux.HandleFunc("POST /admin/groups.json", groups.Create)
	mux.HandleFunc("PUT /admin/groups/{group_id}/members.json", groups.AddMembers)
	mux.HandleFunc("DELETE /admin/groups/{group_id}/members.json", groups.RemoveMembers)
	mux.HandleFunc("PUT /admin/groups/{group_id}/owners.json", groups.AddOwners)
	mux.HandleFunc("DELETE /admin/groups/{group_id}/owners.json", groups.RemoveOwners)
	mux.HandleFunc("DELETE /admin/groups/{group_id}", groups.Delete)
	mux.HandleFunc("PUT /groups/{group_id}", groups.Update)

	// ==================================================================
	// Search
	// ==================================================================
	mux.HandleFunc("GET /search", search.Search)
	mux.HandleFunc("GET /search.json", search.Search)
	mux.HandleFunc("GET /search/query", misc.SearchQuery)
	mux.HandleFunc("POST /search/click", misc.SearchClick)

	// ==================================================================
	// Tags (core + extended)
	// ==================================================================
	mux.HandleFunc("GET /tags.json", tags.List)
	mux.HandleFunc("GET /tags/filter/search", extTags.Search)
	mux.HandleFunc("GET /tags/unused", extTags.Unused)
	mux.HandleFunc("POST /tags", extTags.Create)
	mux.HandleFunc("GET /tags/{tag}/info", extTags.Info)
	mux.HandleFunc("POST /tags/{tag}/synonyms", extTags.AddSynonym)
	mux.HandleFunc("DELETE /tags/{tag}/synonyms/{synonym}", extTags.RemoveSynonym)
	mux.HandleFunc("PUT /tags/{tag}", extTags.Update)
	mux.HandleFunc("DELETE /tags/{tag}", extTags.Delete)
	// Eve app uses /tag/{tag}/l/{variant}.json for tag-filtered topics
	mux.HandleFunc("GET /tag/{tag}/l/latest.json", tags.TopicsByTag)
	mux.HandleFunc("GET /tag/{tag}/l/hot.json", tags.TopicsByTag)
	mux.HandleFunc("GET /tag/{tag}/l/top.json", tags.TopicsByTag)
	mux.HandleFunc("GET /tag/{tag}/l/new.json", tags.TopicsByTag)
	// Eve app uses /tags/c/{slug}/{id}/{tag}/l/{variant}.json for category+tag combos
	mux.HandleFunc("GET /tags/c/{category_slug}/{category_id}/{tag}/l/latest.json", tags.TopicsByCategoryAndTag)
	mux.HandleFunc("GET /tags/c/{category_slug}/{category_id}/{tag}/l/hot.json", tags.TopicsByCategoryAndTag)
	mux.HandleFunc("GET /tags/c/{category_slug}/{category_id}/{tag}/l/top.json", tags.TopicsByCategoryAndTag)
	mux.HandleFunc("GET /tags/c/{category_slug}/{category_id}/{tag}/l/new.json", tags.TopicsByCategoryAndTag)
	mux.HandleFunc("GET /tag/{tag}/notifications", extTags.GetNotifications)
	mux.HandleFunc("PUT /tag/{tag}/notifications", extTags.SetNotifications)
	mux.HandleFunc("GET /tag/{tag}", tags.Show)

	// Tag Groups
	mux.HandleFunc("GET /tag_groups", tagGroups.List)
	mux.HandleFunc("GET /tag_groups.json", tagGroups.List)
	mux.HandleFunc("POST /tag_groups", tagGroups.Create)
	mux.HandleFunc("POST /tag_groups.json", tagGroups.Create)
	mux.HandleFunc("GET /tag_groups/{id}", tagGroups.Show)
	mux.HandleFunc("PUT /tag_groups/{id}", tagGroups.Update)
	mux.HandleFunc("DELETE /tag_groups/{id}", tagGroups.Delete)

	// ==================================================================
	// Badges
	// ==================================================================
	mux.HandleFunc("GET /admin/badges.json", badges.List)
	mux.HandleFunc("POST /admin/badges.json", badges.Create)
	mux.HandleFunc("PUT /admin/badges/{id}", badges.Update)
	mux.HandleFunc("DELETE /admin/badges/{id}", badges.Delete)
	mux.HandleFunc("GET /user-badges/{username}", badges.UserBadges)
	mux.HandleFunc("POST /user_badges", badges.Grant)
	mux.HandleFunc("POST /user_badges.json", badges.Grant)

	// ==================================================================
	// Notifications (core + extended)
	// ==================================================================
	mux.HandleFunc("GET /notifications.json", notifs.List)
	mux.HandleFunc("PUT /notifications/mark-read", extNotifs.MarkRead)
	mux.HandleFunc("GET /notifications/totals", extNotifs.Totals)
	mux.HandleFunc("GET /notifications/{id}", extNotifs.Show)
	mux.HandleFunc("PUT /notifications/{id}", extNotifs.Update)
	mux.HandleFunc("DELETE /notifications/{id}", extNotifs.Delete)

	// ==================================================================
	// Invites
	// ==================================================================
	mux.HandleFunc("POST /invites", invites.Create)
	mux.HandleFunc("POST /invites.json", invites.Create)
	mux.HandleFunc("GET /invites/retrieve.json", invites.Retrieve)
	mux.HandleFunc("PUT /invites/{invite_id}", invites.Update)
	mux.HandleFunc("DELETE /invites", invites.Destroy)
	mux.HandleFunc("POST /invites/destroy-all-expired", invites.DestroyAllExpired)
	mux.HandleFunc("POST /invites/reinvite-all", invites.ResendAll)
	mux.HandleFunc("POST /invites/reinvite", invites.Resend)
	mux.HandleFunc("POST /invite-token/generate", invites.GenerateToken)

	// ==================================================================
	// Uploads (core + extended)
	// ==================================================================
	mux.HandleFunc("POST /uploads", uploads.Create)
	mux.HandleFunc("POST /uploads.json", uploads.Create)
	mux.HandleFunc("GET /uploads/lookup-metadata", extUploads.LookupMetadata)
	mux.HandleFunc("GET /uploads/lookup-urls", extUploads.LookupURLs)
	mux.HandleFunc("POST /uploads/generate-presigned-put", extUploads.GeneratePresignedPut)
	mux.HandleFunc("POST /uploads/complete-external-upload", extUploads.CompleteExternalUpload)
	mux.HandleFunc("POST /uploads/create-multipart", extUploads.CreateMultipart)
	mux.HandleFunc("POST /uploads/batch-presign-multipart-parts", extUploads.BatchPresignParts)
	mux.HandleFunc("POST /uploads/complete-multipart", extUploads.CompleteMultipart)
	mux.HandleFunc("POST /uploads/abort-multipart", extUploads.AbortMultipart)

	// ==================================================================
	// Private Messages (core + extended)
	// ==================================================================
	mux.HandleFunc("GET /topics/private-messages/{username}", pm.Inbox)
	mux.HandleFunc("GET /topics/private-messages-sent/{username}", pm.Sent)
	mux.HandleFunc("GET /topics/private-messages-unread/{username}", extPM.Unread)
	mux.HandleFunc("GET /topics/private-messages-archive/{username}", extPM.Archive)
	mux.HandleFunc("GET /topics/private-messages-new/{username}", extPM.New)
	mux.HandleFunc("GET /topics/private-messages-warnings/{username}", extPM.Warnings)
	mux.HandleFunc("GET /topics/private-messages-group/{username}/{group_name}", extPM.GroupPMs)
	mux.HandleFunc("GET /topics/private-messages-tags/{username}/{tag}", extPM.PMTags)
	mux.HandleFunc("PUT /topics/private-messages/{username}/archive", extPM.MoveToArchive)
	mux.HandleFunc("PUT /topics/private-messages/{username}/move-to-inbox", extPM.MoveToInbox)

	// ==================================================================
	// Admin / Site (core)
	// ==================================================================
	mux.HandleFunc("GET /admin/site_settings.json", admin.GetSiteSettings)
	mux.HandleFunc("PUT /admin/site_settings/{name}", admin.UpdateSiteSetting)
	mux.HandleFunc("GET /admin/dashboard.json", admin.Dashboard)
	mux.HandleFunc("GET /site.json", admin.SiteInfo)
	mux.HandleFunc("GET /session/csrf.json", admin.CSRFToken)

	// ==================================================================
	// Admin Extended
	// ==================================================================
	// Webhooks
	mux.HandleFunc("GET /admin/api/web_hooks.json", extAdmin.ListWebhooks)
	mux.HandleFunc("POST /admin/api/web_hooks.json", extAdmin.CreateWebhook)
	mux.HandleFunc("GET /admin/api/web_hooks/{id}", extAdmin.ShowWebhook)
	mux.HandleFunc("PUT /admin/api/web_hooks/{id}", extAdmin.UpdateWebhook)
	mux.HandleFunc("DELETE /admin/api/web_hooks/{id}", extAdmin.DeleteWebhook)
	mux.HandleFunc("GET /admin/api/web_hooks/{id}/events", extAdmin.WebhookEvents)
	mux.HandleFunc("POST /admin/api/web_hooks/{id}/ping", extAdmin.PingWebhook)

	// Themes
	mux.HandleFunc("GET /admin/themes.json", extAdmin.ListThemes)
	mux.HandleFunc("POST /admin/themes.json", extAdmin.CreateTheme)
	mux.HandleFunc("GET /admin/themes/{id}", extAdmin.ShowTheme)
	mux.HandleFunc("PUT /admin/themes/{id}", extAdmin.UpdateTheme)
	mux.HandleFunc("DELETE /admin/themes/{id}", extAdmin.DeleteTheme)
	mux.HandleFunc("POST /admin/themes/import", extAdmin.ImportTheme)
	mux.HandleFunc("GET /admin/themes/{id}/export", extAdmin.ExportTheme)

	// Color Schemes
	mux.HandleFunc("GET /admin/color_schemes.json", extAdmin.ListColorSchemes)
	mux.HandleFunc("POST /admin/color_schemes.json", extAdmin.CreateColorScheme)
	mux.HandleFunc("PUT /admin/color_schemes/{id}", extAdmin.UpdateColorScheme)
	mux.HandleFunc("DELETE /admin/color_schemes/{id}", extAdmin.DeleteColorScheme)

	// Watched Words
	mux.HandleFunc("GET /admin/customize/watched_words.json", extAdmin.ListWatchedWords)
	mux.HandleFunc("POST /admin/customize/watched_words.json", extAdmin.CreateWatchedWord)
	mux.HandleFunc("DELETE /admin/customize/watched_words/{id}", extAdmin.DeleteWatchedWord)
	mux.HandleFunc("POST /admin/customize/watched_words/upload", extAdmin.UploadWatchedWords)
	mux.HandleFunc("DELETE /admin/customize/watched_words/action/{action}", extAdmin.ClearWatchedWordsAction)

	// Site Texts
	mux.HandleFunc("GET /admin/customize/site_texts.json", extAdmin.ListSiteTexts)
	mux.HandleFunc("GET /admin/customize/site_texts/{id}", extAdmin.ShowSiteText)
	mux.HandleFunc("PUT /admin/customize/site_texts/{id}", extAdmin.UpdateSiteText)
	mux.HandleFunc("DELETE /admin/customize/site_texts/{id}", extAdmin.RevertSiteText)

	// Permalinks
	mux.HandleFunc("GET /admin/permalinks.json", extAdmin.ListPermalinks)
	mux.HandleFunc("POST /admin/permalinks.json", extAdmin.CreatePermalink)
	mux.HandleFunc("PUT /admin/permalinks/{id}", extAdmin.UpdatePermalink)
	mux.HandleFunc("DELETE /admin/permalinks/{id}", extAdmin.DeletePermalink)

	// Staff Action Logs
	mux.HandleFunc("GET /admin/logs/staff_action_logs.json", extAdmin.ListStaffActionLogs)
	mux.HandleFunc("GET /admin/logs/staff_action_logs/{id}/diff", extAdmin.StaffActionLogDiff)

	// Screened Emails/IPs/URLs
	mux.HandleFunc("GET /admin/logs/screened_emails.json", extAdmin.ListScreenedEmails)
	mux.HandleFunc("DELETE /admin/logs/screened_emails/{id}", extAdmin.DeleteScreenedEmail)
	mux.HandleFunc("GET /admin/logs/screened_ip_addresses.json", extAdmin.ListScreenedIPs)
	mux.HandleFunc("POST /admin/logs/screened_ip_addresses.json", extAdmin.CreateScreenedIP)
	mux.HandleFunc("PUT /admin/logs/screened_ip_addresses/{id}", extAdmin.UpdateScreenedIP)
	mux.HandleFunc("DELETE /admin/logs/screened_ip_addresses/{id}", extAdmin.DeleteScreenedIP)
	mux.HandleFunc("GET /admin/logs/screened_urls.json", extAdmin.ListScreenedURLs)

	// Search Logs
	mux.HandleFunc("GET /admin/logs/search_logs.json", extAdmin.ListSearchLogs)
	mux.HandleFunc("GET /admin/logs/search_logs/term.json", extAdmin.SearchLogTerms)

	// Embedding
	mux.HandleFunc("GET /admin/customize/embedding.json", extAdmin.ShowEmbedding)
	mux.HandleFunc("PUT /admin/customize/embedding.json", extAdmin.UpdateEmbedding)
	mux.HandleFunc("POST /admin/embeddable_hosts.json", extAdmin.CreateEmbeddableHost)
	mux.HandleFunc("PUT /admin/embeddable_hosts/{id}", extAdmin.UpdateEmbeddableHost)
	mux.HandleFunc("DELETE /admin/embeddable_hosts/{id}", extAdmin.DeleteEmbeddableHost)

	// Custom User Fields
	mux.HandleFunc("GET /admin/customize/user_fields.json", extAdmin.ListCustomUserFields)
	mux.HandleFunc("POST /admin/customize/user_fields.json", extAdmin.CreateCustomUserField)
	mux.HandleFunc("PUT /admin/customize/user_fields/{id}", extAdmin.UpdateCustomUserField)
	mux.HandleFunc("DELETE /admin/customize/user_fields/{id}", extAdmin.DeleteCustomUserField)

	// Review Queue
	mux.HandleFunc("GET /review.json", extAdmin.ListReviewables)
	mux.HandleFunc("GET /review/count.json", extAdmin.ReviewableCount)
	mux.HandleFunc("GET /review/topics.json", extAdmin.ReviewableTopics)
	mux.HandleFunc("GET /review/settings.json", extAdmin.ReviewSettings)
	mux.HandleFunc("PUT /review/settings.json", extAdmin.UpdateReviewSettings)
	mux.HandleFunc("GET /review/{id}", extAdmin.ShowReviewable)
	mux.HandleFunc("PUT /review/{id}/perform/{action}", extAdmin.PerformReviewAction)
	mux.HandleFunc("PUT /review/{id}", extAdmin.UpdateReviewable)
	mux.HandleFunc("DELETE /review/{id}", extAdmin.DeleteReviewable)

	// Admin Flags
	mux.HandleFunc("GET /admin/config/flags.json", extAdmin.ListFlags)
	mux.HandleFunc("POST /admin/config/flags.json", extAdmin.CreateFlag)
	mux.HandleFunc("PUT /admin/config/flags/{id}", extAdmin.UpdateFlag)
	mux.HandleFunc("DELETE /admin/config/flags/{id}", extAdmin.DeleteFlag)
	mux.HandleFunc("PUT /admin/config/flags/{id}/toggle", extAdmin.ToggleFlag)

	// Reports
	mux.HandleFunc("GET /admin/reports.json", extAdmin.ListReports)
	mux.HandleFunc("GET /admin/reports/{type}", extAdmin.ShowReport)
	mux.HandleFunc("GET /admin/reports/bulk.json", extAdmin.BulkReports)

	// Dashboard sub-routes
	mux.HandleFunc("GET /admin/dashboard/general.json", extAdmin.DashboardGeneral)
	mux.HandleFunc("GET /admin/dashboard/moderation.json", extAdmin.DashboardModeration)
	mux.HandleFunc("GET /admin/dashboard/security.json", extAdmin.DashboardSecurity)
	mux.HandleFunc("GET /admin/dashboard/problems.json", extAdmin.DashboardProblems)

	// Version Check
	mux.HandleFunc("GET /admin/version_check.json", extAdmin.VersionCheck)

	// User Admin ops
	mux.HandleFunc("PUT /admin/users/{id}/silence", extAdmin.SilenceUser)
	mux.HandleFunc("PUT /admin/users/{id}/unsilence", extAdmin.UnsilenceUser)
	mux.HandleFunc("PUT /admin/users/{id}/primary_group", extAdmin.SetPrimaryGroup)
	mux.HandleFunc("DELETE /admin/users/{id}/posts_batch", extAdmin.DeletePostsBatch)
	mux.HandleFunc("POST /admin/users/{id}/merge", extAdmin.MergeUsers)
	mux.HandleFunc("POST /admin/users/{id}/reset_bounce_score", extAdmin.ResetBounceScore)
	mux.HandleFunc("PUT /admin/users/approve-bulk", extAdmin.BulkApproveUsers)
	mux.HandleFunc("DELETE /admin/users/destroy-bulk", extAdmin.BulkDestroyUsers)
	mux.HandleFunc("POST /admin/users/{id}/generate_api_key", extAdmin.GenerateAPIKeyForUser)

	// Impersonation
	mux.HandleFunc("POST /admin/impersonate/{id}", extAdmin.StartImpersonation)
	mux.HandleFunc("DELETE /admin/impersonate", extAdmin.StopImpersonation)

	// Admin Search
	mux.HandleFunc("GET /admin/search/all", extAdmin.AdminSearch)

	// Email Templates
	mux.HandleFunc("GET /admin/customize/email_templates.json", extAdmin.ListEmailTemplates)
	mux.HandleFunc("GET /admin/customize/email_templates/{id}", extAdmin.ShowEmailTemplate)
	mux.HandleFunc("PUT /admin/customize/email_templates/{id}", extAdmin.UpdateEmailTemplate)
	mux.HandleFunc("DELETE /admin/customize/email_templates/{id}", extAdmin.RevertEmailTemplate)

	// Email Style
	mux.HandleFunc("GET /admin/customize/email_style.json", extAdmin.GetEmailStyle)
	mux.HandleFunc("PUT /admin/customize/email_style.json", extAdmin.UpdateEmailStyle)

	// Robots.txt
	mux.HandleFunc("GET /admin/customize/robots.json", extAdmin.GetRobots)
	mux.HandleFunc("PUT /admin/customize/robots.json", extAdmin.UpdateRobots)
	mux.HandleFunc("DELETE /admin/customize/robots.json", extAdmin.ResetRobots)

	// Form Templates
	mux.HandleFunc("GET /admin/customize/form-templates.json", extAdmin.ListFormTemplates)
	mux.HandleFunc("POST /admin/customize/form-templates.json", extAdmin.CreateFormTemplate)
	mux.HandleFunc("PUT /admin/customize/form-templates/{id}", extAdmin.UpdateFormTemplate)
	mux.HandleFunc("DELETE /admin/customize/form-templates/{id}", extAdmin.DeleteFormTemplate)

	// Reseed
	mux.HandleFunc("GET /admin/customize/reseed", extAdmin.GetReseed)
	mux.HandleFunc("POST /admin/customize/reseed", extAdmin.PostReseed)

	// Badge admin extended
	mux.HandleFunc("GET /admin/badge_types.json", extAdmin.BadgeTypes)
	mux.HandleFunc("GET /admin/badge_groupings.json", extAdmin.BadgeGroupings)
	mux.HandleFunc("GET /admin/badges/preview.json", extAdmin.PreviewBadge)

	// Backups (core + extended)
	mux.HandleFunc("GET /admin/backups.json", admin.ListBackups)
	mux.HandleFunc("POST /admin/backups.json", admin.CreateBackup)
	mux.HandleFunc("GET /admin/backups/status", extBackups.Status)
	mux.HandleFunc("GET /admin/backups/logs", extBackups.Logs)
	mux.HandleFunc("GET /admin/backups/is-backup-restore-running", extBackups.IsRunning)
	mux.HandleFunc("DELETE /admin/backups/cancel", extBackups.Cancel)
	mux.HandleFunc("POST /admin/backups/rollback", extBackups.Rollback)
	mux.HandleFunc("PUT /admin/backups/readonly", extBackups.SetReadonly)
	mux.HandleFunc("GET /admin/backups/{filename}/restore", extBackups.Restore)
	mux.HandleFunc("DELETE /admin/backups/{filename}", extBackups.Delete)

	// ==================================================================
	// API Keys
	// ==================================================================
	mux.HandleFunc("GET /admin/api/keys", apiKeys.List)
	mux.HandleFunc("GET /admin/api/keys.json", apiKeys.List)
	mux.HandleFunc("POST /admin/api/keys", apiKeys.Create)
	mux.HandleFunc("POST /admin/api/keys.json", apiKeys.Create)
	mux.HandleFunc("GET /admin/api/keys/scopes", apiKeys.Scopes)
	mux.HandleFunc("POST /admin/api/keys/{id}/revoke", apiKeys.Revoke)
	mux.HandleFunc("POST /admin/api/keys/{id}/undo-revoke", apiKeys.UndoRevoke)
	mux.HandleFunc("DELETE /admin/api/keys/{id}", apiKeys.Delete)

	// ==================================================================
	// Email admin
	// ==================================================================
	mux.HandleFunc("GET /admin/email.json", email.Settings)
	mux.HandleFunc("GET /admin/email/server-settings", email.ServerSettings)
	mux.HandleFunc("GET /admin/email/preview-digest", email.PreviewDigest)
	mux.HandleFunc("POST /admin/email/test", email.Test)
	mux.HandleFunc("GET /admin/email/{filter}", email.List)

	// ==================================================================
	// Polls
	// ==================================================================
	mux.HandleFunc("PUT /polls/vote", polls.Vote)
	mux.HandleFunc("PUT /polls/toggle_status", polls.ToggleStatus)
	mux.HandleFunc("GET /polls/voters.json", polls.Voters)

	// ==================================================================
	// User Actions
	// ==================================================================
	mux.HandleFunc("GET /user_actions.json", userActions.List)

	// ==================================================================
	// SSO
	// ==================================================================
	mux.HandleFunc("POST /admin/users/sync_sso", sso.SyncSSO)
	mux.HandleFunc("POST /admin/users/sync_sso.json", sso.SyncSSO)

	// ==================================================================
	// Session / Auth
	// ==================================================================
	mux.HandleFunc("POST /session", session.Login)
	mux.HandleFunc("GET /session/current.json", session.Current)
	mux.HandleFunc("POST /session/forgot_password", session.ForgotPassword)
	mux.HandleFunc("GET /session/hp", session.Honeypot)
	mux.HandleFunc("GET /session/passkey/challenge", session.PasskeyChallenge)
	mux.HandleFunc("POST /session/passkey/auth", session.PasskeyAuth)
	mux.HandleFunc("POST /session/2fa", session.TwoFactorAuth)
	mux.HandleFunc("GET /session/2fa.json", session.TwoFactorStatus)
	mux.HandleFunc("GET /session/sso", session.SSORedirect)
	mux.HandleFunc("GET /session/sso_provider", session.SSOProvider)
	mux.HandleFunc("GET /session/sso_login", session.SSOLogin)
	mux.HandleFunc("POST /session/email-login/{token}", session.EmailLogin)
	mux.HandleFunc("GET /session/email-login/{token}", session.EmailLoginInfo)
	mux.HandleFunc("DELETE /session/{username}", session.Logout)
	mux.HandleFunc("POST /user-api-key/new", session.NewUserAPIKey)
	mux.HandleFunc("POST /user-api-key", session.CreateUserAPIKey)
	mux.HandleFunc("POST /user-api-key/revoke", session.RevokeUserAPIKey)
	mux.HandleFunc("POST /user-api-key/undo-revoke", session.UndoRevokeUserAPIKey)

	// ==================================================================
	// Misc (undocumented endpoints)
	// ==================================================================
	// Hot/Filter
	mux.HandleFunc("GET /hot.json", misc.HotTopics)
	mux.HandleFunc("GET /filter", misc.FilterTopics)

	// Directory
	mux.HandleFunc("GET /directory_items", misc.DirectoryItems)
	mux.HandleFunc("GET /directory_items.json", misc.DirectoryItems)
	mux.HandleFunc("GET /directory-columns", misc.DirectoryColumns)
	mux.HandleFunc("PUT /edit-directory-columns", misc.EditDirectoryColumns)

	// About
	mux.HandleFunc("GET /about", misc.About)
	mux.HandleFunc("GET /about.json", misc.About)
	mux.HandleFunc("GET /about/live_post_counts", misc.LivePostCounts)

	// Site info
	mux.HandleFunc("GET /site/basic-info.json", misc.SiteBasicInfo)
	mux.HandleFunc("GET /site/statistics.json", misc.SiteStatistics)

	// Drafts
	mux.HandleFunc("GET /drafts", misc.ListDrafts)
	mux.HandleFunc("GET /drafts.json", misc.ListDrafts)
	mux.HandleFunc("POST /drafts", misc.CreateDraft)
	mux.HandleFunc("POST /drafts.json", misc.CreateDraft)
	mux.HandleFunc("GET /drafts/{id}", misc.ShowDraft)
	mux.HandleFunc("DELETE /drafts/{id}", misc.DeleteDraft)

	// Bookmarks (new API)
	mux.HandleFunc("POST /bookmarks", misc.CreateBookmark)
	mux.HandleFunc("POST /bookmarks.json", misc.CreateBookmark)
	mux.HandleFunc("PUT /bookmarks/bulk", misc.BulkBookmarks)
	mux.HandleFunc("PUT /bookmarks/{id}/toggle_pin", misc.ToggleBookmarkPin)
	mux.HandleFunc("PUT /bookmarks/{id}", misc.UpdateBookmark)
	mux.HandleFunc("DELETE /bookmarks/{id}", misc.DeleteBookmark)

	// Published Pages
	mux.HandleFunc("GET /pub/check-slug", misc.CheckPublishedSlug)
	mux.HandleFunc("GET /pub/by-topic/{topic_id}", misc.GetPublishedPage)
	mux.HandleFunc("PUT /pub/by-topic/{topic_id}", misc.UpdatePublishedPage)
	mux.HandleFunc("DELETE /pub/by-topic/{topic_id}", misc.DeletePublishedPage)

	// Sidebar Sections
	mux.HandleFunc("GET /sidebar_sections", misc.ListSidebarSections)
	mux.HandleFunc("POST /sidebar_sections", misc.CreateSidebarSection)
	mux.HandleFunc("PUT /sidebar_sections/reset/{id}", misc.ResetSidebarSection)
	mux.HandleFunc("PUT /sidebar_sections/{id}", misc.UpdateSidebarSection)
	mux.HandleFunc("DELETE /sidebar_sections/{id}", misc.DeleteSidebarSection)

	// Clicks
	mux.HandleFunc("POST /clicks/track", misc.TrackClick)

	// Onebox
	mux.HandleFunc("GET /onebox", misc.Onebox)
	mux.HandleFunc("GET /inline-onebox", misc.InlineOnebox)

	// Slugs
	mux.HandleFunc("POST /slugs", misc.GenerateSlug)

	// Embed
	mux.HandleFunc("GET /embed/topics", misc.EmbedTopics)
	mux.HandleFunc("GET /embed/comments", misc.EmbedComments)
	mux.HandleFunc("GET /embed/count", misc.EmbedCount)
	mux.HandleFunc("GET /embed/info", misc.EmbedInfo)

	// Presence
	mux.HandleFunc("POST /presence/update", misc.PresenceUpdate)
	mux.HandleFunc("GET /presence/get", misc.PresenceGet)

	// User Status / DND
	mux.HandleFunc("GET /user-status", misc.GetUserStatus)
	mux.HandleFunc("PUT /user-status", misc.SetUserStatus)
	mux.HandleFunc("DELETE /user-status", misc.ClearUserStatus)
	mux.HandleFunc("POST /do-not-disturb", misc.EnableDND)
	mux.HandleFunc("DELETE /do-not-disturb", misc.DisableDND)

	// Emojis
	mux.HandleFunc("GET /admin/config/emoji", misc.ListCustomEmojis)
	mux.HandleFunc("POST /admin/config/emoji", misc.CreateCustomEmoji)
	mux.HandleFunc("DELETE /admin/config/emoji/{id}", misc.DeleteCustomEmoji)
	mux.HandleFunc("GET /emojis", misc.ListEmojis)

	// Hashtags
	mux.HandleFunc("GET /hashtags", misc.LookupHashtags)
	mux.HandleFunc("GET /hashtags/search", misc.SearchHashtags)

	// Composer
	mux.HandleFunc("GET /composer/mentions", misc.ComposerMentions)
	mux.HandleFunc("GET /composer_messages", misc.ComposerMessages)

	// Export CSV
	mux.HandleFunc("POST /export_csv/export_entity", misc.ExportCSV)

	// Server Status
	mux.HandleFunc("GET /srv/status", misc.ServerStatus)

	// Permalink Check
	mux.HandleFunc("GET /permalink-check", misc.PermalinkCheck)

	// Pageview
	mux.HandleFunc("POST /pageview", misc.Pageview)

	return mux
}
