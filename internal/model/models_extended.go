// Package model defines data types for undocumented and lesser-known Discourse
// API endpoints.  Field names and JSON tags are derived from the actual
// Discourse Ruby serializers so that SDK clients can deserialise responses
// without modification.
package model

import "time"

// ============================================================================
// Polls (plugin: discourse-poll)
//
// Polls are embedded inside posts.  The post JSON contains a "polls" array
// and an optional "polls_votes" map.
//
// Serializers: PollSerializer, PollOptionSerializer
// ============================================================================

// Poll represents a single poll attached to a post.
type Poll struct {
	ID                 int                    `json:"id"`
	Name               string                 `json:"name"`
	Type               string                 `json:"type"`   // "regular", "multiple", "number", "ranked_choice"
	Status             string                 `json:"status"` // "open" or "closed"
	Public             bool                   `json:"public"`
	Dynamic            bool                   `json:"dynamic"`
	Results            string                 `json:"results"` // "always", "on_vote", "on_close", "staff_only"
	Min                *int                   `json:"min"`
	Max                *int                   `json:"max"`
	Step               *int                   `json:"step"`
	Options            []PollOption           `json:"options"`
	Voters             int                    `json:"voters"`
	Close              *time.Time             `json:"close"`
	PreloadedVoters    map[string][]BasicUser `json:"preloaded_voters,omitempty"`
	ChartType          string                 `json:"chart_type"`
	Groups             string                 `json:"groups,omitempty"`
	Title              string                 `json:"title,omitempty"`
	RankedChoiceOutcome map[string]interface{} `json:"ranked_choice_outcome,omitempty"`
}

// PollOption is a single selectable option inside a poll.
type PollOption struct {
	ID    string `json:"id"`   // digest hash, not an integer
	HTML  string `json:"html"` // rendered option text
	Votes int    `json:"votes"`
}

// PollVote records which options a user voted for on a given poll.
type PollVote struct {
	PollName  string   `json:"poll_name"`
	OptionIDs []string `json:"options"`
}

// PollVotersResponse is returned by GET /polls/voters.json.
type PollVotersResponse struct {
	Voters map[string][]BasicUser `json:"voters"`
}

// ============================================================================
// API Keys (admin)
//
// Endpoint: GET/POST /admin/api/keys
// Serializer: ApiKeySerializer, ApiKeyScopeSerializer
// ============================================================================

// APIKeyScope describes a single scope assigned to an API key.
type APIKeyScope struct {
	Resource          string                 `json:"resource"`
	Action            string                 `json:"action"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	URLs              []string               `json:"urls,omitempty"`
	AllowedParameters map[string][]string    `json:"allowed_parameters,omitempty"`
	Key               string                 `json:"key,omitempty"`
}

// APIKeyRecord is a single admin API key.
type APIKeyRecord struct {
	ID           int          `json:"id"`
	Key          string       `json:"key,omitempty"`          // full key, only on create
	TruncatedKey string       `json:"truncated_key"`          // masked version on list
	Description  string       `json:"description"`
	ScopeMode    string       `json:"scope_mode,omitempty"`   // "global", "granular", "read_only"
	LastUsedAt   *time.Time   `json:"last_used_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	RevokedAt    *time.Time   `json:"revoked_at"`
	User         *BasicUser   `json:"user"`                   // nil when scoped to all users
	CreatedBy    *BasicUser   `json:"created_by,omitempty"`
	APIKeyScopes []APIKeyScope `json:"api_key_scopes,omitempty"`
}

// APIKeyListResponse wraps the list returned by GET /admin/api/keys.json.
type APIKeyListResponse struct {
	Keys []APIKeyRecord `json:"keys"`
}

// ============================================================================
// Email Administration
//
// Endpoint: GET /admin/email.json, GET /admin/email/sent.json, etc.
// Serializer: EmailLogSerializer
// ============================================================================

// EmailSettings is the shape returned by GET /admin/email.json.
type EmailSettings struct {
	DeliveryMethod string                 `json:"delivery_method"`
	Settings       map[string]interface{} `json:"settings"`
}

