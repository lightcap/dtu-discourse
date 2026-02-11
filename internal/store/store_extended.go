package store

import (
	"fmt"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Extended model types -- resources tracked by ExtStore that are not yet
// defined in the model package.  They live here so the extended store is
// self-contained and doesn't require changes to model/models.go.
// ---------------------------------------------------------------------------

type Poll struct {
	ID        int                    `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"` // regular, multiple, number
	Status    string                 `json:"status"`
	TopicID   int                    `json:"topic_id"`
	PostID    int                    `json:"post_id"`
	Options   []PollOption           `json:"options"`
	Voters    int                    `json:"voters"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

type PollOption struct {
	ID    string `json:"id"`
	HTML  string `json:"html"`
	Votes int    `json:"votes"`
}

type APIKeyRecord struct {
	ID          int       `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	UserID      *int      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
}

type EmailLog struct {
	ID          int       `json:"id"`
	ToAddress   string    `json:"to_address"`
	EmailType   string    `json:"email_type"`
	UserID      int       `json:"user_id"`
	SkippedReason *string `json:"skipped_reason,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserAction struct {
	ID             int       `json:"id"`
	ActionType     int       `json:"action_type"`
	UserID         int       `json:"user_id"`
	TargetTopicID  int       `json:"target_topic_id,omitempty"`
	TargetPostID   int       `json:"target_post_id,omitempty"`
	TargetUserID   int       `json:"target_user_id,omitempty"`
	ActingUserID   int       `json:"acting_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type Webhook struct {
	ID             int       `json:"id"`
	PayloadURL     string    `json:"payload_url"`
	ContentType    int       `json:"content_type"` // 1=json, 2=url_encoded
	Secret         string    `json:"secret,omitempty"`
	WildcardWeb    bool      `json:"wildcard_web_hook"`
	VerifyCert     bool      `json:"verify_certificate"`
	Active         bool      `json:"active"`
	EventTypes     []string  `json:"event_types"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Reviewable struct {
	ID             int       `json:"id"`
	Type           string    `json:"type"`
	Status         int       `json:"status"` // 0=pending, 1=approved, 2=rejected
	CreatedByID    int       `json:"created_by_id"`
	TargetID       int       `json:"target_id,omitempty"`
	TargetType     string    `json:"target_type,omitempty"`
	CategoryID     int       `json:"category_id,omitempty"`
	Score          float64   `json:"score"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Theme struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	UserSelectable bool      `json:"user_selectable"`
	Default        bool      `json:"default"`
	Enabled        bool      `json:"enabled"`
	ColorSchemeID  *int      `json:"color_scheme_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ColorScheme struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Enabled   bool             `json:"enabled"`
	IsBase    bool             `json:"is_base"`
	Colors    []ColorEntry     `json:"colors"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type ColorEntry struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

type CustomUserField struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	FieldType    string    `json:"field_type"` // text, confirm, dropdown
	Editable     string    `json:"editable"`
	Required     bool      `json:"required"`
	ShowOnProfile bool     `json:"show_on_profile"`
	ShowOnUserCard bool    `json:"show_on_user_card"`
	Position     int       `json:"position"`
	Options      []string  `json:"options,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TagGroup struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	TagNames   []string  `json:"tag_names"`
	OnePerTopic bool     `json:"one_per_topic"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Draft struct {
	ID        int       `json:"id"`
	DraftKey  string    `json:"draft_key"`
	UserID    int       `json:"user_id"`
	Data      string    `json:"data"`
	Sequence  int       `json:"sequence"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Bookmark struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	BookmarkableID  int       `json:"bookmarkable_id"`
	BookmarkableType string  `json:"bookmarkable_type"`
	Name            string    `json:"name,omitempty"`
	ReminderAt      *time.Time `json:"reminder_at,omitempty"`
	AutoDeletePref  int       `json:"auto_delete_preference"`
	Pinned          bool      `json:"pinned"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WatchedWord struct {
	ID        int       `json:"id"`
	Word      string    `json:"word"`
	Action    int       `json:"action"` // 0=block, 1=censor, 2=require_approval, 3=flag, 4=replace, 5=tag, 6=silence, 7=link
	Replacement *string `json:"replacement,omitempty"`
	CaseSensitive bool  `json:"case_sensitive"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Permalink struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	TopicID     *int      `json:"topic_id,omitempty"`
	PostID      *int      `json:"post_id,omitempty"`
	CategoryID  *int      `json:"category_id,omitempty"`
	ExternalURL *string   `json:"external_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StaffActionLog struct {
	ID             int       `json:"id"`
	ActionType     string    `json:"action"`
	ActingUserID   int       `json:"acting_user_id"`
	TargetUserID   int       `json:"target_user_id,omitempty"`
	Subject        string    `json:"subject,omitempty"`
	Details        string    `json:"details,omitempty"`
	Context        string    `json:"context,omitempty"`
	PreviousValue  string    `json:"previous_value,omitempty"`
	NewValue       string    `json:"new_value,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ScreenedEmail struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	ActionType   int       `json:"action_type"` // 0=block, 1=do_nothing
	MatchCount   int       `json:"match_count"`
	LastMatchAt  *time.Time `json:"last_match_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ScreenedIP struct {
	ID           int       `json:"id"`
	IPAddress    string    `json:"ip_address"`
	ActionType   int       `json:"action_type"` // 0=block, 1=do_nothing, 2=allow_admin
	MatchCount   int       `json:"match_count"`
	LastMatchAt  *time.Time `json:"last_match_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EmbeddableHost struct {
	ID           int       `json:"id"`
	Host         string    `json:"host"`
	CategoryID   int       `json:"category_id"`
	AllowedPaths *string   `json:"allowed_paths,omitempty"`
	ClassName    *string   `json:"class_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SiteText struct {
	ID        string    `json:"id"`
	Value     string    `json:"value"`
	Overridden bool    `json:"overridden"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SidebarSection struct {
	ID        int                `json:"id"`
	Title     string             `json:"title"`
	Public    bool               `json:"public"`
	UserID    int                `json:"user_id"`
	Links     []SidebarLink      `json:"links"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type SidebarLink struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Icon     string `json:"icon"`
	Segment  string `json:"segment"` // primary, secondary
	Position int    `json:"position"`
}

type PublishedPage struct {
	ID        int       `json:"id"`
	TopicID   int       `json:"topic_id"`
	Slug      string    `json:"slug"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomEmoji struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Group     string    `json:"group,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type FormTemplate struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Template  string    `json:"template"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminFlag struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	NameKey     string    `json:"name_key"`
	Description string    `json:"description"`
	AppliesToTopic bool   `json:"applies_to_topic"`
	AppliesToPost  bool   `json:"applies_to_post"`
	Enabled     bool      `json:"enabled"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostRevision struct {
	ID            int                    `json:"id"`
	PostID        int                    `json:"post_id"`
	UserID        int                    `json:"user_id"`
	Number        int                    `json:"current_revision"`
	PreviousRaw   string                 `json:"previous_raw,omitempty"`
	CurrentRaw    string                 `json:"current_raw"`
	PreviousCooked string               `json:"previous_cooked,omitempty"`
	CurrentCooked string                 `json:"current_cooked"`
	CreatedAt     time.Time              `json:"created_at"`
	Modifications map[string]interface{} `json:"body_changes,omitempty"`
}

type UserStatus struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Description string     `json:"description"`
	Emoji       string     `json:"emoji"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
	SetAt       time.Time  `json:"set_at"`
}

// ---------------------------------------------------------------------------
// ExtStore – extends *Store with additional resource collections.
// ---------------------------------------------------------------------------

// ExtStore embeds a *Store and adds maps for extra Discourse resources that
// are not covered by the base store.  It carries its own RWMutex so that
// callers never have to worry about the unexported mu in the embedded Store.
type ExtStore struct {
	*Store

	mu sync.RWMutex

	Polls             map[int]*Poll
	APIKeyRecords     map[int]*APIKeyRecord
	EmailLogs         map[int]*EmailLog
	UserActions       map[int]*UserAction
	Webhooks          map[int]*Webhook
	Reviewables       map[int]*Reviewable
	Themes            map[int]*Theme
	ColorSchemes      map[int]*ColorScheme
	CustomUserFields  map[int]*CustomUserField
	TagGroups         map[int]*TagGroup
	Drafts            map[int]*Draft
	DraftsByKey       map[string]*Draft // "user_id:draft_key" -> draft
	Bookmarks         map[int]*Bookmark
	WatchedWords      map[int]*WatchedWord
	Permalinks        map[int]*Permalink
	StaffActionLogs   map[int]*StaffActionLog
	ScreenedEmails    map[int]*ScreenedEmail
	ScreenedIPs       map[int]*ScreenedIP
	EmbeddableHosts   map[int]*EmbeddableHost
	SiteTexts         map[string]*SiteText
	SidebarSections   map[int]*SidebarSection
	PublishedPages    map[int]*PublishedPage
	CustomEmojis      map[int]*CustomEmoji
	FormTemplates     map[int]*FormTemplate
	AdminFlags        map[int]*AdminFlag
	PostRevisions     map[int]*PostRevision
	PostRevisionsByPost map[int][]*PostRevision // post_id -> revisions
	UserStatuses      map[int]*UserStatus // user_id -> current status

	// ID counters
	NextPollID           int
	NextAPIKeyRecordID   int
	NextEmailLogID       int
	NextUserActionID     int
	NextWebhookID        int
	NextReviewableID     int
	NextThemeID          int
	NextColorSchemeID    int
	NextCustomUserFieldID int
	NextTagGroupID       int
	NextDraftID          int
	NextBookmarkID       int
	NextWatchedWordID    int
	NextPermalinkID      int
	NextStaffActionLogID int
	NextScreenedEmailID  int
	NextScreenedIPID     int
	NextEmbeddableHostID int
	NextSidebarSectionID int
	NextSidebarLinkID    int
	NextPublishedPageID  int
	NextCustomEmojiID    int
	NextFormTemplateID   int
	NextAdminFlagID      int
	NextPostRevisionID   int
	NextUserStatusID     int
}

// NewExtStore creates a fully initialised ExtStore wrapping the given Store.
// All maps are pre-allocated and a small set of seed data is inserted.
func NewExtStore(s *Store) *ExtStore {
	es := &ExtStore{
		Store:               s,
		Polls:               make(map[int]*Poll),
		APIKeyRecords:       make(map[int]*APIKeyRecord),
		EmailLogs:           make(map[int]*EmailLog),
		UserActions:         make(map[int]*UserAction),
		Webhooks:            make(map[int]*Webhook),
		Reviewables:         make(map[int]*Reviewable),
		Themes:              make(map[int]*Theme),
		ColorSchemes:        make(map[int]*ColorScheme),
		CustomUserFields:    make(map[int]*CustomUserField),
		TagGroups:           make(map[int]*TagGroup),
		Drafts:              make(map[int]*Draft),
		DraftsByKey:         make(map[string]*Draft),
		Bookmarks:           make(map[int]*Bookmark),
		WatchedWords:        make(map[int]*WatchedWord),
		Permalinks:          make(map[int]*Permalink),
		StaffActionLogs:     make(map[int]*StaffActionLog),
		ScreenedEmails:      make(map[int]*ScreenedEmail),
		ScreenedIPs:         make(map[int]*ScreenedIP),
		EmbeddableHosts:     make(map[int]*EmbeddableHost),
		SiteTexts:           make(map[string]*SiteText),
		SidebarSections:     make(map[int]*SidebarSection),
		PublishedPages:      make(map[int]*PublishedPage),
		CustomEmojis:        make(map[int]*CustomEmoji),
		FormTemplates:       make(map[int]*FormTemplate),
		AdminFlags:          make(map[int]*AdminFlag),
		PostRevisions:       make(map[int]*PostRevision),
		PostRevisionsByPost: make(map[int][]*PostRevision),
		UserStatuses:        make(map[int]*UserStatus),

		NextPollID:           1,
		NextAPIKeyRecordID:   1,
		NextEmailLogID:       1,
		NextUserActionID:     1,
		NextWebhookID:        1,
		NextReviewableID:     1,
		NextThemeID:          1,
		NextColorSchemeID:    1,
		NextCustomUserFieldID: 1,
		NextTagGroupID:       1,
		NextDraftID:          1,
		NextBookmarkID:       1,
		NextWatchedWordID:    1,
		NextPermalinkID:      1,
		NextStaffActionLogID: 1,
		NextScreenedEmailID:  1,
		NextScreenedIPID:     1,
		NextEmbeddableHostID: 1,
		NextSidebarSectionID: 1,
		NextSidebarLinkID:    1,
		NextPublishedPageID:  1,
		NextCustomEmojiID:    1,
		NextFormTemplateID:   1,
		NextAdminFlagID:      1,
		NextPostRevisionID:   1,
		NextUserStatusID:     1,
	}
	es.seedExtended()
	return es
}

// seedExtended inserts a small amount of realistic data into every collection
// so that the extended store is immediately useful for testing.
func (es *ExtStore) seedExtended() {
	now := time.Now().UTC()

	// --- Themes ---
	es.Themes[1] = &Theme{
		ID: 1, Name: "Default", UserSelectable: true, Default: true,
		Enabled: true, CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	es.Themes[2] = &Theme{
		ID: 2, Name: "Dark", UserSelectable: true, Default: false,
		Enabled: true, CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	es.NextThemeID = 3

	// --- Color Schemes ---
	es.ColorSchemes[1] = &ColorScheme{
		ID: 1, Name: "Light", Enabled: true, IsBase: true,
		Colors: []ColorEntry{
			{Name: "primary", Hex: "222222"},
			{Name: "secondary", Hex: "FFFFFF"},
			{Name: "tertiary", Hex: "0088CC"},
			{Name: "header_background", Hex: "FFFFFF"},
		},
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	es.ColorSchemes[2] = &ColorScheme{
		ID: 2, Name: "Dark", Enabled: true, IsBase: false,
		Colors: []ColorEntry{
			{Name: "primary", Hex: "DDDDDD"},
			{Name: "secondary", Hex: "1E1E1E"},
			{Name: "tertiary", Hex: "6699CC"},
			{Name: "header_background", Hex: "1E1E1E"},
		},
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	es.NextColorSchemeID = 3

	// --- Custom User Fields ---
	es.CustomUserFields[1] = &CustomUserField{
		ID: 1, Name: "Location", Description: "Your city or country",
		FieldType: "text", Editable: "true", Required: false,
		ShowOnProfile: true, ShowOnUserCard: true, Position: 0,
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	es.NextCustomUserFieldID = 2

	// --- Tag Groups ---
	es.TagGroups[1] = &TagGroup{
		ID: 1, Name: "Topic Types", TagNames: []string{"question", "discussion", "announcement"},
		OnePerTopic: true, CreatedAt: now.Add(-20 * 24 * time.Hour), UpdatedAt: now,
	}
	es.NextTagGroupID = 2

	// --- Admin Flags ---
	for i, f := range []struct{ name, key, desc string }{
		{"Spam", "spam", "This post is an advertisement or vandalism"},
		{"Inappropriate", "inappropriate", "This post contains content that is offensive or violates community guidelines"},
		{"Off-Topic", "off_topic", "This post is not relevant to the current discussion"},
		{"Something Else", "notify_moderators", "This post requires staff attention for another reason"},
	} {
		id := i + 1
		es.AdminFlags[id] = &AdminFlag{
			ID: id, Name: f.name, NameKey: f.key, Description: f.desc,
			AppliesToPost: true, AppliesToTopic: true, Enabled: true,
			Position: i, CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
		}
	}
	es.NextAdminFlagID = 5

	// --- Watched Words ---
	es.WatchedWords[1] = &WatchedWord{
		ID: 1, Word: "spam-link.example.com", Action: 0,
		CaseSensitive: false, CreatedAt: now.Add(-10 * 24 * time.Hour), UpdatedAt: now,
	}
	es.NextWatchedWordID = 2

	// --- Embeddable Hosts ---
	es.EmbeddableHosts[1] = &EmbeddableHost{
		ID: 1, Host: "blog.example.com", CategoryID: 1,
		CreatedAt: now.Add(-15 * 24 * time.Hour), UpdatedAt: now,
	}
	es.NextEmbeddableHostID = 2

	// --- Site Texts ---
	es.SiteTexts["js.topic.create"] = &SiteText{
		ID: "js.topic.create", Value: "Create Topic", Overridden: false, UpdatedAt: now,
	}
	es.SiteTexts["js.topic.reply.title"] = &SiteText{
		ID: "js.topic.reply.title", Value: "Reply", Overridden: false, UpdatedAt: now,
	}

	// --- Sidebar Sections ---
	es.SidebarSections[1] = &SidebarSection{
		ID: 1, Title: "Community", Public: true, UserID: 1,
		Links: []SidebarLink{
			{ID: 1, Name: "Everything", Value: "/latest", Icon: "layer-group", Segment: "primary", Position: 0},
			{ID: 2, Name: "My Posts", Value: "/my/activity", Icon: "user", Segment: "primary", Position: 1},
		},
		CreatedAt: now.Add(-20 * 24 * time.Hour), UpdatedAt: now,
	}
	es.NextSidebarSectionID = 2
	es.NextSidebarLinkID = 3

	// --- Custom Emojis ---
	es.CustomEmojis[1] = &CustomEmoji{
		ID: 1, Name: "party_blob", URL: "/uploads/default/custom_emoji/party_blob.gif",
		Group: "", CreatedAt: now.Add(-10 * 24 * time.Hour),
	}
	es.NextCustomEmojiID = 2

	// --- API Key Records ---
	es.APIKeyRecords[1] = &APIKeyRecord{
		ID: 1, Key: "test_api_key", Description: "System test key",
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	es.APIKeyRecords[2] = &APIKeyRecord{
		ID: 2, Key: "admin_api_key", Description: "Admin key",
		UserID: intPtr(1),
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	es.NextAPIKeyRecordID = 3
}

// intPtr is a helper for creating an *int inline.
func intPtr(v int) *int { return &v }

// stringPtr is a helper for creating a *string inline.
func stringPtr(v string) *string { return &v }

// ---------------------------------------------------------------------------
// Poll CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetPoll(id int) (*Poll, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	p, ok := es.Polls[id]
	if !ok {
		return nil, fmt.Errorf("poll not found")
	}
	return p, nil
}

func (es *ExtStore) ListPolls() []Poll {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]Poll, 0, len(es.Polls))
	for _, p := range es.Polls {
		out = append(out, *p)
	}
	return out
}

func (es *ExtStore) CreatePoll(name, pollType string, topicID, postID int, options []PollOption) (*Poll, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	p := &Poll{
		ID: es.NextPollID, Name: name, Type: pollType, Status: "open",
		TopicID: topicID, PostID: postID, Options: options,
		CreatedAt: now, UpdatedAt: now,
	}
	es.Polls[p.ID] = p
	es.NextPollID++
	return p, nil
}

func (es *ExtStore) UpdatePoll(id int, updates map[string]interface{}) (*Poll, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	p, ok := es.Polls[id]
	if !ok {
		return nil, fmt.Errorf("poll not found")
	}
	if v, ok := updates["status"].(string); ok {
		p.Status = v
	}
	if v, ok := updates["name"].(string); ok {
		p.Name = v
	}
	p.UpdatedAt = time.Now().UTC()
	return p, nil
}

func (es *ExtStore) DeletePoll(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.Polls[id]; !ok {
		return fmt.Errorf("poll not found")
	}
	delete(es.Polls, id)
	return nil
}

// ---------------------------------------------------------------------------
// APIKeyRecord CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetAPIKeyRecord(id int) (*APIKeyRecord, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	r, ok := es.APIKeyRecords[id]
	if !ok {
		return nil, fmt.Errorf("api key record not found")
	}
	return r, nil
}

func (es *ExtStore) ListAPIKeyRecords() []APIKeyRecord {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]APIKeyRecord, 0, len(es.APIKeyRecords))
	for _, r := range es.APIKeyRecords {
		out = append(out, *r)
	}
	return out
}

func (es *ExtStore) CreateAPIKeyRecord(description string, userID *int) (*APIKeyRecord, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	key := fmt.Sprintf("dk_%d_%d", es.NextAPIKeyRecordID, now.UnixNano())
	r := &APIKeyRecord{
		ID: es.NextAPIKeyRecordID, Key: key, Description: description,
		UserID: userID, CreatedAt: now, UpdatedAt: now,
	}
	es.APIKeyRecords[r.ID] = r
	es.NextAPIKeyRecordID++
	return r, nil
}

func (es *ExtStore) UpdateAPIKeyRecord(id int, updates map[string]interface{}) (*APIKeyRecord, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	r, ok := es.APIKeyRecords[id]
	if !ok {
		return nil, fmt.Errorf("api key record not found")
	}
	if v, ok := updates["description"].(string); ok {
		r.Description = v
	}
	r.UpdatedAt = time.Now().UTC()
	return r, nil
}

func (es *ExtStore) DeleteAPIKeyRecord(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.APIKeyRecords[id]; !ok {
		return fmt.Errorf("api key record not found")
	}
	delete(es.APIKeyRecords, id)
	return nil
}

// ---------------------------------------------------------------------------
// EmailLog CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetEmailLog(id int) (*EmailLog, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	e, ok := es.EmailLogs[id]
	if !ok {
		return nil, fmt.Errorf("email log not found")
	}
	return e, nil
}

func (es *ExtStore) ListEmailLogs() []EmailLog {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]EmailLog, 0, len(es.EmailLogs))
	for _, e := range es.EmailLogs {
		out = append(out, *e)
	}
	return out
}

func (es *ExtStore) CreateEmailLog(toAddress, emailType string, userID int) (*EmailLog, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	e := &EmailLog{
		ID: es.NextEmailLogID, ToAddress: toAddress, EmailType: emailType,
		UserID: userID, CreatedAt: time.Now().UTC(),
	}
	es.EmailLogs[e.ID] = e
	es.NextEmailLogID++
	return e, nil
}

func (es *ExtStore) DeleteEmailLog(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.EmailLogs[id]; !ok {
		return fmt.Errorf("email log not found")
	}
	delete(es.EmailLogs, id)
	return nil
}

// ---------------------------------------------------------------------------
// UserAction CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetUserAction(id int) (*UserAction, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	a, ok := es.UserActions[id]
	if !ok {
		return nil, fmt.Errorf("user action not found")
	}
	return a, nil
}

