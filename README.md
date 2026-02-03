# football-head-api
**Overview**

This project is a backend service written in Go that analyzes English Premier League (EPL) match data from the 2023–2024 season. The API ingests match-level data from a CSV file and exposes REST endpoints that compute league tables, team statistics, and match insights.
The project focuses on backend fundamentals such as data modeling, aggregation, clean API design, and performance-aware in-memory processing.

**Why Go?**

Go was chosen for its simplicity, strong standard library, and suitability for backend and data-processing workloads. The project avoids heavy frameworks to better understand HTTP servers, concurrency primitives, and clean separation of concerns.

**Features**

- Load and parse EPL match data from CSV
- Compute league tables and team statistics
- Expose RESTful endpoints using net/http

**API Endpoints**

- GET /health
- GET /teams
- GET /league/table
- GET /teams/{team}/stats
- GET /matches?team={team} (In progress...)

**Future Improvements**

- Persist data using a database
- Add caching for computed stats
- Add basic authentication
- Support multiple seasons
