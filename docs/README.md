# docs

API documentation is generated automatically from the route declarations — no
separate codegen step. When `http.docs` is enabled (see
`config/config.example.yml`), the running server serves:

- `GET /openapi.json` — the OpenAPI 3.1 document, built from every endpoint's
  generic request/response documentation callback and its `validate` tags
  (see `internal/pkg/transport/endpoint.go`).
- `GET /docs` — a browsable UI for that document.

To document a new endpoint, attach
`transport.Describe[Request, Response](status)` or
`transport.DescribeNoBody[Request](status)` to its route contribution (for
example `internal/user/service/route.go`). An endpoint with a nil `Document`
callback stays out of the document.

This directory also holds any additional hand-written documents.
