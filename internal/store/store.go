// Package store provides a thread-safe in-memory data store that mimics
// the persistent state of a Discourse instance. All data lives in RAM and
// is pre-seeded on startup so the DTU is immediately usable.
package store

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lightcap/dtu-discourse/internal/model"
)

type Store struct {
	mu sync.RWMutex

	Users         map[int]*model.User
	UsersByName   map[string]*model.User
	UsersByEmail  map[string]*model.User
	UsersByExtID  map[string]*model.User
	NextUserID    int

	Categories      map[int]*model.Category
	CategoriesBySlug map[string]*model.Category
	NextCategoryID  int

	Topics        map[int]*model.Topic
	NextTopicID   int

	Posts         map[int]*model.Post
	PostsByTopic  map[int][]*model.Post
	NextPostID    int

	Groups        map[int]*model.Group
	GroupsByName  map[string]*model.Group
	GroupMembers  map[int][]int // group_id -> user_ids
	GroupOwners   map[int][]int // group_id -> owner user_ids
	NextGroupID   int

	Tags          map[string]*model.Tag
	NextTagID     int

	Badges        map[int]*model.Badge
	UserBadges    map[int][]*model.UserBadge // user_id -> badges
	NextBadgeID   int
	NextUserBadgeID int

	Notifications map[int][]*model.Notification // user_id -> notifications
	NextNotifID   int

	Invites       map[int]*model.Invite
	NextInviteID  int

	Uploads       map[int]*model.Upload
	NextUploadID  int

	SiteSettings  map[string]*model.SiteSetting

	PostActions   map[int]*model.PostAction
	NextPostActionID int

	APIKeys       map[string]string // key -> username
}

func New() *Store {
	s := &Store{
		Users:          make(map[int]*model.User),
		UsersByName:    make(map[string]*model.User),
		UsersByEmail:   make(map[string]*model.User),
		UsersByExtID:   make(map[string]*model.User),
		Categories:     make(map[int]*model.Category),
		CategoriesBySlug: make(map[string]*model.Category),
		Topics:         make(map[int]*model.Topic),
		Posts:          make(map[int]*model.Post),
		PostsByTopic:   make(map[int][]*model.Post),
		Groups:         make(map[int]*model.Group),
		GroupsByName:   make(map[string]*model.Group),
		GroupMembers:   make(map[int][]int),
		GroupOwners:    make(map[int][]int),
		Tags:           make(map[string]*model.Tag),
		Badges:         make(map[int]*model.Badge),
		UserBadges:     make(map[int][]*model.UserBadge),
		Notifications:  make(map[int][]*model.Notification),
		Invites:        make(map[int]*model.Invite),
		Uploads:        make(map[int]*model.Upload),
		SiteSettings:   make(map[string]*model.SiteSetting),
		PostActions:    make(map[int]*model.PostAction),
		APIKeys:        make(map[string]string),
	}
	s.seed()
	return s
}

