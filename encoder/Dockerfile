# docker run -d --network=host transkoder:latest

FROM golang:1.16 as builder
WORKDIR /encoder
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY encoder.go .
RUN CGO_ENABLED=0 go build -o encoder-server


# FROM jrottenberg/ffmpeg:4.1-alpine
# # copy ffmpeg bins
# COPY --from=0 / /
# # FROM <extension of builder>



FROM alpine as production
WORKDIR /encoder
COPY --from=builder /encoder/encoder-server /
COPY --from=mwader/static-ffmpeg:4.4.1 /ffmpeg /usr/local/bin/
COPY /util/Roboto-Regular.ttf .
EXPOSE 8080
CMD ["/encoder-server"]


