// Package model defines data types matching Discourse's API response shapes.
// Field names and JSON tags are chosen for 1:1 compatibility with the official
// Discourse REST API so that SDK clients (discourse_api Ruby gem, pydiscourse,
// discourse-api JS) can deserialise responses without modification.
package model

import "time"

// ---------- Users ----------

type User struct {
	ID               int       `json:"id"`
	Username         string    `json:"username"`
	Name             string    `json:"name"`
	Email            string    `json:"email,omitempty"`
	AvatarTemplate   string    `json:"avatar_template"`
	Active           bool      `json:"active"`
	Admin            bool      `json:"admin"`
	Moderator        bool      `json:"moderator"`
	TrustLevel       int       `json:"trust_level"`
	CreatedAt        time.Time `json:"created_at"`
	LastSeenAt       time.Time `json:"last_seen_at,omitempty"`
	Approved         bool      `json:"approved"`
	Suspended        bool      `json:"suspended,omitempty"`
	SuspendedTill    string    `json:"suspended_till,omitempty"`
	Silenced         bool      `json:"silenced,omitempty"`
	Title            string    `json:"title,omitempty"`
	ExternalID       string    `json:"external_id,omitempty"`
	GroupIDs         []int     `json:"group_ids,omitempty"`
	PrimaryGroupID   *int      `json:"primary_group_id"`
	FlairGroupID     *int      `json:"flair_group_id"`
	UserFields       map[string]interface{} `json:"user_fields,omitempty"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UserListResponse struct {
	Users      []User `json:"users,omitempty"`
	DirectoryItems []DirectoryItem `json:"directory_items,omitempty"`
	TotalRowsDirectoryItems int `json:"total_rows_directory_items,omitempty"`
	LoadMoreURL string `json:"load_more_url,omitempty"`
}

type DirectoryItem struct {
	ID   int  `json:"id"`
	User User `json:"user"`
}

type CreateUserResponse struct {
	Success bool   `json:"success"`
	Active  bool   `json:"active"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

// ---------- Categories ----------

type Category struct {
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	Slug                  string    `json:"slug"`
	Color                 string    `json:"color"`
	TextColor             string    `json:"text_color"`
	Description           string    `json:"description"`
	DescriptionText       string    `json:"description_text"`
	TopicCount            int       `json:"topic_count"`
	PostCount             int       `json:"post_count"`
	Position              int       `json:"position"`
	ParentCategoryID      *int      `json:"parent_category_id"`
	Subcategories         []Category `json:"subcategory_list,omitempty"`
	TopicURL              string    `json:"topic_url"`
	ReadRestricted        bool      `json:"read_restricted"`
	Permission            int       `json:"permission,omitempty"`
	NotificationLevel     *int      `json:"notification_level"`
	CanEdit               bool      `json:"can_edit"`
	TopicTemplate         string    `json:"topic_template"`
	HasChildren           bool      `json:"has_children"`
	NumFeaturedTopics     int       `json:"num_featured_topics"`
	ShowSubcategoryList   bool      `json:"show_subcategory_list"`
	DefaultView           string    `json:"default_view"`
	SubcategoryListStyle  string    `json:"subcategory_list_style"`
	DefaultTopPeriod      string    `json:"default_top_period"`
	MinimumRequiredTags   int       `json:"minimum_required_tags"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
}

type CategoryListResponse struct {
	CategoryList CategoryList `json:"category_list"`
}

type CategoryList struct {
	CanCreateCategory bool       `json:"can_create_category"`
	CanCreateTopic    bool       `json:"can_create_topic"`
	Categories        []Category `json:"categories"`
}

type CategoryResponse struct {
	Category Category `json:"category"`
}

// ---------- Topics ----------

type Topic struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	FancyTitle        string    `json:"fancy_title"`
	Slug              string    `json:"slug"`
	PostsCount        int       `json:"posts_count"`
	ReplyCount        int       `json:"reply_count"`
	HighestPostNumber int       `json:"highest_post_number"`
	CreatedAt         time.Time `json:"created_at"`
	LastPostedAt      time.Time `json:"last_posted_at"`
	Bumped            bool      `json:"bumped"`
	BumpedAt          time.Time `json:"bumped_at"`
	Archetype         string    `json:"archetype"`
	Unseen            bool      `json:"unseen"`
	Pinned            bool      `json:"pinned"`
	Unpinned          *bool     `json:"unpinned"`
	Visible           bool      `json:"visible"`
	Closed            bool      `json:"closed"`
	Archived          bool      `json:"archived"`
	Bookmarked        *bool     `json:"bookmarked"`
	Liked             *bool     `json:"liked"`
	Views             int       `json:"views"`
	LikeCount         int       `json:"like_count"`
	HasSummary        bool      `json:"has_summary"`
	LastPosterUsername string   `json:"last_poster_username"`
	CategoryID        int       `json:"category_id"`
	PinnedGlobally    bool      `json:"pinned_globally"`
	HasAcceptedAnswer bool      `json:"has_accepted_answer"`
	Posters           []Poster  `json:"posters,omitempty"`
	Tags              []string  `json:"tags"`
	ExternalID        string    `json:"external_id,omitempty"`

	// Included when fetching a single topic
	PostStream *PostStream `json:"post_stream,omitempty"`
	Details    *TopicDetails `json:"details,omitempty"`
}

type Poster struct {
	Extras      string `json:"extras,omitempty"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

type TopicDetails struct {
	AutoCloseAt       *time.Time `json:"auto_close_at"`
	CreatedBy         BasicUser  `json:"created_by"`
	LastPoster        BasicUser  `json:"last_poster"`
	Participants      []BasicUser `json:"participants,omitempty"`
	CanEdit           bool       `json:"can_edit"`
	CanInviteTo       bool       `json:"can_invite_to"`
	CanCreatePost     bool       `json:"can_create_post"`
	CanReplyAsNewTopic bool      `json:"can_reply_as_new_topic"`
	CanFlagTopic      bool       `json:"can_flag_topic"`
	NotificationLevel int        `json:"notification_level"`
}

type BasicUser struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name,omitempty"`
	AvatarTemplate string `json:"avatar_template"`
}

type PostStream struct {
	Posts  []Post `json:"posts"`
	Stream []int  `json:"stream"`
}

type TopicListResponse struct {
	Users    []BasicUser `json:"users,omitempty"`
	TopicList TopicList  `json:"topic_list"`
}

type TopicList struct {
	CanCreateTopic bool    `json:"can_create_topic"`
	MoreTopicsURL  string  `json:"more_topics_url,omitempty"`
	PerPage        int     `json:"per_page"`
	TopTags        []string `json:"top_tags,omitempty"`
	Topics         []Topic `json:"topics"`
}

type TopicPostsResponse struct {
	PostStream PostStream `json:"post_stream"`
	ID         int        `json:"id"`
}

// ---------- Posts ----------

type Post struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	AvatarTemplate    string    `json:"avatar_template"`
	CreatedAt         time.Time `json:"created_at"`
	Cooked            string    `json:"cooked"`
	Raw               string    `json:"raw,omitempty"`
	PostNumber        int       `json:"post_number"`
	PostType          int       `json:"post_type"`
	UpdatedAt         time.Time `json:"updated_at"`
	ReplyCount        int       `json:"reply_count"`
	ReplyToPostNumber *int      `json:"reply_to_post_number"`
	QuoteCount        int       `json:"quote_count"`
	AvgTime           *int      `json:"avg_time"`
	Score             float64   `json:"score"`
	Reads             int       `json:"reads"`
	TopicID           int       `json:"topic_id"`
	TopicSlug         string    `json:"topic_slug"`
	TopicTitle        string    `json:"topic_title,omitempty"`
	CategoryID        int       `json:"category_id,omitempty"`
	DisplayUsername    string   `json:"display_username"`
	Version           int       `json:"version"`
	Wiki              bool      `json:"wiki"`
	CanEdit           bool      `json:"can_edit"`
	CanDelete         bool      `json:"can_delete"`
	CanRecover        bool      `json:"can_recover"`
	CanWiki           bool      `json:"can_wiki"`
	UserID            int       `json:"user_id"`
	Hidden            bool      `json:"hidden"`
	TrustLevel        int       `json:"trust_level"`
	Yours             bool      `json:"yours"`
}

