FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /docker-gs-ping
EXPOSE 8080
# Run
CMD ["/docker-gs-ping"]