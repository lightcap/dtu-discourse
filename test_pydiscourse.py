#!/usr/bin/env python3
"""Test DTU Discourse server with the real pydiscourse SDK client."""

import sys
import traceback
from pydiscourse import DiscourseClient

HOST = "http://localhost:4200"
API_KEY = "test_api_key"
API_USER = "system"

client = DiscourseClient(HOST, api_username=API_USER, api_key=API_KEY)
admin_client = DiscourseClient(HOST, api_username="admin", api_key="admin_api_key")

passed = 0
failed = 0
errors = []

def test(name, fn):
    global passed, failed
    try:
        result = fn()
        print(f"  PASS: {name}")
        if result is not None:
            s = str(result)
            if len(s) > 200:
                s = s[:200] + "..."
            print(f"        -> {s}")
        passed += 1
        return result
    except Exception as e:
        print(f"  FAIL: {name}")
        print(f"        Error: {e}")
        traceback.print_exc()
        failed += 1
        errors.append((name, str(e)))
        return None

# =============================================================================
# Categories
# =============================================================================
print("\n== Categories ==")

cats = test("List categories", lambda: client.categories())

test("Create category", lambda: client.create_category(
    name="SDK Test Category",
    color="0088CC",
    text_color="FFFFFF"
))

test("Get category (show)", lambda: client.category(1))

test("Update category", lambda: client.update_category(1, name="Updated General"))

# =============================================================================
# Topics
# =============================================================================
print("\n== Topics ==")

test("Get latest topics", lambda: client.latest_topics())
test("Get hot topics", lambda: client.hot_topics())
test("Get top topics", lambda: client.top_topics())
test("Get new topics", lambda: client.new_topics())

# Create a topic via post creation
new_topic_result = test("Create topic (via create_post)", lambda: client.create_post(
    content="This is a test topic created by pydiscourse SDK",
    title="Pydiscourse SDK Test Topic",
    category_id=1
))

topic_id = None
if new_topic_result and isinstance(new_topic_result, dict):
    topic_id = new_topic_result.get("topic_id")
    print(f"        Created topic_id={topic_id}")

if topic_id:
    # pydiscourse requires slug + topic_id
    test(f"Get topic {topic_id} (slug/id)", lambda: client.topic(
        slug="pydiscourse-sdk-test-topic",
        topic_id=topic_id
    ))

    test(f"Get topic posts for {topic_id}", lambda: client.topic_posts(topic_id=topic_id))

    # update_topic takes a full URL path
    test(f"Update topic {topic_id}", lambda: client.update_topic(
        topic_url=f"/t/-/{topic_id}",
        title="Updated Pydiscourse Topic Title"
    ))

    test(f"Update topic status (close)", lambda: client.update_topic_status(
        topic_id=topic_id,
        status="closed",
        enabled="true"
    ))

    test(f"Invite user to topic {topic_id}", lambda: client.invite_user_to_topic(
        user_email="alice@example.com",
        topic_id=topic_id
    ))

    test(f"Topic timings for {topic_id}", lambda: client.topic_timings(
        topic_id=topic_id,
        time=5000,
        timings={"1": 5000}
    ))

    test(f"Reset bump date for {topic_id}", lambda: client.reset_bump_date(
        topic_id=topic_id
    ))

# =============================================================================
# Posts
# =============================================================================
print("\n== Posts ==")

test("Get latest posts", lambda: client.latest_posts())

if topic_id:
    reply = test("Create reply post", lambda: client.create_post(
        content="This is a reply created by pydiscourse SDK",
        topic_id=topic_id
    ))

    post_id = None
    if reply and isinstance(reply, dict):
        post_id = reply.get("id")
        print(f"        Created post_id={post_id}")

    if post_id:
        # post_by_id uses GET /posts/{id}.json
        test(f"Get post by ID {post_id}", lambda: client.post_by_id(post_id=post_id))

        # post uses GET /t/{topic_id}/{post_id}.json
        test(f"Get post {post_id} in topic", lambda: client.post(
            topic_id=topic_id,
            post_id=post_id
        ))

        test(f"Update post {post_id}", lambda: client.update_post(
            post_id=post_id,
            content="Updated content from pydiscourse SDK"
        ))

        # post_by_number
        test(f"Get post by number", lambda: client.post_by_number(
            topic_id=topic_id,
            post_number=1
        ))