// EmailLog is one row from GET /admin/email/{sent,skipped,bounced,received}.json.
type EmailLog struct {
	ID                       int        `json:"id"`
	To                       string     `json:"to_address,omitempty"`
	CCAddresses              string     `json:"cc_addresses,omitempty"`
	EmailType                string     `json:"email_type"`
	User                     *BasicUser `json:"user,omitempty"`
	PostID                   *int       `json:"post_id"`
	PostURL                  string     `json:"post_url,omitempty"`
	ReplyKey                 string     `json:"reply_key,omitempty"`
	Bounced                  bool       `json:"bounced"`
	HasBounceKey             bool       `json:"has_bounce_key"`
	SMTPTransactionResponse  string     `json:"smtp_transaction_response,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
}

// ============================================================================
// User Actions
//
// Endpoint: GET /user_actions.json?username=…
// Serializer: UserActionSerializer
// ============================================================================

// UserAction is one row from the user activity stream.
type UserAction struct {
	ActionType       int        `json:"action_type"`
	CreatedAt        time.Time  `json:"created_at"`
	AvatarTemplate   string     `json:"avatar_template"`
	ActingAvatarTemplate string `json:"acting_avatar_template,omitempty"`
	Slug             string     `json:"slug"`
	TopicID          int        `json:"topic_id"`
	TargetUserID     *int       `json:"target_user_id"`
	TargetName       string     `json:"target_name,omitempty"`
	TargetUsername   string     `json:"target_username,omitempty"`
	PostNumber       int        `json:"post_number"`
	PostID           *int       `json:"post_id"`
	ReplyToPostNumber *int     `json:"reply_to_post_number"`
	Username         string     `json:"username"`
	Name             string     `json:"name,omitempty"`
	UserID           int        `json:"user_id"`
	ActingUsername   string     `json:"acting_username,omitempty"`
	ActingName       string     `json:"acting_name,omitempty"`
	ActingUserID     *int       `json:"acting_user_id"`
	Title            string     `json:"title"`
	Deleted          bool       `json:"deleted"`
	Hidden           bool       `json:"hidden"`
	PostType         int        `json:"post_type"`
	ActionCode       string     `json:"action_code,omitempty"`
	ActionCodeWho    string     `json:"action_code_who,omitempty"`
	ActionCodePath   string     `json:"action_code_path,omitempty"`
	EditReason       string     `json:"edit_reason,omitempty"`
	CategoryID       int        `json:"category_id"`
	Closed           bool       `json:"closed"`
	Archived         bool       `json:"archived"`
	Excerpt          string     `json:"excerpt,omitempty"`
}

// UserActionListResponse wraps GET /user_actions.json.
type UserActionListResponse struct {
	UserActions []UserAction `json:"user_actions"`
}

// ============================================================================
// Webhooks (admin)
//
// Endpoint: GET /admin/api/web_hooks.json
// Serializer: AdminWebHookSerializer, AdminWebHookEventSerializer
// ============================================================================

// Webhook is an admin-configured outgoing webhook.
type Webhook struct {
	ID                  int        `json:"id"`
	PayloadURL          string     `json:"payload_url"`
	ContentType         int        `json:"content_type"` // 1 = json, 2 = url-encoded
	LastDeliveryStatus  int        `json:"last_delivery_status"`
	Secret              string     `json:"secret,omitempty"`
	WildcardWebHook     bool       `json:"wildcard_web_hook"`
	VerifyCertificate   bool       `json:"verify_certificate"`
	Active              bool       `json:"active"`
	WebHookEventTypeIDs []int      `json:"web_hook_event_type_ids,omitempty"`
	Categories          []Category `json:"categories,omitempty"`
	Tags                []Tag      `json:"tags,omitempty"`
	Groups              []Group    `json:"groups,omitempty"`
	CreatedAt           time.Time  `json:"created_at,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at,omitempty"`
}

// WebhookEvent is one delivery attempt recorded for a webhook.
type WebhookEvent struct {
	ID              int       `json:"id"`
	WebHookID       int       `json:"web_hook_id"`
	RequestURL      string    `json:"request_url"`
	Headers         string    `json:"headers"`
	Payload         string    `json:"payload"`
	Status          int       `json:"status"` // HTTP status code of remote
	ResponseHeaders string    `json:"response_headers"`
	ResponseBody    string    `json:"response_body"`
	Duration        int       `json:"duration"` // milliseconds
	CreatedAt       time.Time `json:"created_at"`
	Redelivering    bool      `json:"redelivering"`
}

// ============================================================================
// Review Queue / Flags
//
// Endpoint: GET /review.json
// Serializer: ReviewableSerializer
// ============================================================================

// ReviewableScore is one score entry on a reviewable item.
type ReviewableScore struct {
	ID            int        `json:"id"`
	Score         float64    `json:"score"`
	AgreedAt      *time.Time `json:"agreed_at"`
	DisagreedAt   *time.Time `json:"disagreed_at"`
	IgnoredAt     *time.Time `json:"ignored_at"`
	CreatedAt     time.Time  `json:"created_at"`
	ReviewableID  int        `json:"reviewable_id"`
	UserID        int        `json:"user_id"`
	ReviewableScoreType int `json:"reviewable_score_type"`
}

// ReviewableItem is a single reviewable entry (flagged post, queued post, user).
type ReviewableItem struct {
	ID                          int               `json:"id"`
	Type                        string            `json:"type"`       // "ReviewableFlaggedPost", "ReviewableQueuedPost", "ReviewableUser"
	TypeSource                  string            `json:"type_source,omitempty"`
	Status                      int               `json:"status"`     // 0=pending, 1=approved, 2=rejected, 3=ignored, 4=deleted
	TopicID                     *int              `json:"topic_id"`
	TopicURL                    string            `json:"topic_url,omitempty"`
	TargetType                  string            `json:"target_type,omitempty"` // "Post", "User", etc.
	TargetID                    *int              `json:"target_id"`
	TargetURL                   string            `json:"target_url,omitempty"`
	TargetCreatedAt             *time.Time        `json:"target_created_at,omitempty"`
	TargetDeletedAt             *time.Time        `json:"target_deleted_at,omitempty"`
	TargetCreatedByTrustLevel   *int              `json:"target_created_by_trust_level,omitempty"`
	TopicTags                   []string          `json:"topic_tags,omitempty"`
	CategoryID                  *int              `json:"category_id"`
	CreatedAt                   time.Time         `json:"created_at"`
	CanEdit                     bool              `json:"can_edit"`
	Score                       float64           `json:"score"`
	Version                     int               `json:"version"`
	CreatedFromFlag             bool              `json:"created_from_flag?"`
	CreatedBy                   *BasicUser        `json:"created_by,omitempty"`
	ReviewableScores            []ReviewableScore `json:"reviewable_scores,omitempty"`
	Payload                     map[string]interface{} `json:"payload,omitempty"` // for queued posts: raw, title, etc.
}

// ReviewableListResponse wraps GET /review.json.
type ReviewableListResponse struct {
	Reviewables      []ReviewableItem `json:"reviewables"`
	Meta             map[string]interface{} `json:"meta,omitempty"`
	TotalRowsReviewables int             `json:"total_rows_reviewables,omitempty"`
}

// ============================================================================
// Themes (admin)
//
// Endpoint: GET /admin/themes.json
// Serializer: ThemeSerializer (extends BasicThemeSerializer), ThemeFieldSerializer
// ============================================================================

// ThemeField is one source file belonging to a theme.
type ThemeField struct {
	Name     string `json:"name"`
	Target   string `json:"target"`   // "common", "desktop", "mobile"
	Value    string `json:"value"`
	Error    string `json:"error,omitempty"`
	TypeID   int    `json:"type_id"`  // 0=html, 1=css, 2=scss, 3=js, 4=yaml
	UploadID *int   `json:"upload_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Filename string `json:"filename,omitempty"`
	Migrated bool   `json:"migrated"`
	FilePath string `json:"file_path,omitempty"`
}

