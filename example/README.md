# srv example

## Run server

```bash
go run ./cmd/srv/main.go
```

## Test GRPC Gateway

```bash
curl http://localhost:8080/api/v1/today
```

The output will be like:
```json
{"today":"2025-07-22", "format":"2006-01-02"}
```