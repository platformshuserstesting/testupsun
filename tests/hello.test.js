'use strict';

const test = require('node:test');
const assert = require('node:assert');
const { createApp } = require('../src/app');

/**
 * Start the app on an ephemeral port, run `fn` with the base URL, then close.
 * @param {(baseUrl: string) => Promise<void>} fn
 */
function withServer(fn) {
  return new Promise((resolve, reject) => {
    const server = createApp();
    server.listen(0, async () => {
      const { port } = server.address();
      try {
        await fn(`http://127.0.0.1:${port}`);
        resolve();
      } catch (err) {
        reject(err);
      } finally {
        server.close();
      }
    });
  });
}

test('GET /hello returns 200 with hello world JSON', async () => {
  await withServer(async (baseUrl) => {
    const res = await fetch(`${baseUrl}/hello`);
    assert.strictEqual(res.status, 200);
    assert.strictEqual(res.headers.get('content-type'), 'application/json');
    const body = await res.json();
    assert.deepStrictEqual(body, { message: 'hello world' });
  });
});

test('unknown route returns 404', async () => {
  await withServer(async (baseUrl) => {
    const res = await fetch(`${baseUrl}/does-not-exist`);
    assert.strictEqual(res.status, 404);
  });
});
