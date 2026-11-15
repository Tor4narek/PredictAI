# PredictAI

PredictAI is a backend platform for company-focused news intelligence. It collects materials from configured public sources, extracts article text, builds structured facts, and lets an AI consultant answer questions using cited evidence and company context.

## Architecture

- **Domain API (.NET 10)** provides registration, JWT authentication, company profiles, source management, one-time material links, and chat endpoints.
- **Domain Worker (.NET 10)** polls RSS, Atom, sitemap, and HTML sources, parses HTML and text-based PDF materials, and submits extracted text to the integration layer.
- **ESB (Go)** manages asynchronous agent calls, chat history, retries, data access, and WebSocket streams between the API and AI runtime.
- **Hermes agents** process company profiles, extract facts from materials, calculate metrics, and provide the conversational interface.
- **PostgreSQL** stores users, company profile JSON, source metadata, call state, and conversation data. **ClickHouse** stores news facts and analytics. **S3-compatible storage** stores generated company context.

The company profile is an extensible JSON object. PredictAI does not require a fixed external company dataset.

## Run locally

Requirements: Docker Compose, .NET 10 SDK, Go 1.24, and an API key for the configured LLM provider.

1. Copy `.env.example` to `.env` and replace the development credentials with local values.
2. Start the databases and object storage, then apply their initial migrations:

   ```sh
   docker compose --env-file .env up -d postgres domain-postgres clickhouse minio
   docker compose --env-file .env --profile migrations run --rm atlas-postgres
   docker compose --env-file .env --profile migrations run --rm atlas-clickhouse
   ```

3. Start the full backend:

   ```sh
   docker compose --env-file .env up --build
   ```

The Domain API listens on `http://localhost:8080`; Swagger is available locally at `/swagger`. Database containers and the API bind to loopback by default. Configure model names and provider credentials in `.env` before starting the agent services.

## Documentation

- [System architecture](predictor/docs/architecture.md)
- [Domain API and worker](predictor/domain/README.md)
- [ESB integration contract](predictor/docs/esb-contracts.md)
- [HTTP, WebSocket, and MCP API references](predictor/docs/api/README.md)
- [Configuration reference](predictor/docs/configuration.md)
- [Development and operations](predictor/docs/development.md)

API descriptions in `predictor/docs/api/` are generated from the Go implementation. Regenerate them with `make docs` after changing the ESB contract.
