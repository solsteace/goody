FROM golang:1.24.6-alpine AS base
WORKDIR /app

FROM base AS devel
WORKDIR ./account
RUN go install github.com/air-verse/air@latest
CMD ["air"]