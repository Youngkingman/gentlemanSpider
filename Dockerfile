FROM golang:alpine as builder

WORKDIR /goapp

ENV GO111MODULE=on
# For Chinese Internet
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o gentlemanSpider .

FROM alpine:latest
WORKDIR /goapp

#Copy the cofig file
COPY . .

#Copy executable from builder
COPY --from=builder /goapp/gentlemanSpider .

CMD [ "./gentlemanSpider" ]