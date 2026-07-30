'use strict';

const test = require('node:test');
const assert = require('node:assert');
const http = require('node:http');
const { createServer } = require('../src/app');

/**
 * Issue an HTTP request against the given server and resolve with the response.
 */
function request(port, options) {
  return new Promise((resolve, reject) => {
    const req = http.request({ host: '127.0.0.1', port, ...options }, (res) => {
      let body = '';
      res.on('data', (chunk) => {
        body += chunk;
      });
      res.on('end', () => resolve({ status: res.statusCode, headers: res.headers, body }));
    });
    req.on('error', reject);
    req.end();
  });
}

async function withServer(fn) {
  const server = createServer();
  await new Promise((resolve) => server.listen(0, resolve));
  const { port } = server.address();
  try {
    await fn(port);
  } finally {
    await new Promise((resolve) => server.close(resolve));
  }
}

test('GET /hello returns hello world', async () => {
  await withServer(async (port) => {
    const res = await request(port, { path: '/hello', method: 'GET' });
    assert.strictEqual(res.status, 200);
    assert.match(res.headers['content-type'], /application\/json/);
    assert.deepStrictEqual(JSON.parse(res.body), { message: 'hello world' });
  });
});

test('GET / returns hello world', async () => {
  await withServer(async (port) => {
    const res = await request(port, { path: '/', method: 'GET' });
    assert.strictEqual(res.status, 200);
    assert.deepStrictEqual(JSON.parse(res.body), { message: 'hello world' });
  });
});

test('POST /hello is not allowed', async () => {
  await withServer(async (port) => {
    const res = await request(port, { path: '/hello', method: 'POST' });
    assert.strictEqual(res.status, 405);
  });
});