// RemoteTheme describes the git remote a theme is synced from.
type RemoteTheme struct {
	ID                int        `json:"id"`
	RemoteURL         string     `json:"remote_url"`
	RemoteVersion     string     `json:"remote_version,omitempty"`
	LocalVersion      string     `json:"local_version,omitempty"`
	CommitsAhead      int        `json:"commits_behind"` // serialized as commits_behind
	Branch            string     `json:"branch,omitempty"`
	LastErrorText     string     `json:"last_error_text,omitempty"`
	IsGit             bool       `json:"is_git"`
	AuthorsString     string     `json:"authors,omitempty"`
	ThemeVersion      string     `json:"theme_version,omitempty"`
	MinDiscourseVersion string   `json:"minimum_discourse_version,omitempty"`
	MaxDiscourseVersion string   `json:"maximum_discourse_version,omitempty"`
	AboutURL          string     `json:"about_url,omitempty"`
	LicenseURL        string     `json:"license_url,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at,omitempty"`
}

// Theme is an admin-managed theme (or component).
type Theme struct {
	// From BasicThemeSerializer
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Default     bool      `json:"default,omitempty"`
	Component   bool      `json:"component"`

	// From ThemeSerializer
	ColorSchemeID     *int            `json:"color_scheme_id"`
	DarkColorSchemeID *int            `json:"dark_color_scheme_id"`
	UserSelectable    bool            `json:"user_selectable"`
	AutoUpdate        bool            `json:"auto_update"`
	RemoteThemeID     *int            `json:"remote_theme_id"`
	Settings          []interface{}   `json:"settings,omitempty"`
	Supported         bool            `json:"supported?"`
	Enabled           bool            `json:"enabled?"`
	DisabledAt        *time.Time      `json:"disabled_at,omitempty"`
	ThemeFields       []ThemeField    `json:"theme_fields,omitempty"`
	ChildThemes       []Theme         `json:"child_themes,omitempty"`
	ParentThemes      []Theme         `json:"parent_themes,omitempty"`
	RemoteTheme       *RemoteTheme    `json:"remote_theme,omitempty"`
	ColorScheme       *ColorScheme    `json:"color_scheme,omitempty"`
	Errors            []string        `json:"errors,omitempty"`
	ScreenshotDarkURL  string         `json:"screenshot_dark_url,omitempty"`
	ScreenshotLightURL string         `json:"screenshot_light_url,omitempty"`
	System            bool            `json:"system,omitempty"`
}

