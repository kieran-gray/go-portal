FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /portal ./cmd/web/*
RUN CGO_ENABLED=0 GOOS=linux go build -o /uploadMockFiles ./cmd/scripts/uploadMockFiles.go

EXPOSE 8080

USER www-data

CMD ["/portal"]