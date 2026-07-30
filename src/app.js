'use strict';

const http = require('http');
const { helloHandler } = require('./routes/hello');

/**
 * Route an incoming request to the appropriate handler.
 *
 * Registers `GET /hello` and `GET /` for the hello world response.
 * Non-GET methods on a known path yield 405; unknown paths yield 404.
 *
 * @param {import('http').IncomingMessage} req
 * @param {import('http').ServerResponse} res
 */
function requestHandler(req, res) {
  const path = (req.url || '').split('?')[0];
  const isHelloPath = path === '/hello' || path === '/';

  if (isHelloPath && req.method === 'GET') {
    helloHandler(req, res);
    return;
  }

  if (isHelloPath) {
    res.statusCode = 405;
    res.setHeader('Content-Type', 'application/json');
    res.end(JSON.stringify({ error: 'Method Not Allowed' }));
    return;
  }

  res.statusCode = 404;
  res.setHeader('Content-Type', 'application/json');
  res.end(JSON.stringify({ error: 'Not Found' }));
}

/**
 * Build an HTTP server bound to the request handler.
 * @returns {import('http').Server}
 */
function createServer() {
  return http.createServer(requestHandler);
}

// Start the server only when run directly (not when imported by tests).
if (require.main === module) {
  const port = process.env.PORT || 3000;
  createServer().listen(port, () => {
    // eslint-disable-next-line no-console
    console.log(`Server listening on port ${port}`);
  });
}

module.exports = { createServer, requestHandler };