type PostResponse struct {
	Post
}

type PostListResponse struct {
	LatestPosts []Post `json:"latest_posts"`
}

// ---------- Groups ----------

type Group struct {
	ID                int       `json:"id"`
	Automatic         bool      `json:"automatic"`
	Name              string    `json:"name"`
	DisplayName       string    `json:"display_name"`
	UserCount         int       `json:"user_count"`
	MentionableLevel  int       `json:"mentionable_level"`
	MessageableLevel  int       `json:"messageable_level"`
	VisibilityLevel   int       `json:"visibility_level"`
	PrimaryGroup      bool      `json:"primary_group"`
	Title             string    `json:"title"`
	GrantTrustLevel   *int      `json:"grant_trust_level"`
	FlairURL          string    `json:"flair_url"`
	FlairBgColor      string    `json:"flair_bg_color"`
	FlairColor        string    `json:"flair_color"`
	BioRaw            string    `json:"bio_raw,omitempty"`
	BioCooked         string    `json:"bio_cooked,omitempty"`
	BioExcerpt        string    `json:"bio_excerpt,omitempty"`
	PublicAdmission   bool      `json:"public_admission"`
	PublicExit        bool      `json:"public_exit"`
	AllowMembershipRequests bool `json:"allow_membership_requests"`
	FullName          string    `json:"full_name"`
	DefaultNotificationLevel int `json:"default_notification_level"`
	MembershipRequestTemplate string `json:"membership_request_template"`
	MemberIDs         []int     `json:"members,omitempty"`
	OwnerIDs          []int     `json:"owners,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

type GroupResponse struct {
	Group Group `json:"group"`
}

type GroupListResponse struct {
	Groups         []Group `json:"groups"`
	TotalRowsGroups int    `json:"total_rows_groups"`
	LoadMoreGroups string  `json:"load_more_groups,omitempty"`
}

type GroupMembersResponse struct {
	Members []BasicUser `json:"members"`
	Owners  []BasicUser `json:"owners"`
	Meta    GroupMembersMeta `json:"meta"`
}

type GroupMembersMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ---------- Search ----------

type SearchResult struct {
	Posts            []Post    `json:"posts"`
	Topics           []Topic   `json:"topics"`
	Users            []BasicUser `json:"users,omitempty"`
	Categories       []Category `json:"categories,omitempty"`
	Tags             []Tag      `json:"tags,omitempty"`
	GroupedSearchResult *GroupedSearchResult `json:"grouped_search_result,omitempty"`
}

type GroupedSearchResult struct {
	MorePosts      *bool `json:"more_posts"`
	MoreUsers      *bool `json:"more_users"`
	MoreCategories *bool `json:"more_categories"`
	PostIDs        []int `json:"post_ids"`
	UserIDs        []int `json:"user_ids"`
	CategoryIDs    []int `json:"category_ids"`
}

// ---------- Tags ----------

type Tag struct {
	ID       int    `json:"id"`
	Name     string `json:"text,omitempty"`
	TagName  string `json:"name,omitempty"`
	Count    int    `json:"count"`
	PMCount  int    `json:"pm_count,omitempty"`
	TargetTag string `json:"target_tag,omitempty"`
}

type TagResponse struct {
	Tag    Tag     `json:"tag,omitempty"`
	Users  []BasicUser `json:"users,omitempty"`
	TopicList TopicList `json:"topic_list,omitempty"`
}

type TagListResponse struct {
	Tags   []Tag  `json:"tags"`
	Extras interface{} `json:"extras,omitempty"`
}

// ---------- Badges ----------

type Badge struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GrantCount  int    `json:"grant_count"`
	AllowTitle  bool   `json:"allow_title"`
	MultipleGrant bool `json:"multiple_grant"`
	Icon        string `json:"icon"`
	Listable    bool   `json:"listable"`
	Enabled     bool   `json:"enabled"`
	BadgeGroupingID int `json:"badge_grouping_id"`
	System      bool   `json:"system"`
	Slug        string `json:"slug,omitempty"`
	BadgeTypeID int    `json:"badge_type_id"`
}

type BadgeListResponse struct {
	Badges []Badge `json:"badges"`
}

type BadgeResponse struct {
	Badge Badge `json:"badge,omitempty"`
	BadgeType BadgeType `json:"badge_type,omitempty"`
}

type BadgeType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserBadge struct {
	ID        int       `json:"id"`
	GrantedAt time.Time `json:"granted_at"`
	BadgeID   int       `json:"badge_id"`
	UserID    int       `json:"user_id"`
	GrantedByID int     `json:"granted_by_id"`
}

type UserBadgeResponse struct {
	Badges      []Badge     `json:"badges"`
	UserBadges  []UserBadge `json:"user_badges"`
}

// ---------- Notifications ----------

type Notification struct {
	ID               int       `json:"id"`
	NotificationType int       `json:"notification_type"`
	Read             bool      `json:"read"`
	CreatedAt        time.Time `json:"created_at"`
	PostNumber       *int      `json:"post_number"`
	TopicID          *int      `json:"topic_id"`
	Slug             string    `json:"slug"`
	Data             NotificationData `json:"data"`
}

type NotificationData struct {
	BadgeID          int    `json:"badge_id,omitempty"`
	BadgeName        string `json:"badge_name,omitempty"`
	BadgeSlug        string `json:"badge_slug,omitempty"`
	TopicTitle       string `json:"topic_title,omitempty"`
	OriginalPostID   int    `json:"original_post_id,omitempty"`
	OriginalPostType int    `json:"original_post_type,omitempty"`
	OriginalUsername string `json:"original_username,omitempty"`
	DisplayUsername   string `json:"display_username,omitempty"`
}

type NotificationListResponse struct {
	Notifications     []Notification `json:"notifications"`
	TotalRowsNotifications int       `json:"total_rows_notifications"`
	SeenNotificationID int           `json:"seen_notification_id"`
	LoadMoreNotifications string     `json:"load_more_notifications,omitempty"`
}

// ---------- Invites ----------

type Invite struct {
	ID          int       `json:"id"`
	Link        string    `json:"link"`
	Email       string    `json:"email,omitempty"`
	MaxRedemptionsAllowed int `json:"max_redemptions_allowed"`
	RedemptionCount int   `json:"redemption_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Expired     bool      `json:"expired"`
	Topics      []Topic   `json:"topics,omitempty"`
	Groups      []Group   `json:"groups,omitempty"`
}