func (s *Store) seed() {
	now := time.Now().UTC()

	// --- API Keys ---
	s.APIKeys["test_api_key"] = "system"
	s.APIKeys["admin_api_key"] = "admin"

	// --- Users ---
	systemUser := &model.User{
		ID: -1, Username: "system", Name: "System",
		Email: "system@localhost", AvatarTemplate: "/letter_avatar_proxy/v4/letter/s/bcef8e/{size}.png",
		Active: true, Admin: true, Moderator: true, TrustLevel: 4,
		CreatedAt: now.Add(-365 * 24 * time.Hour), Approved: true,
	}
	admin := &model.User{
		ID: 1, Username: "admin", Name: "Admin User",
		Email: "admin@example.com", AvatarTemplate: "/letter_avatar_proxy/v4/letter/a/e9a140/{size}.png",
		Active: true, Admin: true, Moderator: true, TrustLevel: 4,
		CreatedAt: now.Add(-30 * 24 * time.Hour), Approved: true,
	}
	user1 := &model.User{
		ID: 2, Username: "alice", Name: "Alice Wonderland",
		Email: "alice@example.com", AvatarTemplate: "/letter_avatar_proxy/v4/letter/a/d0a95e/{size}.png",
		Active: true, Admin: false, Moderator: false, TrustLevel: 2,
		CreatedAt: now.Add(-20 * 24 * time.Hour), Approved: true,
		ExternalID: "ext-alice",
	}
	user2 := &model.User{
		ID: 3, Username: "bob", Name: "Bob Builder",
		Email: "bob@example.com", AvatarTemplate: "/letter_avatar_proxy/v4/letter/b/b4e14e/{size}.png",
		Active: true, Admin: false, Moderator: false, TrustLevel: 1,
		CreatedAt: now.Add(-10 * 24 * time.Hour), Approved: true,
		ExternalID: "ext-bob",
	}
	for _, u := range []*model.User{systemUser, admin, user1, user2} {
		s.Users[u.ID] = u
		s.UsersByName[u.Username] = u
		s.UsersByEmail[u.Email] = u
		if u.ExternalID != "" {
			s.UsersByExtID[u.ExternalID] = u
		}
	}
	s.NextUserID = 4

	// --- Categories ---
	cat1 := &model.Category{
		ID: 1, Name: "General", Slug: "general", Color: "0088CC", TextColor: "FFFFFF",
		Description: "General discussion", DescriptionText: "General discussion",
		TopicCount: 2, PostCount: 3, Position: 0, TopicURL: "/t/about-the-general-category/1",
		CanEdit: true, NumFeaturedTopics: 3, DefaultView: "latest",
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	cat2 := &model.Category{
		ID: 2, Name: "Support", Slug: "support", Color: "ED207B", TextColor: "FFFFFF",
		Description: "Get help here", DescriptionText: "Get help here",
		TopicCount: 1, PostCount: 1, Position: 1, TopicURL: "/t/about-the-support-category/2",
		CanEdit: true, NumFeaturedTopics: 3, DefaultView: "latest",
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	cat3 := &model.Category{
		ID: 3, Name: "Meta", Slug: "meta", Color: "808281", TextColor: "FFFFFF",
		Description: "Discussion about this site", DescriptionText: "Discussion about this site",
		TopicCount: 0, PostCount: 0, Position: 2, TopicURL: "/t/about-the-meta-category/3",
		CanEdit: true, NumFeaturedTopics: 3, DefaultView: "latest",
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	for _, c := range []*model.Category{cat1, cat2, cat3} {
		s.Categories[c.ID] = c
		s.CategoriesBySlug[c.Slug] = c
	}
	s.NextCategoryID = 4

	// --- Topics ---
	topic1 := &model.Topic{
		ID: 1, Title: "Welcome to Discourse", FancyTitle: "Welcome to Discourse",
		Slug: "welcome-to-discourse", PostsCount: 2, ReplyCount: 1, HighestPostNumber: 2,
		CreatedAt: now.Add(-25 * 24 * time.Hour), LastPostedAt: now.Add(-20 * 24 * time.Hour),
		Bumped: true, BumpedAt: now.Add(-20 * 24 * time.Hour),
		Archetype: "regular", Visible: true, Views: 42, LikeCount: 5,
		LastPosterUsername: "alice", CategoryID: 1, Tags: []string{"welcome", "intro"},
	}
	topic2 := &model.Topic{
		ID: 2, Title: "How to use the API", FancyTitle: "How to use the API",
		Slug: "how-to-use-the-api", PostsCount: 1, ReplyCount: 0, HighestPostNumber: 1,
		CreatedAt: now.Add(-15 * 24 * time.Hour), LastPostedAt: now.Add(-15 * 24 * time.Hour),
		Bumped: true, BumpedAt: now.Add(-15 * 24 * time.Hour),
		Archetype: "regular", Visible: true, Views: 15, LikeCount: 2,
		LastPosterUsername: "admin", CategoryID: 1, Tags: []string{"api", "howto"},
	}
	topic3 := &model.Topic{
		ID: 3, Title: "Need help with plugins", FancyTitle: "Need help with plugins",
		Slug: "need-help-with-plugins", PostsCount: 1, ReplyCount: 0, HighestPostNumber: 1,
		CreatedAt: now.Add(-5 * 24 * time.Hour), LastPostedAt: now.Add(-5 * 24 * time.Hour),
		Bumped: true, BumpedAt: now.Add(-5 * 24 * time.Hour),
		Archetype: "regular", Visible: true, Views: 8, LikeCount: 0,
		LastPosterUsername: "bob", CategoryID: 2, Tags: []string{"plugins", "help"},
	}
	for _, t := range []*model.Topic{topic1, topic2, topic3} {
		s.Topics[t.ID] = t
	}
	s.NextTopicID = 4

	// --- Posts ---
	post1 := &model.Post{
		ID: 1, Username: "admin", Name: "Admin User",
		AvatarTemplate: admin.AvatarTemplate,
		CreatedAt: now.Add(-25 * 24 * time.Hour), UpdatedAt: now.Add(-25 * 24 * time.Hour),
		Raw: "Welcome to Discourse! This is your first topic.", Cooked: "<p>Welcome to Discourse! This is your first topic.</p>",
		PostNumber: 1, PostType: 1, TopicID: 1, TopicSlug: "welcome-to-discourse",
		DisplayUsername: "Admin User", Version: 1, UserID: 1, TrustLevel: 4,
		CanEdit: true, CanDelete: true, CanWiki: true,
	}
	post2 := &model.Post{
		ID: 2, Username: "alice", Name: "Alice Wonderland",
		AvatarTemplate: user1.AvatarTemplate,
		CreatedAt: now.Add(-20 * 24 * time.Hour), UpdatedAt: now.Add(-20 * 24 * time.Hour),
		Raw: "Thanks for the warm welcome!", Cooked: "<p>Thanks for the warm welcome!</p>",
		PostNumber: 2, PostType: 1, TopicID: 1, TopicSlug: "welcome-to-discourse",
		DisplayUsername: "Alice Wonderland", Version: 1, UserID: 2, TrustLevel: 2,
		ReplyCount: 0, CanEdit: true, CanDelete: true, CanWiki: true,
	}
	post3 := &model.Post{
		ID: 3, Username: "admin", Name: "Admin User",
		AvatarTemplate: admin.AvatarTemplate,
		CreatedAt: now.Add(-15 * 24 * time.Hour), UpdatedAt: now.Add(-15 * 24 * time.Hour),
		Raw: "Here is a guide on using the Discourse API.", Cooked: "<p>Here is a guide on using the Discourse API.</p>",
		PostNumber: 1, PostType: 1, TopicID: 2, TopicSlug: "how-to-use-the-api",
		DisplayUsername: "Admin User", Version: 1, UserID: 1, TrustLevel: 4,
		CanEdit: true, CanDelete: true, CanWiki: true,
	}
	post4 := &model.Post{
		ID: 4, Username: "bob", Name: "Bob Builder",
		AvatarTemplate: user2.AvatarTemplate,
		CreatedAt: now.Add(-5 * 24 * time.Hour), UpdatedAt: now.Add(-5 * 24 * time.Hour),
		Raw: "Can someone help me install a plugin?", Cooked: "<p>Can someone help me install a plugin?</p>",
		PostNumber: 1, PostType: 1, TopicID: 3, TopicSlug: "need-help-with-plugins",
		DisplayUsername: "Bob Builder", Version: 1, UserID: 3, TrustLevel: 1,
		CanEdit: true, CanDelete: true, CanWiki: true,
	}
	for _, p := range []*model.Post{post1, post2, post3, post4} {
		s.Posts[p.ID] = p
		s.PostsByTopic[p.TopicID] = append(s.PostsByTopic[p.TopicID], p)
	}
	s.NextPostID = 5

	// --- Groups ---
	staffGroup := &model.Group{
		ID: 1, Name: "staff", DisplayName: "Staff", Automatic: true,
		UserCount: 1, VisibilityLevel: 2, PrimaryGroup: false,
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	trustLevel0 := &model.Group{
		ID: 10, Name: "trust_level_0", DisplayName: "Trust Level 0", Automatic: true,
		UserCount: 3, VisibilityLevel: 0, PrimaryGroup: false,
		CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now,
	}
	for _, g := range []*model.Group{staffGroup, trustLevel0} {
		s.Groups[g.ID] = g
		s.GroupsByName[g.Name] = g
	}
	s.GroupMembers[1] = []int{1}
	s.GroupMembers[10] = []int{1, 2, 3}
	s.NextGroupID = 11

	// --- Tags ---
	for _, name := range []string{"welcome", "intro", "api", "howto", "plugins", "help"} {
		s.NextTagID++
		s.Tags[name] = &model.Tag{ID: s.NextTagID, TagName: name, Name: name, Count: 1}
	}

	// --- Badges ---
	s.Badges[1] = &model.Badge{
		ID: 1, Name: "Basic", Description: "Granted when you complete the basic tutorial",
		GrantCount: 2, AllowTitle: false, MultipleGrant: false, Icon: "certificate",
		Listable: true, Enabled: true, BadgeGroupingID: 1, System: true, BadgeTypeID: 3,
	}
	s.Badges[2] = &model.Badge{
		ID: 2, Name: "Member", Description: "Granted when you reach trust level 2",
		GrantCount: 1, AllowTitle: false, MultipleGrant: false, Icon: "certificate",
		Listable: true, Enabled: true, BadgeGroupingID: 1, System: true, BadgeTypeID: 3,
	}
	s.NextBadgeID = 3

	// --- Site Settings ---
	defaults := map[string]interface{}{
		"title":              "DTU Discourse",
		"site_description":   "Digital Twin Universe for Discourse",
		"allow_user_locale":  true,
		"default_locale":     "en",
		"min_topic_title_length": 5,
		"max_topic_title_length": 255,
		"min_post_length":    10,
		"max_post_length":    32000,
		"tagging_enabled":    true,
		"max_tags_per_topic": 5,
	}
	for k, v := range defaults {
		s.SiteSettings[k] = &model.SiteSetting{Setting: k, Value: v, Default: v}
	}
}

// ---------- User Operations ----------

func (s *Store) GetUser(id int) *model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Users[id]
}

func (s *Store) GetUserByUsername(username string) *model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.UsersByName[strings.ToLower(username)]
}

func (s *Store) GetUserByExternalID(extID string) *model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.UsersByExtID[extID]
}

func (s *Store) CreateUser(name, username, email, password string) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lower := strings.ToLower(username)
	if _, exists := s.UsersByName[lower]; exists {
		return nil, fmt.Errorf("username already taken")
	}
	if _, exists := s.UsersByEmail[email]; exists {
		return nil, fmt.Errorf("email already in use")
	}
	u := &model.User{
		ID: s.NextUserID, Username: username, Name: name,
		Email: email, Active: true, TrustLevel: 0, Approved: true,
		AvatarTemplate: fmt.Sprintf("/letter_avatar_proxy/v4/letter/%s/b4e14e/{size}.png", string(lower[0])),
		CreatedAt: time.Now().UTC(),
	}
	s.Users[u.ID] = u
	s.UsersByName[lower] = u
	s.UsersByEmail[email] = u
	s.NextUserID++
	return u, nil
}

