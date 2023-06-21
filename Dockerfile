FROM golang:1.19 AS builder

ARG GIT_TOKEN
ENV GOPRIVATE=github.com/Allen-Career-Institute/
ENV GOSUMDB=off


COPY . /src
WORKDIR /src


RUN apt-get update && apt-get install git
ARG GIT_TOKEN
RUN git config --global url."https://$GIT_TOKEN:@github.com/".insteadOf "https://github.com/"

#Build
RUN go mod download
RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

COPY ./configs/config.yaml /data/conf/config.yaml


RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
WORKDIR /app
EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./go-kratos-sample", "-conf", "/data/conf/config.yaml"]