func (es *ExtStore) ListUserActions(userID int) []UserAction {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]UserAction, 0)
	for _, a := range es.UserActions {
		if a.UserID == userID {
			out = append(out, *a)
		}
	}
	return out
}

func (es *ExtStore) CreateUserAction(actionType, userID, actingUserID int) (*UserAction, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	a := &UserAction{
		ID: es.NextUserActionID, ActionType: actionType,
		UserID: userID, ActingUserID: actingUserID,
		CreatedAt: time.Now().UTC(),
	}
	es.UserActions[a.ID] = a
	es.NextUserActionID++
	return a, nil
}

func (es *ExtStore) DeleteUserAction(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.UserActions[id]; !ok {
		return fmt.Errorf("user action not found")
	}
	delete(es.UserActions, id)
	return nil
}

// ---------------------------------------------------------------------------
// Webhook CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetWebhook(id int) (*Webhook, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	w, ok := es.Webhooks[id]
	if !ok {
		return nil, fmt.Errorf("webhook not found")
	}
	return w, nil
}

func (es *ExtStore) ListWebhooks() []Webhook {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]Webhook, 0, len(es.Webhooks))
	for _, w := range es.Webhooks {
		out = append(out, *w)
	}
	return out
}

func (es *ExtStore) CreateWebhook(payloadURL string, eventTypes []string) (*Webhook, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	w := &Webhook{
		ID: es.NextWebhookID, PayloadURL: payloadURL, ContentType: 1,
		WildcardWeb: false, VerifyCert: true, Active: true,
		EventTypes: eventTypes, CreatedAt: now, UpdatedAt: now,
	}
	es.Webhooks[w.ID] = w
	es.NextWebhookID++
	return w, nil
}

