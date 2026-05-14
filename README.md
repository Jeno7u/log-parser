# log-parser-service
A Go microservice for parsing log folders from `data/`, aggregating topology, storing parsed data in PostgreSQL, and exposing a small REST API.

---

## Requirements

- Go 1.22+
- Docker Compose

---

## Run the Service

1. Copy the environment template:

```bash
cp .env.template .env
```

2. Start the stack:

```bash
docker compose up -d --build
```

3. Send a parse request with a path inside `data/`.

4. Use the returned `log_id` to query the API.

---


## API

Parse folder:

```bash
curl -X POST http://localhost:8081/api/v1/parse/ \
	-H 'Content-Type: application/json' \
	-d '{"path":"data/sample"}'
```

Get log metadata:

```bash
curl http://localhost:8081/api/v1/log/<log_id>
```

Get topology:

```bash
curl http://localhost:8081/api/v1/topology/<log_id>
```

Get node details:

```bash
curl http://localhost:8081/api/v1/node/<node_id>
```

Get node ports:

```bash
curl http://localhost:8081/api/v1/port/<node_id>
```
---

## Postman

Import the collection from `postman/`:

- `postman/log-parser-service.postman_collection.json`
- `postman/log-parser-service.postman_environment.json`

---

## Notes

- Database migrations are applied automatically when the app starts.
- The db container uses `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DB` from `.env`, and the app derives its `DATABASE_URL` from the same values.
- The input folder must contain `ibdiagnet2.db_csv` and `ibdiagnet2.sharp_an_info`. It can be easily adapted to use other names or find files by extension.
- Handler functions can be split into business logic and HTTP for better testability, but the current structure is fine.