func (s *Store) UpdateUser(id int, updates map[string]interface{}) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.Users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	if v, ok := updates["name"].(string); ok {
		u.Name = v
	}
	if v, ok := updates["title"].(string); ok {
		u.Title = v
	}
	if v, ok := updates["trust_level"].(float64); ok {
		u.TrustLevel = int(v)
	}
	if v, ok := updates["active"].(bool); ok {
		u.Active = v
	}
	if v, ok := updates["admin"].(bool); ok {
		u.Admin = v
	}
	if v, ok := updates["moderator"].(bool); ok {
		u.Moderator = v
	}
	if v, ok := updates["suspended"].(bool); ok {
		u.Suspended = v
	}
	return u, nil
}

func (s *Store) DeleteUser(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.Users[id]
	if !ok {
		return fmt.Errorf("user not found")
	}
	delete(s.Users, id)
	delete(s.UsersByName, strings.ToLower(u.Username))
	delete(s.UsersByEmail, u.Email)
	if u.ExternalID != "" {
		delete(s.UsersByExtID, u.ExternalID)
	}
	return nil
}

func (s *Store) ListUsers(listType string) []model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.User
	for _, u := range s.Users {
		if u.ID < 0 {
			continue
		}
		switch listType {
		case "active":
			if u.Active {
				result = append(result, *u)
			}
		case "new":
			result = append(result, *u)
		case "staff":
			if u.Admin || u.Moderator {
				result = append(result, *u)
			}
		case "suspended":
			if u.Suspended {
				result = append(result, *u)
			}
		default:
			result = append(result, *u)
		}
	}
	return result
}

