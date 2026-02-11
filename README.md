# DTU Discourse

A **Digital Twin Universe** (DTU) for [Discourse](https://github.com/discourse/discourse) — a high-fidelity behavioral clone of the Discourse REST API for integration testing.

Built following the [DTU technique](https://factory.strongdm.ai/techniques/dtu): instead of testing against a live Discourse instance (with rate limits, cost, and state management issues), this DTU replicates the API surface with wire-compatible responses so SDK clients work without modification.

## Compatibility Targets

The DTU targets **100% compatibility** with the three most popular public Discourse SDK client libraries:

| SDK | Language | Maintainer | Status |
|-----|----------|------------|--------|
| [discourse_api](https://github.com/discourse/discourse_api) | Ruby | Official (Discourse team) | All endpoints covered |
| [pydiscourse](https://github.com/pydiscourse/pydiscourse) | Python | Community | All endpoints covered |
| [discourse-api](https://github.com/luqin/discourse-api) | JavaScript | Community | All endpoints covered |

## Quick Start

### Run directly

```bash
go run ./cmd/dtu-discourse
# Listening on :4200
```

### Docker

```bash
docker compose up
# or
docker build -t dtu-discourse . && docker run -p 4200:4200 dtu-discourse
```

### Test it

```bash
curl -H "Api-Key: test_api_key" -H "Api-Username: system" http://localhost:4200/latest.json
```

## Authentication

The DTU uses the same header-based authentication as Discourse (post-April 2020):

| Header | Description |
|--------|-------------|
| `Api-Key` | API key for authentication |
| `Api-Username` | Username to act as |

### Pre-configured API keys

| Key | User | Access |
|-----|------|--------|
| `test_api_key` | system | Full access |
| `admin_api_key` | admin | Full access |

## API Coverage

Every endpoint below is implemented and tested against the SDK client libraries.

### Users
- `GET /users/{username}.json` — Get user by username
- `GET /admin/users/{id}.json` — Get user by ID (admin)
- `GET /users/by-external/{external_id}` — Get user by external ID
- `POST /users` — Create user
- `PUT /u/{username}` — Update user
- `PUT /u/{username}/preferences/email` — Update email
- `PUT /u/{username}/preferences/username` — Update username
- `PUT /admin/users/{id}/activate` — Activate user
- `PUT /admin/users/{id}/deactivate` — Deactivate user
- `PUT /admin/users/{id}/trust_level` — Update trust level
- `PUT /admin/users/{id}/grant_admin` — Grant admin
- `PUT /admin/users/{id}/revoke_admin` — Revoke admin
- `PUT /admin/users/{id}/grant_moderation` — Grant moderator
- `PUT /admin/users/{id}/revoke_moderation` — Revoke moderator
- `PUT /admin/users/{id}/suspend` — Suspend user
- `PUT /admin/users/{id}/unsuspend` — Unsuspend user
- `PUT /admin/users/{id}/anonymize` — Anonymize user
- `POST /admin/users/{id}/log_out` — Log out user
- `DELETE /admin/users/{id}.json` — Delete user
- `GET /admin/users/list/{type}.json` — List users
- `GET /users/check_username.json` — Check username availability

### Categories
- `GET /categories.json` — List categories
- `POST /categories` — Create category
- `PUT /categories/{id}` — Update category
- `DELETE /categories/{id}` — Delete category
- `GET /c/{id}/show` — Show category
- `GET /c/{slug}/{id}.json` — List topics in category
- `GET /c/{slug}/l/latest.json` — Latest topics in category
- `GET /c/{slug}/l/top.json` — Top topics in category
- `GET /c/{slug}/l/new.json` — New topics in category
- `POST /categories/reorder` — Reorder categories
- `POST /category/{id}/notifications` — Set notification level

### Topics
- `GET /latest.json` — Latest topics
- `GET /top.json` — Top topics
- `GET /new.json` — New topics
- `GET /t/{id}.json` — Get topic (with post_stream and details)
- `GET /t/external_id/{external_id}` — Get topic by external ID
- `GET /t/{id}/posts.json` — Get topic posts
- `GET /topics/created-by/{username}.json` — Topics by user
- `PUT /t/{id}.json` — Update topic (rename/recategorize)
- `PUT /t/{id}/status` — Update topic status (close/archive/pin)
- `PUT /t/{id}/change-timestamp` — Change timestamp
- `PUT /t/{id}/bookmark.json` — Bookmark topic
- `PUT /t/{id}/remove_bookmarks.json` — Remove bookmark
- `POST /t/{id}/change-owner.json` — Change post owner
- `POST /t/{id}/notifications` — Set notification level
- `POST /t/{id}/invite` — Invite to topic
- `DELETE /t/{id}.json` — Delete topic

### Posts
- `POST /posts` — Create post (or topic when title is provided)
- `GET /posts/{id}.json` — Get post
- `GET /posts.json` — List latest posts
- `PUT /posts/{id}` — Update post
- `DELETE /posts/{id}.json` — Delete post
- `PUT /posts/{id}/wiki` — Toggle wiki status
- `POST /post_actions` — Create post action (like, flag, etc.)
- `DELETE /post_actions/{id}.json` — Delete post action
- `GET /post_action_users.json` — List post action users

### Groups
- `GET /groups.json` — List groups
- `GET /groups/{name}.json` — Get group
- `POST /admin/groups` — Create group
- `PUT /groups/{id}` — Update group
- `DELETE /admin/groups/{id}.json` — Delete group
- `PUT /admin/groups/{id}/members.json` — Add members
- `DELETE /admin/groups/{id}/members.json` — Remove members
- `PUT /admin/groups/{id}/owners.json` — Add owners
- `DELETE /admin/groups/{id}/owners.json` — Remove owners
- `GET /groups/{name}/members.json` — List group members
- `POST /groups/{name}/notifications` — Set notification level

### Search
- `GET /search?q={term}` — Search topics and posts

### Tags
- `GET /tags.json` — List all tags
- `GET /tag/{tag}` — Show tag with topics

### Badges
- `GET /admin/badges.json` — List badges
- `POST /admin/badges.json` — Create badge
- `PUT /admin/badges/{id}.json` — Update badge
- `DELETE /admin/badges/{id}.json` — Delete badge
- `GET /user-badges/{username}.json` — User badges
- `POST /user_badges` — Grant badge to user

### Notifications
- `GET /notifications.json` — List notifications

### Private Messages
- `GET /topics/private-messages/{username}.json` — Inbox
- `GET /topics/private-messages-sent/{username}.json` — Sent messages
- `POST /posts` (with `archetype: private_message`) — Create PM

### Invites
- `POST /invites` — Create invite
- `GET /invites/retrieve.json` — Retrieve invite
- `PUT /invites/{id}` — Update invite
- `DELETE /invites` — Destroy invite
- `POST /invites/destroy-all-expired` — Destroy all expired
- `POST /invites/reinvite-all` — Resend all invites
- `POST /invites/reinvite` — Resend invite
- `POST /invite-token/generate` — Generate disposable token

### Uploads
- `POST /uploads` — Upload file (multipart form)

### SSO
- `POST /admin/users/sync_sso` — Sync SSO record

### Admin / Site
- `GET /admin/site_settings.json` — List site settings
- `PUT /admin/site_settings/{name}.json` — Update site setting
- `GET /admin/backups.json` — List backups
- `POST /admin/backups.json` — Create backup
- `GET /admin/dashboard.json` — Dashboard stats
- `GET /site.json` — Site info (categories, groups, notification types)
- `GET /session/csrf.json` — CSRF token

## Seed Data

The DTU starts with pre-populated data:

- **Users**: system, admin, alice, bob
- **Categories**: General, Support, Meta
- **Topics**: 3 topics with posts
- **Groups**: staff, trust_level_0
- **Tags**: welcome, intro, api, howto, plugins, help
- **Badges**: Basic, Member

## Architecture

```
cmd/dtu-discourse/     — Entry point and route registration
internal/
  model/               — Data types matching Discourse JSON response shapes
  store/               — Thread-safe in-memory data store with seed data
  middleware/           — API key authentication (header-based, post-2020 style)
  handler/             — HTTP handlers for each API resource
```

## Running Tests

```bash
go test ./... -v
```

The test suite covers 89 test cases organized by SDK compatibility target:
- `TestRuby_*` — discourse_api gem compatibility
- `TestPython_*` — pydiscourse compatibility (form-encoded payloads)
- `TestJS_*` — discourse-api JS client compatibility
- `TestLifecycle_*` — Full CRUD lifecycle integration tests

## SDK Client Usage Examples

### Ruby (discourse_api)

```ruby
require 'discourse_api'

client = DiscourseApi::Client.new("http://localhost:4200")
client.api_key = "test_api_key"
client.api_username = "system"

client.latest_topics
client.create_topic(title: "Hello", raw: "World", category: 1)
client.user("admin")
```

### Python (pydiscourse)

```python
from pydiscourse import DiscourseClient

client = DiscourseClient("http://localhost:4200", api_username="system", api_key="test_api_key")

client.latest_topics()
client.create_post("Hello World", category_id=1, content="Post body")
client.user("admin")
```

### JavaScript (discourse-api)

```javascript
const Discourse = require("discourse-api");

const client = new Discourse("http://localhost:4200", "test_api_key", "system");

client.getLatestTopics((err, body) => console.log(body));
client.getUser("admin", (err, body) => console.log(body));
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `4200` | HTTP listen port |