func (es *ExtStore) UpdateWebhook(id int, updates map[string]interface{}) (*Webhook, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	w, ok := es.Webhooks[id]
	if !ok {
		return nil, fmt.Errorf("webhook not found")
	}
	if v, ok := updates["payload_url"].(string); ok {
		w.PayloadURL = v
	}
	if v, ok := updates["active"].(bool); ok {
		w.Active = v
	}
	w.UpdatedAt = time.Now().UTC()
	return w, nil
}

func (es *ExtStore) DeleteWebhook(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.Webhooks[id]; !ok {
		return fmt.Errorf("webhook not found")
	}
	delete(es.Webhooks, id)
	return nil
}

// ---------------------------------------------------------------------------
// Reviewable CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetReviewable(id int) (*Reviewable, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	r, ok := es.Reviewables[id]
	if !ok {
		return nil, fmt.Errorf("reviewable not found")
	}
	return r, nil
}

func (es *ExtStore) ListReviewables() []Reviewable {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]Reviewable, 0, len(es.Reviewables))
	for _, r := range es.Reviewables {
		out = append(out, *r)
	}
	return out
}

func (es *ExtStore) CreateReviewable(reviewType string, createdByID, targetID int, targetType string) (*Reviewable, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	r := &Reviewable{
		ID: es.NextReviewableID, Type: reviewType, Status: 0,
		CreatedByID: createdByID, TargetID: targetID, TargetType: targetType,
		CreatedAt: now, UpdatedAt: now,
	}
	es.Reviewables[r.ID] = r
	es.NextReviewableID++
	return r, nil
}