// ThemeListResponse wraps GET /admin/themes.json.
type ThemeListResponse struct {
	Themes []Theme `json:"themes"`
}

// ============================================================================
// Color Schemes (admin)
//
// Endpoint: GET /admin/color_schemes.json
// Serializer: ColorSchemeSerializer, ColorSchemeColorSerializer
// ============================================================================

// ColorSchemeColor is a named color entry inside a scheme.
type ColorSchemeColor struct {
	Name       string `json:"name"`
	Hex        string `json:"hex"`
	DefaultHex string `json:"default_hex,omitempty"`
	IsAdvanced bool   `json:"is_advanced,omitempty"`
}

// ColorScheme is a selectable color scheme for a theme.
type ColorScheme struct {
	ID               int                `json:"id"`
	Name             string             `json:"name"`
	IsBase           bool               `json:"is_base"`
	BaseSchemeID     string             `json:"base_scheme_id,omitempty"`
	ThemeID          *int               `json:"theme_id,omitempty"`
	ThemeName        string             `json:"theme_name,omitempty"`
	UserSelectable   bool               `json:"user_selectable"`
	IsBuiltinDefault bool               `json:"is_builtin_default,omitempty"`
	Colors           []ColorSchemeColor `json:"colors"`
}

// ============================================================================
// Custom User Fields (admin)
//
// Endpoint: GET /admin/customize/user_fields.json
// Serializer: UserFieldSerializer
// ============================================================================

// CustomUserField describes an admin-created custom profile field.
type CustomUserField struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	FieldType      string   `json:"field_type"` // "text", "confirm", "dropdown", "multiselect"
	Editable       bool     `json:"editable"`
	Required       bool     `json:"required"`
	Requirement    string   `json:"requirement,omitempty"` // "for_all_users", "on_signup", etc.
	ShowOnProfile  bool     `json:"show_on_profile"`
	ShowOnUserCard bool     `json:"show_on_user_card"`
	ShowOnSignup   bool     `json:"show_on_signup,omitempty"`
	Searchable     bool     `json:"searchable"`
	Position       int      `json:"position"`
	Options        []string `json:"options,omitempty"` // for dropdown/multiselect types
}

// ============================================================================
// Tag Groups (admin)
//
// Endpoint: GET /tag_groups.json
// Serializer: TagGroupSerializer
// ============================================================================

// TagGroupTag is a compact tag representation inside a tag group.
type TagGroupTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

// TagGroupPermission is a group permission entry in a tag group.
type TagGroupPermission struct {
	GroupName  string `json:"group_name"`
	PermissionType int `json:"permission_type"`
}

// TagGroup is a named collection of tags with optional restrictions.
type TagGroup struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Tags        []TagGroupTag        `json:"tags,omitempty"`
	ParentTag   *TagGroupTag         `json:"parent_tag,omitempty"`
	OnePerTopic bool                 `json:"one_per_topic"`
	Permissions []TagGroupPermission `json:"permissions,omitempty"`
}

// TagGroupListResponse wraps GET /tag_groups.json.
type TagGroupListResponse struct {
	TagGroups []TagGroup `json:"tag_groups"`
}

// ============================================================================
// Drafts
//
// Endpoint: GET /drafts.json
// Serializer: DraftSerializer
// ============================================================================

// Draft is an auto-saved composer draft.
type Draft struct {
	CreatedAt      time.Time `json:"created_at"`
	DraftKey       string    `json:"draft_key"`       // e.g. "new_topic", "topic_12345"
	Sequence       int       `json:"sequence"`
	DraftUsername   string   `json:"draft_username"`
	AvatarTemplate string    `json:"avatar_template"`
	Data           string    `json:"data"`             // JSON-encoded draft payload
	TopicID        *int      `json:"topic_id,omitempty"`
	Username       string    `json:"username"`
	UsernameLower  string    `json:"username_lower,omitempty"`
	Name           string    `json:"name,omitempty"`
	UserID         int       `json:"user_id"`
	Title          string    `json:"title,omitempty"`
	Slug           string    `json:"slug,omitempty"`
	CategoryID     *int      `json:"category_id,omitempty"`
	Closed         bool      `json:"closed"`
	Archetype      string    `json:"archetype,omitempty"`
	Archived       bool      `json:"archived"`
	Excerpt        string    `json:"excerpt,omitempty"`
}