type InviteResponse struct {
	Invite
}

// ---------- Uploads ----------

type Upload struct {
	ID               int    `json:"id"`
	URL              string `json:"url"`
	OriginalFilename string `json:"original_filename"`
	Filesize         int    `json:"filesize"`
	Width            int    `json:"width,omitempty"`
	Height           int    `json:"height,omitempty"`
	ThumbnailWidth   int    `json:"thumbnail_width,omitempty"`
	ThumbnailHeight  int    `json:"thumbnail_height,omitempty"`
	Extension        string `json:"extension"`
	ShortURL         string `json:"short_url"`
	ShortPath        string `json:"short_path"`
	HumanFilesize    string `json:"human_filesize"`
}

// ---------- Site / Settings ----------

type SiteSetting struct {
	Setting   string      `json:"setting"`
	Value     interface{} `json:"value"`
	Default   interface{} `json:"default,omitempty"`
}

type SiteSettingsResponse struct {
	SiteSettings []SiteSetting `json:"site_settings"`
}

type SiteInfo struct {
	DefaultArchetype string     `json:"default_archetype"`
	NotificationTypes map[string]int `json:"notification_types"`
	PostTypes         map[string]int `json:"post_types"`
	TrustLevels       map[string]int `json:"trust_levels"`
	Groups            []Group    `json:"groups"`
	Categories        []Category `json:"categories"`
}