func (es *ExtStore) UpdateReviewable(id int, updates map[string]interface{}) (*Reviewable, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	r, ok := es.Reviewables[id]
	if !ok {
		return nil, fmt.Errorf("reviewable not found")
	}
	if v, ok := updates["status"].(float64); ok {
		r.Status = int(v)
	}
	r.UpdatedAt = time.Now().UTC()
	return r, nil
}

func (es *ExtStore) DeleteReviewable(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.Reviewables[id]; !ok {
		return fmt.Errorf("reviewable not found")
	}
	delete(es.Reviewables, id)
	return nil
}

// ---------------------------------------------------------------------------
// Theme CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetTheme(id int) (*Theme, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	t, ok := es.Themes[id]
	if !ok {
		return nil, fmt.Errorf("theme not found")
	}
	return t, nil
}

func (es *ExtStore) ListThemes() []Theme {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]Theme, 0, len(es.Themes))
	for _, t := range es.Themes {
		out = append(out, *t)
	}
	return out
}

func (es *ExtStore) CreateTheme(name string, userSelectable bool) (*Theme, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	t := &Theme{
		ID: es.NextThemeID, Name: name, UserSelectable: userSelectable,
		Enabled: true, CreatedAt: now, UpdatedAt: now,
	}
	es.Themes[t.ID] = t
	es.NextThemeID++
	return t, nil
}

func (es *ExtStore) UpdateTheme(id int, updates map[string]interface{}) (*Theme, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	t, ok := es.Themes[id]
	if !ok {
		return nil, fmt.Errorf("theme not found")
	}
	if v, ok := updates["name"].(string); ok {
		t.Name = v
	}
	if v, ok := updates["user_selectable"].(bool); ok {
		t.UserSelectable = v
	}
	if v, ok := updates["enabled"].(bool); ok {
		t.Enabled = v
	}
	if v, ok := updates["default"].(bool); ok {
		t.Default = v
	}
	t.UpdatedAt = time.Now().UTC()
	return t, nil
}

func (es *ExtStore) DeleteTheme(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.Themes[id]; !ok {
		return fmt.Errorf("theme not found")
	}
	delete(es.Themes, id)
	return nil
}

// ---------------------------------------------------------------------------
// ColorScheme CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetColorScheme(id int) (*ColorScheme, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	c, ok := es.ColorSchemes[id]
	if !ok {
		return nil, fmt.Errorf("color scheme not found")
	}
	return c, nil
}

func (es *ExtStore) ListColorSchemes() []ColorScheme {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]ColorScheme, 0, len(es.ColorSchemes))
	for _, c := range es.ColorSchemes {
		out = append(out, *c)
	}
	return out
}

