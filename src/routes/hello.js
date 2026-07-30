'use strict';

/**
 * Handler for the "hello world" endpoint.
 *
 * Responds with HTTP 200 and a JSON body `{ "message": "hello world" }`.
 * The handler is pure and dependency-free so it is trivial to test.
 *
 * @param {import('http').IncomingMessage} req
 * @param {import('http').ServerResponse} res
 */
function helloHandler(req, res) {
  res.statusCode = 200;
  res.setHeader('Content-Type', 'application/json');
  res.end(JSON.stringify({ message: 'hello world' }));
}

module.exports = { helloHandler };
