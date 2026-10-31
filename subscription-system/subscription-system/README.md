# Brightloop Subscription Management System

## Run it

```sh
docker compose up -d --wait   # mysql, mailhog, api (seeds, then serves :8000), web (:5173)
```

- App: http://localhost:5173 · API docs: http://localhost:8000/swagger · Emails: http://localhost:8025 (MailHog)
- Seeded logins (password `Brightloop#2026`): `admin@brightloop.test`, `staff@brightloop.test`, `customer@brightloop.test`
- The schema is applied only when the MySQL volume is first created. To start from a clean database: `docker compose rm -sfv mysql && docker compose up -d --wait`.

Frontend checks: `cd frontend && npm install && npm test && npm run build`.

## API conventions

Success `{ "data", "meta" }` (meta only on lists: `page`, `pageSize` up to 100, `total`); error `{ "error": { "code", "message", "fields" } }`
with `VALIDATION_ERROR` 422, `UNAUTHENTICATED` 401, `FORBIDDEN` 403, `NOT_FOUND` 404 (also for another customer's record), `CONFLICT` 409, `INTERNAL_ERROR` 500.
MySQL 1062 → 409, 3819 → 422, 1213/1205 → retried once by `platform.Tx`, then 409. Money is a decimal string; dates `YYYY-MM-DD`; timestamps UTC.
Roles are enforced per route group in `backend/internal/cmd/routes.go`. Platform contracts: `clock.Today`, `sequence.Next`, `jobs.Enqueue`, `mail.Enqueue`, `audit.Log`, `audit.Update`, `settings.Get`, `auth.CustomerID`.

## Run the backend tests

Needs Docker. Tests run against real MySQL 8 in `golang:1.22`:

```sh
docker compose up -d --wait        # MySQL 8.0 on :3306 (root/root)
./gotest.sh go test ./...          # each test package creates its own database
```

Schema checks: apply `backend/db/migrations/000001_init.up.sql` to an empty database, then `DB=<name> bash backend/db/test_constraints.sh`.

## Billing

`billing.Run(ctx, asOf, dryRun)` (`backend/internal/logic/billing`) is called by the nightly `gcron` job in `backend/internal/cronjob`. Design notes and evidence: [docs/billing-design.md](docs/billing-design.md).

## Credits

The project layout follows [UniBee](https://github.com/UniBee-Billing/unibee) (AGPL-3.0) and the billing design follows [Kill Bill](https://github.com/killbill/killbill) (Apache-2.0). Both were used as references only; this code was written from scratch and copies no source from either project.