func (es *ExtStore) CreateColorScheme(name string, colors []ColorEntry) (*ColorScheme, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	c := &ColorScheme{
		ID: es.NextColorSchemeID, Name: name, Enabled: true,
		Colors: colors, CreatedAt: now, UpdatedAt: now,
	}
	es.ColorSchemes[c.ID] = c
	es.NextColorSchemeID++
	return c, nil
}

func (es *ExtStore) UpdateColorScheme(id int, updates map[string]interface{}) (*ColorScheme, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	c, ok := es.ColorSchemes[id]
	if !ok {
		return nil, fmt.Errorf("color scheme not found")
	}
	if v, ok := updates["name"].(string); ok {
		c.Name = v
	}
	if v, ok := updates["enabled"].(bool); ok {
		c.Enabled = v
	}
	c.UpdatedAt = time.Now().UTC()
	return c, nil
}

func (es *ExtStore) DeleteColorScheme(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.ColorSchemes[id]; !ok {
		return fmt.Errorf("color scheme not found")
	}
	delete(es.ColorSchemes, id)
	return nil
}

// ---------------------------------------------------------------------------
// CustomUserField CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetCustomUserField(id int) (*CustomUserField, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	f, ok := es.CustomUserFields[id]
	if !ok {
		return nil, fmt.Errorf("custom user field not found")
	}
	return f, nil
}

func (es *ExtStore) ListCustomUserFields() []CustomUserField {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]CustomUserField, 0, len(es.CustomUserFields))
	for _, f := range es.CustomUserFields {
		out = append(out, *f)
	}
	return out
}

func (es *ExtStore) CreateCustomUserField(name, description, fieldType string) (*CustomUserField, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	f := &CustomUserField{
		ID: es.NextCustomUserFieldID, Name: name, Description: description,
		FieldType: fieldType, Editable: "true", Position: len(es.CustomUserFields),
		CreatedAt: now, UpdatedAt: now,
	}
	es.CustomUserFields[f.ID] = f
	es.NextCustomUserFieldID++
	return f, nil
}

func (es *ExtStore) UpdateCustomUserField(id int, updates map[string]interface{}) (*CustomUserField, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	f, ok := es.CustomUserFields[id]
	if !ok {
		return nil, fmt.Errorf("custom user field not found")
	}
	if v, ok := updates["name"].(string); ok {
		f.Name = v
	}
	if v, ok := updates["description"].(string); ok {
		f.Description = v
	}
	if v, ok := updates["required"].(bool); ok {
		f.Required = v
	}
	if v, ok := updates["show_on_profile"].(bool); ok {
		f.ShowOnProfile = v
	}
	f.UpdatedAt = time.Now().UTC()
	return f, nil
}

func (es *ExtStore) DeleteCustomUserField(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.CustomUserFields[id]; !ok {
		return fmt.Errorf("custom user field not found")
	}
	delete(es.CustomUserFields, id)
	return nil
}

// ---------------------------------------------------------------------------
// TagGroup CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetTagGroup(id int) (*TagGroup, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	g, ok := es.TagGroups[id]
	if !ok {
		return nil, fmt.Errorf("tag group not found")
	}
	return g, nil
}

func (es *ExtStore) ListTagGroups() []TagGroup {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]TagGroup, 0, len(es.TagGroups))
	for _, g := range es.TagGroups {
		out = append(out, *g)
	}
	return out
}

func (es *ExtStore) CreateTagGroup(name string, tagNames []string) (*TagGroup, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	g := &TagGroup{
		ID: es.NextTagGroupID, Name: name, TagNames: tagNames,
		CreatedAt: now, UpdatedAt: now,
	}
	es.TagGroups[g.ID] = g
	es.NextTagGroupID++
	return g, nil
}

func (es *ExtStore) UpdateTagGroup(id int, updates map[string]interface{}) (*TagGroup, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	g, ok := es.TagGroups[id]
	if !ok {
		return nil, fmt.Errorf("tag group not found")
	}
	if v, ok := updates["name"].(string); ok {
		g.Name = v
	}
	if v, ok := updates["one_per_topic"].(bool); ok {
		g.OnePerTopic = v
	}
	g.UpdatedAt = time.Now().UTC()
	return g, nil
}

func (es *ExtStore) DeleteTagGroup(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.TagGroups[id]; !ok {
		return fmt.Errorf("tag group not found")
	}
	delete(es.TagGroups, id)
	return nil
}

// ---------------------------------------------------------------------------
// Draft CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetDraft(id int) (*Draft, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	d, ok := es.Drafts[id]
	if !ok {
		return nil, fmt.Errorf("draft not found")
	}
	return d, nil
}

func (es *ExtStore) GetDraftByKey(userID int, draftKey string) *Draft {
	es.mu.RLock()
	defer es.mu.RUnlock()
	key := fmt.Sprintf("%d:%s", userID, draftKey)
	return es.DraftsByKey[key]
}

func (es *ExtStore) ListDrafts(userID int) []Draft {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]Draft, 0)
	for _, d := range es.Drafts {
		if d.UserID == userID {
			out = append(out, *d)
		}
	}
	return out
}

func (es *ExtStore) CreateDraft(draftKey string, userID int, data string) (*Draft, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	compositeKey := fmt.Sprintf("%d:%s", userID, draftKey)
	// Upsert: if a draft with this key exists, update it.
	if existing, ok := es.DraftsByKey[compositeKey]; ok {
		existing.Data = data
		existing.Sequence++
		existing.UpdatedAt = now
		return existing, nil
	}
	d := &Draft{
		ID: es.NextDraftID, DraftKey: draftKey, UserID: userID,
		Data: data, Sequence: 0, CreatedAt: now, UpdatedAt: now,
	}
	es.Drafts[d.ID] = d
	es.DraftsByKey[compositeKey] = d
	es.NextDraftID++
	return d, nil
}

func (es *ExtStore) DeleteDraft(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	d, ok := es.Drafts[id]
	if !ok {
		return fmt.Errorf("draft not found")
	}
	compositeKey := fmt.Sprintf("%d:%s", d.UserID, d.DraftKey)
	delete(es.Drafts, id)
	delete(es.DraftsByKey, compositeKey)
	return nil
}

// ---------------------------------------------------------------------------
// Bookmark CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetBookmark(id int) (*Bookmark, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	b, ok := es.Bookmarks[id]
	if !ok {
		return nil, fmt.Errorf("bookmark not found")
	}
	return b, nil
}

func (es *ExtStore) ListBookmarks(userID int) []Bookmark {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]Bookmark, 0)
	for _, b := range es.Bookmarks {
		if b.UserID == userID {
			out = append(out, *b)
		}
	}
	return out
}

func (es *ExtStore) CreateBookmark(userID, bookmarkableID int, bookmarkableType string) (*Bookmark, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	b := &Bookmark{
		ID: es.NextBookmarkID, UserID: userID,
		BookmarkableID: bookmarkableID, BookmarkableType: bookmarkableType,
		CreatedAt: now, UpdatedAt: now,
	}
	es.Bookmarks[b.ID] = b
	es.NextBookmarkID++
	return b, nil
}

func (es *ExtStore) UpdateBookmark(id int, updates map[string]interface{}) (*Bookmark, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	b, ok := es.Bookmarks[id]
	if !ok {
		return nil, fmt.Errorf("bookmark not found")
	}
	if v, ok := updates["name"].(string); ok {
		b.Name = v
	}
	if v, ok := updates["pinned"].(bool); ok {
		b.Pinned = v
	}
	b.UpdatedAt = time.Now().UTC()
	return b, nil
}

