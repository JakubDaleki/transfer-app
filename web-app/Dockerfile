FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
WORKDIR /app/cmd/service
RUN CGO_ENABLED=0 GOOS=linux go build -o /web-app
EXPOSE 8000

CMD ["/web-app"]