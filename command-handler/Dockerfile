FROM golang:1.19
WORKDIR /shared-dependencies
COPY shared-dependencies/ ./
RUN go mod download
WORKDIR /app
COPY command-handler/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /command-handler

CMD ["/command-handler"]