# =============================================================================
# Users
# =============================================================================
print("\n== Users ==")

test("Get user 'alice'", lambda: client.user("alice"))

test("Get user by external ID", lambda: client.user_by_external_id("ext-alice"))

test("List users (admin)", lambda: admin_client.list_users("active"))

test("Get user all (admin)", lambda: admin_client.user_all(2))

test("Get user emails", lambda: client.user_emails("alice"))

test("Get user badges", lambda: client.user_badges("alice"))

test("User actions", lambda: client.user_actions(
    username="alice",
    actions_filter=1
))

# =============================================================================
# Search
# =============================================================================
print("\n== Search ==")

test("Search for 'test'", lambda: client.search("test"))

# =============================================================================
# Groups
# =============================================================================
print("\n== Groups ==")

test("List groups (search)", lambda: client.groups())

test("Get group 'staff'", lambda: client.group("staff"))

test("Get group members", lambda: client.group_members("staff"))

test("Create group (admin)", lambda: admin_client.create_group(
    name="test-sdk-group"
))

# =============================================================================
# Tags
# =============================================================================
print("\n== Tags ==")

test("Create tag group", lambda: client.tag_group(
    name="SDK Test Tags",
    tag_names=["test-tag-1", "test-tag-2"]
))

# =============================================================================
# Badges
# =============================================================================
print("\n== Badges ==")

test("List badges", lambda: admin_client.badges())

# =============================================================================
# Uploads
# =============================================================================
print("\n== Uploads ==")

import tempfile, os
with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as f:
    f.write("test upload content")
    tmp_path = f.name

try:
    test("Upload file", lambda: client.upload_image(tmp_path, "image", True))
finally:
    os.unlink(tmp_path)

# =============================================================================
# Site Info
# =============================================================================
print("\n== Site Info ==")

test("Get site info", lambda: client.get_site_info())

test("About", lambda: client.about())

test("Get site settings (admin)", lambda: admin_client.get_site_settings())

test("Color schemes", lambda: admin_client.color_schemes())

# =============================================================================
# Private Messages
# =============================================================================
print("\n== Private Messages ==")

test("Create PM", lambda: client.create_post(
    content="Test private message from pydiscourse",
    title="SDK PM Test",
    archetype="private_message",
    target_recipients="alice"
))

test("Get user PMs", lambda: client.private_messages("system"))

test("Get unread PMs", lambda: client.private_messages_unread("system"))

# =============================================================================
# Admin
# =============================================================================
print("\n== Admin ==")

test("Admin: list users", lambda: admin_client.list_users("active"))

test("Admin: get user by email", lambda: admin_client.user_by_email("alice@example.com"))

# =============================================================================
# Invites
# =============================================================================
print("\n== Invites ==")

test("Create invite", lambda: client.invite(
    email="testinvite@example.com",
    group_names="staff",
    custom_message="Welcome!"
))

# =============================================================================
# Topics by user
# =============================================================================
print("\n== Topics by User ==")

test("Topics by user", lambda: client.topics_by("system"))

test("Category topics", lambda: client.category_topics(1))

test("Category latest topics", lambda: client.category_latest_topics("general"))

# =============================================================================
# Delete operations
# =============================================================================
print("\n== Delete Operations ==")

if topic_id:
    test(f"Delete topic {topic_id}", lambda: client.delete_topic(topic_id))

# =============================================================================
# Summary
# =============================================================================
print(f"\n{'='*60}")
print(f"RESULTS: {passed} passed, {failed} failed out of {passed+failed} tests")
if errors:
    print(f"\nFailed tests:")
    for name, err in errors:
        print(f"  - {name}: {err}")
print(f"{'='*60}")

sys.exit(0 if failed == 0 else 1)