func (es *ExtStore) DeleteBookmark(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.Bookmarks[id]; !ok {
		return fmt.Errorf("bookmark not found")
	}
	delete(es.Bookmarks, id)
	return nil
}

// ---------------------------------------------------------------------------
// WatchedWord CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetWatchedWord(id int) (*WatchedWord, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	w, ok := es.WatchedWords[id]
	if !ok {
		return nil, fmt.Errorf("watched word not found")
	}
	return w, nil
}

func (es *ExtStore) ListWatchedWords() []WatchedWord {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]WatchedWord, 0, len(es.WatchedWords))
	for _, w := range es.WatchedWords {
		out = append(out, *w)
	}
	return out
}

func (es *ExtStore) CreateWatchedWord(word string, action int) (*WatchedWord, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	w := &WatchedWord{
		ID: es.NextWatchedWordID, Word: word, Action: action,
		CreatedAt: now, UpdatedAt: now,
	}
	es.WatchedWords[w.ID] = w
	es.NextWatchedWordID++
	return w, nil
}

func (es *ExtStore) UpdateWatchedWord(id int, updates map[string]interface{}) (*WatchedWord, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	w, ok := es.WatchedWords[id]
	if !ok {
		return nil, fmt.Errorf("watched word not found")
	}
	if v, ok := updates["word"].(string); ok {
		w.Word = v
	}
	if v, ok := updates["action"].(float64); ok {
		w.Action = int(v)
	}
	if v, ok := updates["case_sensitive"].(bool); ok {
		w.CaseSensitive = v
	}
	w.UpdatedAt = time.Now().UTC()
	return w, nil
}

func (es *ExtStore) DeleteWatchedWord(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.WatchedWords[id]; !ok {
		return fmt.Errorf("watched word not found")
	}
	delete(es.WatchedWords, id)
	return nil
}

// ---------------------------------------------------------------------------
// Permalink CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetPermalink(id int) (*Permalink, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	p, ok := es.Permalinks[id]
	if !ok {
		return nil, fmt.Errorf("permalink not found")
	}
	return p, nil
}

func (es *ExtStore) ListPermalinks() []Permalink {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]Permalink, 0, len(es.Permalinks))
	for _, p := range es.Permalinks {
		out = append(out, *p)
	}
	return out
}

func (es *ExtStore) CreatePermalink(url string, topicID, postID, categoryID *int, externalURL *string) (*Permalink, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	p := &Permalink{
		ID: es.NextPermalinkID, URL: url, TopicID: topicID,
		PostID: postID, CategoryID: categoryID, ExternalURL: externalURL,
		CreatedAt: now, UpdatedAt: now,
	}
	es.Permalinks[p.ID] = p
	es.NextPermalinkID++
	return p, nil
}

func (es *ExtStore) UpdatePermalink(id int, updates map[string]interface{}) (*Permalink, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	p, ok := es.Permalinks[id]
	if !ok {
		return nil, fmt.Errorf("permalink not found")
	}
	if v, ok := updates["url"].(string); ok {
		p.URL = v
	}
	if v, ok := updates["external_url"].(string); ok {
		p.ExternalURL = &v
	}
	p.UpdatedAt = time.Now().UTC()
	return p, nil
}

func (es *ExtStore) DeletePermalink(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.Permalinks[id]; !ok {
		return fmt.Errorf("permalink not found")
	}
	delete(es.Permalinks, id)
	return nil
}

// ---------------------------------------------------------------------------
// StaffActionLog CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetStaffActionLog(id int) (*StaffActionLog, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	l, ok := es.StaffActionLogs[id]
	if !ok {
		return nil, fmt.Errorf("staff action log not found")
	}
	return l, nil
}

func (es *ExtStore) ListStaffActionLogs() []StaffActionLog {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]StaffActionLog, 0, len(es.StaffActionLogs))
	for _, l := range es.StaffActionLogs {
		out = append(out, *l)
	}
	return out
}

func (es *ExtStore) CreateStaffActionLog(actionType string, actingUserID int, details string) (*StaffActionLog, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	l := &StaffActionLog{
		ID: es.NextStaffActionLogID, ActionType: actionType,
		ActingUserID: actingUserID, Details: details,
		CreatedAt: time.Now().UTC(),
	}
	es.StaffActionLogs[l.ID] = l
	es.NextStaffActionLogID++
	return l, nil
}

func (es *ExtStore) DeleteStaffActionLog(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.StaffActionLogs[id]; !ok {
		return fmt.Errorf("staff action log not found")
	}
	delete(es.StaffActionLogs, id)
	return nil
}

// ---------------------------------------------------------------------------
// ScreenedEmail CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetScreenedEmail(id int) (*ScreenedEmail, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	e, ok := es.ScreenedEmails[id]
	if !ok {
		return nil, fmt.Errorf("screened email not found")
	}
	return e, nil
}

func (es *ExtStore) ListScreenedEmails() []ScreenedEmail {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]ScreenedEmail, 0, len(es.ScreenedEmails))
	for _, e := range es.ScreenedEmails {
		out = append(out, *e)
	}
	return out
}

func (es *ExtStore) CreateScreenedEmail(email string, actionType int) (*ScreenedEmail, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	e := &ScreenedEmail{
		ID: es.NextScreenedEmailID, Email: email, ActionType: actionType,
		CreatedAt: now, UpdatedAt: now,
	}
	es.ScreenedEmails[e.ID] = e
	es.NextScreenedEmailID++
	return e, nil
}

func (es *ExtStore) UpdateScreenedEmail(id int, updates map[string]interface{}) (*ScreenedEmail, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	e, ok := es.ScreenedEmails[id]
	if !ok {
		return nil, fmt.Errorf("screened email not found")
	}
	if v, ok := updates["action_type"].(float64); ok {
		e.ActionType = int(v)
	}
	e.UpdatedAt = time.Now().UTC()
	return e, nil
}

func (es *ExtStore) DeleteScreenedEmail(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.ScreenedEmails[id]; !ok {
		return fmt.Errorf("screened email not found")
	}
	delete(es.ScreenedEmails, id)
	return nil
}

// ---------------------------------------------------------------------------
// ScreenedIP CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetScreenedIP(id int) (*ScreenedIP, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	ip, ok := es.ScreenedIPs[id]
	if !ok {
		return nil, fmt.Errorf("screened ip not found")
	}
	return ip, nil
}

func (es *ExtStore) ListScreenedIPs() []ScreenedIP {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]ScreenedIP, 0, len(es.ScreenedIPs))
	for _, ip := range es.ScreenedIPs {
		out = append(out, *ip)
	}
	return out
}

func (es *ExtStore) CreateScreenedIP(ipAddress string, actionType int) (*ScreenedIP, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	ip := &ScreenedIP{
		ID: es.NextScreenedIPID, IPAddress: ipAddress, ActionType: actionType,
		CreatedAt: now, UpdatedAt: now,
	}
	es.ScreenedIPs[ip.ID] = ip
	es.NextScreenedIPID++
	return ip, nil
}

