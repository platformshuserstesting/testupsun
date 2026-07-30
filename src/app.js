'use strict';

const http = require('http');
const hello = require('./routes/hello');

/**
 * Route table mapping `METHOD PATH` to a handler function.
 * @type {Record<string, (req: import('http').IncomingMessage, res: import('http').ServerResponse) => void>}
 */
const routes = {
  'GET /hello': hello,
};

/**
 * Request dispatcher: looks up a handler for the request method/path and
 * invokes it, falling back to a 404 JSON response for unknown routes.
 *
 * @param {import('http').IncomingMessage} req
 * @param {import('http').ServerResponse} res
 */
function requestHandler(req, res) {
  const pathname = new URL(req.url, `http://${req.headers.host || 'localhost'}`).pathname;
  const handler = routes[`${req.method} ${pathname}`];

  if (handler) {
    handler(req, res);
    return;
  }

  res.writeHead(404, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ message: 'not found' }));
}

/**
 * Create (but do not start) the HTTP server.
 * @returns {import('http').Server}
 */
function createApp() {
  return http.createServer(requestHandler);
}

module.exports = { createApp, requestHandler };

// Start the server when run directly (e.g. `node src/app.js`).
if (require.main === module) {
  const port = process.env.PORT || 3000;
  createApp().listen(port, () => {
    // eslint-disable-next-line no-console
    console.log(`Server listening on http://localhost:${port}`);
  });
}