// DraftListResponse wraps GET /drafts.json.
type DraftListResponse struct {
	Drafts []Draft `json:"drafts"`
}

// ============================================================================
// Bookmarks
//
// Endpoint: GET /u/{username}/bookmarks.json
// Serializer: UserBookmarkBaseSerializer
// ============================================================================

// Bookmark is a user-saved bookmark (polymorphic: can target posts, topics, etc.).
type Bookmark struct {
	ID               int        `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Name             string     `json:"name,omitempty"`
	ReminderAt       *time.Time `json:"reminder_at,omitempty"`
	ReminderAtICSStart string   `json:"reminder_at_ics_start,omitempty"`
	ReminderAtICSEnd   string   `json:"reminder_at_ics_end,omitempty"`
	Pinned           bool       `json:"pinned"`
	Title            string     `json:"title,omitempty"`
	FancyTitle       string     `json:"fancy_title,omitempty"`
	Excerpt          string     `json:"excerpt,omitempty"`
	BookmarkableID   int        `json:"bookmarkable_id"`
	BookmarkableType string     `json:"bookmarkable_type"` // "Post", "Topic", etc.
	BookmarkableURL  string     `json:"bookmarkable_url"`
	User             *BasicUser `json:"user,omitempty"`
}

// BookmarkListResponse wraps GET /u/{username}/bookmarks.json.
type BookmarkListResponse struct {
	Bookmarks       []Bookmark `json:"user_bookmark_list"`
	MoreBookmarksURL string    `json:"more_bookmarks_url,omitempty"`
}

// ============================================================================
// Watched Words (admin)
//
// Endpoint: GET /admin/customize/watched_words.json
// Model: WatchedWord
// ============================================================================

// WatchedWord is a word or pattern monitored by Discourse's word filter.
type WatchedWord struct {
	ID               int        `json:"id"`
	Word             string     `json:"word"`
	Action           int        `json:"action"` // 0=block, 1=censor, 2=require_approval, 3=flag, 4=replace, 5=tag, 6=silence, 7=link
	Replacement      string     `json:"replacement,omitempty"`
	CaseSensitive    bool       `json:"case_sensitive"`
	HTML             bool       `json:"html,omitempty"`
	WatchedWordGroupID *int     `json:"watched_word_group_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// WatchedWordListResponse wraps GET /admin/customize/watched_words.json.
// The response is keyed by action name ("block", "censor", etc.).
type WatchedWordListResponse struct {
	Actions  []WatchedWordAction `json:"actions"`
	Words    []WatchedWord       `json:"words,omitempty"`
	Compiled bool                `json:"compiled_regular_expressions,omitempty"`
}

// WatchedWordAction groups words by their action type.
type WatchedWordAction struct {
	ID    int            `json:"id"`
	Name  string         `json:"name"`
	Words []WatchedWord  `json:"words"`
}

// ============================================================================
// Permalinks (admin)
//
// Endpoint: GET /admin/permalinks.json
// Serializer: PermalinkSerializer
// ============================================================================

// Permalink maps an old URL to a new resource.
type Permalink struct {
	ID             int    `json:"id"`
	URL            string `json:"url"`
	TopicID        *int   `json:"topic_id,omitempty"`
	TopicTitle     string `json:"topic_title,omitempty"`
	TopicURL       string `json:"topic_url,omitempty"`
	PostID         *int   `json:"post_id,omitempty"`
	PostURL        string `json:"post_url,omitempty"`
	PostNumber     *int   `json:"post_number,omitempty"`
	PostTopicTitle string `json:"post_topic_title,omitempty"`
	CategoryID     *int   `json:"category_id,omitempty"`
	CategoryName   string `json:"category_name,omitempty"`
	CategoryURL    string `json:"category_url,omitempty"`
	ExternalURL    string `json:"external_url,omitempty"`
	TagID          *int   `json:"tag_id,omitempty"`
	TagName        string `json:"tag_name,omitempty"`
	TagURL         string `json:"tag_url,omitempty"`
	UserID         *int   `json:"user_id,omitempty"`
	UserURL        string `json:"user_url,omitempty"`
	Username       string `json:"username,omitempty"`
}

// PermalinkListResponse wraps GET /admin/permalinks.json.
type PermalinkListResponse struct {
	Permalinks []Permalink `json:"permalinks"`
}

// ============================================================================
// Staff Action Logs (admin)
//
// Endpoint: GET /admin/logs/staff_action_logs.json
// Serializer: UserHistorySerializer
// ============================================================================

// StaffActionLog records an action taken by a staff member.
type StaffActionLog struct {
	ID            int        `json:"id"`
	ActionName    string     `json:"action_name"`
	Details       string     `json:"details,omitempty"`
	Context       string     `json:"context,omitempty"`
	IPAddress     string     `json:"ip_address,omitempty"`
	Email         string     `json:"email,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	Subject       string     `json:"subject,omitempty"`
	PreviousValue string     `json:"previous_value,omitempty"`
	NewValue      string     `json:"new_value,omitempty"`
	TopicID       *int       `json:"topic_id,omitempty"`
	PostID        *int       `json:"post_id,omitempty"`
	CategoryID    *int       `json:"category_id,omitempty"`
	Action        int        `json:"action"`
	CustomType    string     `json:"custom_type,omitempty"`
	ActingUser    *BasicUser `json:"acting_user,omitempty"`
	TargetUser    *BasicUser `json:"target_user,omitempty"`
}

// StaffActionLogListResponse wraps GET /admin/logs/staff_action_logs.json.
type StaffActionLogListResponse struct {
	StaffActionLogs          []StaffActionLog `json:"staff_action_logs"`
	UserHistoryActionsFilter []string         `json:"user_history_actions,omitempty"`
}

// ============================================================================
// Screened Emails / IPs / URLs (admin)
//
// Endpoints:
//   GET /admin/logs/screened_emails.json
//   GET /admin/logs/screened_ip_addresses.json
//   GET /admin/logs/screened_urls.json
//
// Serializers: ScreenedEmailSerializer, ScreenedIpAddressSerializer,
//              ScreenedUrlSerializer
// ============================================================================

// ScreenedEmail is a blocked or watched email pattern.
type ScreenedEmail struct {
	ID          int        `json:"id"`
	Email       string     `json:"email"`
	Action      string     `json:"action"`      // "block", "do_nothing"
	MatchCount  int        `json:"match_count"`
	LastMatchAt *time.Time `json:"last_match_at"`
	CreatedAt   time.Time  `json:"created_at"`
	IPAddress   string     `json:"ip_address,omitempty"`
}

// ScreenedIPAddress is a blocked, allowed, or watched IP range.
type ScreenedIPAddress struct {
	ID          int        `json:"id"`
	IPAddress   string     `json:"ip_address"` // masked CIDR notation
	ActionName  string     `json:"action_name"` // "block", "do_nothing", "allow_admin"
	MatchCount  int        `json:"match_count"`
	LastMatchAt *time.Time `json:"last_match_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ScreenedURL is a blocked URL pattern.
type ScreenedURL struct {
	URL         string     `json:"url"`
	Domain      string     `json:"domain"`
	Action      string     `json:"action"` // "do_nothing", "block"
	MatchCount  int        `json:"match_count"`
	LastMatchAt *time.Time `json:"last_match_at"`
	CreatedAt   time.Time  `json:"created_at"`
	IPAddress   string     `json:"ip_address,omitempty"`
}

// ============================================================================
// Embeddable Hosts (admin)
//
// Endpoint: GET /admin/customize/embedding.json
// Serializer: EmbeddableHostSerializer
// ============================================================================

// EmbeddableHost allows a remote site to embed Discourse comments.
type EmbeddableHost struct {
	ID           int    `json:"id"`
	Host         string `json:"host"`
	AllowedPaths string `json:"allowed_paths,omitempty"`
	ClassName    string `json:"class_name,omitempty"`
	CategoryID   *int   `json:"category_id,omitempty"`
	Tags         []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug,omitempty"`
	} `json:"tags,omitempty"`
	User string `json:"user,omitempty"` // username of the post author
}

// ============================================================================
// Site Texts (admin)
//
// Endpoint: GET /admin/customize/site_texts/{id}.json
// Serializer: SiteTextSerializer
// ============================================================================

// SiteText is a customisable locale string.
type SiteText struct {
	ID                  string `json:"id"` // e.g. "title", "js.topic.create"
	Value               string `json:"value"`
	Status              string `json:"status,omitempty"`
	OldDefault          string `json:"old_default,omitempty"`
	NewDefault          string `json:"new_default,omitempty"`
	InterpolationKeys   string `json:"interpolation_keys,omitempty"`
	HasInterpolationKeys bool  `json:"has_interpolation_keys?,omitempty"`
	Overridden          bool   `json:"overridden?"`
	CanRevert           bool   `json:"can_revert?"`
}

// ============================================================================
// Sidebar Sections
//
// Endpoint: GET /sidebar_sections.json
// Serializer: SidebarSectionSerializer, SidebarUrlSerializer
// ============================================================================

// SidebarLink is one link inside a sidebar section.
type SidebarLink struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Icon     string `json:"icon"`
	External bool   `json:"external"`
	Segment  string `json:"segment,omitempty"`
}

