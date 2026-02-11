#!/usr/bin/env node
/**
 * Test DTU Discourse server with the JavaScript discourse-api npm package.
 */

const Discourse = require('discourse-api');

const HOST = 'http://localhost:4200';
const API_KEY = 'test_api_key';
const API_USER = 'system';
const ADMIN_KEY = 'admin_api_key';
const ADMIN_USER = 'admin';

const client = new Discourse(HOST, API_KEY, API_USER);
const admin = new Discourse(HOST, ADMIN_KEY, ADMIN_USER);

let passed = 0;
let failed = 0;
const errors = [];

function promisify(fn) {
  return new Promise((resolve, reject) => {
    fn((error, body, httpCode) => {
      if (error) {
        reject(error);
        return;
      }
      let parsed = body;
      if (typeof body === 'string') {
        try { parsed = JSON.parse(body); } catch(e) { /* keep as string */ }
      }
      if (httpCode && httpCode >= 400) {
        reject(new Error(`HTTP ${httpCode}: ${typeof parsed === 'object' ? JSON.stringify(parsed).slice(0, 200) : String(parsed).slice(0, 200)}`));
        return;
      }
      resolve(parsed);
    });
  });
}

async function test(name, fn) {
  try {
    const result = await fn();
    console.log(`  PASS: ${name}`);
    if (result != null) {
      const s = typeof result === 'object' ? JSON.stringify(result).slice(0, 200) : String(result).slice(0, 200);
      if (s.length > 10) console.log(`        -> ${s}...`);
    }
    passed++;
    return result;
  } catch(e) {
    console.log(`  FAIL: ${name}`);
    console.log(`        Error: ${e.message}`);
    failed++;
    errors.push([name, e.message]);
    return null;
  }
}

async function main() {
  // ==========================================================================
  // Categories
  // ==========================================================================
  console.log('\n== Categories ==');

  await test('List categories', () => promisify(cb => client.getCategories({}, cb)));

  await test('Create category', () =>
    promisify(cb => client.createCategory('JS SDK Category', 'AABBCC', 'FFFFFF', null, cb)));

  // ==========================================================================
  // Topics
  // ==========================================================================
  console.log('\n== Topics ==');

  let topicResult = await test('Create topic', () =>
    promisify(cb => client.createTopic('JS SDK Test Topic', 'Content from JS discourse-api', 1, cb)));

  let topicId = null;
  let postId = null;
  if (topicResult && typeof topicResult === 'object') {
    topicId = topicResult.topic_id;
    postId = topicResult.id;
    console.log(`        Created topic_id=${topicId}, post_id=${postId}`);
  }

  if (topicId) {
    await test(`Get topic ${topicId}`, () =>
      promisify(cb => client.getTopicAndReplies(topicId, cb)));

    await test(`Update topic ${topicId}`, () =>
      promisify(cb => client.updateTopic('-', topicId, 'Updated JS SDK Topic', 1, cb)));

    await test(`Delete topic ${topicId}`, () =>
      promisify(cb => client.deleteTopic(topicId, cb)));
  }

  // Create another topic for more tests
  topicResult = await test('Create topic 2', () =>
    promisify(cb => client.createTopic('JS SDK Test Topic 2', 'Second topic from JS', 1, cb)));

  topicId = topicResult && topicResult.topic_id;
  postId = topicResult && topicResult.id;
  if (topicId) console.log(`        Created topic_id=${topicId}`);

  // ==========================================================================
  // Posts
  // ==========================================================================
  console.log('\n== Posts ==');

  await test('Get latest posts', () =>
    promisify(cb => client.getLastPostId(cb)));

  if (topicId) {
    const reply = await test('Reply to topic', () =>
      promisify(cb => client.replyToTopic('Reply from JS SDK', topicId, cb)));

    let replyPostId = null;
    if (reply && typeof reply === 'object') {
      replyPostId = reply.id;
      console.log(`        Created reply post_id=${replyPostId}`);
    }

    if (replyPostId) {
      await test(`Get post ${replyPostId}`, () =>
        promisify(cb => client.getPost(replyPostId, cb)));

      await test(`Update post ${replyPostId}`, () =>
        promisify(cb => client.updatePost(replyPostId, 'Updated from JS SDK', 'JS edit', cb)));

      await test('Reply to post', () =>
        promisify(cb => client.replyToPost('Reply to specific post from JS', topicId, 1, cb)));
    }
  }

  // ==========================================================================
  // Users
  // ==========================================================================
  console.log('\n== Users ==');

  await test("Get user 'alice'", () =>
    promisify(cb => client.getUser('alice', cb)));

  await test('Create user', () =>
    promisify(cb => client.createUser('JS Test User', 'jstest@example.com', 'jstest', 'SecureP@ss123', true, cb)));

  await test('Filter users (admin)', () =>
    promisify(cb => admin.filterUsers('alice', cb)));

  await test('Get user activity', () =>
    promisify(cb => client.getUserActivity('alice', 0, cb)));

  await test('Fetch confirmation value (honeypot)', () =>
    promisify(cb => client.fetchConfirmationValue(cb)));

  await test('Activate user (admin)', () =>
    promisify(cb => admin.activateUser(2, 'alice', cb)));

  await test('Approve user (admin)', () =>
    promisify(cb => admin.approveUser(2, 'alice', cb)));

  // ==========================================================================
  // Search
  // ==========================================================================
  console.log('\n== Search ==');

  await test("Search for 'test'", () =>
    promisify(cb => client.search('test', cb)));

  await test("Search for user", () =>
    promisify(cb => client.searchForUser('alice', cb)));

  // ==========================================================================
  // Groups
  // ==========================================================================
  console.log('\n== Groups ==');

  await test('Get group members', () =>
    promisify(cb => client.getGroupMembers('staff', cb)));

  await test('Create group (admin)', () =>
    promisify(cb => admin.createGroup('js-test-group', 0, false, '', false, 1, false, '', true, cb)));

  // ==========================================================================
  // Private Messages
  // ==========================================================================
  console.log('\n== Private Messages ==');

  await test('Create PM', () =>
    promisify(cb => client.createPrivateMessage('JS SDK PM Test', 'PM from JS SDK', 'alice', cb)));

  await test('Get PMs', () =>
    promisify(cb => client.getPrivateMessages('system', cb)));

  await test('Get unread PMs', () =>
    promisify(cb => client.getUnreadPrivateMessages('system', cb)));

  // Login/logout
  console.log('\n== Session ==');

  await test('Login', () =>
    promisify(cb => client.login('alice', 'password123', cb)));

  // ==========================================================================
  // Category latest topics
  // ==========================================================================
  console.log('\n== Category Topics ==');

  await test('Get category latest topics', () =>
    promisify(cb => client.getCategoryLatestTopic('general', {}, cb)));

  // ==========================================================================
  // Cleanup
  // ==========================================================================
  if (topicId) {
    await test(`Delete topic ${topicId}`, () =>
      promisify(cb => client.deleteTopic(topicId, cb)));
  }

  // ==========================================================================
  // Summary
  // ==========================================================================
  console.log();
  console.log('='.repeat(60));
  console.log(`RESULTS: ${passed} passed, ${failed} failed out of ${passed + failed} tests`);
  if (errors.length > 0) {
    console.log('\nFailed tests:');
    errors.forEach(([name, err]) => console.log(`  - ${name}: ${err}`));
  }
  console.log('='.repeat(60));

  process.exit(failed === 0 ? 0 : 1);
}

main().catch(e => {
  console.error('Fatal error:', e);
  process.exit(1);
});