// ListAllUsers returns every user including system.
func (s *Store) ListAllUsers() []model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.User, 0, len(s.Users))
	for _, u := range s.Users {
		result = append(result, *u)
	}
	return result
}

// ---------- Category Operations ----------

func (s *Store) GetCategory(id int) *model.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Categories[id]
}

func (s *Store) GetCategoryBySlug(slug string) *model.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CategoriesBySlug[slug]
}

func (s *Store) ListCategories() []model.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Category, 0, len(s.Categories))
	for _, c := range s.Categories {
		result = append(result, *c)
	}
	return result
}

func (s *Store) CreateCategory(name, slug, color, textColor string) (*model.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}
	if _, exists := s.CategoriesBySlug[slug]; exists {
		return nil, fmt.Errorf("category slug already exists")
	}
	now := time.Now().UTC()
	c := &model.Category{
		ID: s.NextCategoryID, Name: name, Slug: slug,
		Color: color, TextColor: textColor,
		Position: len(s.Categories), CanEdit: true,
		DefaultView: "latest", NumFeaturedTopics: 3,
		CreatedAt: now, UpdatedAt: now,
	}
	if c.Color == "" {
		c.Color = "0088CC"
	}
	if c.TextColor == "" {
		c.TextColor = "FFFFFF"
	}
	s.Categories[c.ID] = c
	s.CategoriesBySlug[slug] = c
	s.NextCategoryID++
	return c, nil
}

func (s *Store) UpdateCategory(id int, updates map[string]interface{}) (*model.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.Categories[id]
	if !ok {
		return nil, fmt.Errorf("category not found")
	}
	if v, ok := updates["name"].(string); ok {
		c.Name = v
	}
	if v, ok := updates["slug"].(string); ok {
		delete(s.CategoriesBySlug, c.Slug)
		c.Slug = v
		s.CategoriesBySlug[v] = c
	}
	if v, ok := updates["color"].(string); ok {
		c.Color = v
	}
	if v, ok := updates["text_color"].(string); ok {
		c.TextColor = v
	}
	if v, ok := updates["description"].(string); ok {
		c.Description = v
		c.DescriptionText = v
	}
	c.UpdatedAt = time.Now().UTC()
	return c, nil
}

func (s *Store) DeleteCategory(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.Categories[id]
	if !ok {
		return fmt.Errorf("category not found")
	}
	delete(s.Categories, id)
	delete(s.CategoriesBySlug, c.Slug)
	return nil
}

// ---------- Topic Operations ----------

func (s *Store) GetTopic(id int) *model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.Topics[id]
	if !ok {
		return nil
	}
	cp := *t
	posts := s.PostsByTopic[id]
	stream := &model.PostStream{
		Posts: make([]model.Post, 0, len(posts)),
		Stream: make([]int, 0, len(posts)),
	}
	for _, p := range posts {
		stream.Posts = append(stream.Posts, *p)
		stream.Stream = append(stream.Stream, p.ID)
	}
	cp.PostStream = stream
	if len(posts) > 0 {
		first := posts[0]
		u := s.Users[first.UserID]
		creator := model.BasicUser{ID: first.UserID, Username: first.Username, AvatarTemplate: first.AvatarTemplate}
		if u != nil {
			creator.Name = u.Name
		}
		cp.Details = &model.TopicDetails{
			CreatedBy:      creator,
			LastPoster:     creator,
			CanEdit:        true,
			CanInviteTo:    true,
			CanCreatePost:  true,
			CanReplyAsNewTopic: true,
			CanFlagTopic:   true,
			NotificationLevel: 1,
		}
	}
	return &cp
}

