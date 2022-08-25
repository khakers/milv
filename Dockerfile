FROM golang:1.19-alpine as builder

ENV BASE_APP_DIR /go/src/github.com/khakers/milv
WORKDIR ${BASE_APP_DIR}

COPY ./ ${BASE_APP_DIR}/

RUN go build -v -o main .
RUN mkdir /app && mv ./main /app/main

FROM alpine:3.16

#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && apk add bash

COPY --from=builder /app /app

ENTRYPOINT ["/app/main"]