'use strict';

/**
 * Handler for `GET /hello`.
 *
 * Responds with HTTP 200 and a small JSON payload `{"message": "hello world"}`.
 *
 * @param {import('http').IncomingMessage} req
 * @param {import('http').ServerResponse} res
 */
function hello(req, res) {
  res.writeHead(200, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ message: 'hello world' }));
}

module.exports = hello;
