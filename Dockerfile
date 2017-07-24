FROM golang:1.7.3 as builder
WORKDIR /go/src/
COPY ./src .
RUN go get -v github.com/zenazn/goji/graceful github.com/rs/cors github.com/dvsekhvalnov/jose2go goji.io goji.io/pat
RUN ls
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app/server app/.
RUN cd app && ls

FROM alpine:latest
EXPOSE 8000
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/server /go/src/app/server.crt /go/src/app/server.key ./
CMD ["./server"]
