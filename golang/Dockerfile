# Build stage: compile a static binary
FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/ip-app .

# Run stage: small image with only the binary
FROM alpine:3.24

WORKDIR /app

COPY --from=build /app/ip-app .

RUN adduser -D -H appuser
USER appuser

EXPOSE 8000

CMD ["./ip-app"]
