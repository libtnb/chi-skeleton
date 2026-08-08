# chi-skeleton

Unlike [fiber-skeleton](https://github.com/libtnb/fiber-skeleton), this skeleton uses the lightest [chi](https://github.com/go-chi/chi) framework on top of net/http.

## Features

- **chi v5** on a hardened `http.Server` (timeouts, body/header limits) with a global middleware stack
- **Database-backed sessions** ([sessions](https://github.com/libtnb/sessions) on a rio-backed store), started for every request (health probes excluded)
- **Compile-time dependency injection** via [libtnb/wire](https://github.com/libtnb/wire): typed provider graphs, generated sequential constructors, rollback on initialization failure and reverse-order cleanup
- **Strongly-typed configuration** ([koanf](https://github.com/knadh/koanf)) with `APP_*` environment overrides, validated at startup
- **Graceful shutdown** on SIGINT/SIGTERM (drains requests and cron jobs) and **zero-downtime upgrades** on SIGHUP ([graceful](https://github.com/libtnb/graceful))
- **Structured logging** with [slog](https://pkg.go.dev/log/slog) on a rotating file writer ([logrotate](https://github.com/libtnb/logrotate)), stdout, or both; request logs share the same logger ([httplog](https://github.com/golang-cz/httplog))
- **Request binding + validation** ([chix](https://github.com/libtnb/chix) + [validator](https://github.com/libtnb/validator)) with a boolean rule DSL and i18n messages; each service holds its validator — no package-global state
- **Typed application errors** (`internal/pkg/apperr`): a closed set of error kinds maps to HTTP statuses in one place, so a module can add error codes without touching any shared file
- **Scheduled jobs** ([cron](https://github.com/libtnb/cron)) with panic recovery and overlap skipping; modules contribute jobs through `registry` without importing the boot wiring
- **rio + SQLite** ([rio](https://github.com/go-rio/rio), swap in MySQL/PostgreSQL freely) with reusable predeclared queries, a prepared-statement cache and versioned migrations ([migrate](https://github.com/go-rio/migrate)): schema as Go code, automatic rollbacks, `migrate status`/`rollback` commands
- **Code generator** (`cmd/gen`) that scaffolds a full CRUD module in one command
- **OpenAPI 3.1 docs from validate tags** — schemas and constraints generated from the validator rules, served with a Scalar UI at `/docs`
- **Tests included**: handler tests against mocked repos ([mockery](https://github.com/vektra/mockery)), data-layer tests on a real migrated SQLite, validate-tag linting, generated-graph tests, and an architecture test that fails CI when a module crosses another module's boundary

## Quick start

Requires Go 1.27rc2. Go 1.27 is still a release candidate; switch the module and
container image to the final release and rerun the complete test/build gate when
it becomes available.

```bash
git clone https://github.com/libtnb/chi-skeleton my-app && cd my-app
make init   # copies config/config.example.yml to config/config.yml
make run    # or `make dev` for hot reload via air
```

The API listens on `:3000` by default: `curl localhost:3000/users`.

## Design

* `cmd` stores the entry point of each application, one directory per binary (`app`, `cli`, `gen`)
* `config` stores the configuration files
* `docs` stores hand-written documentation; the OpenAPI document is generated at runtime
* `internal` stores the application code: one directory per business module plus the shared layers below
* `internal/pkg` stores the contracts shared by every module (transport helpers, apperr, event bus, registry, job)
* `mocks` stores the generated mocks, one package per module (`mocks/user/biz`, `mocks/order/biz`)
* `storage` stores files generated while the application runs (logs, sessions, the SQLite database)
* `web` stores the front-end code of the application
* go.mod and go.sum manage dependencies — including the pinned `tool` directives (wire and mockery)

Each business module (`internal/user`, `internal/order`, ...) follows the three-layer design of [Kratos](https://go-kratos.dev/):

* **biz** holds domain models, repository interfaces and **usecases** — transport-independent business logic
* **data** implements the repositories against the database
* **service** adapts HTTP: binds/validates requests, delegates to usecases, shapes responses

Because usecases are transport-independent, the HTTP handlers, the CLI commands (each module's `service/command.go`) and the cron jobs all share the same business logic instead of each talking to the database on their own.

Wiring follows a typed contribution model. Each module exposes a `Module` from a
`//go:build wireinject` file, contributes routes, CLI commands, jobs,
subscriptions or health checks to named slice types, and explicitly exports the
bindings its parent may use. `internal/app/wire.go` combines the modules and
declares separate App and CLI injectors. `make generate` commits the ordinary Go
constructor code in `internal/app/wire_gen.go`; do not edit that generated file.

The boundaries are enforced, not aspirational: `TestModuleBoundaries` (`internal/app/arch_test.go`) parses the import graph and fails when a module reaches another module past its `biz` package, or any module imports the composition layers. Cross-module needs are expressed as interfaces in the consumer's biz package (see `order/biz.Users`) and adapted over the other module's public usecase in `data` — swap that adapter for an RPC client and the module splits into a service without touching its business logic.

## Configuration

`config/config.yml` is loaded first (override the path with `APP_CONFIG`), then any `APP_*` environment variable wins over the file. A double underscore separates nesting levels:

```bash
APP_HTTP__ADDRESS=:8080 APP_LOG__OUTPUT=stdout ./app
```

Configuration is parsed into a struct and validated at startup — a missing key or a bad value fails fast instead of panicking mid-request.

## Scheduled jobs

Add a job where it belongs — in the module that owns it: one `job.Fn`
(`internal/pkg/job`) provider per job, contributed from the module with
`Contribute[registry.Jobs](NewJob)` and exported as `registry.Jobs` (see
`bootstrap.Heartbeat` for the shape). Specs support an optional seconds field,
`@every 30s` descriptors and per-entry timezones. Jobs receive a
`context.Context` that is cancelled on shutdown; panics are recovered and
overlapping runs are skipped.

## Code generation

```bash
make gen name=article    # or: go run ./cmd/gen article
```

generates the biz entity + repo interface, data repository, service handlers,
typed route documentation, request structs, migration and Wire module for a new
module. Add the generated `<module>.Module` to `ApplicationModule.Include`, then
run `make generate` and `make gen-check`.

## Development

```bash
make help       # list all targets
make generate   # regenerate wire constructors and mocks
make wire-check # verify committed wire output is current
make lint       # golangci-lint
make test       # go test -race with coverage
make build      # static binaries in bin/ with the version injected
```

A `Dockerfile` is included; mount `config/` and `storage/` when running.

## OpenAPI documentation

Every documented endpoint attaches a generic `transport.Describe[Request,
Response]` callback to its route contribution. Schemas, parameters and
constraints are generated from the same `validate` tags that enforce them
([validator/contrib/openapi](https://github.com/libtnb/validator/tree/main/contrib/openapi)) — `min:3 && max:255` becomes `minLength`/`maxLength`, `in:a,b` becomes an enum, and the two can never drift apart. Use
`DescribeNoBody[Request]` when a successful response has no content. With
`http.docs: true` the app serves the OpenAPI 3.1 document at `/openapi.json` and
a [Scalar](https://github.com/scalar/scalar) UI at `/docs`.

## Database queries

Stable Rio query shapes are declared once at package level with `Must` and
reused concurrently. Placeholder values are passed to terminal operations:

```go
var userByIDQuery = rio.From[User]().Where("id = ?").Must()

_, err := userByIDQuery.DeleteAll(ctx, db, id)
```

## Observability

- `/healthz` (liveness) and `/readyz` (readiness, pings the DB) are wired for containers and load balancers; the Dockerfile ships a matching `HEALTHCHECK`.
- Access logs and application logs share one slog logger; every response carries an `X-Request-Id` header and the id is attached to the access log record.
- Set `http.debug_address` (e.g. `127.0.0.1:6060`) to serve `net/http/pprof` and `expvar` on a **separate private port** — profiling in production without exposing it on the API port.
- 404/405 and other framework-level errors leave as JSON in the same shape as the API; 5xx details go to the log, not the client.

## Error model

A usecase creates client-facing errors through `internal/pkg/apperr`:

```go
apperr.Conflict("user.name_taken", "name already taken").In("user").Wrap(ErrNameTaken)
```

The **kind** (conflict, not_found, invalid, ...) is a closed set that `transport.ErrorFrom` maps to an HTTP status; the **code** and public message travel to the client; everything else — stack trace, domain, attributes — goes to the log. Adding a module adds codes, never a new case in shared code. Errors without a kind are unexpected: the client sees a bare 500 and the details stay in the log.

## Serving a frontend

Put your built frontend under `web/` and serve it from the router in `internal/server/server.go` (`NewRouter`):

```go
r.Handle("/*", http.FileServer(http.Dir("./web/dist")))
```

## Graceful lifecycle

| Signal | Behavior |
|---|---|
| SIGINT / SIGTERM | stop accepting connections, drain in-flight requests and cron jobs (30s cap), then close sessions, DB and log writer in reverse dependency order |
| SIGHUP (non-Windows) | zero-downtime binary upgrade via [graceful](https://github.com/libtnb/graceful) |

## Credits

The development of this project refers to the following projects, I would like to express my gratitude:

* [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
* [Kratos](https://go-kratos.dev/)
* [Goravel](https://github.com/goravel/goravel)
* [GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)
* [gin-layout](https://github.com/wannanbigpig/gin-layout)
