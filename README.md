# Hello World API

A minimal HTTP API that returns a "hello world" response. Built with Node.js's
built-in `http` module, so it has no runtime dependencies.

## Requirements

- Node.js 18 or newer (uses the built-in `node:test` runner).

## Running

```bash
npm start
```

The server listens on `PORT` (default `3000`).

## Endpoints

| Method | Path     | Response                                |
| ------ | -------- | --------------------------------------- |
| GET    | `/hello` | `200 OK` — `{"message":"hello world"}`  |
| GET    | `/`      | `200 OK` — `{"message":"hello world"}`  |

Unsupported methods on these paths return `405 Method Not Allowed`; unknown
paths return `404 Not Found`. Responses use `Content-Type: application/json`.

### Example

```bash
curl http://localhost:3000/hello
# {"message":"hello world"}
```

## Testing

```bash
npm test
```
