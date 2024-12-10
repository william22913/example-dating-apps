FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s' -o /home/dating-apps ./main.go
RUN chmod og+w /home/dating-apps 

RUN mkdir -p /home/i18n
RUN mkdir -p /home/sql_mgirations

COPY i18n ./i18n
COPY sql_migrations ./sql_migrations