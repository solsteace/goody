FROM golang:1.24.6-alpine AS base
WORKDIR /app
RUN mkdir ./account ./lib
COPY ./account/go.* ./account
COPY ./lib/go.*  ./lib
RUN go work init && go work use ./account ./lib


FROM alpine/curl:latest AS devtools
WORKDIR .
# goose: 
# - https://pressly.github.io/goose/installation/
# - https://github.com/pressly/goose
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh


FROM base AS devel
COPY --from=devtools /usr/local/bin /
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", "./account/.air.toml"]