// SidebarSection is a group of links in the sidebar navigation.
type SidebarSection struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Links       []SidebarLink `json:"links"`
	Slug        string        `json:"slug"`
	Public      bool          `json:"public"`
	SectionType string        `json:"section_type,omitempty"`
}

// ============================================================================
// Published Pages
//
// Endpoint: GET /pub/by-topic/{topic_id}.json
// Serializer: PublishedPageSerializer
// ============================================================================

// PublishedPage maps a topic to a publicly readable page at /pub/{slug}.
type PublishedPage struct {
	ID     int    `json:"id"` // same as topic_id
	Slug   string `json:"slug"`
	Public bool   `json:"public"`
}

// ============================================================================
// Custom Emojis (admin)
//
// Endpoint: GET /admin/customize/emojis.json
// Model attrs: name, url, group
// ============================================================================

// CustomEmoji is an admin-uploaded custom emoji.
type CustomEmoji struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Group     string `json:"group,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
}

// ============================================================================
// Form Templates (admin)
//
// Endpoint: GET /admin/customize/form-templates.json
// Model: FormTemplate
// ============================================================================

// FormTemplate is a YAML-based form definition that can be attached to a category.
type FormTemplate struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Template  string    `json:"template"` // YAML content
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ============================================================================
// Admin Flags
//
// Endpoint: GET /admin/config/flags.json
// Serializer: FlagSerializer
// ============================================================================

// AdminFlag describes a flag type available for use on posts or topics.
type AdminFlag struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	NameKey          string   `json:"name_key"`
	Description      string   `json:"description"`
	ShortDescription string   `json:"short_description,omitempty"`
	AppliesTo        []string `json:"applies_to,omitempty"` // ["Post"], ["Topic"], ["Post","Topic"]
	Position         int      `json:"position"`
	RequireMessage   bool     `json:"require_message"`
	Enabled          bool     `json:"enabled"`
	IsFlag           bool     `json:"is_flag"`
	IsUsed           bool     `json:"is_used"`
	AutoActionType   *int     `json:"auto_action_type"`
	System           bool     `json:"system"`
}

// ============================================================================
// Extended Directory Items
//
// Endpoint: GET /directory_items.json?period=all
// Serializer: DirectoryItemSerializer
//
// The base DirectoryItem is defined in models.go; this type adds the full set
// of statistics that the serializer emits dynamically.
// ============================================================================

// DirectoryItemExtended adds the statistics that DirectoryItemSerializer
// includes beyond the base ID + User.
type DirectoryItemExtended struct {
	ID               int  `json:"id"` // actually user_id
	User             User `json:"user"`
	LikesReceived    int  `json:"likes_received"`
	LikesGiven       int  `json:"likes_given"`
	TopicsEntered    int  `json:"topics_entered"`
	TopicCount       int  `json:"topic_count"`
	PostCount        int  `json:"post_count"`
	PostsRead        int  `json:"posts_read"`
	DaysVisited      int  `json:"days_visited"`
	TimeRead         int  `json:"time_read,omitempty"` // seconds; only for period=all
	UserFields       map[string]interface{} `json:"user_fields,omitempty"`
}

// ============================================================================
// User Status
//
// Embedded inside user JSON.  Also: PUT /u/{username}/status.json
// Serializer: UserStatusSerializer
// ============================================================================

// UserStatus is the emoji + description status a user can set.
type UserStatus struct {
	Description     string     `json:"description"`
	Emoji           string     `json:"emoji"`
	EndsAt          *time.Time `json:"ends_at"`
	MessageBusLastID int       `json:"message_bus_last_id,omitempty"`
}

// ============================================================================
// Post Revisions
//
// Endpoint: GET /posts/{post_id}/revisions/{revision}.json
// Serializer: PostRevisionSerializer
// ============================================================================

// RevisionDiff holds "before" and "after" for a changed field.
type RevisionDiff struct {
	Inline   string `json:"inline,omitempty"`
	SideBySide string `json:"side_by_side,omitempty"`
	SideBySideMarkdown string `json:"side_by_side_markdown,omitempty"`
}

// PostRevision is a single historical edit of a post.
type PostRevision struct {
	CreatedAt        time.Time     `json:"created_at"`
	PostID           int           `json:"post_id"`
	PreviousHidden   bool          `json:"previous_hidden"`
	CurrentHidden    bool          `json:"current_hidden"`
	FirstRevision    int           `json:"first_revision"`
	PreviousRevision *int          `json:"previous_revision"`
	CurrentRevision  int           `json:"current_revision"`
	NextRevision     *int          `json:"next_revision"`
	LastRevision     int           `json:"last_revision"`
	CurrentVersion   int           `json:"current_version"`
	VersionCount     int           `json:"version_count"`
	Username         string        `json:"username"`
	DisplayUsername  string        `json:"display_username"`
	ActingUserName   string        `json:"acting_user_name,omitempty"`
	AvatarTemplate   string        `json:"avatar_template"`
	EditReason       string        `json:"edit_reason,omitempty"`
	BodyChanges      *RevisionDiff `json:"body_changes,omitempty"`
	TitleChanges     *RevisionDiff `json:"title_changes,omitempty"`
	UserChanges      interface{}   `json:"user_changes,omitempty"`
	TagsChanges      interface{}   `json:"tags_changes,omitempty"`
	CategoryIDChanges interface{}  `json:"category_id_changes,omitempty"`
	WikiChanges      interface{}   `json:"wiki_changes,omitempty"`
	PostTypeChanges  interface{}   `json:"post_type_changes,omitempty"`
	LocaleChanges    interface{}   `json:"locale_changes,omitempty"`
	CanEdit          bool          `json:"can_edit"`
}

// ============================================================================
// User Summary
//
// Endpoint: GET /u/{username}/summary.json
// Serializer: UserSummarySerializer
// ============================================================================

// UserSummaryTopic is a compact topic reference in a user's summary.
type UserSummaryTopic struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	FancyTitle  string `json:"fancy_title,omitempty"`
	Slug        string `json:"slug"`
	LikeCount   int    `json:"like_count,omitempty"`
	PostsCount  int    `json:"posts_count,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	CategoryID  int    `json:"category_id,omitempty"`
}

