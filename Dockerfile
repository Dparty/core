FROM golang:1.20 as build-stage
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /main

FROM ubuntu:22.04 as production-stage
COPY --from=build-stage /main /main
EXPOSE 8080
CMD [ "/main" ]