func (s *Store) ListTopics(filter string) []model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Topic, 0, len(s.Topics))
	for _, t := range s.Topics {
		if t.Archetype == "private_message" {
			continue
		}
		result = append(result, *t)
	}
	return result
}

func (s *Store) TopicsByCategory(categoryID int) []model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Topic
	for _, t := range s.Topics {
		if t.CategoryID == categoryID {
			result = append(result, *t)
		}
	}
	return result
}

func (s *Store) TopicsByUser(username string) []model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Topic
	for _, t := range s.Topics {
		posts := s.PostsByTopic[t.ID]
		if len(posts) > 0 && posts[0].Username == username {
			result = append(result, *t)
		}
	}
	return result
}

func (s *Store) TopicsByTag(tag string) []model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Topic
	for _, t := range s.Topics {
		for _, tg := range t.Tags {
			if tg == tag {
				result = append(result, *t)
				break
			}
		}
	}
	return result
}

func (s *Store) GetTopicByExternalID(extID string) *model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.Topics {
		if t.ExternalID == extID {
			return t
		}
	}
	return nil
}

func (s *Store) CreateTopic(title, raw string, categoryID, userID int, tags []string, archetype string) (*model.Topic, *model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := s.Users[userID]
	if u == nil {
		return nil, nil, fmt.Errorf("user not found")
	}
	now := time.Now().UTC()
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	if archetype == "" {
		archetype = "regular"
	}
	t := &model.Topic{
		ID: s.NextTopicID, Title: title, FancyTitle: title, Slug: slug,
		PostsCount: 1, ReplyCount: 0, HighestPostNumber: 1,
		CreatedAt: now, LastPostedAt: now, Bumped: true, BumpedAt: now,
		Archetype: archetype, Visible: true, CategoryID: categoryID,
		LastPosterUsername: u.Username, Tags: tags,
	}
	if tags == nil {
		t.Tags = []string{}
	}
	s.Topics[t.ID] = t
	s.NextTopicID++

	p := &model.Post{
		ID: s.NextPostID, Username: u.Username, Name: u.Name,
		AvatarTemplate: u.AvatarTemplate,
		CreatedAt: now, UpdatedAt: now, Raw: raw,
		Cooked: "<p>" + raw + "</p>",
		PostNumber: 1, PostType: 1, TopicID: t.ID, TopicSlug: slug,
		DisplayUsername: u.Name, Version: 1, UserID: userID,
		TrustLevel: u.TrustLevel, CanEdit: true, CanDelete: true, CanWiki: true,
	}
	s.Posts[p.ID] = p
	s.PostsByTopic[t.ID] = append(s.PostsByTopic[t.ID], p)
	s.NextPostID++

	if cat, ok := s.Categories[categoryID]; ok {
		cat.TopicCount++
		cat.PostCount++
	}

	return t, p, nil
}

func (s *Store) UpdateTopic(id int, updates map[string]interface{}) (*model.Topic, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.Topics[id]
	if !ok {
		return nil, fmt.Errorf("topic not found")
	}
	if v, ok := updates["title"].(string); ok {
		t.Title = v
		t.FancyTitle = v
	}
	if v, ok := updates["category_id"].(float64); ok {
		t.CategoryID = int(v)
	}
	if v, ok := updates["visible"].(bool); ok {
		t.Visible = v
	}
	return t, nil
}

func (s *Store) DeleteTopic(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.Topics[id]
	if !ok {
		return fmt.Errorf("topic not found")
	}
	if cat, catOk := s.Categories[t.CategoryID]; catOk {
		cat.TopicCount--
		cat.PostCount -= len(s.PostsByTopic[id])
	}
	for _, p := range s.PostsByTopic[id] {
		delete(s.Posts, p.ID)
	}
	delete(s.PostsByTopic, id)
	delete(s.Topics, id)
	return nil
}

func (s *Store) UpdateTopicStatus(id int, status string, enabled bool) (*model.Topic, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.Topics[id]
	if !ok {
		return nil, fmt.Errorf("topic not found")
	}
	switch status {
	case "closed":
		t.Closed = enabled
	case "archived":
		t.Archived = enabled
	case "pinned":
		t.Pinned = enabled
	case "visible":
		t.Visible = enabled
	case "pinned_globally":
		t.PinnedGlobally = enabled
	}
	return t, nil
}

// ---------- Post Operations ----------

func (s *Store) GetPost(id int) *model.Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Posts[id]
}

func (s *Store) ListPosts() []model.Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Post, 0, len(s.Posts))
	for _, p := range s.Posts {
		result = append(result, *p)
	}
	return result
}

func (s *Store) GetTopicPosts(topicID int, postIDs []int) []model.Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	posts := s.PostsByTopic[topicID]
	if len(postIDs) == 0 {
		result := make([]model.Post, 0, len(posts))
		for _, p := range posts {
			result = append(result, *p)
		}
		return result
	}
	idSet := make(map[int]bool, len(postIDs))
	for _, id := range postIDs {
		idSet[id] = true
	}
	var result []model.Post
	for _, p := range posts {
		if idSet[p.ID] {
			result = append(result, *p)
		}
	}
	return result
}

