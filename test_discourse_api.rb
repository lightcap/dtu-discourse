#!/usr/bin/env ruby
# Test DTU Discourse server with the official Ruby discourse_api gem.

require 'discourse_api'

HOST = "http://localhost:4200"
API_KEY = "test_api_key"
API_USER = "system"
ADMIN_KEY = "admin_api_key"
ADMIN_USER = "admin"

client = DiscourseApi::Client.new(HOST)
client.api_key = API_KEY
client.api_username = API_USER

admin = DiscourseApi::Client.new(HOST)
admin.api_key = ADMIN_KEY
admin.api_username = ADMIN_USER

$passed = 0
$failed = 0
$errors = []

def test(name)
  result = yield
  puts "  PASS: #{name}"
  snippet = result.to_s
  puts "        -> #{snippet[0..199]}..." if snippet.length > 10
  $passed += 1
  result
rescue => e
  puts "  FAIL: #{name}"
  puts "        Error: #{e.message}"
  puts "        #{e.backtrace&.first}"
  $failed += 1
  $errors << [name, e.message]
  nil
end

# =============================================================================
# Categories
# =============================================================================
puts "\n== Categories =="

test("List categories") { client.categories }

test("Create category") do
  client.create_category(name: "Ruby SDK Test", color: "FF0000", text_color: "FFFFFF")
end

test("Get category") { client.category(1) }

test("Update category") { client.update_category(id: 1, name: "Updated General") }

# =============================================================================
# Topics
# =============================================================================
puts "\n== Topics =="

test("Get latest topics") { client.latest_topics }
test("Get top topics") { client.top_topics }
test("Get new topics") { client.new_topics }

topic_result = test("Create topic") do
  client.create_topic(
    title: "Ruby SDK Test Topic",
    raw: "This is a topic created by the Ruby discourse_api gem",
    category: 1
  )
end

topic_id = topic_result && topic_result["topic_id"]
puts "        Created topic_id=#{topic_id}" if topic_id

if topic_id
  test("Get topic #{topic_id}") { client.topic(topic_id) }

  test("Get topic posts #{topic_id}") { client.topic_posts(topic_id) }

  test("Rename topic #{topic_id}") do
    client.rename_topic(topic_id, "Updated Ruby SDK Topic Title")
  end

  test("Update topic status (close)") do
    client.update_topic_status(topic_id, status: "closed", enabled: "true")
  end

  test("Bookmark topic #{topic_id}") { client.bookmark_topic(topic_id) }

  test("Remove topic bookmark #{topic_id}") { client.remove_topic_bookmark(topic_id) }

  test("Edit topic timestamp #{topic_id}") do
    client.edit_topic_timestamp(topic_id, Time.now.to_i)
  end

  test("Set topic notification level") do
    client.topic_set_user_notification_level(topic_id, notification_level: 3)
  end

  test("Topics by user") { client.topics_by("system") }
end

# =============================================================================
# Posts
# =============================================================================
puts "\n== Posts =="

test("Get latest posts") { client.posts }

post_id = nil
if topic_id
  reply = test("Create reply post") do
    client.create_post(
      topic_id: topic_id,
      raw: "This is a reply created by the Ruby discourse_api gem"
    )
  end

  post_id = reply && reply["id"]
  puts "        Created post_id=#{post_id}" if post_id

  if post_id
    test("Get post #{post_id}") { client.get_post(post_id) }

    test("Edit post #{post_id}") do
      client.edit_post(post_id, "Updated content from Ruby SDK")
    end

    test("Wikify post #{post_id}") { client.wikify_post(post_id) }
  end
end

# =============================================================================
# Users
# =============================================================================
puts "\n== Users =="

test("Get user 'alice'") { client.user("alice") }

test("Get user by external ID") { client.by_external_id("ext-alice") }

test("Check username") { client.check_username("newuser123") }

test("List users (admin)") { admin.list_users("active") }

test("User SSO (admin)") { admin.user_sso(2) }

test("Create user") do
  client.create_user(
    name: "Ruby Test User",
    email: "rubytest@example.com",
    password: "SecureP@ss123",
    username: "rubytest"
  )
end

test("Update email") { client.update_email("alice", "newalice@example.com") }

test("Update username") { client.update_username("rubytest", "rubytest2") }

test("Update trust level (admin)") do
  admin.update_trust_level(2, level: 3)
end

test("Activate user (admin)") { admin.activate(2) }

test("Grant admin (admin)") { admin.grant_admin(2) }

test("Revoke admin (admin)") { admin.revoke_admin(2) }

test("Grant moderation (admin)") { admin.grant_moderation(2) }

test("Revoke moderation (admin)") { admin.revoke_moderation(2) }

test("Suspend user (admin)") do
  admin.suspend(2, Time.now.to_s, "test suspension")
end

test("Unsuspend user (admin)") { admin.unsuspend(2) }

# =============================================================================
# Search
# =============================================================================
puts "\n== Search =="

test("Search for 'test'") { client.search("test") }

# =============================================================================
# Groups
# =============================================================================
puts "\n== Groups =="

test("List groups") { client.groups }

test("Get group 'staff'") { client.group("staff") }

test("Get group members") { client.group_members("staff") }

test("Create group (admin)") do
  admin.create_group(name: "ruby-test-group")
end

# =============================================================================
# Badges
# =============================================================================
puts "\n== Badges =="

test("List badges") { admin.badges }

test("User badges") { client.user_badges("alice") }

# =============================================================================
# Notifications
# =============================================================================
puts "\n== Notifications =="

test("Get notifications") { client.notifications }

# =============================================================================
# Uploads
# =============================================================================
puts "\n== Uploads =="

require 'tempfile'
tmp = Tempfile.new(['test', '.txt'])
tmp.write("test upload from Ruby SDK")
tmp.close

test("Upload file") do
  client.upload_file(file: tmp.path)
end

tmp.unlink

# =============================================================================
# Private Messages
# =============================================================================
puts "\n== Private Messages =="

test("Create PM") do
  client.create_pm(
    title: "Ruby SDK PM Test",
    raw: "Test PM from Ruby discourse_api gem",
    target_recipients: "alice"
  )
end

test("Get user PMs") { client.private_messages("system") }

test("Get sent PMs") { client.sent_private_messages("system") }

# =============================================================================
# Invites
# =============================================================================
puts "\n== Invites =="

test("Create invite") do
  client.invite_user(email: "rubyinvite@example.com")
end

if topic_id
  test("Invite to topic") do
    client.invite_to_topic(topic_id, email: "rubyinvite2@example.com")
  end
end

# =============================================================================
# Admin
# =============================================================================
puts "\n== Admin =="

test("Get category latest topics") do
  client.category_latest_topics(category_slug: "general")
end

# =============================================================================
# Delete Operations
# =============================================================================
puts "\n== Delete Operations =="

if post_id
  test("Delete post #{post_id}") { client.delete_post(post_id) }
end

if topic_id
  test("Delete topic #{topic_id}") { client.delete_topic(topic_id) }
end

# =============================================================================
# Summary
# =============================================================================
puts
puts "=" * 60
puts "RESULTS: #{$passed} passed, #{$failed} failed out of #{$passed + $failed} tests"
if $errors.any?
  puts "\nFailed tests:"
  $errors.each { |name, err| puts "  - #{name}: #{err}" }
end
puts "=" * 60

exit($failed == 0 ? 0 : 1)
