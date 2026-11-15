# Company profile and metric contract

The company profile is Markdown at `companies/{company_id}/profile.md` in S3. It is generated from the extensible JSON profile supplied by Domain. ESB is the only component that accesses S3, PostgreSQL, and ClickHouse.

## Metric settings

PostgreSQL is the source of truth for company-specific metric settings. Each row stores a database-generated opaque `metric_id`, `company_id`, `type_code`, `name`, `description`, `metric_type`, `characteristic`, `unit`, and structured `parameters`.

At most one row per `(company_id, type_code)` is stored. `metric_type` must match the PostgreSQL `company_metric_type` enum. Use a stable descriptive `type_code` for each measurable business indicator, a concise unit when applicable, and `parameters.source_fields` to list the exact profile keys that justify the setting. Do not include field values in `parameters`. The profile may contain any application-defined fields, so derive indicators only from data actually supplied. Do not create placeholder rows.

## Profile save and retry

`profile_revision` is a positive, monotonically increasing revision supplied by Domain. `save_company_profile_and_metrics` takes the complete Markdown profile and complete metric snapshot for that revision. ESB replaces the company's metric settings in PostgreSQL, then writes the stable S3 profile key. Omitted metric rows are removed; an empty array removes all rows. An older revision cannot overwrite newer data. Repeating the same revision and payload is idempotent, including repair after a partial S3 failure.

On a transient failure, retry once with the exact same company ID, revision, Markdown, and metric settings. Do not invent IDs or increment the revision for a retry.

## Analytics boundary

ClickHouse stores news, company-independent facts, fact-to-company links, observations, and generated analyses. An observation references a database-generated `metric_id`; it is opaque and must not encode company or metric data. News owns publication time, URL, and material kind used by analytics. Totals are calculated when requested and are not persisted.

The Domain request to `POST /call?agent=valera` uses `companyId`, `profileRevision`, and `company`. ESB normalizes it before invoking Hermes as `company_id`, `company_data`, and `profile_revision`.