// UserSummaryReply is a reply entry in a user's summary.
type UserSummaryReply struct {
	PostNumber int    `json:"post_number"`
	LikeCount  int    `json:"like_count"`
	TopicID    int    `json:"topic_id"`
	CreatedAt  string `json:"created_at,omitempty"`
}

// UserSummaryLink is a link shared by a user.
type UserSummaryLink struct {
	URL        string `json:"url"`
	Title      string `json:"title,omitempty"`
	Clicks     int    `json:"clicks"`
	TopicID    int    `json:"topic_id,omitempty"`
	PostNumber int    `json:"post_number,omitempty"`
}

// UserWithCount is a user reference accompanied by a count (likes, replies).
type UserWithCount struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name,omitempty"`
	AvatarTemplate string `json:"avatar_template"`
	Count          int    `json:"count"`
}

// CategoryWithCounts is a category accompanied by topic/post counts in a user summary.
type CategoryWithCounts struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Color      string `json:"color"`
	TopicCount int    `json:"topic_count"`
	PostCount  int    `json:"post_count"`
}

// UserSummary is the aggregated activity overview of a user.
type UserSummary struct {
	LikesGiven           int                  `json:"likes_given"`
	LikesReceived        int                  `json:"likes_received"`
	TopicsEntered        int                  `json:"topics_entered"`
	PostsReadCount       int                  `json:"posts_read_count"`
	DaysVisited          int                  `json:"days_visited"`
	TopicCount           int                  `json:"topic_count"`
	PostCount            int                  `json:"post_count"`
	TimeRead             int                  `json:"time_read"`        // seconds
	RecentTimeRead       int                  `json:"recent_time_read"` // seconds
	BookmarkCount        int                  `json:"bookmark_count"`
	CanSeeSummaryStats   bool                 `json:"can_see_summary_stats"`
	CanSeeUserActions    bool                 `json:"can_see_user_actions"`
	Topics               []UserSummaryTopic   `json:"topics,omitempty"`
	Replies              []UserSummaryReply   `json:"replies,omitempty"`
	Links                []UserSummaryLink    `json:"links,omitempty"`
	MostLikedByUsers     []UserWithCount      `json:"most_liked_by_users,omitempty"`
	MostLikedUsers       []UserWithCount      `json:"most_liked_users,omitempty"`
	MostRepliedToUsers   []UserWithCount      `json:"most_replied_to_users,omitempty"`
	Badges               []UserBadge          `json:"badges,omitempty"`
	TopCategories        []CategoryWithCounts `json:"top_categories,omitempty"`
}

// UserSummaryResponse wraps GET /u/{username}/summary.json.
type UserSummaryResponse struct {
	UserSummary UserSummary `json:"user_summary"`
}