func (s *Store) CreatePost(topicID int, raw string, userID int, replyTo *int) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.Topics[topicID]
	if !ok {
		return nil, fmt.Errorf("topic not found")
	}
	u := s.Users[userID]
	if u == nil {
		return nil, fmt.Errorf("user not found")
	}
	now := time.Now().UTC()
	t.PostsCount++
	t.HighestPostNumber++
	t.ReplyCount++
	t.LastPostedAt = now
	t.BumpedAt = now
	t.LastPosterUsername = u.Username

	p := &model.Post{
		ID: s.NextPostID, Username: u.Username, Name: u.Name,
		AvatarTemplate: u.AvatarTemplate,
		CreatedAt: now, UpdatedAt: now, Raw: raw,
		Cooked: "<p>" + raw + "</p>",
		PostNumber: t.HighestPostNumber, PostType: 1, TopicID: topicID,
		TopicSlug: t.Slug, DisplayUsername: u.Name, Version: 1,
		UserID: userID, TrustLevel: u.TrustLevel,
		ReplyToPostNumber: replyTo,
		CanEdit: true, CanDelete: true, CanWiki: true,
	}
	s.Posts[p.ID] = p
	s.PostsByTopic[topicID] = append(s.PostsByTopic[topicID], p)
	s.NextPostID++

	if cat, catOk := s.Categories[t.CategoryID]; catOk {
		cat.PostCount++
	}

	return p, nil
}

func (s *Store) UpdatePost(id int, raw string) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.Posts[id]
	if !ok {
		return nil, fmt.Errorf("post not found")
	}
	p.Raw = raw
	p.Cooked = "<p>" + raw + "</p>"
	p.Version++
	p.UpdatedAt = time.Now().UTC()
	return p, nil
}

func (s *Store) DeletePost(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.Posts[id]
	if !ok {
		return fmt.Errorf("post not found")
	}
	if t, tOk := s.Topics[p.TopicID]; tOk {
		t.PostsCount--
		if cat, catOk := s.Categories[t.CategoryID]; catOk {
			cat.PostCount--
		}
	}
	posts := s.PostsByTopic[p.TopicID]
	for i, tp := range posts {
		if tp.ID == id {
			s.PostsByTopic[p.TopicID] = append(posts[:i], posts[i+1:]...)
			break
		}
	}
	delete(s.Posts, id)
	return nil
}

func (s *Store) WikifyPost(id int, wiki bool) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.Posts[id]
	if !ok {
		return nil, fmt.Errorf("post not found")
	}
	p.Wiki = wiki
	return p, nil
}

// ---------- Group Operations ----------

func (s *Store) GetGroup(nameOrID interface{}) *model.Group {
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch v := nameOrID.(type) {
	case string:
		return s.GroupsByName[v]
	case int:
		return s.Groups[v]
	}
	return nil
}

func (s *Store) ListGroups() []model.Group {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Group, 0, len(s.Groups))
	for _, g := range s.Groups {
		result = append(result, *g)
	}
	return result
}

func (s *Store) CreateGroup(name string, attrs map[string]interface{}) (*model.Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.GroupsByName[name]; exists {
		return nil, fmt.Errorf("group name already exists")
	}
	now := time.Now().UTC()
	g := &model.Group{
		ID: s.NextGroupID, Name: name, DisplayName: name,
		CreatedAt: now, UpdatedAt: now,
	}
	if v, ok := attrs["visibility_level"].(float64); ok {
		g.VisibilityLevel = int(v)
	}
	if v, ok := attrs["full_name"].(string); ok {
		g.FullName = v
	}
	s.Groups[g.ID] = g
	s.GroupsByName[name] = g
	s.GroupMembers[g.ID] = []int{}
	s.NextGroupID++
	return g, nil
}

func (s *Store) UpdateGroup(id int, updates map[string]interface{}) (*model.Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, ok := s.Groups[id]
	if !ok {
		return nil, fmt.Errorf("group not found")
	}
	if v, ok := updates["name"].(string); ok {
		delete(s.GroupsByName, g.Name)
		g.Name = v
		s.GroupsByName[v] = g
	}
	if v, ok := updates["full_name"].(string); ok {
		g.FullName = v
	}
	g.UpdatedAt = time.Now().UTC()
	return g, nil
}

func (s *Store) DeleteGroup(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, ok := s.Groups[id]
	if !ok {
		return fmt.Errorf("group not found")
	}
	delete(s.Groups, id)
	delete(s.GroupsByName, g.Name)
	delete(s.GroupMembers, id)
	delete(s.GroupOwners, id)
	return nil
}

func (s *Store) AddGroupMembers(groupID int, userIDs []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, ok := s.Groups[groupID]
	if !ok {
		return fmt.Errorf("group not found")
	}
	members := s.GroupMembers[groupID]
	existing := make(map[int]bool, len(members))
	for _, id := range members {
		existing[id] = true
	}
	for _, id := range userIDs {
		if !existing[id] {
			members = append(members, id)
			g.UserCount++
		}
	}
	s.GroupMembers[groupID] = members
	return nil
}

