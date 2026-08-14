FROM golang:1.26-alpine AS build

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

RUN GOBIN=/tmp/bin go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -tags=seeder -o /out/seeder ./src/database/seeders/...

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata postgresql-client && addgroup -S app && adduser -S app -G app

COPY --from=build /out/api /app/api
COPY --from=build /out/seeder /app/seeder
COPY --from=build /tmp/bin/migrate /usr/local/bin/migrate
COPY src/database/migrations /app/migrations
COPY docker/app-entrypoint.sh /app/app-entrypoint.sh

RUN chmod +x /app/api /app/seeder /app/app-entrypoint.sh && chown -R app:app /app

USER app

EXPOSE 8080

ENTRYPOINT ["/app/app-entrypoint.sh"]
