FROM golang:1.16.10-alpine3.13 as build
ENV GOPROXY https://goproxy.cn
WORKDIR /opt/
COPY . /opt/
RUN go build .


FROM alpine:3.14.0 as prod
WORKDIR /eat-and-go
COPY --from=build /opt/eat-and-go .
RUN mkdir -p /var/log/eat-and-go/

EXPOSE 8004
ENTRYPOINT ["./eat-and-go"]
