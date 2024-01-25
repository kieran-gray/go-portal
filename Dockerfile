FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /portal ./cmd/web/*
RUN CGO_ENABLED=0 GOOS=linux go build -o /uploadMockFiles ./cmd/scripts/uploadMockFiles.go

FROM gcr.io/distroless/base-debian12 as prod

WORKDIR /

COPY /ui /ui
COPY /config /config
COPY --from=builder /portal /portal
COPY --from=builder /uploadMockFiles /uploadMockFiles

EXPOSE 8080

USER nonroot:nonroot

CMD ["/portal"]

FROM prod as dev

COPY /data /data

USER root

CMD ["/portal"]