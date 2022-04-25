# docker build -t swr.cn-east-3.myhuaweicloud.com/wx-xiaoyang-public/golang:runner-demacia .
FROM swr.cn-east-3.myhuaweicloud.com/wx-xiaoyang-public/golang:runner-base

COPY .gitlab/* /home/
COPY go.sum /home/
COPY go.mod /home/
WORKDIR /home
RUN go mod download