// ---------- Private Messages ----------

type PrivateMessageListResponse struct {
	Users    []BasicUser `json:"users,omitempty"`
	TopicList TopicList  `json:"topic_list"`
}

// ---------- Backups ----------

type Backup struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Link     string `json:"link,omitempty"`
}

type BackupListResponse []Backup

// ---------- Post Actions ----------

type PostAction struct {
	ID         int `json:"id"`
	PostID     int `json:"post_id"`
	PostActionTypeID int `json:"post_action_type_id"`
}

type PostActionResponse struct {
	PostAction
}

// ---------- SSO ----------

type SSORecord struct {
	ExternalID   string `json:"external_id"`
	ExternalUsername string `json:"external_username,omitempty"`
	ExternalEmail string `json:"external_email,omitempty"`
	ExternalName  string `json:"external_name,omitempty"`
	ExternalAvatarURL string `json:"external_avatar_url,omitempty"`
}

// ---------- Generic / Shared ----------

type SuccessResponse struct {
	Success string `json:"success,omitempty"`
}

type ErrorResponse struct {
	Errors    []string `json:"errors,omitempty"`
	ErrorType string   `json:"error_type,omitempty"`
}

type StatusResponse struct {
	Status string `json:"status,omitempty"`
	Topic  *Topic `json:"topic,omitempty"`
}
