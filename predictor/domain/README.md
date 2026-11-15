# Domain API and Worker

Domain contains the user-facing backend. The API handles accounts, company profiles, source links, and chat. The worker polls configured sources, extracts readable text from HTML and text-based PDFs, and forwards it to ESB. Raw material text is not persisted by Domain.

## Local development

Copy `.env.example` to `.env` in the repository root and set the required keys. Start the dependencies and Domain services with:

```sh
docker compose --env-file .env up --build domain-api domain-worker
```

The API listens on `http://localhost:8080`, health check is `GET /health`, and Swagger UI is `/swagger` in development.

## Company profiles

Create a company with an extensible JSON profile, for example:

```json
{
  "profile": {
    "name": "Example Company",
    "industry": "Logistics",
    "markets": ["Northern Europe", "Central Asia"],
    "priorities": ["delivery reliability", "cost control"]
  }
}
```

The profile schema is owned by each deployment and can evolve without a fixed dataset. Domain stores it as JSONB and sends each revision to the profile agent through ESB.

## Responsibilities

- User registration, login, and JWT issuance
- One company per user, with editable JSON profile and revision checks
- CRUD for recurring news sources and one-time material URLs
- Background polling and extraction of HTML, RSS, Atom, sitemap, and text-based PDF materials
- Retry-safe delivery of company profiles and extracted material to ESB
- Chat history and analysis responses read from ESB
