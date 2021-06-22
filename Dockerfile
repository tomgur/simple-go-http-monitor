FROM golang:1.12.0-alpine3.9
RUN apk add --update --no-cache curl bash git
RUN mkdir /monitor
ADD . /monitor
WORKDIR /monitor
RUN go build -o monitor .
CMD ["/monitor/monitor"]