func (s *Store) RemoveGroupMembers(groupID int, userIDs []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Groups[groupID]; !ok {
		return fmt.Errorf("group not found")
	}
	toRemove := make(map[int]bool, len(userIDs))
	for _, id := range userIDs {
		toRemove[id] = true
	}
	members := s.GroupMembers[groupID]
	var kept []int
	for _, id := range members {
		if !toRemove[id] {
			kept = append(kept, id)
		}
	}
	s.GroupMembers[groupID] = kept
	s.Groups[groupID].UserCount = len(kept)
	return nil
}

func (s *Store) AddGroupOwners(groupID int, userIDs []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Groups[groupID]; !ok {
		return fmt.Errorf("group not found")
	}
	owners := s.GroupOwners[groupID]
	existing := make(map[int]bool, len(owners))
	for _, id := range owners {
		existing[id] = true
	}
	for _, id := range userIDs {
		if !existing[id] {
			owners = append(owners, id)
		}
	}
	s.GroupOwners[groupID] = owners
	return nil
}

func (s *Store) RemoveGroupOwners(groupID int, userIDs []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Groups[groupID]; !ok {
		return fmt.Errorf("group not found")
	}
	toRemove := make(map[int]bool, len(userIDs))
	for _, id := range userIDs {
		toRemove[id] = true
	}
	owners := s.GroupOwners[groupID]
	var kept []int
	for _, id := range owners {
		if !toRemove[id] {
			kept = append(kept, id)
		}
	}
	s.GroupOwners[groupID] = kept
	return nil
}

func (s *Store) GetGroupMembers(groupID int) ([]model.BasicUser, []model.BasicUser) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var members, owners []model.BasicUser
	for _, uid := range s.GroupMembers[groupID] {
		if u, ok := s.Users[uid]; ok {
			members = append(members, model.BasicUser{ID: u.ID, Username: u.Username, Name: u.Name, AvatarTemplate: u.AvatarTemplate})
		}
	}
	for _, uid := range s.GroupOwners[groupID] {
		if u, ok := s.Users[uid]; ok {
			owners = append(owners, model.BasicUser{ID: u.ID, Username: u.Username, Name: u.Name, AvatarTemplate: u.AvatarTemplate})
		}
	}
	return members, owners
}

// ---------- Search Operations ----------

func (s *Store) Search(term string) model.SearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	lower := strings.ToLower(term)
	var posts []model.Post
	var topics []model.Topic
	for _, p := range s.Posts {
		if strings.Contains(strings.ToLower(p.Raw), lower) || strings.Contains(strings.ToLower(p.Cooked), lower) {
			posts = append(posts, *p)
		}
	}
	for _, t := range s.Topics {
		if strings.Contains(strings.ToLower(t.Title), lower) {
			topics = append(topics, *t)
		}
	}
	return model.SearchResult{Posts: posts, Topics: topics}
}

// ---------- Tag Operations ----------

func (s *Store) GetTag(name string) *model.Tag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Tags[name]
}

func (s *Store) ListTags() []model.Tag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Tag, 0, len(s.Tags))
	for _, t := range s.Tags {
		result = append(result, *t)
	}
	return result
}

// ---------- Badge Operations ----------

func (s *Store) ListBadges() []model.Badge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Badge, 0, len(s.Badges))
	for _, b := range s.Badges {
		result = append(result, *b)
	}
	return result
}

func (s *Store) CreateBadge(name, description string, badgeTypeID int) (*model.Badge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b := &model.Badge{
		ID: s.NextBadgeID, Name: name, Description: description,
		BadgeTypeID: badgeTypeID, Enabled: true, Listable: true, System: false,
	}
	s.Badges[b.ID] = b
	s.NextBadgeID++
	return b, nil
}

func (s *Store) UpdateBadge(id int, updates map[string]interface{}) (*model.Badge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.Badges[id]
	if !ok {
		return nil, fmt.Errorf("badge not found")
	}
	if v, ok := updates["name"].(string); ok {
		b.Name = v
	}
	if v, ok := updates["description"].(string); ok {
		b.Description = v
	}
	return b, nil
}

func (s *Store) DeleteBadge(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Badges[id]; !ok {
		return fmt.Errorf("badge not found")
	}
	delete(s.Badges, id)
	return nil
}

func (s *Store) GrantUserBadge(userID, badgeID int) (*model.UserBadge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Badges[badgeID]; !ok {
		return nil, fmt.Errorf("badge not found")
	}
	ub := &model.UserBadge{
		ID: s.NextUserBadgeID, GrantedAt: time.Now().UTC(),
		BadgeID: badgeID, UserID: userID, GrantedByID: 1,
	}
	s.UserBadges[userID] = append(s.UserBadges[userID], ub)
	s.Badges[badgeID].GrantCount++
	s.NextUserBadgeID++
	return ub, nil
}

