# Billing design: what UniBee does, what we take from Kill Bill

Billing only. Sources were read-only clones taken 2026-10-31: `UniBee-Billing/unibee` @ `bcdd257`, `unibee-api` @ `0f5e55d`, `killbill/killbill` @ `52e2f55`, `killbill/killbill-commons` @ `9f64499`. Runtime evidence is from UniBee `unibee/api:v1.9.0`, run locally with its own docker-compose (db, redis, api; `ENV=local` so its test clock works), with one merchant, a monthly 10.00 USD plan and a wire-transfer gateway.

## UniBee, observed running

| Question | What we saw (DB rows / API replies) |
|---|---|
| Creating a subscription | `POST /merchant/subscription/create_submit` writes the `subscription` row (status 1, pending), one `invoice` (period = `current_period_start..end`, `unique_id` = `invoice_id`) and one `payment`, all at once. `dunning_time` = period end − 3 days |
| Unpaid first invoice | After the test clock passed the period end, the subscription went to status 9 and the invoice to status 5 (cancelled). No renewal |
| Paid, clock past period end | One renewal invoice: `period_start` = previous `period_end`, status 2 (processing). The subscription went to status 7 and its `current_period_*` did **not** move until payment. The renewal invoice's `billing_cycle_anchor` was 0 and its `create_time` was wall-clock time, not test-clock time |
| Walking again later in the same period | No second invoice |
| Six test-clock walks fired at once | One renewal invoice. The per-subscription Redis lock held in this test |
| Moving the clock to the same time or back | Rejected: "The Subscription Has Walk To The TestClock Exceed The New One" |
| Invoice lines | One JSON text column `invoice.lines`. Each line has `amount`, `tax`, `discountAmount`, `taxPercentage` in cents. `periodStart`/`periodEnd` inside the line were `0`; the period appears only in the description text ("1 * Monthly 10 (2026-11-30-2026-12-31)") |

## UniBee, confirmed in source

- **Scheduling.** Two `gcron.AddSingleton(ctx, "@every 10s", …)` jobs (`internal/cronjob/cronjob.go:30,38`).
- **Per-subscription lock.** `TryLock(key, 60)` is `SET key 1 NX EX 60` (`utility/jedis.go:9`), taken in `billingcycle/cycle/cycle.go:62`. `ReleaseLock` is a plain `DEL` with no owner check (`jedis.go:16`). If the work outlives 60 s, a worker can delete another worker's lock. That was not reproduced here: the six concurrent walks finished well inside 60 s.
- **Whether to invoice** is decided in code from the latest invoice. `needInvoiceGenerate` is a local flag in `cycle.go:102–118` (the findings file called it a function). The only unique key on `invoice` is `unique_id`, so nothing in the database stops a second invoice for the same subscription and period.
- **Test clock.** Stored per subscription and forward-only. Capped at 60 days (2 years for yearly plans). Refused when `ENV=prod` (`billingcycle/testclock/testclock.go`). The `/system` route group has CORS and response middleware but no auth (`internal/cmd/cmd.go:378`); only the prod check guards the clock.

## What replaces each part (Kill Bill ideas, our code)

| UniBee | Ours | Kill Bill source of the idea (confirmed) |
|---|---|---|
| Invoice lines as JSON text, cents, period in a string | `invoice_items` rows typed `recurring` / `discount` / `tax`; discounts and taxes point at their base item via `linked_item_id`; `service_start`/`service_end` columns; `DECIMAL(12,2)` + `shopspring/decimal` | `invoice/.../invoice/ddl.sql`: `invoice_items` with `type`, `start_date`, `end_date`, `linked_item_id` |
| "Look at the latest invoice" decides | Compute every period that should be billed up to `asOf`, subtract the periods in `billing_records`, invoice only the difference. `UNIQUE (subscription_id, period_start)` makes a duplicate impossible even if the code is wrong | `FixedAndRecurringInvoiceItemGenerator.java:118` ("proposed items are ALL items since beginning of time"), then merged with existing items |
| No bound | Refuse more than 31 periods per subscription per run, and log it | `safetyBounds` (`…Generator.java:477`); `maxDailyNumberOfItemsSafetyBound` default 15 (`InvoiceConfig.java:104`) |
| Preview endpoint only | `billing.Run(ctx, asOf, dryRun)`: a dry run computes the same items and writes only the `billing_runs` row | `InvoiceDispatcher` dry-run arguments (per the findings file, K7) |
| Redis lock with expiry | A MySQL row lock `FOR UPDATE SKIP LOCKED` on the subscription, held only for its transaction; retry once on 1213/1205 | Kill Bill takes a DB-backed per-account lock (K6). `SKIP LOCKED` is our choice, not Kill Bill's |
| Redis streams (`go-redismq`) | `job_queue` table: `Enqueue` (dedupe key), `Claim` (`SKIP LOCKED`), `Done`, `Fail` (retry, then `failed`) | Queue tables with `processing_state`, `processing_owner`, `processing_available_date`, `error_count` (`util/.../ddl.sql:171–174`). Kill Bill claims with `UPDATE … WHERE processing_owner IS NULL` (`killbill-commons queue/.../QueueSqlDao.sql.stg:172`), not `SKIP LOCKED` |
| `merchant_operation_log` | `audit_log` row per billing change, written in the same transaction | `audit_log` (K9) |
| Per-subscription test clock | One global clock in `company_settings.test_clock_date`, dev-only, forward-only | `ClockMock` behind a test endpoint (K10) |

## Not observed

- UniBee's behaviour when work outlives the 60 s lock (the race in U7).
- Kill Bill was not run.
- REA process capture is unavailable on native Windows (`capability_unavailable`), so the runtime evidence is curl plus SQL output (kept outside the repo), not REA evidence records.