func (es *ExtStore) UpdateScreenedIP(id int, updates map[string]interface{}) (*ScreenedIP, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	ip, ok := es.ScreenedIPs[id]
	if !ok {
		return nil, fmt.Errorf("screened ip not found")
	}
	if v, ok := updates["action_type"].(float64); ok {
		ip.ActionType = int(v)
	}
	ip.UpdatedAt = time.Now().UTC()
	return ip, nil
}

func (es *ExtStore) DeleteScreenedIP(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.ScreenedIPs[id]; !ok {
		return fmt.Errorf("screened ip not found")
	}
	delete(es.ScreenedIPs, id)
	return nil
}

// ---------------------------------------------------------------------------
// EmbeddableHost CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetEmbeddableHost(id int) (*EmbeddableHost, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	h, ok := es.EmbeddableHosts[id]
	if !ok {
		return nil, fmt.Errorf("embeddable host not found")
	}
	return h, nil
}

func (es *ExtStore) ListEmbeddableHosts() []EmbeddableHost {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]EmbeddableHost, 0, len(es.EmbeddableHosts))
	for _, h := range es.EmbeddableHosts {
		out = append(out, *h)
	}
	return out
}

func (es *ExtStore) CreateEmbeddableHost(host string, categoryID int) (*EmbeddableHost, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	h := &EmbeddableHost{
		ID: es.NextEmbeddableHostID, Host: host, CategoryID: categoryID,
		CreatedAt: now, UpdatedAt: now,
	}
	es.EmbeddableHosts[h.ID] = h
	es.NextEmbeddableHostID++
	return h, nil
}

func (es *ExtStore) UpdateEmbeddableHost(id int, updates map[string]interface{}) (*EmbeddableHost, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	h, ok := es.EmbeddableHosts[id]
	if !ok {
		return nil, fmt.Errorf("embeddable host not found")
	}
	if v, ok := updates["host"].(string); ok {
		h.Host = v
	}
	if v, ok := updates["category_id"].(float64); ok {
		h.CategoryID = int(v)
	}
	h.UpdatedAt = time.Now().UTC()
	return h, nil
}

func (es *ExtStore) DeleteEmbeddableHost(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.EmbeddableHosts[id]; !ok {
		return fmt.Errorf("embeddable host not found")
	}
	delete(es.EmbeddableHosts, id)
	return nil
}

// ---------------------------------------------------------------------------
// SiteText CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetSiteText(id string) (*SiteText, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	t, ok := es.SiteTexts[id]
	if !ok {
		return nil, fmt.Errorf("site text not found")
	}
	return t, nil
}

func (es *ExtStore) ListSiteTexts() []SiteText {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]SiteText, 0, len(es.SiteTexts))
	for _, t := range es.SiteTexts {
		out = append(out, *t)
	}
	return out
}

func (es *ExtStore) UpdateSiteText(id, value string) (*SiteText, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	t, ok := es.SiteTexts[id]
	if !ok {
		// Create if it doesn't exist.
		t = &SiteText{ID: id}
		es.SiteTexts[id] = t
	}
	t.Value = value
	t.Overridden = true
	t.UpdatedAt = time.Now().UTC()
	return t, nil
}

func (es *ExtStore) DeleteSiteText(id string) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.SiteTexts[id]; !ok {
		return fmt.Errorf("site text not found")
	}
	delete(es.SiteTexts, id)
	return nil
}

// ---------------------------------------------------------------------------
// SidebarSection CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetSidebarSection(id int) (*SidebarSection, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	s, ok := es.SidebarSections[id]
	if !ok {
		return nil, fmt.Errorf("sidebar section not found")
	}
	return s, nil
}

func (es *ExtStore) ListSidebarSections() []SidebarSection {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]SidebarSection, 0, len(es.SidebarSections))
	for _, s := range es.SidebarSections {
		out = append(out, *s)
	}
	return out
}

func (es *ExtStore) CreateSidebarSection(title string, public bool, userID int, links []SidebarLink) (*SidebarSection, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	// Assign IDs to links
	for i := range links {
		links[i].ID = es.NextSidebarLinkID
		es.NextSidebarLinkID++
	}
	s := &SidebarSection{
		ID: es.NextSidebarSectionID, Title: title, Public: public,
		UserID: userID, Links: links, CreatedAt: now, UpdatedAt: now,
	}
	es.SidebarSections[s.ID] = s
	es.NextSidebarSectionID++
	return s, nil
}

func (es *ExtStore) UpdateSidebarSection(id int, updates map[string]interface{}) (*SidebarSection, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	s, ok := es.SidebarSections[id]
	if !ok {
		return nil, fmt.Errorf("sidebar section not found")
	}
	if v, ok := updates["title"].(string); ok {
		s.Title = v
	}
	if v, ok := updates["public"].(bool); ok {
		s.Public = v
	}
	s.UpdatedAt = time.Now().UTC()
	return s, nil
}

func (es *ExtStore) DeleteSidebarSection(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.SidebarSections[id]; !ok {
		return fmt.Errorf("sidebar section not found")
	}
	delete(es.SidebarSections, id)
	return nil
}

// ---------------------------------------------------------------------------
// PublishedPage CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetPublishedPage(id int) (*PublishedPage, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	p, ok := es.PublishedPages[id]
	if !ok {
		return nil, fmt.Errorf("published page not found")
	}
	return p, nil
}

func (es *ExtStore) GetPublishedPageBySlug(slug string) *PublishedPage {
	es.mu.RLock()
	defer es.mu.RUnlock()
	for _, p := range es.PublishedPages {
		if p.Slug == slug {
			return p
		}
	}
	return nil
}

func (es *ExtStore) ListPublishedPages() []PublishedPage {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]PublishedPage, 0, len(es.PublishedPages))
	for _, p := range es.PublishedPages {
		out = append(out, *p)
	}
	return out
}

func (es *ExtStore) CreatePublishedPage(topicID int, slug string, public bool) (*PublishedPage, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	p := &PublishedPage{
		ID: es.NextPublishedPageID, TopicID: topicID, Slug: slug,
		Public: public, CreatedAt: now, UpdatedAt: now,
	}
	es.PublishedPages[p.ID] = p
	es.NextPublishedPageID++
	return p, nil
}

func (es *ExtStore) UpdatePublishedPage(id int, updates map[string]interface{}) (*PublishedPage, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	p, ok := es.PublishedPages[id]
	if !ok {
		return nil, fmt.Errorf("published page not found")
	}
	if v, ok := updates["slug"].(string); ok {
		p.Slug = v
	}
	if v, ok := updates["public"].(bool); ok {
		p.Public = v
	}
	p.UpdatedAt = time.Now().UTC()
	return p, nil
}

func (es *ExtStore) DeletePublishedPage(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.PublishedPages[id]; !ok {
		return fmt.Errorf("published page not found")
	}
	delete(es.PublishedPages, id)
	return nil
}

// ---------------------------------------------------------------------------
// CustomEmoji CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetCustomEmoji(id int) (*CustomEmoji, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	e, ok := es.CustomEmojis[id]
	if !ok {
		return nil, fmt.Errorf("custom emoji not found")
	}
	return e, nil
}

