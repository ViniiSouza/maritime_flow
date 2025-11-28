# Resource Manager API

REST API that exposes towers, vehicles and structures stored in PostgreSQL.

## Prerequisites

- Node.js >= 18
- PostgreSQL instance and credentials

## Setup

```bash
cd server
npm install
cp .env.example .env
# edit .env with your DATABASE_URL
npm run dev
```

The API will boot on `http://localhost:4000` by default.

## Endpoints

| Method | Path | Description |
| --- | --- | --- |
| GET | `/api/towers` | List towers |
| GET | `/api/towers/:id` | Fetch a tower |
| POST | `/api/towers` | Create a tower (`name`, `latitude`, `longitude`, `is_leader`) |
| DELETE | `/api/towers/:id` | Remove a tower |
| GET | `/api/vehicles` | List vehicles |
| GET | `/api/vehicles/:id` | Fetch a vehicle |
| POST | `/api/vehicles` | Create a vehicle (`name`, `type`, `latitude`, `longitude`) |
| DELETE | `/api/vehicles/:id` | Remove a vehicle |
| GET | `/api/structures` | List structures |
| GET | `/api/structures/:id` | Fetch a structure |
| POST | `/api/structures` | Create a structure (`name`, `type`, `latitude`, `longitude`) |
| DELETE | `/api/structures/:id` | Remove a structure (slots are deleted together) |

## Example table definitions

```sql
CREATE TABLE towers (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  latitude NUMERIC NOT NULL,
  longitude NUMERIC NOT NULL,
  is_leader BOOLEAN DEFAULT FALSE
);

CREATE TABLE vehicles (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('Helicopter', 'Ship')),
  latitude NUMERIC NOT NULL,
  longitude NUMERIC NOT NULL
);

CREATE TABLE structures (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('Platform', 'Central')),
  latitude NUMERIC NOT NULL,
  longitude NUMERIC NOT NULL
);
```

Only issue `DELETE` requests against entities you intend to remove from both the API and UI.
