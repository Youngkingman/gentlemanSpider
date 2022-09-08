FROM golang:alpine as builder

WORKDIR /go/src/app

ENV GO111MODULE=on
# For Chinese Internet
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o gentlemanSpider .

FROM alpine:latest
WORKDIR /root/

#Copy the cofig file
COPY . .

#Copy executable from builder
COPY --from=builder /go/src/app/gentlemanSpider .

EXPOSE 8080

CMD [ "./gentlemanSpider" ]