func (s *Store) GetUserBadges(username string) ([]model.Badge, []model.UserBadge) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u := s.UsersByName[strings.ToLower(username)]
	if u == nil {
		return nil, nil
	}
	ubs := s.UserBadges[u.ID]
	var badges []model.Badge
	seen := make(map[int]bool)
	var userBadges []model.UserBadge
	for _, ub := range ubs {
		userBadges = append(userBadges, *ub)
		if !seen[ub.BadgeID] {
			if b, ok := s.Badges[ub.BadgeID]; ok {
				badges = append(badges, *b)
				seen[ub.BadgeID] = true
			}
		}
	}
	return badges, userBadges
}

// ---------- Notification Operations ----------

func (s *Store) GetNotifications(userID int) []model.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	notifs := s.Notifications[userID]
	result := make([]model.Notification, 0, len(notifs))
	for _, n := range notifs {
		result = append(result, *n)
	}
	return result
}

// ---------- Invite Operations ----------

func (s *Store) CreateInvite(email string, groupIDs []int, topicID *int) (*model.Invite, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	inv := &model.Invite{
		ID: s.NextInviteID, Email: email,
		Link: fmt.Sprintf("/invites/%d", s.NextInviteID),
		MaxRedemptionsAllowed: 1, CreatedAt: now, UpdatedAt: now,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
	}
	s.Invites[inv.ID] = inv
	s.NextInviteID++
	return inv, nil
}

// ---------- Upload Operations ----------

func (s *Store) CreateUpload(filename, ext string, filesize int) *model.Upload {
	s.mu.Lock()
	defer s.mu.Unlock()
	up := &model.Upload{
		ID: s.NextUploadID, OriginalFilename: filename, Extension: ext,
		Filesize: filesize,
		URL: fmt.Sprintf("/uploads/default/original/1X/%d.%s", s.NextUploadID, ext),
		ShortURL: fmt.Sprintf("upload://%d.%s", s.NextUploadID, ext),
		ShortPath: fmt.Sprintf("/uploads/short-url/%d.%s", s.NextUploadID, ext),
		HumanFilesize: fmt.Sprintf("%d Bytes", filesize),
	}
	s.Uploads[up.ID] = up
	s.NextUploadID++
	return up
}

// ---------- Site Settings Operations ----------

func (s *Store) GetSiteSettings() []model.SiteSetting {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.SiteSetting, 0, len(s.SiteSettings))
	for _, ss := range s.SiteSettings {
		result = append(result, *ss)
	}
	return result
}

func (s *Store) UpdateSiteSetting(name string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ss, ok := s.SiteSettings[name]; ok {
		ss.Value = value
	} else {
		s.SiteSettings[name] = &model.SiteSetting{Setting: name, Value: value}
	}
}

// ---------- Post Action Operations ----------

func (s *Store) CreatePostAction(postID, actionTypeID int) (*model.PostAction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Posts[postID]; !ok {
		return nil, fmt.Errorf("post not found")
	}
	pa := &model.PostAction{
		ID: s.NextPostActionID, PostID: postID, PostActionTypeID: actionTypeID,
	}
	s.PostActions[pa.ID] = pa
	s.NextPostActionID++

	if actionTypeID == 2 { // like
		s.Posts[postID].Score++
		if t, ok := s.Topics[s.Posts[postID].TopicID]; ok {
			t.LikeCount++
		}
	}
	return pa, nil
}

func (s *Store) DeletePostAction(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.PostActions[id]; !ok {
		return fmt.Errorf("post action not found")
	}
	delete(s.PostActions, id)
	return nil
}

// ---------- Private Message Operations ----------

func (s *Store) GetPrivateMessages(username string) []model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Topic
	for _, t := range s.Topics {
		if t.Archetype == "private_message" {
			posts := s.PostsByTopic[t.ID]
			for _, p := range posts {
				if p.Username == username {
					result = append(result, *t)
					break
				}
			}
		}
	}
	return result
}

func (s *Store) GetSentPrivateMessages(username string) []model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Topic
	for _, t := range s.Topics {
		if t.Archetype == "private_message" {
			posts := s.PostsByTopic[t.ID]
			if len(posts) > 0 && posts[0].Username == username {
				result = append(result, *t)
			}
		}
	}
	return result
}

// ---------- SSO Operations ----------

func (s *Store) SyncSSO(externalID, email, username, name string) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if u, ok := s.UsersByExtID[externalID]; ok {
		u.Email = email
		u.Username = username
		u.Name = name
		return u, nil
	}

	lower := strings.ToLower(username)
	u := &model.User{
		ID: s.NextUserID, Username: username, Name: name, Email: email,
		ExternalID: externalID, Active: true, TrustLevel: 0, Approved: true,
		AvatarTemplate: fmt.Sprintf("/letter_avatar_proxy/v4/letter/%s/b4e14e/{size}.png", string(lower[0])),
		CreatedAt: time.Now().UTC(),
	}
	s.Users[u.ID] = u
	s.UsersByName[lower] = u
	s.UsersByEmail[email] = u
	s.UsersByExtID[externalID] = u
	s.NextUserID++
	return u, nil
}

// ---------- Auth ----------

func (s *Store) ValidateAPIKey(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	username, ok := s.APIKeys[key]
	return username, ok
}
