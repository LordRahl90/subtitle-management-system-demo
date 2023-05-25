FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o translations ./cmd/


FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/translations translations


EXPOSE 8080