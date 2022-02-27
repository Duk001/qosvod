FROM golang:1.16 as builder
WORKDIR /backend
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY main.go .
COPY pkg pkg
RUN ls -r
RUN CGO_ENABLED=0 go build -o backend-server

FROM alpine as production
WORKDIR /backend
COPY --from=builder /backend/backend-server /
EXPOSE 8080
CMD ["/backend-server"]