func (es *ExtStore) ListCustomEmojis() []CustomEmoji {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]CustomEmoji, 0, len(es.CustomEmojis))
	for _, e := range es.CustomEmojis {
		out = append(out, *e)
	}
	return out
}

func (es *ExtStore) CreateCustomEmoji(name, url, group string) (*CustomEmoji, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	e := &CustomEmoji{
		ID: es.NextCustomEmojiID, Name: name, URL: url, Group: group,
		CreatedAt: time.Now().UTC(),
	}
	es.CustomEmojis[e.ID] = e
	es.NextCustomEmojiID++
	return e, nil
}

func (es *ExtStore) DeleteCustomEmoji(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.CustomEmojis[id]; !ok {
		return fmt.Errorf("custom emoji not found")
	}
	delete(es.CustomEmojis, id)
	return nil
}

// ---------------------------------------------------------------------------
// FormTemplate CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetFormTemplate(id int) (*FormTemplate, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	f, ok := es.FormTemplates[id]
	if !ok {
		return nil, fmt.Errorf("form template not found")
	}
	return f, nil
}

func (es *ExtStore) ListFormTemplates() []FormTemplate {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]FormTemplate, 0, len(es.FormTemplates))
	for _, f := range es.FormTemplates {
		out = append(out, *f)
	}
	return out
}

func (es *ExtStore) CreateFormTemplate(name, template string) (*FormTemplate, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	f := &FormTemplate{
		ID: es.NextFormTemplateID, Name: name, Template: template,
		CreatedAt: now, UpdatedAt: now,
	}
	es.FormTemplates[f.ID] = f
	es.NextFormTemplateID++
	return f, nil
}

func (es *ExtStore) UpdateFormTemplate(id int, updates map[string]interface{}) (*FormTemplate, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	f, ok := es.FormTemplates[id]
	if !ok {
		return nil, fmt.Errorf("form template not found")
	}
	if v, ok := updates["name"].(string); ok {
		f.Name = v
	}
	if v, ok := updates["template"].(string); ok {
		f.Template = v
	}
	f.UpdatedAt = time.Now().UTC()
	return f, nil
}

func (es *ExtStore) DeleteFormTemplate(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.FormTemplates[id]; !ok {
		return fmt.Errorf("form template not found")
	}
	delete(es.FormTemplates, id)
	return nil
}

// ---------------------------------------------------------------------------
// AdminFlag CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetAdminFlag(id int) (*AdminFlag, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	f, ok := es.AdminFlags[id]
	if !ok {
		return nil, fmt.Errorf("admin flag not found")
	}
	return f, nil
}

func (es *ExtStore) ListAdminFlags() []AdminFlag {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]AdminFlag, 0, len(es.AdminFlags))
	for _, f := range es.AdminFlags {
		out = append(out, *f)
	}
	return out
}

func (es *ExtStore) CreateAdminFlag(name, nameKey, description string) (*AdminFlag, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	now := time.Now().UTC()
	f := &AdminFlag{
		ID: es.NextAdminFlagID, Name: name, NameKey: nameKey,
		Description: description, AppliesToPost: true, AppliesToTopic: true,
		Enabled: true, Position: len(es.AdminFlags),
		CreatedAt: now, UpdatedAt: now,
	}
	es.AdminFlags[f.ID] = f
	es.NextAdminFlagID++
	return f, nil
}

func (es *ExtStore) UpdateAdminFlag(id int, updates map[string]interface{}) (*AdminFlag, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	f, ok := es.AdminFlags[id]
	if !ok {
		return nil, fmt.Errorf("admin flag not found")
	}
	if v, ok := updates["name"].(string); ok {
		f.Name = v
	}
	if v, ok := updates["description"].(string); ok {
		f.Description = v
	}
	if v, ok := updates["enabled"].(bool); ok {
		f.Enabled = v
	}
	f.UpdatedAt = time.Now().UTC()
	return f, nil
}

func (es *ExtStore) DeleteAdminFlag(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.AdminFlags[id]; !ok {
		return fmt.Errorf("admin flag not found")
	}
	delete(es.AdminFlags, id)
	return nil
}

// ---------------------------------------------------------------------------
// PostRevision CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetPostRevision(id int) (*PostRevision, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	r, ok := es.PostRevisions[id]
	if !ok {
		return nil, fmt.Errorf("post revision not found")
	}
	return r, nil
}

func (es *ExtStore) ListPostRevisions(postID int) []PostRevision {
	es.mu.RLock()
	defer es.mu.RUnlock()
	revs := es.PostRevisionsByPost[postID]
	out := make([]PostRevision, 0, len(revs))
	for _, r := range revs {
		out = append(out, *r)
	}
	return out
}

func (es *ExtStore) CreatePostRevision(postID, userID int, previousRaw, currentRaw string) (*PostRevision, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	existing := es.PostRevisionsByPost[postID]
	number := len(existing) + 1
	r := &PostRevision{
		ID: es.NextPostRevisionID, PostID: postID, UserID: userID,
		Number: number, PreviousRaw: previousRaw, CurrentRaw: currentRaw,
		PreviousCooked: "<p>" + previousRaw + "</p>",
		CurrentCooked:  "<p>" + currentRaw + "</p>",
		CreatedAt: time.Now().UTC(),
	}
	es.PostRevisions[r.ID] = r
	es.PostRevisionsByPost[postID] = append(es.PostRevisionsByPost[postID], r)
	es.NextPostRevisionID++
	return r, nil
}

func (es *ExtStore) DeletePostRevision(id int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	r, ok := es.PostRevisions[id]
	if !ok {
		return fmt.Errorf("post revision not found")
	}
	// Remove from the per-post slice.
	revs := es.PostRevisionsByPost[r.PostID]
	for i, rev := range revs {
		if rev.ID == id {
			es.PostRevisionsByPost[r.PostID] = append(revs[:i], revs[i+1:]...)
			break
		}
	}
	delete(es.PostRevisions, id)
	return nil
}

// ---------------------------------------------------------------------------
// UserStatus CRUD
// ---------------------------------------------------------------------------

func (es *ExtStore) GetUserStatus(userID int) (*UserStatus, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	s, ok := es.UserStatuses[userID]
	if !ok {
		return nil, fmt.Errorf("user status not found")
	}
	return s, nil
}

func (es *ExtStore) ListUserStatuses() []UserStatus {
	es.mu.RLock()
	defer es.mu.RUnlock()
	out := make([]UserStatus, 0, len(es.UserStatuses))
	for _, s := range es.UserStatuses {
		out = append(out, *s)
	}
	return out
}

func (es *ExtStore) SetUserStatus(userID int, description, emoji string, endsAt *time.Time) (*UserStatus, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	s := &UserStatus{
		ID: es.NextUserStatusID, UserID: userID,
		Description: description, Emoji: emoji,
		EndsAt: endsAt, SetAt: time.Now().UTC(),
	}
	es.UserStatuses[userID] = s
	es.NextUserStatusID++
	return s, nil
}

func (es *ExtStore) DeleteUserStatus(userID int) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, ok := es.UserStatuses[userID]; !ok {
		return fmt.Errorf("user status not found")
	}
	delete(es.UserStatuses, userID)
	return nil
}

// Ensure stringPtr is used (prevents "declared and not used" if only intPtr is
// referenced during compilation).
var _ = stringPtr
